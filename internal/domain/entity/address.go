package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type Address struct {
	BaseModel

	UserID id.UUID `gorm:"index;not null"`

	Street   string `gorm:"size:200;not null"`
	City     string `gorm:"size:100;not null"`
	State    string `gorm:"size:100;not null"`
	Country  string `gorm:"size:100;not null"`
	Number   int32
	Postcode string `gorm:"size:20;not null"`

	IsDefault bool `gorm:"default:false"`
}

func NewAddress(userID id.UUID, street, city, state, country, postCode string, number int32) *Address {
	return &Address{
		UserID:   userID,
		Street:   street,
		City:     city,
		State:    state,
		Country:  country,
		Number:   number,
		Postcode: postCode,
	}
}
