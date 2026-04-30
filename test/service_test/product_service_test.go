package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dproduct"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newProductService(
	productRepo *mocks.ProductRepositoryMock,
	userRepo *mocks.UserRepositoryMock,
	sellerRepo *mocks.SellerRepositoryMock,
) *service.ProductService {
	return service.NewProductService(
		logger.NewLogger(logger.TRACE),
		productRepo,
		userRepo,
		sellerRepo,
		mapper.NewProductMapper(),
	)
}

func newProductUser() (*entity.User, *entity.Seller) {
	userID := id.NewUUID()
	seller := &entity.Seller{ID: id.NewUUID(), UserID: userID, DisplayName: "Store"}
	user := &entity.User{
		ID:       userID,
		IsSeller: true,
		Seller:   seller,
	}
	return user, seller
}

func newProduct(sellerID id.UUID) *entity.Product {
	return &entity.Product{
		ID:       id.NewUUID(),
		Name:     "Test Product",
		PriceBTC: 0.01,
		Stock:    100,
		SellerID: sellerID,
	}
}

func TestProductService_CreateProduct(t *testing.T) {
	productRepo := new(mocks.ProductRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newProductService(productRepo, userRepo, sellerRepo)

	user, seller := newProductUser()
	product := newProduct(seller.ID)

	userRepo.On("FindByIDWithSeller", user.ID).Return(user, nil)
	productRepo.On("Create", mock.AnythingOfType("*entity.Product")).Return(product, nil)

	req := &dproduct.CreateProduct{Name: "Test Product", PriceBTC: 0.01, Stock: 100}
	res, err := svc.CreateProduct(user, req)

	assert.NoError(t, err)
	assert.Equal(t, product.Name, res.Name)
	assert.Equal(t, seller.ID.String(), res.SellerID)
	userRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct_NotSeller(t *testing.T) {
	productRepo := new(mocks.ProductRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newProductService(productRepo, userRepo, sellerRepo)

	user := &entity.User{ID: id.NewUUID(), IsSeller: false, Seller: nil}
	userRepo.On("FindByIDWithSeller", user.ID).Return(user, nil)

	req := &dproduct.CreateProduct{Name: "Test", PriceBTC: 0.01, Stock: 1}
	_, err := svc.CreateProduct(user, req)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestProductService_GetProductByID(t *testing.T) {
	productRepo := new(mocks.ProductRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newProductService(productRepo, userRepo, sellerRepo)

	_, seller := newProductUser()
	product := newProduct(seller.ID)
	productRepo.On("FindByID", product.ID).Return(product, nil)

	res, err := svc.GetProductByID(product.ID)

	assert.NoError(t, err)
	assert.Equal(t, product.ID.String(), res.ID)
	productRepo.AssertExpectations(t)
}

func TestProductService_GetProductByID_NotFound(t *testing.T) {
	productRepo := new(mocks.ProductRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newProductService(productRepo, userRepo, sellerRepo)

	ghostID := id.NewUUID()
	productRepo.On("FindByID", ghostID).Return(nil, domainErr.ErrNotFound("Product", nil))

	_, err := svc.GetProductByID(ghostID)

	assert.Error(t, err)
	productRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct(t *testing.T) {
	productRepo := new(mocks.ProductRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newProductService(productRepo, userRepo, sellerRepo)

	user, seller := newProductUser()
	product := newProduct(seller.ID)
	updated := *product
	updated.Name = "Updated Name"
	updated.PriceBTC = 0.02

	userRepo.On("FindByIDWithSeller", user.ID).Return(user, nil)
	productRepo.On("FindByID", product.ID).Return(product, nil)
	productRepo.On("Updates", mock.AnythingOfType("*entity.Product")).Return(&updated, nil)

	req := &dproduct.UpdateProduct{Name: "Updated Name", PriceBTC: 0.02, Stock: 50}
	res, err := svc.UpdateProduct(user, req, product.ID)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", res.Name)
	userRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_NotOwner(t *testing.T) {
	productRepo := new(mocks.ProductRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newProductService(productRepo, userRepo, sellerRepo)

	user1, seller1 := newProductUser()
	user2, _ := newProductUser()
	product := newProduct(seller1.ID)

	userRepo.On("FindByIDWithSeller", user2.ID).Return(user2, nil)
	productRepo.On("FindByID", product.ID).Return(product, nil)

	req := &dproduct.UpdateProduct{Name: "Hacked"}
	_, err := svc.UpdateProduct(user2, req, product.ID)

	assert.Error(t, err)
	_ = user1
	userRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestProductService_DeleteProduct(t *testing.T) {
	productRepo := new(mocks.ProductRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newProductService(productRepo, userRepo, sellerRepo)

	user, seller := newProductUser()
	product := newProduct(seller.ID)

	userRepo.On("FindByIDWithSeller", user.ID).Return(user, nil)
	productRepo.On("FindByID", product.ID).Return(product, nil)
	productRepo.On("DeleteByID", product.ID).Return(nil)

	err := svc.DeleteProduct(user, product.ID)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestProductService_DeleteProduct_NotOwner(t *testing.T) {
	productRepo := new(mocks.ProductRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newProductService(productRepo, userRepo, sellerRepo)

	user1, seller1 := newProductUser()
	user2, _ := newProductUser()
	product := newProduct(seller1.ID)

	userRepo.On("FindByIDWithSeller", user2.ID).Return(user2, nil)
	productRepo.On("FindByID", product.ID).Return(product, nil)

	err := svc.DeleteProduct(user2, product.ID)

	assert.Error(t, err)
	_ = user1
	userRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}
