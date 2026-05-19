package entity

import (
	"fmt"
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
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
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

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

	now := time.Now()
	return &Order{
		ID:        id.NewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
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

func (o *Order) IsOwnedBy(userID id.UUID) bool {
	return o.UserID == userID
}

func (o *Order) EnsureOwnedBy(userID id.UUID) error {
	if !o.IsOwnedBy(userID) {
		return errors.ErrForbidden(fmt.Errorf("order does not belong to user"))
	}
	return nil
}

func (o *Order) CanBePaidBy(userID id.UUID) error {
	if !o.IsOwnedBy(userID) {
		return errors.ErrForbidden(fmt.Errorf("order does not belong to user"))
	}
	if o.Status != OrderStatusPending {
		return errors.ErrConflict("Order", fmt.Errorf("order is not in pending status"))
	}
	return nil
}

func (o *Order) CancelPending() error {
	if o.Status != OrderStatusPending {
		return errors.ErrConflict("Order", fmt.Errorf("order cannot be cancelled in current status"))
	}

	o.Status = OrderStatusCancelled
	return nil
}

func (o *Order) SetOrderIDForItems() {
	for i := range o.Items {
		o.Items[i].OrderID = o.ID
	}
}
