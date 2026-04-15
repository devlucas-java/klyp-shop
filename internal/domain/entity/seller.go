package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type Seller struct {
	BaseModel

	UserID id.UUID `gorm:"uniqueIndex"`

	DisplayName string
	Bio         string

	Products []Product
}
