package database

import (
	"context"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"gorm.io/gorm"
)

type ShoppingCartDB struct {
	DB  *gorm.DB
	log *logger.Logger
}

func NewShoppingCartDB(db *gorm.DB, log *logger.Logger) repository.ShoppingCartRepository {
	return &ShoppingCartDB{
		DB:  db,
		log: log,
	}
}

func (s *ShoppingCartDB) FindByUserID(userID id.UUID) (*entity.ShoppingCart, error) {

	var cart entity.ShoppingCart

	err := s.DB.WithContext(context.Background()).Preload("Items.Product").
		Where("user_id = ?", userID).
		First(&cart).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		s.log.Errorf("Failed to get shopping cart for user %s: %v", userID, err)
		return nil, err
	}
	return &cart, nil
}

func (s *ShoppingCartDB) Create(cart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	err := s.DB.WithContext(context.Background()).Create(cart).Error
	if err != nil {
		s.log.Errorf("Failed to create shopping cart: %v", err)
		return nil, err
	}
	return cart, nil
}

func (s *ShoppingCartDB) Updates(cart *entity.ShoppingCart) (*entity.ShoppingCart, error) {

	err := s.DB.WithContext(context.Background()).Session(&gorm.Session{FullSaveAssociations: true}).Updates(cart).Error
	if err != nil {
		s.log.Errorf("Failed to update shopping cart %s: %v", cart.ID, err)
		return nil, err
	}
	return cart, nil
}

func (s *ShoppingCartDB) DeleteByID(uuid id.UUID) error {
	err := s.DB.WithContext(context.Background()).Where("id = ?", uuid).Delete(&entity.ShoppingCart{}).Error
	if err != nil {
		s.log.Errorf("Failed to delete shopping cart %s: %v", uuid, err)
		return err
	}
	return nil
}

func (s *ShoppingCartDB) FindByID(uuid id.UUID) (*entity.ShoppingCart, error) {

	var cart entity.ShoppingCart

	err := s.DB.WithContext(context.Background()).Preload("Items.Product").
		Where("id = ?", uuid).
		First(&cart).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		s.log.Errorf("Failed to get shopping cart %s: %v", uuid, err)
		return nil, err
	}
	return &cart, nil
}

func (s *ShoppingCartDB) Save(cart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	err := s.DB.WithContext(context.Background()).Session(&gorm.Session{FullSaveAssociations: true}).Save(cart).Error
	if err != nil {
		s.log.Errorf("Failed to save shopping cart %s: %v", cart.ID, err)
		return nil, err
	}
	return cart, nil
}
