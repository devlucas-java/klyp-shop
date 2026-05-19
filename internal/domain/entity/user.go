package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
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

// HasRole verifica se o usuário possui a role informada.
func (u *User) HasRole(role enums.Role) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// VerifyPassword valida a senha fornecida contra o hash armazenado.
func (u *User) VerifyPassword(password string) (bool, error) {
	match, err := password_encoder.Match(password, u.Password)
	if err != nil {
		return false, errors.ErrInternal("failed to verify password", err)
	}
	return match, nil
}

// ChangePassword valida a senha atual e aplica a nova senha com hash.
func (u *User) ChangePassword(currentPassword, newPassword string) error {
	match, err := u.VerifyPassword(currentPassword)
	if err != nil {
		return err
	}
	if !match {
		return errors.ErrInvalidCredentials(nil)
	}

	hash, err := password_encoder.Encoder(newPassword)
	if err != nil {
		return errors.ErrInternal("failed to encode password", err)
	}

	u.Password = hash
	return nil
}

// ChangeName atualiza o nome do usuário se não for vazio.
func (u *User) ChangeName(name string) {
	if name != "" {
		u.Name = name
	}
}

// ChangeEmail atualiza o email do usuário se não for vazio.
func (u *User) ChangeEmail(email string) {
	if email != "" {
		u.Email = email
	}
}

// ChangeUsername atualiza o username do usuário se não for vazio.
func (u *User) ChangeUsername(username string) {
	if username != "" {
		u.Username = username
	}
}

// EnsureSeller retorna erro se o usuário não for um seller ativo.
func (u *User) EnsureSeller() error {
	if !u.IsSeller || u.Seller == nil {
		return errors.ErrNotFound("Seller", nil)
	}
	return nil
}

// MarkAsSeller marca o usuário como seller.
func (u *User) MarkAsSeller() error {
	if u.IsSeller {
		return errors.ErrConflict("Seller", nil)
	}
	u.IsSeller = true
	return nil
}

// UnmarkAsSeller remove o status de seller do usuário.
func (u *User) UnmarkAsSeller() error {
	if !u.IsSeller {
		return errors.ErrConflict("Seller", nil)
	}
	u.IsSeller = false
	return nil
}

// PromoteToAdmin aplica a role de admin ao usuário.
// Pré-condição: validar com UserPolicy.CanPromoteToAdmin antes de chamar.
func (u *User) PromoteToAdmin() {
	u.Roles = []enums.Role{enums.ADMIN}
}

// DemoteToUser remove roles especiais e retorna o usuário ao estado padrão.
// Pré-condição: validar com UserPolicy.CanDemoteToUser antes de chamar.
func (u *User) DemoteToUser() {
	u.IsSeller = false
	u.Roles = []enums.Role{enums.USER}
}
