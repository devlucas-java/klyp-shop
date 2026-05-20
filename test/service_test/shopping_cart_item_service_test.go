package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	cartDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/cart"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newShoppingCartItemService(cartRepo *mocks.ShoppingCartRepositoryMock, productRepo *mocks.ProductRepositoryMock) *service.ShoppingCartItemService {
	return service.NewShoppingCartItemService(
		logger.NewLogger(logger.TRACE),
		cartRepo,
		productRepo,
		mapper.NewShoppingCartMapper(),
	)
}

func TestShoppingCartItemService_AddItem(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	product := &entity.Product{ID: id.NewUUID(), PriceBTC: 0.05}
	cart := &entity.ShoppingCart{ID: id.NewUUID(), UserID: user.ID, Items: []*entity.ShoppingCartItem{}}

	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)
	productRepo.On("FindByID", product.ID).Return(product, nil)
	cartRepo.On("Save", mock.AnythingOfType("*entity.ShoppingCart")).Return(cart, nil)

	req := &cartDTO.AddShoppingCartItemRequest{ProductID: product.ID.String(), Quantity: 2}
	res, err := svc.AddItem(user, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user.ID.String(), res.UserID)
	assert.Equal(t, 1, len(res.Items))
	assert.Equal(t, 0.1, res.TotalBTC)
	cartRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestShoppingCartItemService_UpdateItem(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	item := &entity.ShoppingCartItem{ID: id.NewUUID(), CartID: id.NewUUID(), ProductID: id.NewUUID(), Quantity: 2, PriceBTC: 0.05}
	cart := &entity.ShoppingCart{ID: item.CartID, UserID: user.ID, Items: []*entity.ShoppingCartItem{item}}

	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)
	cartRepo.On("Save", mock.AnythingOfType("*entity.ShoppingCart")).Return(cart, nil)

	req := &cartDTO.UpdateShoppingCartItemRequest{Quantity: 4}
	res, err := svc.UpdateItem(user, item.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res.Items))
	assert.Equal(t, 4, res.Items[0].Quantity)
	assert.Equal(t, 0.2, res.TotalBTC)
	cartRepo.AssertExpectations(t)
}

func TestShoppingCartItemService_RemoveItem(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newShoppingCartItemService(cartRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	item := &entity.ShoppingCartItem{ID: id.NewUUID(), CartID: id.NewUUID(), ProductID: id.NewUUID(), Quantity: 1, PriceBTC: 0.05}
	cart := &entity.ShoppingCart{ID: item.CartID, UserID: user.ID, Items: []*entity.ShoppingCartItem{item}}

	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)
	cartRepo.On("DeleteByID", cart.ID).Return(nil)

	err := svc.RemoveItem(user, item.ID)

	assert.NoError(t, err)
	cartRepo.AssertExpectations(t)
}
