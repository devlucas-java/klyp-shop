package mocks

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/mock"
)

type OrderItemRepositoryMock struct {
	mock.Mock
}

func (m *OrderItemRepositoryMock) Create(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	args := m.Called(orderItem)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.OrderItem), args.Error(1)
}

func (m *OrderItemRepositoryMock) Save(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	args := m.Called(orderItem)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.OrderItem), args.Error(1)
}

func (m *OrderItemRepositoryMock) Update(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	args := m.Called(orderItem)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.OrderItem), args.Error(1)
}

func (m *OrderItemRepositoryMock) Updates(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	args := m.Called(orderItem)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.OrderItem), args.Error(1)
}

func (m *OrderItemRepositoryMock) FindByID(id id.UUID) (*entity.OrderItem, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.OrderItem), args.Error(1)
}

func (m *OrderItemRepositoryMock) FindByOrder(orderID id.UUID) ([]*entity.OrderItem, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.OrderItem), args.Error(1)
}

func (m *OrderItemRepositoryMock) DeleteByID(id id.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
