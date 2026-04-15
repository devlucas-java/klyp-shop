package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type Product struct {
	BaseModel

	Name        string
	Description string

	PriceBTC float64

	Stock int

	SellerID id.UUID `gorm:"index"`
	Seller   Seller

	Reviews []Review
}
