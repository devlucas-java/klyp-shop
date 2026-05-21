package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

const featuredProductDB = "featured_product_db.FeaturedProductDB"

type FeaturedProductDB struct {
	db *gorm.DB
}

func NewFeaturedProductDB(db *gorm.DB) repository.FeaturedProductRepository {
	return &FeaturedProductDB{db: db}
}

func (r *FeaturedProductDB) Add(featured *entity.FeaturedProduct) (*entity.FeaturedProduct, error) {
	if err := r.db.WithContext(context.Background()).Create(featured).Error; err != nil {
		return nil, apperrors.HandlePgError(featuredProductDB+".add", err)
	}
	return featured, nil
}

func (r *FeaturedProductDB) Remove(sellerID, productID id.UUID) error {
	err := r.db.WithContext(context.Background()).
		Where("seller_id = ? AND product_id = ?", sellerID, productID).
		Delete(&entity.FeaturedProduct{}).Error
	if err != nil {
		return apperrors.HandlePgError(featuredProductDB+".remove", err)
	}
	return nil
}

func (r *FeaturedProductDB) FindAll() ([]*entity.FeaturedProduct, error) {
	var featured []*entity.FeaturedProduct
	err := r.db.WithContext(context.Background()).
		Preload("Product").
		Order("seller_id, position asc").
		Find(&featured).Error
	if err != nil {
		return nil, apperrors.HandlePgError(featuredProductDB+".find_all", err)
	}
	return featured, nil
}

func (r *FeaturedProductDB) FindBySellerID(sellerID id.UUID) ([]*entity.FeaturedProduct, error) {
	var featured []*entity.FeaturedProduct
	err := r.db.WithContext(context.Background()).
		Preload("Product").
		Where("seller_id = ?", sellerID).
		Order("position asc").
		Find(&featured).Error
	if err != nil {
		return nil, apperrors.HandlePgError(featuredProductDB+".find_by_seller_id", err)
	}
	return featured, nil
}

func (r *FeaturedProductDB) FindBySellerIDAndProductID(sellerID, productID id.UUID) (*entity.FeaturedProduct, error) {
	var featured entity.FeaturedProduct
	err := r.db.WithContext(context.Background()).
		Where("seller_id = ? AND product_id = ?", sellerID, productID).
		First(&featured).Error
	if err != nil {
		return nil, apperrors.HandlePgError(featuredProductDB+".find_by_seller_id_and_product_id", err)
	}
	return &featured, nil
}

func (r *FeaturedProductDB) CountBySellerID(sellerID id.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(context.Background()).Model(&entity.FeaturedProduct{}).
		Where("seller_id = ?", sellerID).
		Count(&count).Error
	if err != nil {
		return 0, apperrors.HandlePgError(featuredProductDB+".count_by_seller_id", err)
	}
	return count, nil
}

func (r *FeaturedProductDB) UpdatePosition(sellerID, productID id.UUID, position int) error {
	err := r.db.WithContext(context.Background()).Model(&entity.FeaturedProduct{}).
		Where("seller_id = ? AND product_id = ?", sellerID, productID).
		Update("position", position).Error
	if err != nil {
		return apperrors.HandlePgError(featuredProductDB+".update_position", err)
	}
	return nil
}
