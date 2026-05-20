package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

const MaxAddressesPerUser = 3

type AddressPolicy struct{}

func NewAddressPolicy() *AddressPolicy {
	return &AddressPolicy{}
}

func (p *AddressPolicy) CanCreate(existing []*entity.Address) error {
	if len(existing) >= MaxAddressesPerUser {
		return errors.ErrUnprocessable("maximum number of addresses (3) reached", nil)
	}
	return nil
}

func (p *AddressPolicy) CanModify(address *entity.Address, userID id.UUID) error {
	if address.UserID != userID {
		return errors.ErrForbidden(nil)
	}
	return nil
}
