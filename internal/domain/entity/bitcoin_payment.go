package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type Payment struct {
	BaseModel

	OrderID id.UUID

	WalletAddress string
	TxHash        string `gorm:"index"`

	AmountBTC float64

	Status string // pending, confirmed, failed
}
