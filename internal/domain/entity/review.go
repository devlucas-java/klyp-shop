package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type Review struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID    id.UUID `gorm:"index;not null"`
	ProductID id.UUID `gorm:"index;not null"`

	Rating  int    `gorm:"check:rating >= 1 AND rating <= 5;not null"`
	Comment string `gorm:"size:1000"`
}

func NewReview(userID, productID id.UUID, rating int, comment string) *Review {
	now := time.Now()
	return &Review{
		ID:        id.NewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		ProductID: productID,
		Rating:    rating,
		Comment:   comment,
	}
}
