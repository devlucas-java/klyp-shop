package service_test

import (
	"context"
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		policy.NewOrderPolicy(),
	)
}

// ── GetOrderItems ─────────────────────────────────────────────────────────────

func TestOrderItemService_GetOrderItems(t *testing.T) {
	orderItemRepo := new(mocks.OrderItemRepositoryMock)
	orderRepo := new(mocks.OrderRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderItemService(orderItemRepo, orderRepo, productRepo)

	orderID := id.NewUUID()
	items := []*entity.OrderItem{
		{ID: id.NewUUID(), OrderID: orderID, ProductID: id.NewUUID(), Quantity: 2, PriceBTC: 100},
	}

	orderRepo.On("FindByID", mock.Anything, orderID).Return(&entity.Order{ID: orderID}, nil)
	orderItemRepo.On("FindByOrder", orderID).Return(items, nil)

	res, err := svc.GetOrderItems(context.Background(), orderID)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, items[0].ID.String(), res[0].ID)
	assert.Equal(t, int64(200), res[0].Subtotal) // 100 * 2
	orderRepo.AssertExpectations(t)
	orderItemRepo.AssertExpectations(t)
}

func TestOrderItemService_GetOrderItems_Empty(t *testing.T) {
	orderItemRepo := new(mocks.OrderItemRepositoryMock)
	orderRepo := new(mocks.OrderRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderItemService(orderItemRepo, orderRepo, productRepo)

	orderID := id.NewUUID()

	orderRepo.On("FindByID", mock.Anything, orderID).Return(&entity.Order{ID: orderID}, nil)
	orderItemRepo.On("FindByOrder", orderID).Return([]*entity.OrderItem{}, nil)

	res, err := svc.GetOrderItems(context.Background(), orderID)

	assert.NoError(t, err)
	assert.Empty(t, res)
	orderRepo.AssertExpectations(t)
	orderItemRepo.AssertExpectations(t)
}

func TestOrderItemService_GetOrderItems_OrderNotFound(t *testing.T) {
	orderItemRepo := new(mocks.OrderItemRepositoryMock)
	orderRepo := new(mocks.OrderRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderItemService(orderItemRepo, orderRepo, productRepo)

	orderID := id.NewUUID()
	orderRepo.On("FindByID", mock.Anything, orderID).Return(nil, apperrors.NotFound("order", nil))

	_, err := svc.GetOrderItems(context.Background(), orderID)

	assert.Error(t, err)
	orderItemRepo.AssertNotCalled(t, "FindByOrder")
	orderRepo.AssertExpectations(t)
}

// ── GetOrderItem ──────────────────────────────────────────────────────────────

func TestOrderItemService_GetOrderItem(t *testing.T) {
	orderItemRepo := new(mocks.OrderItemRepositoryMock)
	orderRepo := new(mocks.OrderRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderItemService(orderItemRepo, orderRepo, productRepo)

	orderID := id.NewUUID()
	itemID := id.NewUUID()
	item := &entity.OrderItem{ID: itemID, OrderID: orderID, ProductID: id.NewUUID(), Quantity: 3, PriceBTC: 200}

	orderRepo.On("FindByID", mock.Anything, orderID).Return(&entity.Order{ID: orderID}, nil)
	orderItemRepo.On("FindByID", itemID).Return(item, nil)

	res, err := svc.GetOrderItem(context.Background(), orderID, itemID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, itemID.String(), res.ID)
	assert.Equal(t, int64(600), res.Subtotal) // 200 * 3
	orderRepo.AssertExpectations(t)
	orderItemRepo.AssertExpectations(t)
}

func TestOrderItemService_GetOrderItem_WrongOrder(t *testing.T) {
	orderItemRepo := new(mocks.OrderItemRepositoryMock)
	orderRepo := new(mocks.OrderRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderItemService(orderItemRepo, orderRepo, productRepo)

	orderID := id.NewUUID()
	otherOrderID := id.NewUUID()
	itemID := id.NewUUID()
	// item belongs to otherOrderID, not orderID
	item := &entity.OrderItem{ID: itemID, OrderID: otherOrderID, ProductID: id.NewUUID(), Quantity: 1, PriceBTC: 100}

	orderRepo.On("FindByID", mock.Anything, orderID).Return(&entity.Order{ID: orderID}, nil)
	orderItemRepo.On("FindByID", itemID).Return(item, nil)

	_, err := svc.GetOrderItem(context.Background(), orderID, itemID)

	assert.Error(t, err)
	orderRepo.AssertExpectations(t)
	orderItemRepo.AssertExpectations(t)
}

func TestOrderItemService_GetOrderItem_NotFound(t *testing.T) {
	orderItemRepo := new(mocks.OrderItemRepositoryMock)
	orderRepo := new(mocks.OrderRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderItemService(orderItemRepo, orderRepo, productRepo)

	orderID := id.NewUUID()
	itemID := id.NewUUID()

	orderRepo.On("FindByID", mock.Anything, orderID).Return(&entity.Order{ID: orderID}, nil)
	orderItemRepo.On("FindByID", itemID).Return(nil, apperrors.NotFound("order_item", nil))

	_, err := svc.GetOrderItem(context.Background(), orderID, itemID)

	assert.Error(t, err)
	orderRepo.AssertExpectations(t)
	orderItemRepo.AssertExpectations(t)
}

// ── CalculateOrderTotal ───────────────────────────────────────────────────────

func TestOrderItemService_CalculateOrderTotal(t *testing.T) {
	orderItemRepo := new(mocks.OrderItemRepositoryMock)
	orderRepo := new(mocks.OrderRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderItemService(orderItemRepo, orderRepo, productRepo)

	orderID := id.NewUUID()
	items := []*entity.OrderItem{
		{ID: id.NewUUID(), OrderID: orderID, Quantity: 2, PriceBTC: 500},
		{ID: id.NewUUID(), OrderID: orderID, Quantity: 1, PriceBTC: 300},
	}

	orderItemRepo.On("FindByOrder", orderID).Return(items, nil)

	total, err := svc.CalculateOrderTotal(orderID)

	assert.NoError(t, err)
	assert.Equal(t, int64(1300), total) // (500*2) + (300*1)
	orderItemRepo.AssertExpectations(t)
}
