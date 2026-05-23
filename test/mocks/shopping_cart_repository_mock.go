package mocks

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/mock"
)

type ShoppingCartRepositoryMock struct {
	mock.Mock
}

func (m *ShoppingCartRepositoryMock) FindByUserID(userID id.UUID) (*entity.ShoppingCart, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShoppingCart), args.Error(1)
}

func (m *ShoppingCartRepositoryMock) FindByID(cartID id.UUID) (*entity.ShoppingCart, error) {
	args := m.Called(cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShoppingCart), args.Error(1)
}

func (m *ShoppingCartRepositoryMock) Create(cart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	args := m.Called(cart)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShoppingCart), args.Error(1)
}

func (m *ShoppingCartRepositoryMock) Save(cart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	args := m.Called(cart)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShoppingCart), args.Error(1)
}

func (m *ShoppingCartRepositoryMock) DeleteByID(uuid id.UUID) error {
	args := m.Called(uuid)
	return args.Error(0)
}

func (m *ShoppingCartRepositoryMock) FindCartsByProductID(productID id.UUID) ([]*entity.ShoppingCart, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ShoppingCart), args.Error(1)
}
