package database

import (
	"context"
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"gorm.io/gorm"
)

type OrderDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewOrderDB(db *gorm.DB, log *logger.Logger) repository.OrderRepository {
	return &OrderDB{db: db, log: log}
}

func (o *OrderDB) Create(order *entity.Order) (*entity.Order, error) {
	if err := o.db.WithContext(context.Background()).Create(order).Error; err != nil {
		return nil, errors.HandlePgError(o.log, err, "failed to create order")
	}
	return order, nil
}

func (o *OrderDB) Save(order *entity.Order) (*entity.Order, error) {
	if err := o.db.WithContext(context.Background()).Where("id = ?", order.ID).Save(order).Error; err != nil {
		return nil, errors.HandlePgError(o.log, err, "failed to save order")
	}
	return order, nil
}

func (o *OrderDB) Updates(order *entity.Order) (*entity.Order, error) {
	if err := o.db.WithContext(context.Background()).Model(order).Where("id = ?", order.ID).Updates(order).Error; err != nil {
		return nil, errors.HandlePgError(o.log, err, "failed to update order")
	}
	return order, nil
}

func (o *OrderDB) Update(order *entity.Order) (*entity.Order, error) {
	saved, err := o.Save(order)
	if err != nil {
		return nil, errors.HandlePgError(o.log, err, "erros in update order")
	}
	return saved, nil
}

func (o *OrderDB) FindByID(orderID id.UUID) (*entity.Order, error) {
	var order entity.Order
	err := o.db.WithContext(context.Background()).Preload("Items").Preload("User").Preload("Address").Preload("BitcoinPayment").First(&order, "id = ?", orderID).Error
	if err != nil {
		return nil, errors.HandlePgError(o.log, err, "failed to find order")
	}
	return &order, nil
}

func (o *OrderDB) FindByUser(userID id.UUID) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := o.db.WithContext(context.Background()).Preload("Items").Preload("User").Preload("Address").Preload("BitcoinPayment").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, errors.HandlePgError(o.log, err, "failed to find orders")
	}
	return orders, nil
}

func (o *OrderDB) FindAll() ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := o.db.WithContext(context.Background()).
		Preload("Items").
		Preload("User").
		Preload("Address").
		Preload("BitcoinPayment").
		Find(&orders).Error; err != nil {
		return nil, errors.HandlePgError(o.log, err, "failed to find orders")
	}
	return orders, nil
}

func (o *OrderDB) FindAllWithDetails() ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := o.db.WithContext(context.Background()).
		Preload("Items.Product.Seller").
		Preload("User").
		Preload("Address").
		Preload("BitcoinPayment").
		Order("created_at desc").
		Find(&orders).Error; err != nil {
		return nil, errors.HandlePgError(o.log, err, "failed to find orders with details")
	}
	return orders, nil
}

func (o *OrderDB) FindAllPaginated(page, size int, status string) ([]*entity.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	var total int64
	countQ := o.db.WithContext(context.Background()).Model(&entity.Order{})
	if status != "" {
		countQ = countQ.Where("status = ?", status)
	}
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("OrderDB.FindAllPaginated count: %w", err)
	}

	var orders []*entity.Order
	q := o.db.WithContext(context.Background()).
		Preload("Items.Product.Seller").
		Preload("User").
		Preload("Address").
		Preload("BitcoinPayment").
		Order("created_at desc").
		Limit(size).
		Offset((page - 1) * size)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find paginated orders: %w", err)
	}
	return orders, total, nil
}

func (o *OrderDB) FindBySellerIDPaginated(sellerID id.UUID, page, size int, status string) ([]*entity.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	baseQ := o.db.WithContext(context.Background()).
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.seller_id = ?", sellerID)
	if status != "" {
		baseQ = baseQ.Where("orders.status = ?", status)
	}

	var total int64
	if err := baseQ.Model(&entity.Order{}).Distinct("orders.id").Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("OrderDB.FindBySellerIDPaginated count: %w", err)
	}

	var orders []*entity.Order
	q := o.db.WithContext(context.Background()).
		Preload("Items.Product").
		Preload("User").
		Preload("Address").
		Preload("BitcoinPayment").
		Joins("JOIN order_items oi2 ON oi2.order_id = orders.id").
		Joins("JOIN products p2 ON p2.id = oi2.product_id").
		Where("p2.seller_id = ?", sellerID).
		Group("orders.id").
		Order("orders.created_at desc").
		Limit(size).
		Offset((page - 1) * size)
	if status != "" {
		q = q.Where("orders.status = ?", status)
	}
	if err := q.Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find paginated seller orders: %w", err)
	}
	return orders, total, nil
}

func (o *OrderDB) FindBySellerID(sellerID id.UUID) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := o.db.WithContext(context.Background()).
		Preload("Items.Product").
		Preload("User").
		Preload("Address").
		Preload("BitcoinPayment").
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.seller_id = ?", sellerID).
		Group("orders.id").
		Order("orders.created_at desc").
		Find(&orders).Error; err != nil {
		return nil, errors.HandlePgError(o.log, err, "failed to find orders by seller")
	}
	return orders, nil
}

func (o *OrderDB) DeleteByID(orderID id.UUID) error {
	if err := o.db.WithContext(context.Background()).Delete(&entity.Order{}, "id = ?", orderID).Error; err != nil {
		return errors.HandlePgError(o.log, err, "failed to delete order")
	}
	return nil
}
