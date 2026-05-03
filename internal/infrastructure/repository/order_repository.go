package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type OrderRepository interface {
	Create(order *entity.Order) (*entity.Order, error)
	Save(order *entity.Order) (*entity.Order, error)
	Update(order *entity.Order) (*entity.Order, error)
	Updates(order *entity.Order) (*entity.Order, error)
	FindByID(id id.UUID) (*entity.Order, error)
	FindByUser(userID id.UUID) ([]*entity.Order, error)
	FindAll() ([]*entity.Order, error)
	DeleteByID(id id.UUID) error
}
