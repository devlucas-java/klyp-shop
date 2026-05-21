package database

import (
	"context"

	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"gorm.io/gorm"
)

const commentDB = "comment_db.CommentDB"

type CommentDB struct {
	db *gorm.DB
}

func NewCommentDB(db *gorm.DB) repository.CommentRepository {
	return &CommentDB{db: db}
}

func (c *CommentDB) Create(comment *entity.Comment) (*entity.Comment, error) {
	if err := c.db.WithContext(context.Background()).Create(comment).Error; err != nil {
		return nil, apperrors.HandlePgError(commentDB+".create", err)
	}
	return comment, nil
}

func (c *CommentDB) Save(comment *entity.Comment) (*entity.Comment, error) {
	if err := c.db.WithContext(context.Background()).Where("id = ?", comment.ID).Save(comment).Error; err != nil {
		return nil, apperrors.HandlePgError(commentDB+".save", err)
	}
	return comment, nil
}

func (c *CommentDB) Updates(comment *entity.Comment) (*entity.Comment, error) {
	if err := c.db.WithContext(context.Background()).Model(comment).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return nil, apperrors.HandlePgError(commentDB+".updates", err)
	}
	return comment, nil
}

func (c *CommentDB) Update(comment *entity.Comment) (*entity.Comment, error) {
	return c.Save(comment)
}

func (c *CommentDB) FindByID(commentID id.UUID) (*entity.Comment, error) {
	var comment entity.Comment
	err := c.db.WithContext(context.Background()).First(&comment, "id = ?", commentID).Error
	if err != nil {
		return nil, apperrors.HandlePgError(commentDB+".find_by_id", err)
	}
	return &comment, nil
}

func (c *CommentDB) FindByUser(userID id.UUID) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	if err := c.db.WithContext(context.Background()).Where("user_id = ?", userID).Find(&comments).Error; err != nil {
		return nil, apperrors.HandlePgError(commentDB+".find_by_user", err)
	}
	return comments, nil
}

func (c *CommentDB) FindByProduct(productID id.UUID) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	if err := c.db.WithContext(context.Background()).Where("product_id = ?", productID).Find(&comments).Error; err != nil {
		return nil, apperrors.HandlePgError(commentDB+".find_by_product", err)
	}
	return comments, nil
}

func (c *CommentDB) DeleteByID(commentID id.UUID) error {
	if err := c.db.WithContext(context.Background()).Delete(&entity.Comment{}, "id = ?", commentID).Error; err != nil {
		return apperrors.HandlePgError(commentDB+".delete_by_id", err)
	}
	return nil
}
