package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
)

// SellerPolicy contém as regras de negócio para vendedores.
type SellerPolicy struct{}

func NewSellerPolicy() *SellerPolicy {
	return &SellerPolicy{}
}

// CanBecomeSeller verifica se o usuário pode se tornar seller.
func (p *SellerPolicy) CanBecomeSeller(user *entity.User) error {
	if user.IsSeller {
		return errors.ErrConflict("Seller", nil)
	}
	return nil
}

// CanManage verifica se o usuário autenticado é o seller.
func (p *SellerPolicy) CanManage(user *entity.User) error {
	return user.EnsureSeller()
}
