package entity

import (
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type Address struct {
	BaseModel

	UserID id.UUID `gorm:"index"`

	Street  string
	City    string
	State   string
	Country string

	Number   int32
	Postcode string
}

func NewAddress(street, city, state, country, postCode string, number int32) *Address {
	return &Address{
		Street:   street,
		City:     city,
		State:    state,
		Country:  country,
		Number:   number,
		Postcode: postCode,
	}
}
