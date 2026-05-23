package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type ShoppingCartItemDB struct {
	db *gorm.DB
}

func NewShoppingCartItemDB(db *gorm.DB) repository.ShoppingCartItemRepository {
	return &ShoppingCartItemDB{db: db}
}

func (s *ShoppingCartItemDB) FindByID(itemID id.UUID) (*entity.ShoppingCartItem, error) {
	var item entity.ShoppingCartItem
	err := s.db.WithContext(context.Background()).
		Preload("Product").
		Where("id = ?", itemID).
		First(&item).Error
	if err != nil {
		return nil, apperrors.HandlePgError("shopping_cart_item", err)
	}
	return &item, nil
}

func (s *ShoppingCartItemDB) FindByCartID(cartID id.UUID) ([]*entity.ShoppingCartItem, error) {
	var items []*entity.ShoppingCartItem
	err := s.db.WithContext(context.Background()).
		Preload("Product").
		Where("cart_id = ?", cartID).
		Find(&items).Error
	if err != nil {
		return nil, apperrors.HandlePgError("shopping_cart_item", err)
	}
	return items, nil
}

func (s *ShoppingCartItemDB) Create(item *entity.ShoppingCartItem) (*entity.ShoppingCartItem, error) {
	err := s.db.WithContext(context.Background()).Create(item).Error
	if err != nil {
		return nil, apperrors.HandlePgError("shopping_cart_item", err)
	}
	return item, nil
}

func (s *ShoppingCartItemDB) Save(item *entity.ShoppingCartItem) (*entity.ShoppingCartItem, error) {
	err := s.db.WithContext(context.Background()).Save(item).Error
	if err != nil {
		return nil, apperrors.HandlePgError("shopping_cart_item", err)
	}
	return item, nil
}

func (s *ShoppingCartItemDB) DeleteByID(itemID id.UUID) error {
	err := s.db.WithContext(context.Background()).
		Where("id = ?", itemID).
		Delete(&entity.ShoppingCartItem{}).Error
	if err != nil {
		return apperrors.HandlePgError("shopping_cart_item", err)
	}
	return nil
}
