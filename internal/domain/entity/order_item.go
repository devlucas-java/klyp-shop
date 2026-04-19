package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type OrderItem struct {
	BaseModel

	OrderID   id.UUID `gorm:"index;not null"`
	ProductID id.UUID `gorm:"index;not null"`
	Product   Product

	Quantity int     `gorm:"not null"`
	PriceBTC float64 `gorm:"not null"`
}

func NewOrderItem(productID id.UUID, quantity int, priceBTC float64) *OrderItem {
	return &OrderItem{
		ProductID: productID,
		Quantity:  quantity,
		PriceBTC:  priceBTC,
	}
}

func (oi *OrderItem) Subtotal() float64 {
	return oi.PriceBTC * float64(oi.Quantity)
}
