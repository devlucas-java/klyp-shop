package entity

import (
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type Review struct {
	BaseModel

	UserID    id.UUID `gorm:"index;not null"`
	ProductID id.UUID `gorm:"index;not null"`

	Rating  int    `gorm:"check:rating >= 1 AND rating <= 5;not null"`
	Comment string `gorm:"size:1000"`
}

func NewReview(userID, productID id.UUID, rating int, comment string) *Review {
	return &Review{
		UserID:    userID,
		ProductID: productID,
		Rating:    rating,
		Comment:   comment,
	}
}
