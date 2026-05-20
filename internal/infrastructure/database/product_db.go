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

type ProductDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewProductDB(db *gorm.DB, log *logger.Logger) repository.ProductRepository {
	return &ProductDB{db: db, log: log}
}

func (r *ProductDB) Create(product *entity.Product) (*entity.Product, error) {
	if err := r.db.WithContext(context.Background()).Create(product).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to create product")
	}
	return product, nil
}

func (r *ProductDB) Save(product *entity.Product) (*entity.Product, error) {
	if err := r.db.WithContext(context.Background()).Where("id = ?", product.ID).Save(product).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to save product")
	}
	return product, nil
}

func (r *ProductDB) Updates(product *entity.Product) (*entity.Product, error) {
	if err := r.db.WithContext(context.Background()).Model(product).Where("id = ?", product.ID).Updates(product).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to update product")
	}
	var saved entity.Product
	if err := r.db.WithContext(context.Background()).First(&saved, "id = ?", product.ID).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to reload product")
	}
	return &saved, nil
}

func (r *ProductDB) FindByID(productID id.UUID) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(context.Background()).Preload("Reviews").First(&product, "id = ?", productID).Error
	if err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to find product")
	}
	return &product, nil
}

func (r *ProductDB) FindBySellerID(sellerID id.UUID, page, size int) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(context.Background()).
		Preload("Reviews").
		Where("seller_id = ?", sellerID).
		Limit(size).
		Offset((page - 1) * size).
		Order("created_at desc").
		Find(&products).Error
	if err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to find products by seller")
	}
	return products, nil
}

func (r *ProductDB) Search(page, size int, order, search string, categories []string) ([]*entity.Product, error) {
	var products []*entity.Product

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	query := r.db.WithContext(context.Background()).
		Limit(size).
		Offset((page - 1) * size).
		Order("created_at " + order)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if len(categories) > 0 {
		query = query.Where("categories ?| array[?]", categories)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, errors.HandlePgError(r.log, err, "failed to search products")
	}
	return products, nil
}

func (r *ProductDB) DeleteByID(productID id.UUID) error {
	if err := r.db.WithContext(context.Background()).Delete(&entity.Product{}, "id = ?", productID).Error; err != nil {
		return errors.HandlePgError(r.log, err, "failed to delete product")
	}
	return nil
}
