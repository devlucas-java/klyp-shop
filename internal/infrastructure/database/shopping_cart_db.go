package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type ShoppingCartDB struct {
	db *gorm.DB
}

func NewShoppingCartDB(db *gorm.DB) repository.ShoppingCartRepository {
	return &ShoppingCartDB{db: db}
}

func (s *ShoppingCartDB) FindByUserID(userID id.UUID) (*entity.ShoppingCart, error) {
	var cart entity.ShoppingCart
	err := s.db.WithContext(context.Background()).
		Preload("Items.Product").
		Where("user_id = ?", userID).
		First(&cart).Error
	if err != nil {
		return nil, apperrors.HandlePgError("shopping_cart", err)
	}
	return &cart, nil
}

func (s *ShoppingCartDB) FindByID(cartID id.UUID) (*entity.ShoppingCart, error) {
	var cart entity.ShoppingCart
	err := s.db.WithContext(context.Background()).
		Preload("Items.Product").
		Where("id = ?", cartID).
		First(&cart).Error
	if err != nil {
		return nil, apperrors.HandlePgError("shopping_cart", err)
	}
	return &cart, nil
}

func (s *ShoppingCartDB) FindCartsByProductID(productID id.UUID) ([]*entity.ShoppingCart, error) {
	var carts []*entity.ShoppingCart
	err := s.db.WithContext(context.Background()).
		Preload("Items.Product").
		Joins("JOIN shopping_cart_items ON shopping_cart_items.cart_id = shopping_carts.id").
		Where("shopping_cart_items.product_id = ?", productID).
		Find(&carts).Error
	if err != nil {
		return nil, apperrors.HandlePgError("shopping_cart", err)
	}
	return carts, nil
}

func (s *ShoppingCartDB) Create(cart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	err := s.db.WithContext(context.Background()).Create(cart).Error
	if err != nil {
		return nil, apperrors.HandlePgError("shopping_cart", err)
	}
	return cart, nil
}

func (s *ShoppingCartDB) Save(cart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	err := s.db.WithContext(context.Background()).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(cart).Error
	if err != nil {
		return nil, apperrors.HandlePgError("shopping_cart", err)
	}
	return cart, nil
}

func (s *ShoppingCartDB) DeleteByID(cartID id.UUID) error {
	err := s.db.WithContext(context.Background()).
		Where("id = ?", cartID).
		Delete(&entity.ShoppingCart{}).Error
	if err != nil {
		return apperrors.HandlePgError("shopping_cart", err)
	}
	return nil
}
