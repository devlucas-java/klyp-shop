package mapper

import (
	productDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/product"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type ProductMaper struct{}

func NewProductMapper() *ProductMaper {
	return &ProductMaper{}
}

func (p *ProductMaper) ProductToProductResponse(product *entity.Product) *productDTO.ProductResponse {
	return &productDTO.ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		PriceBTC:    product.PriceBTC,
		Stock:       product.Stock,
		SellerID:    product.SellerID.String(),
		Categories:  product.Categories,
	}
}
