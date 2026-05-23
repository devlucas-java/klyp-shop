package service

import (
	"context"
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/product"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/pkg/pagination"
)

type ProductService struct {
	log               *logger.Logger
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
	sellerRepository  repository.SellerRepository
	cartRepository    repository.ShoppingCartRepository
	productMapper     *mapper.ProductMaper
	productPolicy     *policy.ProductPolicy
}

func NewProductService(
	log *logger.Logger,
	productRepository repository.ProductRepository,
	userRepository repository.UserRepository,
	sellerRepository repository.SellerRepository,
	productMapper *mapper.ProductMaper,
	cartRepository repository.ShoppingCartRepository,
) *ProductService {
	return &ProductService{
		log:               log,
		productRepository: productRepository,
		userRepository:    userRepository,
		sellerRepository:  sellerRepository,
		cartRepository:    cartRepository,
		productMapper:     productMapper,
		productPolicy:     policy.NewProductPolicy(),
	}
}

func (s *ProductService) CreateProduct(auth *entity.User, req *product.CreateProduct) (*product.ProductResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, err
	}

	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	prod, err := entity.NewProduct(req.Name, req.Description, req.PriceBTC, req.Stock, req.Categories)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	prod.SellerID = user.Seller.ID

	saved, err := s.productRepository.Create(prod)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return s.productMapper.ProductToProductResponse(saved), nil
}

func (s *ProductService) UpdateProduct(auth *entity.User, req *product.UpdateProduct, productID id.UUID) (*product.ProductResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, err
	}

	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	prod, err := s.productRepository.FindByID(productID)
	if err != nil {
		return nil, err
	}

	if err := s.productPolicy.CanManage(prod, user.Seller.ID); err != nil {
		return nil, err
	}

	if err := prod.UpdateDetails(req.Name, req.Description, req.PriceBTC, req.Stock, req.Categories); err != nil {
		return nil, err
	}

	saved, err := s.productRepository.Updates(prod)
	if err != nil {
		return nil, err
	}

	return s.productMapper.ProductToProductResponse(saved), nil
}

func (s *ProductService) DeleteProduct(auth *entity.User, productID id.UUID) error {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return err
	}

	if err := user.EnsureSeller(); err != nil {
		return err
	}

	prod, err := s.productRepository.FindByID(productID)
	if err != nil {
		return err
	}

	if err := s.productPolicy.CanManage(prod, user.Seller.ID); err != nil {
		return err
	}

	affectedCarts, err := s.cartRepository.FindCartsByProductID(productID)
	if err != nil {
		s.log.Warnf("could not find carts with product %s: %v", productID, err)
	}

	if err := s.productRepository.DeleteByID(prod.ID); err != nil {
		return err
	}

	for _, c := range affectedCarts {
		updated, err := s.cartRepository.FindByID(c.ID)
		if err != nil {
			s.log.Warnf("could not reload cart %s after product delete: %v", c.ID, err)
			continue
		}
		updated.RecalculateTotal()
		if _, err := s.cartRepository.Save(updated); err != nil {
			s.log.Warnf("could not save recalculated cart %s: %v", c.ID, err)
		}
	}

	return nil
}

func (s *ProductService) GetProductByID(productID id.UUID) (*product.ProductResponse, error) {
	prod, err := s.productRepository.FindByID(productID)
	if err != nil {
		return nil, err
	}

	return s.productMapper.ProductToProductResponse(prod), nil
}

func (s *ProductService) ListProducts(page, size int) ([]*product.ProductResponse, error) {
	products, err := s.productRepository.Search(page, size, "name", "", []string{})
	if err != nil {
		return nil, err
	}

	responses := make([]*product.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = s.productMapper.ProductToProductResponse(p)
	}
	return responses, nil
}

func (s *ProductService) GetProductsBySeller(sellerID id.UUID, page, size int) ([]*product.ProductResponse, error) {
	products, err := s.productRepository.FindBySellerID(sellerID, page, size)
	if err != nil {
		return nil, err
	}

	responses := make([]*product.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = s.productMapper.ProductToProductResponse(p)
	}
	return responses, nil
}

func (s *ProductService) SetTop10(ctx context.Context, auth *entity.User, productID id.UUID) ([]*product.ProductResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, err
	}

	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	prod, err := s.productRepository.FindByID(productID)
	if err != nil {
		return nil, err
	}

	if err := s.productPolicy.CanManage(prod, user.Seller.ID); err != nil {
		return nil, err
	}

	size, err := s.productRepository.CountTop10BySellerID(user.Seller.ID)
	if err != nil {
		return nil, err
	}

	if size >= 10 {
		return nil, apperrors.BadRequest(fmt.Sprintf("you already have %d products in the top 10 list", size), nil)
	}

	if err := prod.AddTop10(size); err != nil {
		return nil, err
	}

	if _, err := s.productRepository.Updates(prod); err != nil {
		return nil, apperrors.Internal(err)
	}

	return s.GetProductsBySeller(user.Seller.ID, 1, 10)
}

func (s *ProductService) SearchProducts(ctx context.Context, inputPagination pagination.InputPagination, categories []string) ([]*product.ProductResponse, error) {
	products, err := s.productRepository.Search(inputPagination.Page, inputPagination.Size, inputPagination.Search, inputPagination.Search, categories)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	responses := make([]*product.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = s.productMapper.ProductToProductResponse(p)
	}
	return responses, nil
}
