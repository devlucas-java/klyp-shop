package order

import "github.com/devlucas-java/klyp-shop/internal/domain/apperrors"

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

func (r *UpdateOrderStatusRequest) Validate() error {
	if r.Status == "" {
		return apperrors.BadRequest("status is required", nil)
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"delivered": true,
		"cancelled": true,
	}

	if !validStatuses[r.Status] {
		return apperrors.BadRequest("invalid status value", nil)
	}

	return nil
}
