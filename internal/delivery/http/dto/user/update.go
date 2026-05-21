package user

import "github.com/devlucas-java/klyp-shop/internal/domain/apperrors"

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (r *UpdateUserRequest) Validate() error {
	if r.Name == "" && r.Email == "" && r.Username == "" {
		return apperrors.Validation("at least one field (name, email or username) must be provided")
	}
	if r.Name != "" && len(r.Name) > 120 {
		return apperrors.Validation("name must not exceed 120 characters")
	}
	if r.Email != "" && len(r.Email) > 120 {
		return apperrors.Validation("email must not exceed 120 characters")
	}
	if r.Username != "" && len(r.Username) > 120 {
		return apperrors.Validation("username must not exceed 120 characters")
	}
	return nil
}
