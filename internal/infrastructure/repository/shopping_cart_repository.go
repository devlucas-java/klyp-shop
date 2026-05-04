package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ShoppingCartRepository interface {
	FindByUserID(userID id.UUID) (*entity.ShoppingCart, error)
	Create(cart *entity.ShoppingCart) (*entity.ShoppingCart, error)
	Updates(cart *entity.ShoppingCart) (*entity.ShoppingCart, error)
	Save(cart *entity.ShoppingCart) (*entity.ShoppingCart, error)
	DeleteByID(uuid id.UUID) error
}
