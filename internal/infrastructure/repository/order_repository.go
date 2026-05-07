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
	FindAllWithDetails() ([]*entity.Order, error)
	FindAllPaginated(page, size int, status string) ([]*entity.Order, int64, error)
	FindBySellerID(sellerID id.UUID) ([]*entity.Order, error)
	FindBySellerIDPaginated(sellerID id.UUID, page, size int, status string) ([]*entity.Order, int64, error)
	DeleteByID(id id.UUID) error
}
