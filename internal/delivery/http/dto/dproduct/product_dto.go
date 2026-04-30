package dproduct

type ProductResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	PriceBTC    float64  `json:"price_btc"`
	Stock       int      `json:"stock"`
	SellerID    string   `json:"seller_id"`
	Categories  []string `json:"categories"`
}
