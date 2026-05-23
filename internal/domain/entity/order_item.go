package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type OrderItem struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	OrderID   id.UUID `gorm:"index;not null"`
	ProductID id.UUID `gorm:"index;not null"`
	Product   Product `gorm:"foreignKey:ProductID"`

	Quantity int   `gorm:"not null"`
	PriceBTC int64 `gorm:"not null"`
}

func NewOrderItem(productID id.UUID, quantity int, priceBTC int64) (*OrderItem, error) {
	if quantity <= 0 {
		return nil, apperrors.BadRequest("quantity must be greater than zero", nil)
	}

	now := time.Now()
	return &OrderItem{
		ID:        id.NewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
		ProductID: productID,
		Quantity:  quantity,
		PriceBTC:  priceBTC,
	}, nil
}

func (oi *OrderItem) Subtotal() int64 {
	return oi.PriceBTC * int64(oi.Quantity)
}
