package mocks

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) Create(product *entity.Product) (*entity.Product, error) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) Save(product *entity.Product) (*entity.Product, error) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) Updates(product *entity.Product) (*entity.Product, error) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) FindByID(productID id.UUID) (*entity.Product, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) FindBySellerID(sellerID id.UUID, page, size int) ([]*entity.Product, error) {
	args := m.Called(sellerID, page, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) Search(page, size int, order, search string, categories []string) ([]*entity.Product, error) {
	args := m.Called(page, size, order, search, categories)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) CountTop10BySellerID(sellerID id.UUID) (int64, error) {
	args := m.Called(sellerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *ProductRepositoryMock) DeleteByID(productID id.UUID) error {
	args := m.Called(productID)
	return args.Error(0)
}
