package database

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
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
	err := r.db.Create(product).Error
	return product, err
}

func (r *ProductDB) FindByID(productID id.UUID) (*entity.Product, error) {
	var product entity.Product

	err := r.db.
		Preload("Seller").
		Preload("Reviews").
		First(&product, "id = ?", productID).Error

	return &product, err
}

func (r *ProductDB) Search(page, size int, order, search string, categories []string) ([]*entity.Product, error) {
	var products []*entity.Product

	offset := (page - 1) * size

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	query := r.db.
		Preload("Seller").
		Limit(size).
		Offset(offset).
		Order("created_at " + order)

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ?", like)
	}

	if len(categories) > 0 {
		query = query.Where("categories ?| array[?]", categories)
	}

	err := query.Find(&products).Error
	return products, err
}

func (r *ProductDB) FindBySellerID(sellerID id.UUID, page, size int) ([]*entity.Product, error) {
	var products []*entity.Product

	offset := (page - 1) * size

	err := r.db.
		Preload("Reviews").
		Where("seller_id = ?", sellerID).
		Limit(size).
		Offset(offset).
		Order("created_at desc").
		Find(&products).Error

	return products, err
}

func (r *ProductDB) Update(product *entity.Product) (*entity.Product, error) {
	err := r.db.Save(&product).Error
	return product, err
}

func (r *ProductDB) Delete(productID id.UUID) error {
	return r.db.Delete(&entity.Product{}, "id = ?", productID).Error
}
