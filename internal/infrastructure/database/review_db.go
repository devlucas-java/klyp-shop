package database

import (
	"fmt"

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
	if err := r.db.Create(review).Error; err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}
	return review, nil
}

func (r *ReviewDB) Save(review *entity.Review) (*entity.Review, error) {
	if err := r.db.Where("id = ?", review.ID).Save(review).Error; err != nil {
		return nil, fmt.Errorf("failed to save review: %w", err)
	}
	return review, nil
}

func (r *ReviewDB) Updates(review *entity.Review) (*entity.Review, error) {
	if err := r.db.Model(review).Where("id = ?", review.ID).Updates(review).Error; err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}
	return review, nil
}

func (r *ReviewDB) Update(review *entity.Review) (*entity.Review, error) {
	return r.Save(review)
}

func (r *ReviewDB) FindByProductID(productID id.UUID) ([]*entity.Review, error) {
	var reviews []*entity.Review
	if err := r.db.Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		return nil, fmt.Errorf("failed to find reviews: %w", err)
	}
	return reviews, nil
}

func (r *ReviewDB) DeleteByID(reviewID id.UUID) error {
	if err := r.db.Delete(&entity.Review{}, "id = ?", reviewID).Error; err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}
