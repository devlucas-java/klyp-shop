package auth

import "github.com/devlucas-java/klyp-shop/internal/domain/apperrors"

type VerifyPasswordRequest struct {
	Password string `json:"password"`
}

func (r *VerifyPasswordRequest) Validate() error {
	if r.Password == "" {
		return apperrors.BadRequest("password is required", nil)
	}
	return nil
}
