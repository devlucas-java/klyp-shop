package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
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
	if err := r.db.WithContext(context.Background()).Create(review).Error; err != nil {
		return nil, apperrors.HandlePgError("review", err)
	}
	return review, nil
}

func (r *ReviewDB) Save(review *entity.Review) (*entity.Review, error) {
	if err := r.db.WithContext(context.Background()).Where("id = ?", review.ID).Save(review).Error; err != nil {
		return nil, apperrors.HandlePgError("review", err)
	}
	return review, nil
}

func (r *ReviewDB) Updates(review *entity.Review) (*entity.Review, error) {
	if err := r.db.WithContext(context.Background()).Model(review).Where("id = ?", review.ID).Updates(review).Error; err != nil {
		return nil, apperrors.HandlePgError("review", err)
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
		return nil, apperrors.HandlePgError("review", err)
	}
	return &review, nil
}

func (r *ReviewDB) FindByUser(userID id.UUID) ([]*entity.Review, error) {
	var reviews []*entity.Review
	if err := r.db.WithContext(context.Background()).Where("user_id = ?", userID).Find(&reviews).Error; err != nil {
		return nil, apperrors.HandlePgError("review", err)
	}
	return reviews, nil
}

func (r *ReviewDB) FindByProductID(productID id.UUID) ([]*entity.Review, error) {
	var reviews []*entity.Review
	if err := r.db.WithContext(context.Background()).Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		return nil, apperrors.HandlePgError("review", err)
	}
	return reviews, nil
}

func (r *ReviewDB) DeleteByID(reviewID id.UUID) error {
	if err := r.db.WithContext(context.Background()).Delete(&entity.Review{}, "id = ?", reviewID).Error; err != nil {
		return apperrors.HandlePgError("review", err)
	}
	return nil
}
