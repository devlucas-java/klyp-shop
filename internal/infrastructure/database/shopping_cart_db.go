package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"gorm.io/gorm"
)

const shoppingCartDB = "shopping_cart_db.ShoppingCartDB"

type ShoppingCartDB struct {
	DB  *gorm.DB
	log *logger.Logger
}

func NewShoppingCartDB(db *gorm.DB) repository.ShoppingCartRepository {
	return &ShoppingCartDB{DB: db}
}

func (s *ShoppingCartDB) FindByUserID(userID id.UUID) (*entity.ShoppingCart, error) {
	var cart entity.ShoppingCart
	err := s.DB.WithContext(context.Background()).Preload("Items.Product").
		Where("user_id = ?", userID).
		First(&cart).Error
	if err != nil {
		return nil, apperrors.HandlePgError(shoppingCartDB+".find_by_user_id", err)
	}
	return &cart, nil
}

func (s *ShoppingCartDB) Create(cart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	err := s.DB.WithContext(context.Background()).Create(cart).Error
	if err != nil {
		return nil, apperrors.HandlePgError(shoppingCartDB+".create", err)
	}
	return cart, nil
}

func (s *ShoppingCartDB) Updates(cart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	err := s.DB.WithContext(context.Background()).Session(&gorm.Session{FullSaveAssociations: true}).Updates(cart).Error
	if err != nil {
		return nil, apperrors.HandlePgError(shoppingCartDB+".updates", err)
	}
	return cart, nil
}

func (s *ShoppingCartDB) DeleteByID(uuid id.UUID) error {
	err := s.DB.WithContext(context.Background()).Where("id = ?", uuid).Delete(&entity.ShoppingCart{}).Error
	if err != nil {
		return apperrors.HandlePgError(shoppingCartDB+".delete_by_id", err)
	}
	return nil
}

func (s *ShoppingCartDB) FindByID(uuid id.UUID) (*entity.ShoppingCart, error) {
	var cart entity.ShoppingCart
	err := s.DB.WithContext(context.Background()).Preload("Items.Product").
		Where("id = ?", uuid).
		First(&cart).Error
	if err != nil {
		return nil, apperrors.HandlePgError(shoppingCartDB+".find_by_id", err)
	}
	return &cart, nil
}

func (s *ShoppingCartDB) Save(cart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	err := s.DB.WithContext(context.Background()).Session(&gorm.Session{FullSaveAssociations: true}).Save(cart).Error
	if err != nil {
		return nil, apperrors.HandlePgError(shoppingCartDB+".save", err)
	}
	return cart, nil
}
