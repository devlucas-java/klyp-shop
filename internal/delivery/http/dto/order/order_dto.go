package order

import dorderitem "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/order_item"

type OrderResponse struct {
	ID        string                         `json:"id"`
	UserID    string                         `json:"user_id"`
	AddressID string                         `json:"address_id"`
	Status    string                         `json:"status"`
	TotalBTC  float64                        `json:"total_btc"`
	Items     []dorderitem.OrderItemResponse `json:"items"`
	CreatedAt string                         `json:"created_at,omitempty"`
	UpdatedAt string                         `json:"updated_at,omitempty"`
}
