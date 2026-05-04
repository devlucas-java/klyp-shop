package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ShoppingCart struct {
	ID        id.UUID             `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time           `gorm:"autoCreateTime"`
	UpdatedAt time.Time           `gorm:"autoUpdateTime"`
	UserID    id.UUID             `gorm:"index;not null"`
	Items     []*ShoppingCartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;"`
	TotalBTC  float64             `gorm:"type:decimal(18,8);not null"`
}
