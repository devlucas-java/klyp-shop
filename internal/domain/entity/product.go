package entity

import (
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type Product struct {
	BaseModel

	Name        string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`

	PriceBTC float64 `gorm:"not null"`

	Stock int `gorm:"default:0"`

	SellerID id.UUID `gorm:"index;not null"`
	Seller   Seller

	Reviews    []Review
	Categories []string `gorm:"serializer:json"`
}
