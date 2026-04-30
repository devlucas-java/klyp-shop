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

	if !user.IsSeller || user.Seller == nil {
		return nil, errors.ErrForbidden(fmt.Errorf("user is not a seller"))
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

	if product.SellerID != user.Seller.ID {
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

	if product.SellerID != user.Seller.ID {
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
