package mocks

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/mock"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) Create(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Save(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Update(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Updates(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindByID(ctx context.Context, id id.UUID) (*entity.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) FindAllPaginated(ctx context.Context, page, size int, status string) ([]*entity.Order, int64, error) {
	args := m.Called(ctx, page, size, status)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Order), args.Get(1).(int64), args.Error(2)
}

func (m *OrderRepositoryMock) FindBySellerIDPaginated(ctx context.Context, sellerID id.UUID, page, size int, status string) ([]*entity.Order, int64, error) {
	args := m.Called(ctx, sellerID, page, size, status)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Order), args.Get(1).(int64), args.Error(2)
}

func (m *OrderRepositoryMock) FindByUserIDPaginated(ctx context.Context, userID id.UUID, page, size int, status string) ([]*entity.Order, int64, error) {
	args := m.Called(ctx, userID, page, size, status)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Order), args.Get(1).(int64), args.Error(2)
}

func (m *OrderRepositoryMock) DeleteByID(ctx context.Context, id id.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
