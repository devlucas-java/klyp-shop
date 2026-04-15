package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type Order struct {
	BaseModel

	UserID id.UUID
	User   User

	Status string // pending, paid, shipped, delivered

	TotalBTC float64

	Items []OrderItem
}
