package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ShoppingCartItemRepository interface {
	FindByID(itemID id.UUID) (*entity.ShoppingCartItem, error)
	FindByCartID(cartID id.UUID) ([]*entity.ShoppingCartItem, error)
	Create(item *entity.ShoppingCartItem) (*entity.ShoppingCartItem, error)
	Save(item *entity.ShoppingCartItem) (*entity.ShoppingCartItem, error)
	DeleteByID(itemID id.UUID) error
}
