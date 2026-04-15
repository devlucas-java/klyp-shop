package database

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type UserDB struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) repository.UserRepository {
	return &UserDB{DB: db}
}

func (r *UserDB) Save(user *entity.User) (*entity.User, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserDB) Updates(user *entity.User) (*entity.User, error) {

	err := r.DB.
		Model(&entity.User{}).
		Where("id =? ", user.ID).
		Updates(user).Error

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserDB) FindByID(id id.UUID) (*entity.User, error) {
	var user entity.User

	err := r.DB.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserDB) FindByEmailOrUsername(str string) (*entity.User, error) {
	var user entity.User

	err := r.DB.
		Where("email = ? OR username = ?", str, str).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
