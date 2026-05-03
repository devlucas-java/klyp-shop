package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type OrderItemRepository interface {
	Create(orderItem *entity.OrderItem) (*entity.OrderItem, error)
	Save(orderItem *entity.OrderItem) (*entity.OrderItem, error)
	Update(orderItem *entity.OrderItem) (*entity.OrderItem, error)
	Updates(orderItem *entity.OrderItem) (*entity.OrderItem, error)
	FindByID(id id.UUID) (*entity.OrderItem, error)
	FindByOrder(orderID id.UUID) ([]*entity.OrderItem, error)
	DeleteByID(id id.UUID) error
}
