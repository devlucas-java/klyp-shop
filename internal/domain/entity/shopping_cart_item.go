package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/pkg/id"
)

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

func NewShoppingCartItem(cartID, productID id.UUID, quantity int, priceBTC float64) *ShoppingCartItem {
	now := time.Now()
	return &ShoppingCartItem{
		ID:        id.NewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
		CartID:    cartID,
		ProductID: productID,
		Quantity:  quantity,
		PriceBTC:  priceBTC,
	}
}

func (item *ShoppingCartItem) Subtotal() float64 {
	return item.PriceBTC * float64(item.Quantity)
}
