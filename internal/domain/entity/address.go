package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type Address struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID id.UUID `gorm:"index;not null"`

	Street   string `gorm:"size:200;not null"`
	City     string `gorm:"size:100;not null"`
	State    string `gorm:"size:100;not null"`
	Country  string `gorm:"size:100;not null"`
	Number   int32
	Postcode string `gorm:"size:20;not null"`
}

func NewAddress(userID id.UUID, street, city, state, country, postCode string, number int32) *Address {
	now := time.Now()
	return &Address{
		ID:        id.NewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		Street:    street,
		City:      city,
		State:     state,
		Country:   country,
		Number:    number,
		Postcode:  postCode,
	}
}
