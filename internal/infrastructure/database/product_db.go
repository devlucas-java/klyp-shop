package database

import (
	"context"
	"errors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type ProductDB struct {
	db *gorm.DB
}

func NewProductDB(db *gorm.DB) repository.ProductRepository {
	return &ProductDB{db: db}
}

func (r *ProductDB) Create(product *entity.Product) (*entity.Product, error) {
	if err := r.db.WithContext(context.Background()).Create(product).Error; err != nil {
		return nil, domainErr.ErrDatabase("failed to create product", err)
	}
	return product, nil
}

func (r *ProductDB) Save(product *entity.Product) (*entity.Product, error) {
	if err := r.db.WithContext(context.Background()).Where("id = ?", product.ID).Save(product).Error; err != nil {
		return nil, domainErr.ErrDatabase("failed to save product", err)
	}
	return product, nil
}

func (r *ProductDB) Updates(product *entity.Product) (*entity.Product, error) {
	if err := r.db.WithContext(context.Background()).Model(product).Where("id = ?", product.ID).Updates(product).Error; err != nil {
		return nil, domainErr.ErrDatabase("failed to update product", err)
	}
	var saved entity.Product
	if err := r.db.WithContext(context.Background()).First(&saved, "id = ?", product.ID).Error; err != nil {
		return nil, domainErr.ErrDatabase("failed to reload product", err)
	}
	return &saved, nil
}

func (r *ProductDB) FindByID(productID id.UUID) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(context.Background()).Preload("Reviews").First(&product, "id = ?", productID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("Product", err)
		}
		return nil, domainErr.ErrDatabase("failed to find product", err)
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
		return nil, domainErr.ErrDatabase("failed to find products by seller", err)
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
		return nil, domainErr.ErrDatabase("failed to search products", err)
	}
	return products, nil
}

func (r *ProductDB) DeleteByID(productID id.UUID) error {
	if err := r.db.WithContext(context.Background()).Delete(&entity.Product{}, "id = ?", productID).Error; err != nil {
		return domainErr.ErrDatabase("failed to delete product", err)
	}
	return nil
}
