package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
	FindByID(id id.UUID) (*entity.User, error)
	FindByEmailOrUsername(str string) (*entity.User, error)
	DeleteByID(id id.UUID) error
}
