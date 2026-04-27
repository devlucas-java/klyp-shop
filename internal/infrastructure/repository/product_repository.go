package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ProductRepository interface {
	Create(product *entity.Product) (*entity.Product, error)
	Update(product *entity.Product) (*entity.Product, error)
	FindByID(id id.UUID) (*entity.Product, error)
	FindBySellerID(sellerID id.UUID, page, size int) ([]*entity.Product, error)
	Search(page, size int, order, search string, categories []string) ([]*entity.Product, error)
	DeleteByID(id id.UUID) error
}
