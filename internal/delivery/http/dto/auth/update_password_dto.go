package auth

import "github.com/devlucas-java/klyp-shop/internal/domain/apperrors"

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (r *UpdatePasswordRequest) Validate() error {
	if r.CurrentPassword == "" {
		return apperrors.BadRequest("current_password is required", nil)
	}
	if r.NewPassword == "" {
		return apperrors.BadRequest("new_password is required", nil)
	}
	if len(r.NewPassword) < 6 {
		return apperrors.BadRequest("new_password must be at least 6 characters", nil)
	}
	if r.CurrentPassword == r.NewPassword {
		return apperrors.BadRequest("new_password must differ from current_password", nil)
	}
	return nil
}
