package entity

import (
	"github.com/devlucas-java/klyp-shop/pkg/password_encoder"
)

type User struct {
	BaseModel

	Name     string `gorm:"size:120;not null"`
	Email    string `gorm:"size:120;uniqueIndex;not null"`
	Username string `gorm:"size:120;uniqueIndex;not null"`
	Password string `gorm:"size:255;not null"`
	TxHash   string `gorm:"index"`

	IsSeller bool `gorm:"default:false"`

	Addresses []Address
	Roles     []Role `gorm:"many2many:user_roles"`
}

func NewUser(name, email, username, pass string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Username: username,
		Password: password_encoder.Encoder(pass),
	}
}
