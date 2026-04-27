package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type Comment struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID    id.UUID `gorm:"index;not null"`
	ProductID id.UUID `gorm:"index;not null"`

	Content string `gorm:"size:2000;not null"`
}

func NewComment(userID, productID id.UUID, content string) *Comment {
	now := time.Now()
	return &Comment{
		ID:        id.NewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		ProductID: productID,
		Content:   content,
	}
}
