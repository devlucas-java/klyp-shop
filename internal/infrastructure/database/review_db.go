package database

import (
	"errors"
	"fmt"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"gorm.io/gorm"
)

type ReviewDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewReviewDB(db *gorm.DB, log *logger.Logger) repository.ReviewRepository {
	return &ReviewDB{db: db, log: log}
}

func (r *ReviewDB) Create(review *entity.Review) (*entity.Review, error) {
	if err := r.db.Create(review).Error; err != nil {
		r.log.Errorf("ReviewDB.Create: %v", err)
		return nil, fmt.Errorf("failed to create review: %w", err)
	}
	return review, nil
}

func (r *ReviewDB) Save(review *entity.Review) (*entity.Review, error) {
	if err := r.db.Where("id = ?", review.ID).Save(review).Error; err != nil {
		r.log.Errorf("ReviewDB.Save %s: %v", review.ID, err)
		return nil, fmt.Errorf("failed to save review: %w", err)
	}
	return review, nil
}

func (r *ReviewDB) Updates(review *entity.Review) (*entity.Review, error) {
	if err := r.db.Model(review).Where("id = ?", review.ID).Updates(review).Error; err != nil {
		r.log.Errorf("ReviewDB.Updates %s: %v", review.ID, err)
		return nil, fmt.Errorf("failed to update review: %w", err)
	}
	return review, nil
}

func (r *ReviewDB) Update(review *entity.Review) (*entity.Review, error) {
	return r.Save(review)
}

func (r *ReviewDB) FindByID(reviewID id.UUID) (*entity.Review, error) {
	var review entity.Review
	err := r.db.First(&review, "id = ?", reviewID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("Review", err)
		}
		r.log.Errorf("ReviewDB.FindByID %s: %v", reviewID, err)
		return nil, fmt.Errorf("failed to find review: %w", err)
	}
	return &review, nil
}

func (r *ReviewDB) FindByUser(userID id.UUID) ([]*entity.Review, error) {
	var reviews []*entity.Review
	if err := r.db.Where("user_id = ?", userID).Find(&reviews).Error; err != nil {
		r.log.Errorf("ReviewDB.FindByUser %s: %v", userID, err)
		return nil, fmt.Errorf("failed to find reviews: %w", err)
	}
	return reviews, nil
}

func (r *ReviewDB) FindByProductID(productID id.UUID) ([]*entity.Review, error) {
	var reviews []*entity.Review
	if err := r.db.Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		r.log.Errorf("ReviewDB.FindByProductID %s: %v", productID, err)
		return nil, fmt.Errorf("failed to find reviews: %w", err)
	}
	return reviews, nil
}

func (r *ReviewDB) DeleteByID(reviewID id.UUID) error {
	if err := r.db.Delete(&entity.Review{}, "id = ?", reviewID).Error; err != nil {
		r.log.Errorf("ReviewDB.DeleteByID %s: %v", reviewID, err)
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}
