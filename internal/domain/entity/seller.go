package entity

import "github.com/devlucas-java/klyp-shop/pkg/id"

type Seller struct {
	BaseModel

	UserID id.UUID `gorm:"uniqueIndex;not null"`
	User   User

	DisplayName string `gorm:"size:120;not null"`
	Bio         string `gorm:"size:500"`

	Products []Product
}

func NewSeller(userID id.UUID, displayName, bio string) *Seller {
	return &Seller{
		UserID:      userID,
		DisplayName: displayName,
		Bio:         bio,
	}
}

func (s *Seller) UpdateInfo(displayName, bio string) {
	if displayName != "" {
		s.DisplayName = displayName
	}
	if bio != "" {
		s.Bio = bio
	}
}
