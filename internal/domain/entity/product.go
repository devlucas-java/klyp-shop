package entity

import (
	"strconv"
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

const productEntity = "product_entity.Product"

type Product struct {
	ID        id.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`

	PriceBTC float64 `gorm:"not null"`

	Stock int `gorm:"default:0"`

	SellerID id.UUID `gorm:"index;not null"`

	IsTop10 bool `gorm:"default:false"`

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
		return nil, apperrors.BadRequest(productEntity+".new_product: price must be greater than zero", nil)
	}
	if stock < 0 {
		return nil, apperrors.BadRequest(productEntity+".new_product: stock cannot be negative", nil)
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

func (p *Product) AddTop10(size int64) error {
	if p.IsTop10 {
		return apperrors.BadRequest(productEntity+".add_top10: product is already in top 10", nil)
	}
	if size >= 10 {
		return apperrors.BadRequest(productEntity+".add_top10: seller already has "+strconv.FormatInt(size, 10)+" products in top 10", nil)
	}
	p.IsTop10 = true
	return nil
}

// UpdateDetails aplica as alterações de nome, descrição, preço, estoque e categorias.
func (p *Product) UpdateDetails(name, description string, priceBTC float64, stock int, categories []string) error {
	if priceBTC <= 0 {
		return apperrors.BadRequest(productEntity+".update_details: price must be greater than zero", nil)
	}
	if stock < 0 {
		return apperrors.BadRequest(productEntity+".update_details: stock cannot be negative", nil)
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
		return apperrors.BadRequest(productEntity+".decrement_stock: quantity must be greater than zero", nil)
	}
	if p.Stock < quantity {
		return apperrors.Unprocessable(productEntity+".decrement_stock: insufficient stock", nil)
	}
	p.Stock -= quantity
	return nil
}
