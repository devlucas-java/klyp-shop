package entity

import (
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type Product struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`

	PriceBTC float64 `gorm:"not null"`

	Stock int `gorm:"default:0"`

	SellerID id.UUID `gorm:"index;not null"`

	Reviews    []Review
	Categories []string `gorm:"serializer:json"`
}

func NewProduct(
	name string,
	description string,
	priceBTC float64,
	stock int,
	categories []string,
) (*Product, error) {
	if priceBTC <= 0 {
		return nil, errors.ErrBadRequest("price must be greater than zero", nil)
	}
	if stock < 0 {
		return nil, errors.ErrBadRequest("stock cannot be negative", nil)
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

// IsOwnedBy verifica se o produto pertence ao seller informado.
func (p *Product) IsOwnedBy(sellerID id.UUID) bool {
	return p.SellerID == sellerID
}

// UpdateDetails aplica as alterações de nome, descrição, preço, estoque e categorias.
func (p *Product) UpdateDetails(name, description string, priceBTC float64, stock int, categories []string) error {
	if priceBTC <= 0 {
		return errors.ErrBadRequest("price must be greater than zero", nil)
	}
	if stock < 0 {
		return errors.ErrBadRequest("stock cannot be negative", nil)
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

// DecrementStock reduz o estoque após uma venda.
func (p *Product) DecrementStock(quantity int) error {
	if quantity <= 0 {
		return errors.ErrBadRequest("quantity must be greater than zero", nil)
	}
	if p.Stock < quantity {
		return errors.ErrUnprocessable("insufficient stock", nil)
	}
	p.Stock -= quantity
	return nil
}
