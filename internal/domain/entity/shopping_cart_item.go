package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

const shoppingCartItemEntity = "shopping_cart_item_entity.ShoppingCartItem"

type ShoppingCartItem struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	CartID    id.UUID `gorm:"index;not null"`
	ProductID id.UUID `gorm:"index;not null"`
	Product   Product `gorm:"foreignKey:ProductID"`

	Quantity int     `gorm:"not null"`
	PriceBTC float64 `gorm:"not null"`
}

func NewShoppingCartItem(cartID, productID id.UUID, quantity int, priceBTC float64) (*ShoppingCartItem, error) {
	if quantity <= 0 {
		return nil, apperrors.BadRequest(shoppingCartItemEntity+".new_shopping_cart_item: quantity must be greater than zero", nil)
	}

	now := time.Now()
	return &ShoppingCartItem{
		ID:        id.NewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
		PriceBTC:  priceBTC,
	}, nil
}

func (item *ShoppingCartItem) SetQuantity(quantity int) error {
	if quantity <= 0 {
		return apperrors.BadRequest(shoppingCartItemEntity+".set_quantity: quantity must be greater than zero", nil)
	}
	item.Quantity = quantity
	item.UpdatedAt = time.Now()
	return nil
}

func (item *ShoppingCartItem) Subtotal() float64 {
	return item.PriceBTC * float64(item.Quantity)
}
