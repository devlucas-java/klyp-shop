package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type OrderItem struct {
	BaseModel

	OrderID   id.UUID
	ProductID id.UUID

	Quantity int
	PriceBTC float64
}
