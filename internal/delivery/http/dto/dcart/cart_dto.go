package dcart

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
