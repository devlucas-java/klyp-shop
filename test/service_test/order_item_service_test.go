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

func newOrderItemService(
	orderItemRepo *mocks.OrderItemRepositoryMock,
	orderRepo *mocks.OrderRepositoryMock,
	productRepo *mocks.ProductRepositoryMock,
) *service.OrderItemService {
	return service.NewOrderItemService(
		logger.NewLogger(logger.TRACE),
		orderItemRepo,
		orderRepo,
		productRepo,
		mapper.NewOrderMapper(),
	)
}

func TestOrderItemService_GetOrderItems(t *testing.T) {
	orderItemRepo := new(mocks.OrderItemRepositoryMock)
	orderRepo := new(mocks.OrderRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderItemService(orderItemRepo, orderRepo, productRepo)

	orderID := id.NewUUID()
	items := []*entity.OrderItem{{ID: id.NewUUID(), OrderID: orderID, ProductID: id.NewUUID(), Quantity: 2, PriceBTC: 0.1}}

	orderRepo.On("FindByID", orderID).Return(&entity.Order{ID: orderID}, nil)
	orderItemRepo.On("FindByOrder", orderID).Return(items, nil)

	res, err := svc.GetOrderItems(orderID)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, items[0].ID.String(), res[0].ID)
	orderRepo.AssertExpectations(t)
	orderItemRepo.AssertExpectations(t)
}

func TestOrderItemService_GetOrderItem(t *testing.T) {
	orderItemRepo := new(mocks.OrderItemRepositoryMock)
	orderRepo := new(mocks.OrderRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderItemService(orderItemRepo, orderRepo, productRepo)

	orderID := id.NewUUID()
	itemID := id.NewUUID()
	item := &entity.OrderItem{ID: itemID, OrderID: orderID, ProductID: id.NewUUID(), Quantity: 3, PriceBTC: 0.2}

	orderRepo.On("FindByID", orderID).Return(&entity.Order{ID: orderID}, nil)
	orderItemRepo.On("FindByID", itemID).Return(item, nil)

	res, err := svc.GetOrderItem(orderID, itemID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, itemID.String(), res.ID)
	assert.InDelta(t, 0.6, res.Subtotal, 1e-9)
	orderRepo.AssertExpectations(t)
	orderItemRepo.AssertExpectations(t)
}
