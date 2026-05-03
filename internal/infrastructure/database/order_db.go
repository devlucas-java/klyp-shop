package database

import (
	"errors"
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
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
	if err := o.db.Create(order).Error; err != nil {
		o.log.Errorf("OrderDB.Create: %v", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return order, nil
}

func (o *OrderDB) Save(order *entity.Order) (*entity.Order, error) {
	if err := o.db.Where("id = ?", order.ID).Save(order).Error; err != nil {
		o.log.Errorf("OrderDB.Save %s: %v", order.ID, err)
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	return order, nil
}

func (o *OrderDB) Updates(order *entity.Order) (*entity.Order, error) {
	if err := o.db.Model(order).Where("id = ?", order.ID).Updates(order).Error; err != nil {
		o.log.Errorf("OrderDB.Updates %s: %v", order.ID, err)
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return order, nil
}

func (o *OrderDB) Update(order *entity.Order) (*entity.Order, error) {
	return o.Save(order)
}

func (o *OrderDB) FindByID(orderID id.UUID) (*entity.Order, error) {
	var order entity.Order
	err := o.db.Preload("Items").Preload("User").Preload("Address").Preload("BitcoinPayment").First(&order, "id = ?", orderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("Order", err)
		}
		o.log.Errorf("OrderDB.FindByID %s: %v", orderID, err)
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	return &order, nil
}

func (o *OrderDB) FindByUser(userID id.UUID) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := o.db.Preload("Items").Preload("User").Preload("Address").Preload("BitcoinPayment").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		o.log.Errorf("OrderDB.FindByUser %s: %v", userID, err)
		return nil, fmt.Errorf("failed to find orders: %w", err)
	}
	return orders, nil
}

func (o *OrderDB) FindAll() ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := o.db.Preload("Items").Preload("User").Preload("Address").Preload("BitcoinPayment").Find(&orders).Error; err != nil {
		o.log.Errorf("OrderDB.FindAll: %v", err)
		return nil, fmt.Errorf("failed to find orders: %w", err)
	}
	return orders, nil
}

func (o *OrderDB) DeleteByID(orderID id.UUID) error {
	if err := o.db.Delete(&entity.Order{}, "id = ?", orderID).Error; err != nil {
		o.log.Errorf("OrderDB.DeleteByID %s: %v", orderID, err)
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
