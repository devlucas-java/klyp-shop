package database

import (
	"context"
	"errors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"gorm.io/gorm"
)

type OrderItemDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewOrderItemDB(db *gorm.DB, log *logger.Logger) repository.OrderItemRepository {
	return &OrderItemDB{db: db, log: log}
}

func (oi *OrderItemDB) Create(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	if err := oi.db.WithContext(context.Background()).Create(orderItem).Error; err != nil {
		oi.log.Errorf("OrderItemDB.Create: %v", err)
		return nil, domainErr.ErrDatabase("failed to create order item", err)
	}
	return orderItem, nil
}

func (oi *OrderItemDB) Save(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	if err := oi.db.WithContext(context.Background()).Where("id = ?", orderItem.ID).Save(orderItem).Error; err != nil {
		oi.log.Errorf("OrderItemDB.Save %s: %v", orderItem.ID, err)
		return nil, domainErr.ErrDatabase("failed to save order item", err)
	}
	return orderItem, nil
}

func (oi *OrderItemDB) Updates(orderItem *entity.OrderItem) (*entity.OrderItem, error) {
	if err := oi.db.WithContext(context.Background()).Model(orderItem).Where("id = ?", orderItem.ID).Updates(orderItem).Error; err != nil {
		oi.log.Errorf("OrderItemDB.Updates %s: %v", orderItem.ID, err)
		return nil, domainErr.ErrDatabase("failed to update order item", err)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("OrderItem", err)
		}
		oi.log.Errorf("OrderItemDB.FindByID %s: %v", orderItemID, err)
		return nil, domainErr.ErrDatabase("failed to find order item", err)
	}
	return &orderItem, nil
}

func (oi *OrderItemDB) FindByOrder(orderID id.UUID) ([]*entity.OrderItem, error) {
	var orderItems []*entity.OrderItem
	if err := oi.db.WithContext(context.Background()).Preload("Product").Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
		oi.log.Errorf("OrderItemDB.FindByOrder %s: %v", orderID, err)
		return nil, domainErr.ErrDatabase("failed to find order items", err)
	}
	return orderItems, nil
}

func (oi *OrderItemDB) DeleteByID(orderItemID id.UUID) error {
	if err := oi.db.WithContext(context.Background()).Delete(&entity.OrderItem{}, "id = ?", orderItemID).Error; err != nil {
		oi.log.Errorf("OrderItemDB.DeleteByID %s: %v", orderItemID, err)
		return domainErr.ErrDatabase("failed to delete order item", err)
	}
	return nil
}
