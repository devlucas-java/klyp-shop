package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type SellerPolicy struct{}

func NewSellerPolicy() *SellerPolicy {
	return &SellerPolicy{}
}

func (p *SellerPolicy) CanBecomeSeller(user *entity.User) error {
	if user.IsSeller {
		return apperrors.Conflict("user is already a seller", nil)
	}
	return nil
}

func (p *SellerPolicy) CanManage(user *entity.User) error {
	return user.EnsureSeller()
}
