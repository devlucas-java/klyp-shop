package mapper

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dproduct"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type ProductMaper struct {
}

func NewProductMapper() *ProductMaper {
	return &ProductMaper{}
}

func (p *ProductMaper) ProductToProductResponse(product *entity.Product) *dproduct.ProductResponse {
	return &dproduct.ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		PriceBTC:    product.PriceBTC,
		Stock:       product.Stock,
		SellerID:    product.SellerID.String(),
		Categories:  product.Categories,
	}
}

func (p *ProductMaper) CreateProductToProduct(dto *dproduct.CreateProduct) *entity.Product {
	return entity.NewProduct(dto.Name, dto.Description, dto.PriceBTC, dto.Stock, dto.Categories)
}

func (p *ProductMaper) UpdateProductToProduct(dto *dproduct.UpdateProduct) *entity.Product {
	return entity.NewProduct(dto.Name, dto.Description, dto.PriceBTC, dto.Stock, dto.Categories)
}
