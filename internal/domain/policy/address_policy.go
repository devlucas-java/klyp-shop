package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

const MaxAddressesPerUser = 3

type AddressPolicy struct{}

func NewAddressPolicy() *AddressPolicy {
	return &AddressPolicy{}
}

func (p *AddressPolicy) CanCreate(existing []*entity.Address) error {
	if len(existing) >= MaxAddressesPerUser {
		return apperrors.Policy("you have reached the maximum number of addresses allowed (3)")
	}
	return nil
}

func (p *AddressPolicy) CanModify(address *entity.Address, userID id.UUID) error {
	if address.UserID != userID {
		return apperrors.Unauthorized("address_policy.AddressPolicy.can_modify: address does not belong to user", nil)
	}
	return nil
}
