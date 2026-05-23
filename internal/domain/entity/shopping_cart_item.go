package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ShoppingCartItem struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	CartID    id.UUID `gorm:"index;not null"`
	ProductID id.UUID `gorm:"index;not null"`
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`

	Quantity int   `gorm:"not null"`
	PriceBTC int64 `gorm:"not null"`
}

func NewShoppingCartItem(cartID, productID id.UUID, quantity int, priceBTC int64) (*ShoppingCartItem, error) {
	if quantity <= 0 {
		return nil, apperrors.BadRequest("quantity must be greater than zero", nil)
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
		return apperrors.BadRequest("quantity must be greater than zero", nil)
	}
	item.Quantity = quantity
	item.UpdatedAt = time.Now()
	return nil
}

func (item *ShoppingCartItem) Subtotal() int64 {
	return item.PriceBTC * int64(item.Quantity)
}
