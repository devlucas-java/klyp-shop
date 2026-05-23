package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type OrderItemDB struct {
	db *gorm.DB
}

func NewOrderItemDB(db *gorm.DB) repository.OrderItemRepository {
	return &OrderItemDB{db: db}
}

func (oi *OrderItemDB) Create(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	if err := oi.db.WithContext(context.Background()).Create(orderItem).Error; err != nil {
		return nil, apperrors.HandlePgError("order_item", err)
	}
	return orderItem, nil
}

func (oi *OrderItemDB) Save(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	if err := oi.db.WithContext(context.Background()).Where("id = ?", orderItem.ID).Save(orderItem).Error; err != nil {
		return nil, apperrors.HandlePgError("order_item", err)
	}
	return orderItem, nil
}

func (oi *OrderItemDB) Updates(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	if err := oi.db.WithContext(context.Background()).Model(orderItem).Where("id = ?", orderItem.ID).Updates(orderItem).Error; err != nil {
		return nil, apperrors.HandlePgError("order_item", err)
	}
	return orderItem, nil
}

func (oi *OrderItemDB) Update(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	return oi.Save(orderItem)
}

func (oi *OrderItemDB) FindByID(orderItemID id.UUID) (*entity.OrderItem, error) {
	var orderItem entity.OrderItem
	err := oi.db.WithContext(context.Background()).Preload("Product").First(&orderItem, "id = ?", orderItemID).Error
	if err != nil {
		return nil, apperrors.HandlePgError("order_item", err)
	}
	return &orderItem, nil
}

func (oi *OrderItemDB) FindByOrder(orderID id.UUID) ([]*entity.OrderItem, error) {
	var orderItems []*entity.OrderItem
	if err := oi.db.WithContext(context.Background()).Preload("Product").Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		return nil, apperrors.HandlePgError("order_item", err)
	}
	return orderItems, nil
}

func (oi *OrderItemDB) DeleteByID(orderItemID id.UUID) error {
	if err := oi.db.WithContext(context.Background()).Delete(&entity.OrderItem{}, "id = ?", orderItemID).Error; err != nil {
		return apperrors.HandlePgError("order_item", err)
	}
	return nil
}
