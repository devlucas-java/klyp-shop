package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type DashboardDB struct {
	db *gorm.DB
}

func NewDashboardDB(db *gorm.DB) repository.DashboardRepository {
	return &DashboardDB{db: db}
}

func (d *DashboardDB) CountOrdersByStatusForSeller(sellerID id.UUID) ([]repository.OrderStatusCount, error) {
	var rows []repository.OrderStatusCount
	err := d.db.WithContext(context.Background()).Raw(`
		SELECT o.status, COUNT(DISTINCT o.id) AS count
		FROM orders o
		INNER JOIN order_items oi ON oi.order_id = o.id
		INNER JOIN products p ON p.id = oi.product_id
		WHERE p.seller_id = ?
		GROUP BY o.status
	`, sellerID).Scan(&rows).Error
	if err != nil {
		return nil, handlePgError(err, "failed to count orders by status for seller")
	}
	return rows, nil
}

func (d *DashboardDB) SumRevenueForSeller(sellerID id.UUID) (float64, error) {
	var revenue float64
	err := d.db.WithContext(context.Background()).Raw(`
		SELECT COALESCE(SUM(oi.price_btc * oi.quantity), 0)
		FROM order_items oi
		INNER JOIN products p ON p.id = oi.product_id
		INNER JOIN orders o ON o.id = oi.order_id
		WHERE p.seller_id = ?
		  AND o.status IN ('paid', 'delivered')
	`, sellerID).Scan(&revenue).Error
	if err != nil {
		return 0, handlePgError(err, "failed to sum revenue for seller")
	}
	return revenue, nil
}

func (d *DashboardDB) CountProductsForSeller(sellerID id.UUID) (int64, error) {
	var count int64
	err := d.db.WithContext(context.Background()).Raw(`SELECT COUNT(*) FROM products WHERE seller_id = ?`, sellerID).Scan(&count).Error
	if err != nil {
		return 0, handlePgError(err, "failed to count products for seller")
	}
	return count, nil
}

func (d *DashboardDB) AvgRatingForSeller(sellerID id.UUID) (float64, int64, error) {
	type result struct {
		Avg   float64
		Total int64
	}
	var r result
	err := d.db.WithContext(context.Background()).Raw(`
		SELECT COALESCE(AVG(rv.rating), 0) AS avg, COUNT(rv.id) AS total
		FROM reviews rv
		INNER JOIN products p ON p.id = rv.product_id
		WHERE p.seller_id = ?
	`, sellerID).Scan(&r).Error
	if err != nil {
		return 0, 0, handlePgError(err, "failed to calculate average rating for seller")
	}
	return r.Avg, r.Total, nil
}

func (d *DashboardDB) TopProductsForSeller(sellerID id.UUID, limit int) ([]repository.ProductSalesRow, error) {
	var rows []repository.ProductSalesRow
	err := d.db.WithContext(context.Background()).Raw(`
		SELECT
			p.id          AS product_id,
			p.name        AS name,
			p.stock       AS stock,
			COALESCE(SUM(oi.quantity), 0)                    AS total_sold,
			COALESCE(SUM(oi.price_btc * oi.quantity), 0)    AS revenue_btc
		FROM products p
		LEFT JOIN order_items oi ON oi.product_id = p.id
		LEFT JOIN orders o ON o.id = oi.order_id AND o.status IN ('paid', 'delivered')
		WHERE p.seller_id = ?
		GROUP BY p.id, p.name, p.stock
		ORDER BY total_sold DESC
		LIMIT ?
	`, sellerID, limit).Scan(&rows).Error
	if err != nil {
		return nil, handlePgError(err, "failed to get top products for seller")
	}
	return rows, nil
}

func (d *DashboardDB) CountAllUsers() (int64, error) {
	var count int64
	err := d.db.WithContext(context.Background()).Raw(`SELECT COUNT(*) FROM users`).Scan(&count).Error
	if err != nil {
		return 0, handlePgError(err, "failed to count all users")
	}
	return count, nil
}

func (d *DashboardDB) CountAllSellers() (int64, error) {
	var count int64
	err := d.db.WithContext(context.Background()).Raw(`SELECT COUNT(*) FROM sellers`).Scan(&count).Error
	if err != nil {
		return 0, handlePgError(err, "failed to count all sellers")
	}
	return count, nil
}

func (d *DashboardDB) CountAllProducts() (int64, error) {
	var count int64
	err := d.db.WithContext(context.Background()).Raw(`SELECT COUNT(*) FROM products`).Scan(&count).Error
	if err != nil {
		return 0, handlePgError(err, "failed to count all products")
	}
	return count, nil
}

func (d *DashboardDB) CountAllOrdersByStatus() ([]repository.OrderStatusCount, error) {
	var rows []repository.OrderStatusCount
	err := d.db.WithContext(context.Background()).Raw(`
		SELECT status, COUNT(*) AS count
		FROM orders
		GROUP BY status
	`).Scan(&rows).Error
	if err != nil {
		return nil, handlePgError(err, "failed to count orders by status")
	}
	return rows, nil
}

func (d *DashboardDB) SumTotalRevenue() (float64, error) {
	var revenue float64
	err := d.db.WithContext(context.Background()).Raw(`
		SELECT COALESCE(SUM(oi.price_btc * oi.quantity), 0)
		FROM order_items oi
		INNER JOIN orders o ON o.id = oi.order_id
		WHERE o.status IN ('paid', 'delivered')
	`).Scan(&revenue).Error
	if err != nil {
		return 0, handlePgError(err, "failed to sum total revenue")
	}
	return revenue, nil
}

func (d *DashboardDB) TopSellersByRevenue(limit int) ([]repository.SellerRevenueRow, error) {
	var rows []repository.SellerRevenueRow
	err := d.db.WithContext(context.Background()).Raw(`
		SELECT
			s.id           AS seller_id,
			s.display_name AS display_name,
			COUNT(DISTINCT o.id)                             AS total_orders,
			COALESCE(SUM(oi.price_btc * oi.quantity), 0)    AS revenue_btc,
			COALESCE(SUM(oi.quantity), 0)                   AS total_sold
		FROM sellers s
		LEFT JOIN products p ON p.seller_id = s.id
		LEFT JOIN order_items oi ON oi.product_id = p.id
		LEFT JOIN orders o ON o.id = oi.order_id AND o.status IN ('paid', 'delivered')
		GROUP BY s.id, s.display_name
		ORDER BY revenue_btc DESC
		LIMIT ?
	`, limit).Scan(&rows).Error
	if err != nil {
		return nil, handlePgError(err, "failed to get top sellers by revenue")
	}
	return rows, nil
}
