package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type DashboardRepository interface {
	CountOrdersByStatusForSeller(sellerID id.UUID) (int, error)
	SumRevenueForSeller(sellerID id.UUID) (int64, error)
	CountProductsForSeller(sellerID id.UUID) (int64, error)
	AvgRatingForSeller(sellerID id.UUID) (int64, int64, error)
	TopProductsForSeller(sellerID id.UUID, limit int) ([]entity.Product, error)

	CountAllUsers() (int64, error)
	CountAllSellers() (int64, error)
	CountAllProducts() (int64, error)
	CountAllOrdersByStatus() (int, error)
	SumTotalRevenue() (int64, error)
	TopSellersByRevenue(limit int) ([]entity.Seller, error)
}
