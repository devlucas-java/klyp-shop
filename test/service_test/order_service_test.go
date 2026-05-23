package service_test

import (
	"context"
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order"
	dorderitem "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order_item"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/pkg/pagination"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newOrderService(
	orderRepo *mocks.OrderRepositoryMock,
	userRepo *mocks.UserRepositoryMock,
	addressRepo *mocks.AddressRepositoryMock,
	productRepo *mocks.ProductRepositoryMock,
) *service.OrderService {
	return service.NewOrderService(
		logger.NewLogger(logger.TRACE),
		orderRepo,
		userRepo,
		addressRepo,
		productRepo,
		mapper.NewOrderMapper(),
		mocks.NewTestMetric(),
	)
}

// ── CreateOrder ───────────────────────────────────────────────────────────────

func TestOrderService_CreateOrder(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	address := &entity.Address{ID: id.NewUUID(), UserID: user.ID}
	product := &entity.Product{ID: id.NewUUID(), PriceBTC: 500}

	orderItem, _ := entity.NewOrderItem(product.ID, 2, product.PriceBTC)
	createdOrder := entity.NewOrder(user.ID, address.ID, []entity.OrderItem{*orderItem})

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addressRepo.On("FindByID", address.ID).Return(address, nil)
	productRepo.On("FindByID", product.ID).Return(product, nil)
	orderRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(createdOrder, nil)

	req := &order.CreateOrderRequest{
		AddressID: address.ID.String(),
		Items:     []dorderitem.OrderItemRequest{{ProductID: product.ID.String(), Quantity: 2}},
	}

	res, err := svc.CreateOrder(context.Background(), user, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Items, 1)
	assert.Equal(t, int64(1000), res.TotalBTC) // 500 * 2
	assert.Equal(t, address.ID.String(), res.AddressID)
	assert.Equal(t, string(entity.OrderStatusPending), res.Status)
	orderRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	addressRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestOrderService_CreateOrder_AddressNotOwned(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	otherUserID := id.NewUUID()
	address := &entity.Address{ID: id.NewUUID(), UserID: otherUserID} // belongs to someone else

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addressRepo.On("FindByID", address.ID).Return(address, nil)

	req := &order.CreateOrderRequest{
		AddressID: address.ID.String(),
		Items:     []dorderitem.OrderItemRequest{{ProductID: id.NewUUID().String(), Quantity: 1}},
	}

	_, err := svc.CreateOrder(context.Background(), user, req)

	assert.Error(t, err)
	orderRepo.AssertNotCalled(t, "Create")
	userRepo.AssertExpectations(t)
	addressRepo.AssertExpectations(t)
}

func TestOrderService_CreateOrder_ProductNotFound(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	address := &entity.Address{ID: id.NewUUID(), UserID: user.ID}
	ghostProductID := id.NewUUID()

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addressRepo.On("FindByID", address.ID).Return(address, nil)
	productRepo.On("FindByID", ghostProductID).Return(nil, apperrors.NotFound("product", nil))

	req := &order.CreateOrderRequest{
		AddressID: address.ID.String(),
		Items:     []dorderitem.OrderItemRequest{{ProductID: ghostProductID.String(), Quantity: 1}},
	}

	_, err := svc.CreateOrder(context.Background(), user, req)

	assert.Error(t, err)
	orderRepo.AssertNotCalled(t, "Create")
}

// ── GetOrder ──────────────────────────────────────────────────────────────────

func TestOrderService_GetOrder(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	ord := &entity.Order{
		ID:        id.NewUUID(),
		UserID:    user.ID,
		AddressID: id.NewUUID(),
		Status:    entity.OrderStatusPending,
		TotalBTC:  0,
		Items:     []entity.OrderItem{},
	}

	orderRepo.On("FindByID", mock.Anything, ord.ID).Return(ord, nil)

	res, err := svc.GetOrder(context.Background(), user, ord.ID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, ord.ID.String(), res.ID)
	assert.Equal(t, string(entity.OrderStatusPending), res.Status)
	orderRepo.AssertExpectations(t)
}

func TestOrderService_GetOrder_NotOwned(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	ord := &entity.Order{
		ID:     id.NewUUID(),
		UserID: id.NewUUID(), // different user
		Status: entity.OrderStatusPending,
	}

	orderRepo.On("FindByID", mock.Anything, ord.ID).Return(ord, nil)

	_, err := svc.GetOrder(context.Background(), user, ord.ID)

	assert.Error(t, err)
	orderRepo.AssertExpectations(t)
}

func TestOrderService_GetOrder_NotFound(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	ghostID := id.NewUUID()

	orderRepo.On("FindByID", mock.Anything, ghostID).Return(nil, apperrors.NotFound("order", nil))

	_, err := svc.GetOrder(context.Background(), user, ghostID)

	assert.Error(t, err)
	orderRepo.AssertExpectations(t)
}

// ── ListUserOrders ────────────────────────────────────────────────────────────

func TestOrderService_ListUserOrders(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	orders := []*entity.Order{
		{ID: id.NewUUID(), UserID: user.ID, Status: entity.OrderStatusPending, Items: []entity.OrderItem{}},
		{ID: id.NewUUID(), UserID: user.ID, Status: entity.OrderStatusPaid, Items: []entity.OrderItem{}},
	}

	orderRepo.On("FindByUserIDPaginated", mock.Anything, user.ID, 1, 10, "").
		Return(orders, int64(2), nil)

	res, err := svc.ListUserOrders(context.Background(), user, pagination.InputPagination{Page: 1, Size: 10})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Items, 2)
	assert.Equal(t, int64(2), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.TotalPages)
	orderRepo.AssertExpectations(t)
}

// ── CancelOrder ───────────────────────────────────────────────────────────────

func TestOrderService_CancelOrder(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	ord := &entity.Order{
		ID:     id.NewUUID(),
		UserID: user.ID,
		Status: entity.OrderStatusPending,
	}

	orderRepo.On("FindByID", mock.Anything, ord.ID).Return(ord, nil)
	orderRepo.On("Updates", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(ord, nil)

	err := svc.CancelOrder(context.Background(), user, ord.ID)

	assert.NoError(t, err)
	assert.Equal(t, entity.OrderStatusCancelled, ord.Status)
	orderRepo.AssertExpectations(t)
}

func TestOrderService_CancelOrder_AlreadyPaid(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	ord := &entity.Order{
		ID:     id.NewUUID(),
		UserID: user.ID,
		Status: entity.OrderStatusPaid,
	}

	orderRepo.On("FindByID", mock.Anything, ord.ID).Return(ord, nil)

	err := svc.CancelOrder(context.Background(), user, ord.ID)

	assert.Error(t, err)
	orderRepo.AssertNotCalled(t, "Updates")
	orderRepo.AssertExpectations(t)
}

func TestOrderService_CancelOrder_NotOwned(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	ord := &entity.Order{
		ID:     id.NewUUID(),
		UserID: id.NewUUID(), // different user
		Status: entity.OrderStatusPending,
	}

	orderRepo.On("FindByID", mock.Anything, ord.ID).Return(ord, nil)

	err := svc.CancelOrder(context.Background(), user, ord.ID)

	assert.Error(t, err)
	orderRepo.AssertNotCalled(t, "Updates")
	orderRepo.AssertExpectations(t)
}
