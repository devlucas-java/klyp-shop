package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type SellerRepository interface {
	Create(seller *entity.Seller) (*entity.Seller, error)
	Update(seller *entity.Seller) (*entity.Seller, error)
	FindByID(id id.UUID) (*entity.Seller, error)
	Find(page, size int, order, search string) ([]*entity.Seller, error)
	DeleteByID(id id.UUID) error
}
