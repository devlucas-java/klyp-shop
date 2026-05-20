package btcpay

type WebhookPayload struct {
	DeliveryID         string `json:"deliveryId"`
	WebhookID          string `json:"webhookId"`
	OriginalDeliveryID string `json:"originalDeliveryId"`
	IsRedelivery       bool   `json:"isRedelivery"`
	Type               string `json:"type"`
	Timestamp          int64  `json:"timestamp"`
	StoreID            string `json:"storeId"`
	InvoiceID          string `json:"invoiceId"`
	Metadata           struct {
		OrderID string `json:"orderId"`
	} `json:"metadata"`
}
