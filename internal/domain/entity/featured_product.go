package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

const MaxFeaturedProducts = 10

type FeaturedProduct struct {
	ID        id.UUID   `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	SellerID  id.UUID `gorm:"uniqueIndex:idx_featured_seller_product;not null"`
	ProductID id.UUID `gorm:"uniqueIndex:idx_featured_seller_product;not null"`

	Position int `gorm:"not null;check:position >= 1 AND position <= 10"`

	Product Product `gorm:"foreignKey:ProductID;references:ID"`
}

func NewFeaturedProduct(sellerID, productID id.UUID, position int) (*FeaturedProduct, error) {
	if position < 1 || position > 10 {
		return nil, errors.ErrBadRequest("position must be between 1 and 10", nil)
	}

	now := time.Now()
	return &FeaturedProduct{
		ID:        id.NewUUID(),
		CreatedAt: now,
		UpdatedAt: now,
		SellerID:  sellerID,
		ProductID: productID,
		Position:  position,
	}, nil
}

func (f *FeaturedProduct) SetPosition(position int) error {
	if position < 1 || position > 10 {
		return errors.ErrBadRequest("position must be between 1 and 10", nil)
	}
	f.Position = position
	f.UpdatedAt = time.Now()
	return nil
}
