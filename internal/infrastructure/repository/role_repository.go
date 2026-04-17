package repository

import "github.com/devlucas-java/klyp-shop/internal/domain/entity"

type RoleRepository interface {
	Create(role *entity.Role) error
	FindByName(name string) (*entity.Role, error)
	FindAll() ([]entity.Role, error)
}
