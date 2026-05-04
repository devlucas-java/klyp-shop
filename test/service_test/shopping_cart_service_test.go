package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
)

func newShoppingCartService(cartRepo *mocks.ShoppingCartRepositoryMock) *service.ShoppingCartService {
	return service.NewShoppingCartService(
		logger.NewLogger(logger.TRACE),
		cartRepo,
		mapper.NewShoppingCartMapper(),
	)
}

func TestShoppingCartService_GetCart_WhenNoCartExists_ReturnsEmptyCart(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	svc := newShoppingCartService(cartRepo)

	user := &entity.User{ID: id.NewUUID()}
	cartRepo.On("FindByUserID", user.ID).Return(nil, nil)

	res, err := svc.GetCart(user)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user.ID.String(), res.UserID)
	assert.Empty(t, res.Items)
	assert.Equal(t, 0.0, res.TotalBTC)
	cartRepo.AssertExpectations(t)
}

func TestShoppingCartService_ClearCart_WhenCartExists_DeletesCart(t *testing.T) {
	cartRepo := new(mocks.ShoppingCartRepositoryMock)
	svc := newShoppingCartService(cartRepo)

	user := &entity.User{ID: id.NewUUID()}
	cart := &entity.ShoppingCart{ID: id.NewUUID(), UserID: user.ID}

	cartRepo.On("FindByUserID", user.ID).Return(cart, nil)
	cartRepo.On("DeleteByID", cart.ID).Return(nil)

	err := svc.ClearCart(user)

	assert.NoError(t, err)
	cartRepo.AssertExpectations(t)
}
