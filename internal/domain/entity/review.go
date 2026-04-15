package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type Review struct {
	BaseModel

	UserID    id.UUID
	ProductID id.UUID

	Rating int `gorm:"check:rating >= 1 AND rating <= 5"`

	Comment string
}
