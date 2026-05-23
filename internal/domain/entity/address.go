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

func (a *Address) ChangeStreet(str string) {
	if str == "" {
		return
	}
	a.Street = str
	a.UpdatedAt = time.Now()
}
func (a *Address) ChangeCity(city string) {
	if city == "" {
		return
	}
	a.City = city
	a.UpdatedAt = time.Now()
}

func (a *Address) ChangeState(state string) {
	if state == "" {
		return
	}
	a.State = state
	a.UpdatedAt = time.Now()
}

func (a *Address) ChangeCountry(country string) {
	if country == "" {
		return
	}
	a.Country = country
	a.UpdatedAt = time.Now()
}

func (a *Address) ChangeNumber(number int32) {
	if number <= 0 {
		return
	}
	a.Number = number
	a.UpdatedAt = time.Now()
}

func (a *Address) ChangePostcode(postcode string) {
	if postcode == "" {
		return
	}
	a.Postcode = postcode
	a.UpdatedAt = time.Now()
}
