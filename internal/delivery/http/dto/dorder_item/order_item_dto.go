package dorderitem

import "github.com/devlucas-java/klyp-shop/internal/domain/errors"

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	PriceBTC  float64 `json:"price_btc"`
	Subtotal  float64 `json:"subtotal"`
}

func (r *OrderItemRequest) Validate() error {
	if r.ProductID == "" {
		return errors.ErrBadRequest("product_id is required", nil)
	}
	if r.Quantity <= 0 {
		return errors.ErrBadRequest("quantity must be greater than 0", nil)
	}
	return nil
}
