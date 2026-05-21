package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

const userDB = "user_db.UserDB"

type UserDB struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) repository.UserRepository {
	return &UserDB{db: db}
}

func (r *UserDB) Create(user *entity.User) (*entity.User, error) {
	if err := r.db.WithContext(context.Background()).Create(user).Error; err != nil {
		return nil, apperrors.HandlePgError(userDB+".create", err)
	}
	return user, nil
}

func (r *UserDB) Save(user *entity.User) (*entity.User, error) {
	if err := r.db.WithContext(context.Background()).
		Model(user).
		Session(&gorm.Session{FullSaveAssociations: false}).
		Select("name", "email", "username", "password", "is_seller", "roles", "updated_at").
		Save(user).Error; err != nil {
		return nil, apperrors.HandlePgError(userDB+".save", err)
	}
	return user, nil
}

func (r *UserDB) Updates(user *entity.User) (*entity.User, error) {
	if err := r.db.WithContext(context.Background()).Model(user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, apperrors.HandlePgError(userDB+".updates", err)
	}
	return user, nil
}

func (r *UserDB) Update(user *entity.User) (*entity.User, error) {
	return r.Save(user)
}

func (r *UserDB) DeleteByID(userID id.UUID) error {
	if err := r.db.WithContext(context.Background()).Where("id = ?", userID).Delete(&entity.User{}).Error; err != nil {
		return apperrors.HandlePgError(userDB+".delete_by_id", err)
	}
	return nil
}

func (r *UserDB) FindByID(userID id.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(context.Background()).First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, apperrors.HandlePgError(userDB+".find_by_id", err)
	}
	return &user, nil
}

func (r *UserDB) FindByIDWithSeller(userID id.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(context.Background()).Preload("Seller").First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, apperrors.HandlePgError(userDB+".find_by_id_with_seller", err)
	}
	return &user, nil
}

func (r *UserDB) FindByEmailOrUsername(str string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(context.Background()).Where("email = ? OR username = ?", str, str).First(&user).Error
	if err != nil {
		return nil, apperrors.HandlePgError(userDB+".find_by_email_or_username", err)
	}
	return &user, nil
}

func (r *UserDB) ExistsUserByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(context.Background()).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, apperrors.HandlePgError(userDB+".exists_user_by_email", err)
	}
	return count > 0, nil
}

func (r *UserDB) ExistsUserByUserName(username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(context.Background()).Model(&entity.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, apperrors.HandlePgError(userDB+".exists_user_by_username", err)
	}
	return count > 0, nil
}
