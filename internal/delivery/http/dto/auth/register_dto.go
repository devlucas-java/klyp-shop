package auth

import "github.com/devlucas-java/klyp-shop/internal/domain/apperrors"

type RegisterDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *RegisterDTO) Validate() error {
	if r.Name == "" {
		return apperrors.BadRequest("name is required", nil)
	}
	if len(r.Name) > 120 {
		return apperrors.BadRequest("name must not exceed 120 characters", nil)
	}
	if r.Email == "" {
		return apperrors.BadRequest("email is required", nil)
	}
	if len(r.Email) > 120 {
		return apperrors.BadRequest("email must not exceed 120 characters", nil)
	}
	if r.Username == "" {
		return apperrors.BadRequest("username is required", nil)
	}
	if len(r.Username) > 120 {
		return apperrors.BadRequest("username must not exceed 120 characters", nil)
	}
	if r.Password == "" {
		return apperrors.BadRequest("password is required", nil)
	}
	if len(r.Password) < 6 {
		return apperrors.BadRequest("password must be at least 6 characters", nil)
	}
	return nil
}
