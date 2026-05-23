package entity

import (
	"fmt"
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type Product struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`

	PriceBTC int64 `gorm:"not null"`

	Stock int `gorm:"default:0"`

	SellerID id.UUID `gorm:"index;not null"`

	IsTop10 bool `gorm:"default:false"`

	Reviews    []Review
	Categories []string `gorm:"serializer:json"`
}

func NewProduct(
	name string,
	description string,
	priceBTC int64,
	stock int,
	categories []string,
) (*Product, error) {
	if priceBTC <= 0 {
		return nil, apperrors.BadRequest("price must be greater than zero", nil)
	}
	if stock < 0 {
		return nil, apperrors.BadRequest("stock cannot be negative", nil)
	}
	now := time.Now()
	return &Product{
		ID:          id.NewUUID(),
		CreatedAt:   now,
		UpdatedAt:   now,
		Name:        name,
		Description: description,
		PriceBTC:    priceBTC,
		Stock:       stock,
		Categories:  categories,
	}, nil
}

func (p *Product) IsOwnedBy(sellerID id.UUID) bool {
	return p.SellerID == sellerID
}

func (p *Product) AddTop10(size int64) error {
	if p.IsTop10 {
		return apperrors.BadRequest("this product is already in the top 10 list", nil)
	}
	if size >= 10 {
		return apperrors.BadRequest(fmt.Sprintf("you already have %d products in the top 10 list", size), nil)
	}
	p.IsTop10 = true
	return nil
}

func (p *Product) UpdateDetails(name, description string, priceBTC int64, stock int, categories []string) error {
	if priceBTC <= 0 {
		return apperrors.BadRequest("price must be greater than zero", nil)
	}
	if stock < 0 {
		return apperrors.BadRequest("stock cannot be negative", nil)
	}
	if name != "" {
		p.Name = name
	}
	if description != "" {
		p.Description = description
	}
	p.PriceBTC = priceBTC
	p.Stock = stock
	if len(categories) > 0 {
		p.Categories = categories
	}
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) DecrementStock(quantity int) error {
	if quantity <= 0 {
		return apperrors.BadRequest("quantity must be greater than zero", nil)
	}
	if p.Stock < quantity {
		return apperrors.Unprocessable("insufficient stock for this product", nil)
	}
	p.Stock -= quantity
	return nil
}
