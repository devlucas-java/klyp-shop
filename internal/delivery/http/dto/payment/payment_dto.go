package payment

type InvoiceResponse struct {
	PaymentID     string `json:"payment_id"`
	OrderID       string `json:"order_id"`
	AmountSats    int64  `json:"amount_sats"`
	Status        string `json:"status"`
	WalletAddress string `json:"wallet_address"`
	CheckoutURL   string `json:"checkout_url,omitempty"`
	InvoiceID     string `json:"invoice_id,omitempty"`
}
