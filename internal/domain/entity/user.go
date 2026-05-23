package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
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
		return nil, apperrors.Internal(err)
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

func (u *User) VerifyPassword(password string) (bool, error) {
	match, err := password_encoder.Match(password, u.Password)
	if err != nil {
		return false, apperrors.Internal(err)
	}
	return match, nil
}

func (u *User) ChangePassword(currentPassword, newPassword string) error {
	match, err := u.VerifyPassword(currentPassword)
	if err != nil {
		return err
	}
	if !match {
		return apperrors.InvalidCredentials(nil)
	}

	hash, err := password_encoder.Encoder(newPassword)
	if err != nil {
		return apperrors.Internal(err)
	}

	u.Password = hash
	return nil
}

func (u *User) ChangeName(name string) {
	if name != "" {
		u.Name = name
	}
}

func (u *User) ChangeEmail(email string) {
	if email != "" {
		u.Email = email
	}
}

func (u *User) ChangeUsername(username string) {
	if username != "" {
		u.Username = username
	}
}

func (u *User) EnsureSeller() error {
	if !u.IsSeller || u.Seller == nil {
		return apperrors.Unprocessable("this action requires a seller account", nil)
	}
	return nil
}

func (u *User) MarkAsSeller() error {
	if u.IsSeller {
		return apperrors.Conflict("user is already a seller", nil)
	}
	u.IsSeller = true
	return nil
}

func (u *User) UnmarkAsSeller() error {
	if !u.IsSeller {
		return apperrors.Conflict("user is not a seller", nil)
	}
	u.IsSeller = false
	return nil
}

func (u *User) ChangerToAdmin() {
	u.Roles = []enums.Role{enums.ADMIN}
}

func (u *User) ChangerToSuperAdmin() {
	u.Roles = []enums.Role{enums.ADMIN, enums.USER, enums.SELLER}
}

func (u *User) ChangerToSeller() {
	u.Roles = []enums.Role{enums.SELLER}
}

func (u *User) ChangerToUser() {
	u.IsSeller = false
	u.Roles = []enums.Role{enums.USER}
}
