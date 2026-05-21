package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

const orderDB = "order_db.OrderDB"

type OrderDB struct {
	db *gorm.DB
}

func NewOrderDB(db *gorm.DB) repository.OrderRepository {
	return &OrderDB{db: db}
}

func (o *OrderDB) Create(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	tx := o.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, apperrors.HandlePgError(orderDB+".create", tx.Error)
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, apperrors.HandlePgError(orderDB+".create", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apperrors.HandlePgError(orderDB+".create", err)
	}

	return order, nil
}

func (o *OrderDB) Save(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	tx := o.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, apperrors.HandlePgError(orderDB+".save", tx.Error)
	}

	if err := tx.Where("id = ?", order.ID).Save(order).Error; err != nil {
		tx.Rollback()
		return nil, apperrors.HandlePgError(orderDB+".save", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apperrors.HandlePgError(orderDB+".save", err)
	}

	return order, nil
}

func (o *OrderDB) Updates(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	tx := o.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, apperrors.HandlePgError(orderDB+".updates", tx.Error)
	}

	if err := tx.Model(order).Where("id = ?", order.ID).Updates(order).Error; err != nil {
		tx.Rollback()
		return nil, apperrors.HandlePgError(orderDB+".updates", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apperrors.HandlePgError(orderDB+".updates", err)
	}

	return order, nil
}

func (o *OrderDB) Update(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	saved, err := o.Save(ctx, order)
	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (o *OrderDB) FindByID(ctx context.Context, orderID id.UUID) (*entity.Order, error) {
	var order entity.Order
	err := o.db.WithContext(ctx).
		Preload("Items").
		Preload("User").
		Preload("Address").
		Preload("BitcoinPayment").
		First(&order, "id = ?", orderID).Error
	if err != nil {
		return nil, apperrors.HandlePgError(orderDB+".find_by_id", err)
	}
	return &order, nil
}

func (o *OrderDB) FindAllPaginated(ctx context.Context, page, size int, status string) ([]*entity.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	var total int64
	countQ := o.db.WithContext(ctx).Model(&entity.Order{})
	if status != "" {
		countQ = countQ.Where("status = ?", status)
	}
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, apperrors.HandlePgError(orderDB+".find_all_paginated", err)
	}

	var orders []*entity.Order
	q := o.db.WithContext(ctx).
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
		return nil, 0, apperrors.HandlePgError(orderDB+".find_all_paginated", err)
	}
	return orders, total, nil
}

func (o *OrderDB) FindBySellerIDPaginated(ctx context.Context, sellerID id.UUID, page, size int, status string) ([]*entity.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	baseQ := o.db.WithContext(ctx).
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("products.seller_id = ?", sellerID)
	if status != "" {
		baseQ = baseQ.Where("orders.status = ?", status)
	}

	var total int64
	if err := baseQ.Model(&entity.Order{}).Distinct("orders.id").Count(&total).Error; err != nil {
		return nil, 0, apperrors.HandlePgError(orderDB+".find_by_seller_id_paginated", err)
	}

	var orders []*entity.Order
	q := o.db.WithContext(ctx).
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
		return nil, 0, apperrors.HandlePgError(orderDB+".find_by_seller_id_paginated", err)
	}
	return orders, total, nil
}

func (o *OrderDB) FindByUserIDPaginated(ctx context.Context, userID id.UUID, page, size int, status string) ([]*entity.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	baseQ := o.db.WithContext(ctx).
		Model(&entity.Order{}).
		Distinct("orders.id").
		Where("orders.users_id = ?", userID)
	if status != "" {
		baseQ = baseQ.Where("orders.status = ?", status)
	}

	var total int64
	if err := baseQ.Count(&total).Error; err != nil {
		return nil, 0, apperrors.HandlePgError(orderDB+".find_by_user_id_paginated", err)
	}

	var orders []*entity.Order
	q := o.db.WithContext(ctx).
		Preload("Items.Product").
		Preload("User").
		Preload("Address").
		Preload("BitcoinPayment").
		Where("orders.users_id = ?", userID).
		Group("orders.id").
		Order("orders.created_at desc").
		Limit(size).
		Offset((page - 1) * size)
	if status != "" {
		q = q.Where("orders.status = ?", status)
	}
	if err := q.Find(&orders).Error; err != nil {
		return nil, 0, apperrors.HandlePgError(orderDB+".find_by_user_id_paginated", err)
	}
	return orders, total, nil
}

func (o *OrderDB) DeleteByID(ctx context.Context, orderID id.UUID) error {
	tx := o.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return apperrors.HandlePgError(orderDB+".delete_by_id", tx.Error)
	}

	if err := tx.Delete(&entity.Order{}, "id = ?", orderID).Error; err != nil {
		tx.Rollback()
		return apperrors.HandlePgError(orderDB+".delete_by_id", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return apperrors.HandlePgError(orderDB+".delete_by_id", err)
	}

	return nil
}
