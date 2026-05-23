package mocks

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/mock"
)

type ShoppingCartItemRepositoryMock struct {
	mock.Mock
}

func (m *ShoppingCartItemRepositoryMock) FindByID(itemID id.UUID) (*entity.ShoppingCartItem, error) {
	args := m.Called(itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShoppingCartItem), args.Error(1)
}

func (m *ShoppingCartItemRepositoryMock) FindByCartID(cartID id.UUID) ([]*entity.ShoppingCartItem, error) {
	args := m.Called(cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ShoppingCartItem), args.Error(1)
}

func (m *ShoppingCartItemRepositoryMock) Create(item *entity.ShoppingCartItem) (*entity.ShoppingCartItem, error) {
	args := m.Called(item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShoppingCartItem), args.Error(1)
}

func (m *ShoppingCartItemRepositoryMock) Save(item *entity.ShoppingCartItem) (*entity.ShoppingCartItem, error) {
	args := m.Called(item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ShoppingCartItem), args.Error(1)
}

func (m *ShoppingCartItemRepositoryMock) DeleteByID(itemID id.UUID) error {
	args := m.Called(itemID)
	return args.Error(0)
}
