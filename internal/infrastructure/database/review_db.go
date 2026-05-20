package database

import (
	"context"
	"errors"
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
	if err := r.db.WithContext(context.Background()).Create(review).Error; err != nil {
		r.log.Errorf("ReviewDB.Create: %v", err)
		return nil, handlePgError(err, "failed to create review")
	}
	return review, nil
}

func (r *ReviewDB) Save(review *entity.Review) (*entity.Review, error) {
	if err := r.db.WithContext(context.Background()).Where("id = ?", review.ID).Save(review).Error; err != nil {
		r.log.Errorf("ReviewDB.Save %s: %v", review.ID, err)
		return nil, handlePgError(err, "failed to save review")
	}
	return review, nil
}

func (r *ReviewDB) Updates(review *entity.Review) (*entity.Review, error) {
	if err := r.db.WithContext(context.Background()).Model(review).Where("id = ?", review.ID).Updates(review).Error; err != nil {
		r.log.Errorf("ReviewDB.Updates %s: %v", review.ID, err)
		return nil, handlePgError(err, "failed to update review")
	}
	return review, nil
}

func (r *ReviewDB) Update(review *entity.Review) (*entity.Review, error) {
	return r.Save(review)
}

func (r *ReviewDB) FindByID(reviewID id.UUID) (*entity.Review, error) {
	var review entity.Review
	err := r.db.WithContext(context.Background()).First(&review, "id = ?", reviewID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("Review", err)
		}
		r.log.Errorf("ReviewDB.FindByID %s: %v", reviewID, err)
		return nil, handlePgError(err, "failed to find review")
	}
	return &review, nil
}

func (r *ReviewDB) FindByUser(userID id.UUID) ([]*entity.Review, error) {
	var reviews []*entity.Review
	if err := r.db.WithContext(context.Background()).Where("user_id = ?", userID).Find(&reviews).Error; err != nil {
		r.log.Errorf("ReviewDB.FindByUser %s: %v", userID, err)
		return nil, handlePgError(err, "failed to find reviews")
	}
	return reviews, nil
}

func (r *ReviewDB) FindByProductID(productID id.UUID) ([]*entity.Review, error) {
	var reviews []*entity.Review
	if err := r.db.WithContext(context.Background()).Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		r.log.Errorf("ReviewDB.FindByProductID %s: %v", productID, err)
		return nil, handlePgError(err, "failed to find reviews")
	}
	return reviews, nil
}

func (r *ReviewDB) DeleteByID(reviewID id.UUID) error {
	if err := r.db.WithContext(context.Background()).Delete(&entity.Review{}, "id = ?", reviewID).Error; err != nil {
		r.log.Errorf("ReviewDB.DeleteByID %s: %v", reviewID, err)
		return handlePgError(err, "failed to delete review")
	}
	return nil
}
