package repository

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Save(user *entity.User) (*entity.User, error)
	Updates(user *entity.User) (*entity.User, error)
	FindByID(id uuid.UUID) (*entity.User, error)
	FindByEmailOrUsername(str string) (*entity.User, error)
}
