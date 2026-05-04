package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/password_encoder"
)

type User struct {
	ID        id.UUID   `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Name      string    `gorm:"size:120;not null"`
	Email     string    `gorm:"size:120;uniqueIndex;not null"`
	Username  string    `gorm:"size:120;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`

	IsSeller bool `gorm:"default:false"`

	Seller *Seller `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ShoppingCart ShoppingCart `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Reviews      []Review     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Addresses []Address    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Orders    []Order      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Roles     []enums.Role `gorm:"type:json;serializer:json"`
}

func NewUser(name, email, username, pass string) (*User, error) {
	hash, err := password_encoder.Encoder(pass)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &User{
		ID:        id.NewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Email:     email,
		Username:  username,
		Password:  hash,
		Roles:     []enums.Role{enums.USER},
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
