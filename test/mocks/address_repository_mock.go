package mocks

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/mock"
)

type AddressRepositoryMock struct {
	mock.Mock
}

func (m *AddressRepositoryMock) Create(address *entity.Address) (*entity.Address, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Address), args.Error(1)
}

func (m *AddressRepositoryMock) Save(address *entity.Address) (*entity.Address, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Address), args.Error(1)
}

func (m *AddressRepositoryMock) Update(address *entity.Address) (*entity.Address, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Address), args.Error(1)
}

func (m *AddressRepositoryMock) Updates(address *entity.Address) (*entity.Address, error) {
	args := m.Called(address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Address), args.Error(1)
}

func (m *AddressRepositoryMock) FindByID(addressID id.UUID) (*entity.Address, error) {
	args := m.Called(addressID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Address), args.Error(1)
}

func (m *AddressRepositoryMock) FindByUser(userID id.UUID) ([]*entity.Address, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Address), args.Error(1)
}

func (m *AddressRepositoryMock) DeleteByID(addressID id.UUID) error {
	args := m.Called(addressID)
	return args.Error(0)
}
