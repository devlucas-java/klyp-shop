package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type AddressRepository interface {
	Create(address *entity.Address) (*entity.Address, error)
	Update(address *entity.Address) (*entity.Address, error)
	FindByID(id id.UUID) (*entity.Address, error)
	FindByUser(userID id.UUID) ([]*entity.Address, error)
	Delete(id id.UUID) error
}
