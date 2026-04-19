package entity

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/pkg/password_encoder"
)

type User struct {
	BaseModel

	Name     string `gorm:"size:120;not null"`
	Email    string `gorm:"size:120;uniqueIndex;not null"`
	Username string `gorm:"size:120;uniqueIndex;not null"`
	Password string `gorm:"size:255;not null"`

	IsSeller bool `gorm:"default:false"`

	Addresses []Address
	Orders    []Order
	Roles     []enums.Role `gorm:"type:json;serializer:json"`
}

func NewUser(name, email, username, pass string) (*User, error) {
	hash, err := password_encoder.Encoder(pass)
	if err != nil {
		return nil, err
	}
	return &User{
		Name:     name,
		Email:    email,
		Username: username,
		Password: hash,
		Roles:    []enums.Role{enums.USER},
	}, nil
}

func (u *User) HasRole(role enums.Role) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}
