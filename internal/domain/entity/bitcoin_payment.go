package entity

import (
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusConfirmed PaymentStatus = "confirmed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type BitcoinPayment struct {
	BaseModel

	OrderID id.UUID `gorm:"uniqueIndex;not null"`

	WalletAddress string        `gorm:"not null"`
	TxHash        string        `gorm:"index"`
	AmountBTC     float64       `gorm:"not null"`
	Status        PaymentStatus `gorm:"default:'pending'"`
}

func NewBitcoinPayment(orderID id.UUID, walletAddress string, amountBTC float64) *BitcoinPayment {
	return &BitcoinPayment{
		OrderID:       orderID,
		WalletAddress: walletAddress,
		AmountBTC:     amountBTC,
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
