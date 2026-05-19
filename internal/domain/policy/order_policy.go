package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

// OrderPolicy contém as regras de negócio para pedidos.
type OrderPolicy struct{}

func NewOrderPolicy() *OrderPolicy {
	return &OrderPolicy{}
}

// CanView verifica se o usuário pode visualizar o pedido.
func (p *OrderPolicy) CanView(order *entity.Order, userID id.UUID) error {
	return order.EnsureOwnedBy(userID)
}

// CanCancel verifica se o pedido pode ser cancelado pelo usuário.
func (p *OrderPolicy) CanCancel(order *entity.Order, userID id.UUID) error {
	if err := order.EnsureOwnedBy(userID); err != nil {
		return err
	}
	if order.Status != entity.OrderStatusPending {
		return errors.ErrConflict("Order", nil)
	}
	return nil
}

// CanPay verifica se o usuário pode pagar o pedido.
func (p *OrderPolicy) CanPay(order *entity.Order, userID id.UUID) error {
	return order.CanBePaidBy(userID)
}

// AddressBelongsToUser verifica se o endereço pertence ao usuário.
func (p *OrderPolicy) AddressBelongsToUser(address *entity.Address, userID id.UUID) error {
	if address.UserID != userID {
		return errors.ErrForbidden(nil)
	}
	return nil
}
