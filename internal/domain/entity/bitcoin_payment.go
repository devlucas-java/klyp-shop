package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusConfirmed PaymentStatus = "confirmed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

// BitcoinPayment representa um pagamento Bitcoin.
// AmountSats é o valor em satoshis (1 BTC = 100_000_000 satoshis).
type BitcoinPayment struct {
	ID        id.UUID   `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	OrderID id.UUID `gorm:"uniqueIndex;not null"`

	WalletAddress string        `gorm:"not null"`
	TxHash        string        `gorm:"index"`
	AmountSats    int64         `gorm:"not null"`
	Status        PaymentStatus `gorm:"default:'pending'"`
}

func NewBitcoinPayment(orderID id.UUID, walletAddress string, amountSats int64) *BitcoinPayment {
	now := time.Now()
	return &BitcoinPayment{
		ID:            id.NewUUID(),
		CreatedAt:     now,
		UpdatedAt:     now,
		OrderID:       orderID,
		WalletAddress: walletAddress,
		AmountSats:    amountSats,
		Status:        PaymentStatusPending,
	}
}

func (p *BitcoinPayment) Confirm(txHash string) {
	p.TxHash = txHash
	p.Status = PaymentStatusConfirmed
}

func (p *BitcoinPayment) Fail() {
	p.Status = PaymentStatusFailed
}

func (p *BitcoinPayment) IsConfirmed() bool {
	return p.Status == PaymentStatusConfirmed
}
