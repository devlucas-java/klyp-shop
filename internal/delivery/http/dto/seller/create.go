package seller

import "github.com/devlucas-java/klyp-shop/internal/domain/apperrors"

type CreateSeller struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}

func (r *CreateSeller) Validate() error {
	if len(r.DisplayName) < 3 || len(r.DisplayName) > 120 {
		return apperrors.BadRequest("display_name must be between 3 and 120 characters", nil)
	}
	if len(r.Bio) > 500 {
		return apperrors.BadRequest("bio must not exceed 500 characters", nil)
	}
	return nil
}
