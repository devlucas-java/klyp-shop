package service

import (
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dproduct"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type ProductService struct {
	log               *logger.Logger
	productRepository repository.ProductRepository
	userReposiory     repository.UserRepository
	sellerRepository  repository.SellerRepository
	productMapper     *mapper.ProductMaper
}

func NewProductService(l *logger.Logger, pr repository.ProductRepository, ur repository.UserRepository, sr repository.SellerRepository, pm *mapper.ProductMaper) *ProductService {
	return &ProductService{
		log:               l,
		productRepository: pr,
		userReposiory:     ur,
		sellerRepository:  sr,
		productMapper:     pm,
	}
}

func (s *ProductService) CreateProduct(auth *entity.User, req *dproduct.CreateProduct) (*dproduct.ProductResponse, error) {
	user, err := s.userReposiory.FindByIDWithSeller(auth.ID)
	if err != nil {
		s.log.Errorf("failed to find user by id: %v", err)
		return nil, err
	}

	if err := user.EnsureSeller(); err != nil {
		return nil, err
	}

	product := s.productMapper.CreateProductToProduct(req)
	product.SellerID = user.Seller.ID

	saved, err := s.productRepository.Create(product)
	if err != nil {
		s.log.Errorf("failed to create product: %v", err)
		return nil, err
	}

	return s.productMapper.ProductToProductResponse(saved), err
}

func (s *ProductService) UpdateProduct(auth *entity.User, req *dproduct.UpdateProduct, id id.UUID) (*dproduct.ProductResponse, error) {

	user, err := s.userReposiory.FindByIDWithSeller(auth.ID)
	if err != nil {
		s.log.Errorf("failed to find user by id: %v", err)
		return nil, err
	}

	product, err := s.productRepository.FindByID(id)
	if err != nil {
		s.log.Errorf("failed to find product by id: %v", err)
		return nil, err
	}

	if !product.IsOwnedBy(user.Seller.ID) {
		s.log.Errorf("user is not the owner of the product")
		return nil, errors.ErrUnauthorized(fmt.Errorf("user is not the owner of the product"))
	}

	productReq := s.productMapper.UpdateProductToProduct(req)
	productReq.ID = product.ID

	saved, err := s.productRepository.Updates(productReq)
	if err != nil {
		s.log.Errorf("failed to update product: %v", err)
		return nil, err
	}

	return s.productMapper.ProductToProductResponse(saved), nil
}

func (s *ProductService) DeleteProduct(auth *entity.User, id id.UUID) error {

	user, err := s.userReposiory.FindByIDWithSeller(auth.ID)
	if err != nil {
		s.log.Errorf("failed to find user by id: %v", err)
		return err
	}

	product, err := s.productRepository.FindByID(id)
	if err != nil {
		s.log.Errorf("failed to find product by id: %v", err)
		return err
	}

	if !product.IsOwnedBy(user.Seller.ID) {
		s.log.Errorf("user is not the owner of the product")
		return errors.ErrUnauthorized(fmt.Errorf("user is not the owner of the product"))
	}

	err = s.productRepository.DeleteByID(product.ID)
	if err != nil {
		s.log.Errorf("failed to delete product: %v", err)
		return err
	}

	return nil
}

func (s *ProductService) GetProductByID(id id.UUID) (*dproduct.ProductResponse, error) {

	product, err := s.productRepository.FindByID(id)
	if err != nil {
		s.log.Errorf("failed to find product by id: %v", err)
		return nil, err
	}

	return s.productMapper.ProductToProductResponse(product), nil
}

// ListProducts retorna una lista paginada de productos
func (s *ProductService) ListProducts(page, size int) ([]*dproduct.ProductResponse, error) {
	s.log.Infof("Listing products with pagination: page=%d, size=%d", page, size)

	products, err := s.productRepository.Search(page, size, "name", "", []string{})
	if err != nil {
		s.log.Errorf("failed to list products: %v", err)
		return nil, err
	}

	responses := make([]*dproduct.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = s.productMapper.ProductToProductResponse(product)
	}

	return responses, nil
}

// GetProductsBySeller retorna los productos de un vendedor específico
func (s *ProductService) GetProductsBySeller(sellerID id.UUID, page, size int) ([]*dproduct.ProductResponse, error) {
	s.log.Infof("Getting products for seller %s: page=%d, size=%d", sellerID, page, size)

	products, err := s.productRepository.FindBySellerID(sellerID, page, size)
	if err != nil {
		s.log.Errorf("failed to get products for seller: %v", err)
		return nil, err
	}

	responses := make([]*dproduct.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = s.productMapper.ProductToProductResponse(product)
	}

	return responses, nil
}

// SearchProducts busca productos por criterios
func (s *ProductService) SearchProducts(page, size int, search string, categories []string) ([]*dproduct.ProductResponse, error) {
	s.log.Infof("Searching products: page=%d, size=%d, search=%s, categories=%v", page, size, search, categories)

	products, err := s.productRepository.Search(page, size, "name", search, categories)
	if err != nil {
		s.log.Errorf("failed to search products: %v", err)
		return nil, err
	}

	responses := make([]*dproduct.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = s.productMapper.ProductToProductResponse(product)
	}

	return responses, nil
}
