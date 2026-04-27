package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type Product struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`

	PriceBTC float64 `gorm:"not null"`

	Stock int `gorm:"default:0"`

	SellerID id.UUID `gorm:"index;not null"`
	Seller   Seller

	Reviews    []Review
	Categories []string `gorm:"serializer:json"`
}

func NewProduct(
	name string,
	description string,
	priceBTC float64,
	stock int,
	sellerID id.UUID,
	categories []string,
) *Product {
	now := time.Now()
	return &Product{
		ID:          id.NewUUID(),
		CreatedAt:   now,
		UpdatedAt:   now,
		Name:        name,
		Description: description,
		PriceBTC:    priceBTC,
		Stock:       stock,
		SellerID:    sellerID,
		Categories:  categories,
	}
}
