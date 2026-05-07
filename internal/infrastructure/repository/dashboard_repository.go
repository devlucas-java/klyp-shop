package repository

import "github.com/devlucas-java/klyp-shop/pkg/id"

type DashboardOrderFilter struct {
	Status string
	Page   int
	Size   int
}

type OrderStatusCount struct {
	Status string
	Count  int64
}

type ProductSalesRow struct {
	ProductID  string
	Name       string
	TotalSold  int64
	RevenueBTC float64
	Stock      int
}

type SellerRevenueRow struct {
	SellerID    string
	DisplayName string
	TotalOrders int64
	RevenueBTC  float64
	TotalSold   int64
}

type DashboardRepository interface {
	CountOrdersByStatusForSeller(sellerID id.UUID) ([]OrderStatusCount, error)
	SumRevenueForSeller(sellerID id.UUID) (float64, error)
	CountProductsForSeller(sellerID id.UUID) (int64, error)
	AvgRatingForSeller(sellerID id.UUID) (float64, int64, error)
	TopProductsForSeller(sellerID id.UUID, limit int) ([]ProductSalesRow, error)

	CountAllUsers() (int64, error)
	CountAllSellers() (int64, error)
	CountAllProducts() (int64, error)
	CountAllOrdersByStatus() ([]OrderStatusCount, error)
	SumTotalRevenue() (float64, error)
	TopSellersByRevenue(limit int) ([]SellerRevenueRow, error)
}
