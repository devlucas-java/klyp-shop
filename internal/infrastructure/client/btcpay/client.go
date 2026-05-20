package btcpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	storeID    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, storeID, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		storeID: storeID,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// CreateInvoice cria uma invoice no BTCPay. amountSats é o valor em satoshis.
func (c *Client) CreateInvoice(orderID string, amountSats int64) (*InvoiceResponse, error) {
	req := CreateInvoiceRequest{
		Amount:   amountSats,
		Currency: "BTC",
	}
	req.Metadata.OrderID = orderID
	req.Checkout.SpeedPolicy = "MediumSpeed"
	req.Checkout.ExpirationMinutes = 60

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("btcpay: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/stores/%s/invoices", c.baseURL, c.storeID)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("btcpay: build request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "token "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("btcpay: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("btcpay: server error %d: %s", resp.StatusCode, string(raw))
	}

	var invoice InvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, fmt.Errorf("btcpay: decode response: %w", err)
	}

	return &invoice, nil
}

func (c *Client) GetInvoice(invoiceID string) (*InvoiceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/stores/%s/invoices/%s", c.baseURL, c.storeID, invoiceID)
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("btcpay: build request: %w", err)
	}

	httpReq.Header.Set("Authorization", "token "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("btcpay: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("btcpay: server error %d: %s", resp.StatusCode, string(raw))
	}

	var invoice InvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, fmt.Errorf("btcpay: decode response: %w", err)
	}

	return &invoice, nil
}
