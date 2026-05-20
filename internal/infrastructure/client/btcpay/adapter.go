package btcpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/infrastructure/client/port"
)

// BTCPayAdapter implementa port.PaymentGateway usando o BTCPay Server.
// É o único lugar do projeto que conhece os tipos internos do pacote btcpay.
// Não importa nada de service, handler ou dto — apenas a port que a aplicação define.
type BTCPayAdapter struct {
	client *Client
	secret string
}

// NewBTCPayAdapter cria um adapter pronto para uso.
// webhookSecret pode ser vazio — nesse caso a verificação de assinatura é ignorada.
func NewBTCPayAdapter(baseURL, storeID, apiKey, webhookSecret string) *BTCPayAdapter {
	return &BTCPayAdapter{
		client: NewClient(baseURL, storeID, apiKey),
		secret: webhookSecret,
	}
}

// CreateInvoice cria uma invoice no BTCPay e mapeia para o tipo da aplicação.
// amountSats é o valor em satoshis.
func (a *BTCPayAdapter) CreateInvoice(orderID string, amountSats int64) (*port.InvoiceResult, error) {
	resp, err := a.client.CreateInvoice(orderID, amountSats)
	if err != nil {
		return nil, fmt.Errorf("btcpay adapter: create invoice: %w", err)
	}
	return &port.InvoiceResult{
		ID:           resp.ID,
		Status:       resp.Status,
		CheckoutLink: resp.CheckoutLink,
		AmountSats:   resp.Amount,
	}, nil
}

// GetInvoice consulta o status de uma invoice e mapeia para o tipo da aplicação.
func (a *BTCPayAdapter) GetInvoice(invoiceID string) (*port.InvoiceResult, error) {
	resp, err := a.client.GetInvoice(invoiceID)
	if err != nil {
		return nil, fmt.Errorf("btcpay adapter: get invoice: %w", err)
	}
	return &port.InvoiceResult{
		ID:           resp.ID,
		Status:       resp.Status,
		CheckoutLink: resp.CheckoutLink,
		AmountSats:   resp.Amount,
	}, nil
}

// ParseWebhook valida a assinatura HMAC-SHA256, faz o parse do payload BTCPay
// e retorna um WebhookEvent normalizado para a aplicação.
func (a *BTCPayAdapter) ParseWebhook(rawBody []byte, signature string) (*port.WebhookEvent, error) {
	if !a.verifySignature(rawBody, signature) {
		return nil, fmt.Errorf("btcpay adapter: invalid webhook signature")
	}

	var payload WebhookPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, fmt.Errorf("btcpay adapter: parse webhook payload: %w", err)
	}

	return &port.WebhookEvent{
		Type:      payload.Type,
		InvoiceID: payload.InvoiceID,
		OrderID:   payload.Metadata.OrderID,
	}, nil
}

func (a *BTCPayAdapter) verifySignature(body []byte, signature string) bool {
	if a.secret == "" {
		return true
	}
	mac := hmac.New(sha256.New, []byte(a.secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
