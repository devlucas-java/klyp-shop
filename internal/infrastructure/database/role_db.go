package database

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *entity.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) FindByName(name string) (*entity.Role, error) {
	var role entity.Role

	err := r.db.
		Preload("Authorities").
		First(&role, "name = ?", name).Error

	return &role, err
}

func (r *roleRepository) FindAll() ([]entity.Role, error) {
	var roles []entity.Role

	err := r.db.
		Preload("Authorities").
		Find(&roles).Error

	return roles, err
}
