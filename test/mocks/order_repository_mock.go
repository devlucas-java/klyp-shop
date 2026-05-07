package mocks

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/mock"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) Create(order *entity.Order) (*entity.Order, error) {
	args := m.Called(order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Save(order *entity.Order) (*entity.Order, error) {
	args := m.Called(order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Update(order *entity.Order) (*entity.Order, error) {
	args := m.Called(order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Updates(order *entity.Order) (*entity.Order, error) {
	args := m.Called(order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindByID(id id.UUID) (*entity.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindByUser(userID id.UUID) ([]*entity.Order, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindAll() ([]*entity.Order, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindAllWithDetails() ([]*entity.Order, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindAllPaginated(page, size int, status string) ([]*entity.Order, int64, error) {
	args := m.Called(page, size, status)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Order), args.Get(1).(int64), args.Error(2)
}

func (m *OrderRepositoryMock) FindBySellerID(sellerID id.UUID) ([]*entity.Order, error) {
	args := m.Called(sellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindBySellerIDPaginated(sellerID id.UUID, page, size int, status string) ([]*entity.Order, int64, error) {
	args := m.Called(sellerID, page, size, status)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Order), args.Get(1).(int64), args.Error(2)
}

func (m *OrderRepositoryMock) DeleteByID(id id.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
