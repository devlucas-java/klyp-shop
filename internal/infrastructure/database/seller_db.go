package database

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
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
		return nil, err
	}
	return seller, nil
}

func (r *SellerDB) Update(seller *entity.Seller) (*entity.Seller, error) {
	err := r.db.Save(seller).Error
	return seller, err
}

func (r *SellerDB) DeleteByID(id id.UUID) error {
	return r.db.Where("id = ? ", id).Delete(&entity.Seller{}).Error
}

func (r *SellerDB) FindByID(sellerID id.UUID) (*entity.Seller, error) {
	var seller entity.Seller
	err := r.db.
		First(&seller, "id = ?", sellerID).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

func (r *SellerDB) Find(page, size int, order, search string) ([]*entity.Seller, error) {
	var sellers []*entity.Seller

	offset := (page - 1) * size

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	query := r.db.
		Limit(size).
		Offset(offset).
		Order("created_at " + order)

	if search != "" {
		query = query.Where("display_name LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&sellers).Error; err != nil {
		return nil, err
	}
	return sellers, nil
}
