package mocks

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/mock"
)

type SellerRepositoryMock struct {
	mock.Mock
}

func (m *SellerRepositoryMock) Create(seller *entity.Seller) (*entity.Seller, error) {
	args := m.Called(seller)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Seller), args.Error(1)
}

func (m *SellerRepositoryMock) Save(seller *entity.Seller) (*entity.Seller, error) {
	args := m.Called(seller)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Seller), args.Error(1)
}

func (m *SellerRepositoryMock) Updates(seller *entity.Seller) (*entity.Seller, error) {
	args := m.Called(seller)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Seller), args.Error(1)
}

func (m *SellerRepositoryMock) FindByID(sellerID id.UUID) (*entity.Seller, error) {
	args := m.Called(sellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Seller), args.Error(1)
}

func (m *SellerRepositoryMock) Find(page, size int, order, search string) ([]*entity.Seller, error) {
	args := m.Called(page, size, order, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Seller), args.Error(1)
}

func (m *SellerRepositoryMock) DeleteByID(sellerID id.UUID) error {
	args := m.Called(sellerID)
	return args.Error(0)
}
