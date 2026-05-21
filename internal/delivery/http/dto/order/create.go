package order

import (
	dorderitem "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order_item"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
)

type CreateOrderRequest struct {
	AddressID string                        `json:"address_id"`
	Items     []dorderitem.OrderItemRequest `json:"items"`
}

func (r *CreateOrderRequest) Validate() error {
	if r.AddressID == "" {
		return apperrors.BadRequest("address_id is required", nil)
	}

	if len(r.Items) == 0 {
		return apperrors.BadRequest("at least one item must be provided", nil)
	}

	for _, item := range r.Items {
		if item.ProductID == "" {
			return apperrors.BadRequest("product_id is required for all items", nil)
		}
		if item.Quantity <= 0 {
			return apperrors.BadRequest("quantity must be greater than 0", nil)
		}
	}

	return nil
}
