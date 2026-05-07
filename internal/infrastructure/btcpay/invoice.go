package btcpay

type CreateInvoiceRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Metadata struct {
		OrderID string `json:"orderId"`
	} `json:"metadata"`
	Checkout struct {
		SpeedPolicy       string `json:"speedPolicy"`
		ExpirationMinutes int    `json:"expirationMinutes"`
	} `json:"checkout"`
}

type InvoiceResponse struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	CheckoutLink   string `json:"checkoutLink"`
	ExpirationTime int64  `json:"expirationTime"`
	Metadata       struct {
		OrderID string `json:"orderId"`
	} `json:"metadata"`
}
