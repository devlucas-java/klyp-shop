package database

import (
	"errors"
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
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
		r.log.Errorf("UserDB.Create: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (r *UserDB) Save(user *entity.User) (*entity.User, error) {
	if err := r.db.Model(user).Where("id = ?", user.ID).
		Select("name", "email", "username", "password", "is_seller", "roles", "updated_at").
		Updates(map[string]interface{}{
			"name":       user.Name,
			"email":      user.Email,
			"username":   user.Username,
			"password":   user.Password,
			"is_seller":  user.IsSeller,
			"roles":      user.Roles,
			"updated_at": user.UpdatedAt,
		}).Error; err != nil {
		r.log.Errorf("UserDB.Save %s: %v", user.ID, err)
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	return user, nil
}

func (r *UserDB) Updates(user *entity.User) (*entity.User, error) {
	if err := r.db.Model(user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		r.log.Errorf("UserDB.Updates %s: %v", user.ID, err)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

func (r *UserDB) Update(user *entity.User) (*entity.User, error) {
	return r.Save(user)
}

func (r *UserDB) DeleteByID(userID id.UUID) error {
	if err := r.db.Where("id = ?", userID).Delete(&entity.User{}).Error; err != nil {
		r.log.Errorf("UserDB.DeleteByID %s: %v", userID, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *UserDB) FindByID(userID id.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("User", err)
		}
		r.log.Errorf("UserDB.FindByID %s: %v", userID, err)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

func (r *UserDB) FindByIDWithSeller(userID id.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Seller").First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("User", err)
		}
		r.log.Errorf("UserDB.FindByIDWithSeller %s: %v", userID, err)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

func (r *UserDB) FindByEmailOrUsername(str string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ? OR username = ?", str, str).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("User", err)
		}
		r.log.Errorf("UserDB.FindByEmailOrUsername %s: %v", str, err)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}
