package database

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"gorm.io/gorm"

	"github.com/devlucas-java/klyp-shop/pkg/id"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) repository.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *entity.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) FindByProduct(productID id.UUID) ([]entity.Review, error) {
	var reviews []entity.Review

	err := r.db.
		Where("product_id = ?", productID).
		Find(&reviews).Error

	return reviews, err
}

func (r *reviewRepository) Delete(reviewID id.UUID) error {
	return r.db.Delete(&entity.Review{}, "id = ?", reviewID).Error
}
