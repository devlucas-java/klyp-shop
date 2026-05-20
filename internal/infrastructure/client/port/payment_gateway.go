package port

// InvoiceResult é o resultado de uma invoice criada ou consultada no gateway de pagamento.
// AmountSats é o valor em satoshis (1 BTC = 100_000_000 satoshis).
type InvoiceResult struct {
	ID           string
	Status       string
	CheckoutLink string
	AmountSats   int64
}

// WebhookEvent representa um evento normalizado recebido via webhook.
// O adapter de infra é responsável por fazer o parse do payload bruto e preencher esta struct.
type WebhookEvent struct {
	Type      string
	InvoiceID string
	OrderID   string
}

// PaymentGateway é a port de saída (driven port) definida pela aplicação.
// Qualquer gateway de pagamento (BTCPay, Strike, etc.) deve implementar esta interface.
// O PaymentService depende apenas daqui — nunca de um client concreto.
type PaymentGateway interface {
	// CreateInvoice cria uma invoice. amountSats é o valor em satoshis.
	CreateInvoice(orderID string, amountSats int64) (*InvoiceResult, error)
	GetInvoice(invoiceID string) (*InvoiceResult, error)
	ParseWebhook(rawBody []byte, signature string) (*WebhookEvent, error)
}
