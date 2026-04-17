package database

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SellerDB struct {
	DB *gorm.DB
}

func NewSellerDB(db *gorm.DB) repository.SellerRepository {
	return &SellerDB{DB: db}
}

func (s *SellerDB) Save(seller *entity.Seller) (*entity.Seller, error) {
	err := s.DB.Create(seller).Error
	if err != nil {
		return nil, err
	}
	return seller, nil
}

func (s *SellerDB) Updates(seller *entity.Seller) (*entity.Seller, error) {
	err := s.DB.
		Model(&entity.Seller{}).
		Where("id = ?", seller.ID).
		Updates(&seller).Error

	if err != nil {
		return nil, err
	}
	return seller, nil
}
func (s *SellerDB) FindByID(id uuid.UUID) (*entity.Seller, error) {
	var seller entity.Seller

	err := s.DB.First(&seller, id).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

func (s *SellerDB) Find(page, size int, order, search string) ([]*entity.Seller, error) {

	var sellers []*entity.Seller
	offset := (page - 1) * size

	if search != "" {
		search = "%" + search + "%"
	}

	if order != "desc" && order != "asc" {
		order = "desc"
	}

	query := s.DB.
		Limit(size).
		Offset(offset).
		Order("created_at " + order)

	if search != "" {
		query = query.Where("display_name LIKE ?", search)
	}

	err := query.Find(&sellers).Error

	if err != nil {
		return nil, err
	}

	return sellers, nil
}
