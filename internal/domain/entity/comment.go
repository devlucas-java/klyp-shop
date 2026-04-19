package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type Comment struct {
	BaseModel

	UserID    id.UUID `gorm:"index;not null"`
	ProductID id.UUID `gorm:"index;not null"`

	Content string `gorm:"size:2000;not null"`
}

func NewComment(userID, productID id.UUID, content string) *Comment {
	return &Comment{
		UserID:    userID,
		ProductID: productID,
		Content:   content,
	}
}
