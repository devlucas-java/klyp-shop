package database

import (
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
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

	if err := r.db.Create(user).Error; err != nil {
		r.log.Errorf("Database error creating duser: %v", err)
		return nil, fmt.Errorf("failed to create duser in database: %w", err)
	}

	return user, nil
}

func (r *UserDB) Update(user *entity.User) (*entity.User, error) {

	if err := r.db.Save(user).Error; err != nil {
		r.log.Errorf("Database error updating duser %s: %v", user.ID, err)
		return nil, fmt.Errorf("failed to update duser in database: %w", err)
	}

	return user, nil
}

func (r *UserDB) DeleteByID(userID id.UUID) error {

	err := r.db.Where("id = ?", userID).Delete(&entity.User{}).Error

	if err != nil {
		r.log.Errorf("Database error deleting duser %s: %v", userID, err)
		return fmt.Errorf("failed to delete duser from database: %w", err)
	}
	return nil
}

func (r *UserDB) FindByID(userID id.UUID) (*entity.User, error) {

	var user entity.User

	err := r.db.First(&user, "id = ?", userID).Error
	if err != nil {
		r.log.Errorf("Database error finding duser by ID %s: %v", userID, err)
		return nil, fmt.Errorf("failed to find duser by ID: %w", err)
	}
	return &user, nil
}

func (r *UserDB) FindByEmailOrUsername(str string) (*entity.User, error) {

	var user entity.User
	err := r.db.
		Where("email = ? OR username = ?", str, str).
		First(&user).Error

	if err != nil {
		r.log.Errorf("Database error finding duser by email or username %s: %v", str, err)
		return nil, fmt.Errorf("failed to find duser by email or username: %w", err)
	}
	return &user, nil
}
