package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ReviewRepository interface {
	Create(review *entity.Review) (*entity.Review, error)
	Update(review *entity.Review) (*entity.Review, error)
	FindByProductID(productID id.UUID) ([]*entity.Review, error)
	Delete(review id.UUID) error
}
