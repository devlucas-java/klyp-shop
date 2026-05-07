package mocks

import (
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/mock"
)

type DashboardRepositoryMock struct {
	mock.Mock
}

func (m *DashboardRepositoryMock) CountOrdersByStatusForSeller(sellerID id.UUID) ([]repository.OrderStatusCount, error) {
	args := m.Called(sellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.OrderStatusCount), args.Error(1)
}

func (m *DashboardRepositoryMock) SumRevenueForSeller(sellerID id.UUID) (float64, error) {
	args := m.Called(sellerID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *DashboardRepositoryMock) CountProductsForSeller(sellerID id.UUID) (int64, error) {
	args := m.Called(sellerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *DashboardRepositoryMock) AvgRatingForSeller(sellerID id.UUID) (float64, int64, error) {
	args := m.Called(sellerID)
	return args.Get(0).(float64), args.Get(1).(int64), args.Error(2)
}

func (m *DashboardRepositoryMock) TopProductsForSeller(sellerID id.UUID, limit int) ([]repository.ProductSalesRow, error) {
	args := m.Called(sellerID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.ProductSalesRow), args.Error(1)
}

func (m *DashboardRepositoryMock) CountAllUsers() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *DashboardRepositoryMock) CountAllSellers() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *DashboardRepositoryMock) CountAllProducts() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *DashboardRepositoryMock) CountAllOrdersByStatus() ([]repository.OrderStatusCount, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.OrderStatusCount), args.Error(1)
}

func (m *DashboardRepositoryMock) SumTotalRevenue() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *DashboardRepositoryMock) TopSellersByRevenue(limit int) ([]repository.SellerRevenueRow, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.SellerRevenueRow), args.Error(1)
}
