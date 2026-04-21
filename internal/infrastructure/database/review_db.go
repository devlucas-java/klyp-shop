package database

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

type ReviewDB struct {
	db *gorm.DB
}

func NewReviewDB(db *gorm.DB) repository.ReviewRepository {
	return &ReviewDB{db: db}
}

func (r *ReviewDB) Create(review *entity.Review) (*entity.Review, error) {
	err := r.db.Create(review).Error
	return review, err
}

func (r *ReviewDB) Update(review *entity.Review) (*entity.Review, error) {
	err := r.db.Save(review).Error
	return review, err
}

func (r *ReviewDB) FindByProductID(productID id.UUID) ([]*entity.Review, error) {
	var reviews []*entity.Review

	err := r.db.
		Where("product_id = ?", productID).
		Find(&reviews).Error

	return reviews, err
}

func (r *ReviewDB) Delete(reviewID id.UUID) error {
	return r.db.Delete(&entity.Review{}, "id = ?", reviewID).Error
}
