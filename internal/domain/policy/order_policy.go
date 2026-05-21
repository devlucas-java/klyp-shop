package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

const orderPolicy = "order_policy.OrderPolicy"

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
		return apperrors.Conflict(orderPolicy+".can_cancel: only pending orders can be cancelled", nil)
	}
	return nil
}

func (p *OrderPolicy) CanPay(order *entity.Order, userID id.UUID) error {
	return order.CanBePaidBy(userID)
}

func (p *OrderPolicy) AddressBelongsToUser(address *entity.Address, userID id.UUID) error {
	if address.UserID != userID {
		return apperrors.Forbidden(orderPolicy+".address_belongs_to_user: address does not belong to user", nil)
	}
	return nil
}
