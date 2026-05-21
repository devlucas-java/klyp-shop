package port

type InvoiceResult struct {
	ID           string
	Status       string
	CheckoutLink string
	AmountSats   int64
}

type WebhookEvent struct {
	Type      string
	InvoiceID string
	OrderID   string
}

type PaymentGateway interface {
	CreateInvoice(orderID string, amountSats int64) (*InvoiceResult, error)
	GetInvoice(invoiceID string) (*InvoiceResult, error)
	ParseWebhook(rawBody []byte, signature string) (*WebhookEvent, error)
}
