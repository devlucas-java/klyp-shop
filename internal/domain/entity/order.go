package entity

import (
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	BaseModel

	UserID id.UUID `gorm:"index;not null"`
	User   User

	AddressID id.UUID `gorm:"index;not null"`
	Address   Address

	Status   OrderStatus `gorm:"default:'pending'"`
	TotalBTC float64

	Items          []OrderItem
	BitcoinPayment *BitcoinPayment
}

func NewOrder(userID, addressID id.UUID, items []OrderItem) *Order {
	var total float64
	for _, item := range items {
		total += item.PriceBTC * float64(item.Quantity)
	}

	return &Order{
		UserID:    userID,
		AddressID: addressID,
		Status:    OrderStatusPending,
		TotalBTC:  total,
		Items:     items,
	}
}

func (o *Order) MarkAsPaid() {
	o.Status = OrderStatusPaid
}

func (o *Order) MarkAsShipped() {
	o.Status = OrderStatusShipped
}

func (o *Order) MarkAsDelivered() {
	o.Status = OrderStatusDelivered
}

func (o *Order) Cancel() {
	o.Status = OrderStatusCancelled
}

func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}
