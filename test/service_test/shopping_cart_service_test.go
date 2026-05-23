package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newShoppingCartService(cartRepo *mocks.ShoppingCartRepositoryMock) *service.ShoppingCartService {
	return service.NewShoppingCartService(
		logger.NewLogger(logger.TRACE),
		cartRepo,
		mapper.NewShoppingCartMapper(),
	)
}

func TestShoppingCartService_GetCart(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	svc := newShoppingCartService(cartRepo)

	user := &entity.User{ID: id.NewUUID()}
	cart := &entity.ShoppingCart{
		ID:     id.NewUUID(),
		UserID: user.ID,
		Items:  []*entity.ShoppingCartItem{},
	}

	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)

	res, err := svc.GetCart(user)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, cart.ID.String(), res.ID)
	assert.Equal(t, user.ID.String(), res.UserID)
	assert.Empty(t, res.Items)
	assert.Equal(t, int64(0), res.TotalBTC)
	cartRepo.AssertExpectations(t)
}

func TestShoppingCartService_GetCart_NotFound(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	svc := newShoppingCartService(cartRepo)

	user := &entity.User{ID: id.NewUUID()}
	cartRepo.On("FindByUserID", user.ID).Return(nil, apperrors.NotFound("shopping_cart", nil))

	_, err := svc.GetCart(user)

	assert.Error(t, err)
	cartRepo.AssertExpectations(t)
}

func TestShoppingCartService_ClearCart(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	svc := newShoppingCartService(cartRepo)

	user := &entity.User{ID: id.NewUUID()}
	cart := &entity.ShoppingCart{ID: id.NewUUID(), UserID: user.ID}
	newCart := entity.NewShoppingCart(user.ID)

	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)
	cartRepo.On("DeleteByID", cart.ID).Return(nil)
	cartRepo.On("Create", mock.AnythingOfType("*entity.ShoppingCart")).Return(newCart, nil)

	err := svc.ClearCart(user)

	assert.NoError(t, err)
	cartRepo.AssertCalled(t, "FindByUserID", user.ID)
	cartRepo.AssertCalled(t, "DeleteByID", cart.ID)
	cartRepo.AssertExpectations(t)
}
