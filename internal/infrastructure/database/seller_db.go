package database

import (
	"errors"
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type SellerDB struct {
	db *gorm.DB
}

func NewSellerDB(db *gorm.DB) repository.SellerRepository {
	return &SellerDB{db: db}
}

func (r *SellerDB) Create(seller *entity.Seller) (*entity.Seller, error) {
	if err := r.db.Create(seller).Error; err != nil {
		return nil, fmt.Errorf("failed to create seller: %w", err)
	}
	return seller, nil
}

func (r *SellerDB) Save(seller *entity.Seller) (*entity.Seller, error) {
	if err := r.db.Where("id = ?", seller.ID).Save(seller).Error; err != nil {
		return nil, fmt.Errorf("failed to save seller: %w", err)
	}
	return seller, nil
}

func (r *SellerDB) Updates(seller *entity.Seller) (*entity.Seller, error) {
	if err := r.db.Model(seller).Where("id = ?", seller.ID).Updates(seller).Error; err != nil {
		return nil, fmt.Errorf("failed to update seller: %w", err)
	}
	var saved entity.Seller
	if err := r.db.First(&saved, "id = ?", seller.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload seller: %w", err)
	}
	return &saved, nil
}

func (r *SellerDB) FindByID(sellerID id.UUID) (*entity.Seller, error) {
	var seller entity.Seller
	err := r.db.First(&seller, "id = ?", sellerID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("Seller", err)
		}
		return nil, fmt.Errorf("failed to find seller: %w", err)
	}
	return &seller, nil
}

func (r *SellerDB) Find(page, size int, order, search string) ([]*entity.Seller, error) {
	var sellers []*entity.Seller

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	query := r.db.
		Limit(size).
		Offset((page - 1) * size).
		Order("created_at " + order)

	if search != "" {
		query = query.Where("display_name LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&sellers).Error; err != nil {
		return nil, fmt.Errorf("failed to list sellers: %w", err)
	}
	return sellers, nil
}

func (r *SellerDB) DeleteByID(sellerID id.UUID) error {
	result := r.db.Where("id = ?", sellerID).Delete(&entity.Seller{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete seller: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domainErr.ErrNotFound("Seller", nil)
	}
	return nil
}
