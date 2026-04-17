package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type ReviewRepository interface {
	Create(review *entity.Review) error
	FindByProduct(productID id.UUID) ([]entity.Review, error)
	Delete(id id.UUID) error
}
