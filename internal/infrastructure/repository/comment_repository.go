package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type CommentRepository interface {
	Create(comment *entity.Comment) (*entity.Comment, error)
	Save(comment *entity.Comment) (*entity.Comment, error)
	Update(comment *entity.Comment) (*entity.Comment, error)
	Updates(comment *entity.Comment) (*entity.Comment, error)
	FindByID(id id.UUID) (*entity.Comment, error)
	FindByUser(userID id.UUID) ([]*entity.Comment, error)
	FindByProduct(productID id.UUID) ([]*entity.Comment, error)
	DeleteByID(id id.UUID) error
}