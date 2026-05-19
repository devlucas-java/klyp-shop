package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dproduct"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type ProductService struct {
	log               *logger.Logger
	productRepository repository.ProductRepository
	userRepository    repository.UserRepository
	sellerRepository  repository.SellerRepository
	productMapper     *mapper.ProductMaper
	productPolicy     *policy.ProductPolicy
}

func NewProductService(
	log *logger.Logger,
	productRepository repository.ProductRepository,
	userRepository repository.UserRepository,
	sellerRepository repository.SellerRepository,
	productMapper *mapper.ProductMaper,
) *ProductService {
	return &ProductService{
		log:               log,
		productRepository: productRepository,
		userRepository:    userRepository,
		sellerRepository:  sellerRepository,
		productMapper:     productMapper,
		productPolicy:     policy.NewProductPolicy(),
	}
}

func (s *ProductService) CreateProduct(auth *entity.User, req *dproduct.CreateProduct) (*dproduct.ProductResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}

	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	product, err := entity.NewProduct(req.Name, req.Description, req.PriceBTC, req.Stock, req.Categories)
	if err != nil {
		return nil, err
	}
	product.SellerID = user.Seller.ID

	saved, err := s.productRepository.Create(product)
	if err != nil {
		s.log.Errorf("failed to create product: %v", err)
		return nil, errors.ErrDatabase("failed to create product", err)
	}

	return s.productMapper.ProductToProductResponse(saved), nil
}

func (s *ProductService) UpdateProduct(auth *entity.User, req *dproduct.UpdateProduct, productID id.UUID) (*dproduct.ProductResponse, error) {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return nil, errors.ErrNotFound("User", err)
	}

	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return nil, errors.ErrNotFound("Product", err)
	}

	if err := s.productPolicy.CanManage(product, user.Seller.ID); err != nil {
		return nil, err
	}

	if err := product.UpdateDetails(req.Name, req.Description, req.PriceBTC, req.Stock, req.Categories); err != nil {
		return nil, err
	}

	saved, err := s.productRepository.Updates(product)
	if err != nil {
		s.log.Errorf("failed to update product: %v", err)
		return nil, errors.ErrDatabase("failed to update product", err)
	}

	return s.productMapper.ProductToProductResponse(saved), nil
}

func (s *ProductService) DeleteProduct(auth *entity.User, productID id.UUID) error {
	user, err := s.userRepository.FindByIDWithSeller(auth.ID)
	if err != nil {
		return errors.ErrNotFound("User", err)
	}

	if err := user.EnsureSeller(); err != nil {
		return err
	}

	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return errors.ErrNotFound("Product", err)
	}

	if err := s.productPolicy.CanManage(product, user.Seller.ID); err != nil {
		return err
	}

	if err := s.productRepository.DeleteByID(product.ID); err != nil {
		s.log.Errorf("failed to delete product: %v", err)
		return errors.ErrDatabase("failed to delete product", err)
	}

	return nil
}

func (s *ProductService) GetProductByID(productID id.UUID) (*dproduct.ProductResponse, error) {
	product, err := s.productRepository.FindByID(productID)
	if err != nil {
		return nil, errors.ErrNotFound("Product", err)
	}

	return s.productMapper.ProductToProductResponse(product), nil
}

func (s *ProductService) ListProducts(page, size int) ([]*dproduct.ProductResponse, error) {
	products, err := s.productRepository.Search(page, size, "name", "", []string{})
	if err != nil {
		return nil, errors.ErrDatabase("failed to list products", err)
	}

	responses := make([]*dproduct.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = s.productMapper.ProductToProductResponse(p)
	}
	return responses, nil
}

func (s *ProductService) GetProductsBySeller(sellerID id.UUID, page, size int) ([]*dproduct.ProductResponse, error) {
	products, err := s.productRepository.FindBySellerID(sellerID, page, size)
	if err != nil {
		return nil, errors.ErrDatabase("failed to get products for seller", err)
	}

	responses := make([]*dproduct.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = s.productMapper.ProductToProductResponse(p)
	}
	return responses, nil
}

func (s *ProductService) SearchProducts(page, size int, search string, categories []string) ([]*dproduct.ProductResponse, error) {
	products, err := s.productRepository.Search(page, size, "name", search, categories)
	if err != nil {
		return nil, errors.ErrDatabase("failed to search products", err)
	}

	responses := make([]*dproduct.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = s.productMapper.ProductToProductResponse(p)
	}
	return responses, nil
}
