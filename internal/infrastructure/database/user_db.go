package database

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"gorm.io/gorm"
)

type userRepository struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewUserRepository(db *gorm.DB, log *logger.Logger) repository.UserRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) Create(user *entity.User) (*entity.User, error) {
	r.log.Tracef("Creating new user in database: %s", user.ID)
	if err := r.db.Create(user).Error; err != nil {
		r.log.Errorf("Database error creating user: %v", err)
		return nil, errors.Wrap("DB_CREATE_ERROR", "failed to create user in database", 500, err)
	}
	return user, nil
}

func (r *userRepository) Update(user *entity.User) (*entity.User, error) {
	r.log.Tracef("Updating user in database: %s", user.ID)
	if err := r.db.Save(user).Error; err != nil {
		r.log.Errorf("Database error updating user %s: %v", user.ID, err)
		return nil, errors.Wrap("DB_UPDATE_ERROR", "failed to update user in database", 500, err)
	}
	return user, nil
}

func (r *userRepository) Delete(userID id.UUID) error {
	r.log.Tracef("Deleting user from database: %s", userID)
	err := r.db.Where("id = ?", userID).Delete(&entity.User{}).Error
	if err != nil {
		r.log.Errorf("Database error deleting user %s: %v", userID, err)
		return errors.Wrap("DB_DELETE_ERROR", "failed to delete user from database", 500, err)
	}
	return nil
}

func (r *userRepository) FindByID(userID id.UUID) (*entity.User, error) {
	r.log.Tracef("Finding user by ID: %s", userID)
	var user entity.User
	err := r.db.First(&user, "id = ?", userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		r.log.Errorf("Database error finding user by ID %s: %v", userID, err)
		return nil, errors.Wrap("DB_FIND_ERROR", "failed to find user by ID", 500, err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmailOrUsername(str string) (*entity.User, error) {
	r.log.Tracef("Finding user by email or username: %s", str)
	var user entity.User
	err := r.db.
		Where("email = ? OR username = ?", str, str).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		r.log.Errorf("Database error finding user by email or username %s: %v", str, err)
		return nil, errors.Wrap("DB_FIND_ERROR", "failed to find user by email or username", 500, err)
	}
	return &user, nil
}
