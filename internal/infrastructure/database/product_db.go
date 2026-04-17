package database

import (
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"gorm.io/gorm"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *entity.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(productID id.UUID) (*entity.Product, error) {
	var product entity.Product

	err := r.db.
		Preload("Seller").
		Preload("Reviews").
		First(&product, "id = ?", productID).Error

	return &product, err
}

func (r *productRepository) FindAll() ([]entity.Product, error) {
	var products []entity.Product

	err := r.db.
		Preload("Seller").
		Find(&products).Error

	return products, err
}

func (r *productRepository) Update(product *entity.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(productID id.UUID) error {
	return r.db.Delete(&entity.Product{}, "id = ?", productID).Error
}
