package policy

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
)

type UserPolicy struct{}

func NewUserPolicy() *UserPolicy {
	return &UserPolicy{}
}

func (p *UserPolicy) CanPromoteToAdmin(user *entity.User) error {
	if user.IsSeller {
		return apperrors.Unprocessable("a seller account cannot be promoted to admin", nil)
	}
	if user.HasRole(enums.ADMIN) {
		return apperrors.Unprocessable("user is already an admin", nil)
	}
	return nil
}

func (p *UserPolicy) CanDemoteToUser(user *entity.User) error {
	if user.IsSeller {
		return apperrors.Unprocessable("a seller account cannot be demoted", nil)
	}
	if user.HasRole(enums.USER) {
		return apperrors.Unprocessable("user is already a regular user", nil)
	}
	return nil
}
