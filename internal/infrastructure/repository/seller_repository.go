package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type SellerRepository interface {
	Save(seller *entity.Seller) (*entity.Seller, error)
	Updates(seller *entity.Seller) (*entity.Seller, error)
	FindByID(ID id.UUID) (*entity.Seller, error)
	Find(page, size int, order, search string) ([]*entity.Seller, error)
}
