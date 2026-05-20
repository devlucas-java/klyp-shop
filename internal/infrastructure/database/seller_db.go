package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"gorm.io/gorm"
)

type SellerDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewSellerDB(db *gorm.DB, log *logger.Logger) repository.SellerRepository {
	return &SellerDB{db: db, log: log}
}

func (r *SellerDB) Create(seller *entity.Seller) (*entity.Seller, error) {
	if err := r.db.WithContext(context.Background()).Create(seller).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to create seller")
	}
	return seller, nil
}

func (r *SellerDB) Save(seller *entity.Seller) (*entity.Seller, error) {
	if err := r.db.WithContext(context.Background()).Where("id = ?", seller.ID).Save(seller).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to save seller")
	}
	return seller, nil
}

func (r *SellerDB) Updates(seller *entity.Seller) (*entity.Seller, error) {
	if err := r.db.WithContext(context.Background()).Model(seller).Where("id = ?", seller.ID).Updates(seller).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to update seller")
	}
	var saved entity.Seller
	if err := r.db.WithContext(context.Background()).First(&saved, "id = ?", seller.ID).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to reload seller")
	}
	return &saved, nil
}

func (r *SellerDB) FindByID(sellerID id.UUID) (*entity.Seller, error) {
	var seller entity.Seller
	err := r.db.WithContext(context.Background()).First(&seller, "id = ?", sellerID).Error
	if err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to find seller")
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

	query := r.db.WithContext(context.Background()).
		Limit(size).
		Offset((page - 1) * size).
		Order("created_at " + order)

	if search != "" {
		query = query.Where("display_name LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&sellers).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to list sellers")
	}
	return sellers, nil
}

func (r *SellerDB) DeleteByID(sellerID id.UUID) error {
	err := r.db.WithContext(context.Background()).Where("id = ?", sellerID).Delete(&entity.Seller{}).Error
	if err != nil {
		return errors.HandlePgError(r.log, err, "failed to delete seller")
	}
	return nil
}
