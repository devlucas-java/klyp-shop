package dorder

import "github.com/devlucas-java/klyp-shop/internal/domain/errors"

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

func (r *UpdateOrderStatusRequest) Validate() error {
	if r.Status == "" {
		return errors.ErrBadRequest("status is required", nil)
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"delivered": true,
		"cancelled": true,
	}

	if !validStatuses[r.Status] {
		return errors.ErrBadRequest("invalid status value", nil)
	}

	return nil
}
