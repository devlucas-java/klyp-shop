package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type FeaturedProductRepository interface {
	Add(featured *entity.FeaturedProduct) (*entity.FeaturedProduct, error)
	Remove(sellerID, productID id.UUID) error
	FindAll() ([]*entity.FeaturedProduct, error)
	FindBySellerID(sellerID id.UUID) ([]*entity.FeaturedProduct, error)
	FindBySellerIDAndProductID(sellerID, productID id.UUID) (*entity.FeaturedProduct, error)
	CountBySellerID(sellerID id.UUID) (int64, error)
	UpdatePosition(sellerID, productID id.UUID, position int) error
}
