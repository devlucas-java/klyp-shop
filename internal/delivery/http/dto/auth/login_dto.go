package auth

import "github.com/devlucas-java/klyp-shop/internal/domain/apperrors"

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Login == "" {
		return apperrors.BadRequest("login is required", nil)
	}
	if r.Password == "" {
		return apperrors.BadRequest("password is required", nil)
	}
	return nil
}
