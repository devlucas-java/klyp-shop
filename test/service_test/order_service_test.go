package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dorder"
	dorderitem "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dorder_item"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
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
	)
}

func TestOrderService_CreateOrder(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	address := &entity.Address{ID: id.NewUUID(), UserID: user.ID}
	product := &entity.Product{ID: id.NewUUID(), PriceBTC: 0.5}

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addressRepo.On("FindByID", address.ID).Return(address, nil)
	productRepo.On("FindByID", product.ID).Return(product, nil)

	orderItem, _ := entity.NewOrderItem(product.ID, 2, product.PriceBTC)
	createdOrder := entity.NewOrder(user.ID, address.ID, []entity.OrderItem{*orderItem})
	orderRepo.On("Create", mock.AnythingOfType("*entity.Order")).Return(createdOrder, nil)

	req := &dorder.CreateOrderRequest{
		AddressID: address.ID.String(),
		Items:     []dorderitem.OrderItemRequest{{ProductID: product.ID.String(), Quantity: 2}},
	}

	res, err := svc.CreateOrder(user, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res.Items))
	assert.Equal(t, 1.0, res.TotalBTC)
	assert.Equal(t, address.ID.String(), res.AddressID)
	orderRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	addressRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestOrderService_GetOrder(t *testing.T) {
	orderRepo := new(mocks.OrderRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	addressRepo := new(mocks.AddressRepositoryMock)
	productRepo := new(mocks.ProductRepositoryMock)
	svc := newOrderService(orderRepo, userRepo, addressRepo, productRepo)

	user := &entity.User{ID: id.NewUUID()}
	order := &entity.Order{ID: id.NewUUID(), UserID: user.ID, AddressID: id.NewUUID(), Status: entity.OrderStatusPending, TotalBTC: 0.0, Items: []entity.OrderItem{}}

	orderRepo.On("FindByID", order.ID).Return(order, nil)

	res, err := svc.GetOrder(user, order.ID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, order.ID.String(), res.ID)
	orderRepo.AssertExpectations(t)
}
