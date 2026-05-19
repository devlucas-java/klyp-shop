package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type FeaturedProductDB struct {
	db *gorm.DB
}

func NewFeaturedProductDB(db *gorm.DB) repository.FeaturedProductRepository {
	return &FeaturedProductDB{db: db}
}

func (r *FeaturedProductDB) Add(featured *entity.FeaturedProduct) (*entity.FeaturedProduct, error) {
	if err := r.db.WithContext(context.Background()).Create(featured).Error; err != nil {
		return nil, domainErr.ErrDatabase("failed to add featured product", err)
	}
	return featured, nil
}

func (r *FeaturedProductDB) Remove(sellerID, productID id.UUID) error {
	result := r.db.WithContext(context.Background()).
		Where("seller_id = ? AND product_id = ?", sellerID, productID).
		Delete(&entity.FeaturedProduct{})
	if result.Error != nil {
		return fmt.Errorf("failed to remove featured product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domainErr.ErrNotFound("FeaturedProduct", nil)
	}
	return nil
}

func (r *FeaturedProductDB) FindBySellerID(sellerID id.UUID) ([]*entity.FeaturedProduct, error) {
	var featured []*entity.FeaturedProduct
	err := r.db.WithContext(context.Background()).
		Preload("Product").
		Where("seller_id = ?", sellerID).
		Order("position asc").
		Find(&featured).Error
	if err != nil {
		return nil, domainErr.ErrDatabase("failed to find featured products", err)
	}
	return featured, nil
}

func (r *FeaturedProductDB) FindBySellerIDAndProductID(sellerID, productID id.UUID) (*entity.FeaturedProduct, error) {
	var featured entity.FeaturedProduct
	err := r.db.WithContext(context.Background()).
		Where("seller_id = ? AND product_id = ?", sellerID, productID).
		First(&featured).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("FeaturedProduct", err)
		}
		return nil, domainErr.ErrDatabase("failed to find featured product", err)
	}
	return &featured, nil
}

func (r *FeaturedProductDB) CountBySellerID(sellerID id.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(context.Background()).Model(&entity.FeaturedProduct{}).
		Where("seller_id = ?", sellerID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count featured products: %w", err)
	}
	return count, nil
}

func (r *FeaturedProductDB) UpdatePosition(sellerID, productID id.UUID, position int) error {
	err := r.db.WithContext(context.Background()).Model(&entity.FeaturedProduct{}).
		Where("seller_id = ? AND product_id = ?", sellerID, productID).
		Update("position", position).Error
	if err != nil {
		return domainErr.ErrDatabase("failed to update featured product position", err)
	}
	return nil
}
