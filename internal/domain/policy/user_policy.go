package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
)

// UserPolicy contém as regras de negócio para gerenciamento de usuários.
type UserPolicy struct{}

func NewUserPolicy() *UserPolicy {
	return &UserPolicy{}
}

// CanPromoteToAdmin verifica se o usuário pode ser promovido a admin.
func (p *UserPolicy) CanPromoteToAdmin(user *entity.User) error {
	if user.IsSeller {
		return errors.ErrInvalidRole("seller cannot be promoted to admin", nil)
	}
	if user.HasRole(enums.ADMIN) {
		return errors.ErrInvalidRole("user is already an admin", nil)
	}
	return nil
}

// CanDemoteToUser verifica se o usuário pode ser rebaixado a user comum.
func (p *UserPolicy) CanDemoteToUser(user *entity.User) error {
	if user.IsSeller {
		return errors.ErrInvalidRole("seller cannot be demoted", nil)
	}
	if user.HasRole(enums.USER) {
		return errors.ErrInvalidRole("user is already a regular user", nil)
	}
	return nil
}
