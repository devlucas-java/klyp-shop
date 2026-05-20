package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
)

type SellerPolicy struct{}

func NewSellerPolicy() *SellerPolicy {
	return &SellerPolicy{}
}

func (p *SellerPolicy) CanBecomeSeller(user *entity.User) error {
	if user.IsSeller {
		return errors.ErrConflict("Seller", nil)
	}
	return nil
}

func (p *SellerPolicy) CanManage(user *entity.User) error {
	return user.EnsureSeller()
}
