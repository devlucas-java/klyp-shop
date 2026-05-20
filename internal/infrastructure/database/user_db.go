package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"gorm.io/gorm"
)

type UserDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewUserDB(db *gorm.DB, log *logger.Logger) repository.UserRepository {
	return &UserDB{db: db, log: log}
}

func (r *UserDB) Create(user *entity.User) (*entity.User, error) {

	if err := r.db.WithContext(context.Background()).Create(user).Error; err != nil {

		return nil, errors.HandlePgError(r.log, err, "failed to create user")
	}
	return user, nil
}

func (r *UserDB) Save(user *entity.User) (*entity.User, error) {

	if err := r.db.WithContext(context.Background()).
		Model(user).
		Session(&gorm.Session{FullSaveAssociations: false}).
		Select("name", "email", "username", "password", "is_seller", "roles", "updated_at").
		Save(user).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to save user")
	}
	return user, nil
}

func (r *UserDB) Updates(user *entity.User) (*entity.User, error) {

	if err := r.db.WithContext(context.Background()).Model(user).Where("id = ?", user.ID).Updates(user).Error; err != nil {

		return nil, errors.HandlePgError(r.log, err, "failed to update user")
	}
	return user, nil
}

func (r *UserDB) Update(user *entity.User) (*entity.User, error) {
	return r.Save(user)
}

func (r *UserDB) DeleteByID(userID id.UUID) error {

	if err := r.db.WithContext(context.Background()).Where("id = ?", userID).Delete(&entity.User{}).Error; err != nil {

		return errors.HandlePgError(r.log, err, "failed to delete user")
	}
	return nil
}

func (r *UserDB) FindByID(userID id.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(context.Background()).First(&user, "id = ?", userID).Error

	if err != nil {

		return nil, errors.HandlePgError(r.log, err, "failed to find user")
	}
	return &user, nil
}

func (r *UserDB) FindByIDWithSeller(userID id.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(context.Background()).Preload("Seller").First(&user, "id = ?", userID).Error

	if err != nil {

		return nil, errors.HandlePgError(r.log, err, "failed to find user")
	}
	return &user, nil
}

func (r *UserDB) FindByEmailOrUsername(str string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(context.Background()).Where("email = ? OR username = ?", str, str).First(&user).Error

	if err != nil {

		return nil, errors.HandlePgError(r.log, err, "failed to find user")
	}
	return &user, nil
}

func (r *UserDB) ExistsUserByEmail(email string) (bool, error) {
	var count int64

	if err := r.db.WithContext(context.Background()).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error; err != nil {

		return false, errors.HandlePgError(r.log, err, "failed to check existing email")
	}

	return count > 0, nil
}

func (r *UserDB) ExistsUserByUserName(username string) (bool, error) {
	var count int64

	if err := r.db.WithContext(context.Background()).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error; err != nil {

		return false, errors.HandlePgError(r.log, err, "failed to check existing username")
	}

	return count > 0, nil
}
