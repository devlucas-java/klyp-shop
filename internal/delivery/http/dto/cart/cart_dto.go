package cart

import "github.com/devlucas-java/klyp-shop/internal/domain/apperrors"

type AddShoppingCartItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type UpdateShoppingCartItemRequest struct {
	Quantity int `json:"quantity"`
}

type ShoppingCartItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	PriceBTC  float64 `json:"price_btc"`
	Subtotal  float64 `json:"subtotal"`
}

type ShoppingCartResponse struct {
	ID        string                     `json:"id"`
	UserID    string                     `json:"user_id"`
	TotalBTC  float64                    `json:"total_btc"`
	Items     []ShoppingCartItemResponse `json:"items"`
	CreatedAt string                     `json:"created_at,omitempty"`
	UpdatedAt string                     `json:"updated_at,omitempty"`
}

func (r *AddShoppingCartItemRequest) Validate() error {
	if r.ProductID == "" {
		return apperrors.Validation("product_id is required")
	}
	if r.Quantity <= 0 {
		return apperrors.Validation("quantity must be greater than 0")
	}
	return nil
}

func (r *UpdateShoppingCartItemRequest) Validate() error {
	if r.Quantity <= 0 {
		return apperrors.Validation("quantity must be greater than 0")
	}
	return nil
}
