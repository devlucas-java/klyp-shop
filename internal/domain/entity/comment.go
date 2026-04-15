package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type Comment struct {
	BaseModel

	UserID    id.UUID
	ProductID id.UUID

	Content string
}
