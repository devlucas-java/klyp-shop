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

// IsOwnedBy verifica se o pedido pertence ao usuário informado.
func (o *Order) IsOwnedBy(userID id.UUID) bool {
	return o.UserID == userID
}

// EnsureOwnedBy retorna erro se o pedido não pertencer ao usuário.
func (o *Order) EnsureOwnedBy(userID id.UUID) error {
	if !o.IsOwnedBy(userID) {
		return errors.ErrForbidden(fmt.Errorf("order does not belong to user"))
	}
	return nil
}

// IsPending retorna true se o pedido estiver em status pendente.
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// CanBePaidBy verifica se o usuário pode pagar este pedido.
func (o *Order) CanBePaidBy(userID id.UUID) error {
	if err := o.EnsureOwnedBy(userID); err != nil {
		return err
	}
	if o.Status != OrderStatusPending {
		return errors.ErrConflict("Order", fmt.Errorf("order is not in pending status"))
	}
	return nil
}

// MarkAsPaid transiciona o pedido para o status pago.
func (o *Order) MarkAsPaid() {
	o.Status = OrderStatusPaid
	o.UpdatedAt = time.Now()
}

// MarkAsShipped transiciona o pedido para o status enviado.
func (o *Order) MarkAsShipped() {
	o.Status = OrderStatusShipped
	o.UpdatedAt = time.Now()
}

// MarkAsDelivered transiciona o pedido para o status entregue.
func (o *Order) MarkAsDelivered() {
	o.Status = OrderStatusDelivered
	o.UpdatedAt = time.Now()
}

// CancelPending cancela o pedido apenas se estiver pendente.
func (o *Order) CancelPending() error {
	if o.Status != OrderStatusPending {
		return errors.ErrConflict("Order", fmt.Errorf("only pending orders can be cancelled"))
	}
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

// SetOrderIDForItems propaga o ID do pedido para todos os itens.
func (o *Order) SetOrderIDForItems() {
	for i := range o.Items {
		o.Items[i].OrderID = o.ID
	}
}
