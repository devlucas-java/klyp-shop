package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type OrderPolicy struct{}

func NewOrderPolicy() *OrderPolicy {
	return &OrderPolicy{}
}

func (p *OrderPolicy) CanView(order *entity.Order, userID id.UUID) error {
	return order.EnsureOwnedBy(userID)
}

func (p *OrderPolicy) CanCancel(order *entity.Order, userID id.UUID) error {
	if err := order.EnsureOwnedBy(userID); err != nil {
		return err
	}
	if order.Status != entity.OrderStatusPending {
		return errors.ErrConflict("Order", nil)
	}
	return nil
}

func (p *OrderPolicy) CanPay(order *entity.Order, userID id.UUID) error {
	return order.CanBePaidBy(userID)
}

func (p *OrderPolicy) AddressBelongsToUser(address *entity.Address, userID id.UUID) error {
	if address.UserID != userID {
		return errors.ErrForbidden(nil)
	}
	return nil
}
