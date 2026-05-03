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

type CommentDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCommentDB(db *gorm.DB, log *logger.Logger) repository.CommentRepository {
	return &CommentDB{db: db, log: log}
}

func (c *CommentDB) Create(comment *entity.Comment) (*entity.Comment, error) {
	if err := c.db.Create(comment).Error; err != nil {
		c.log.Errorf("CommentDB.Create: %v", err)
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}
	return comment, nil
}

func (c *CommentDB) Save(comment *entity.Comment) (*entity.Comment, error) {
	if err := c.db.Where("id = ?", comment.ID).Save(comment).Error; err != nil {
		c.log.Errorf("CommentDB.Save %s: %v", comment.ID, err)
		return nil, fmt.Errorf("failed to save comment: %w", err)
	}
	return comment, nil
}

func (c *CommentDB) Updates(comment *entity.Comment) (*entity.Comment, error) {
	if err := c.db.Model(comment).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		c.log.Errorf("CommentDB.Updates %s: %v", comment.ID, err)
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}
	return comment, nil
}

func (c *CommentDB) Update(comment *entity.Comment) (*entity.Comment, error) {
	return c.Save(comment)
}

func (c *CommentDB) FindByID(commentID id.UUID) (*entity.Comment, error) {
	var comment entity.Comment
	err := c.db.First(&comment, "id = ?", commentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("Comment", err)
		}
		c.log.Errorf("CommentDB.FindByID %s: %v", commentID, err)
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}
	return &comment, nil
}

func (c *CommentDB) FindByUser(userID id.UUID) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	if err := c.db.Where("user_id = ?", userID).Find(&comments).Error; err != nil {
		c.log.Errorf("CommentDB.FindByUser %s: %v", userID, err)
		return nil, fmt.Errorf("failed to find comments: %w", err)
	}
	return comments, nil
}

func (c *CommentDB) FindByProduct(productID id.UUID) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	if err := c.db.Where("product_id = ?", productID).Find(&comments).Error; err != nil {
		c.log.Errorf("CommentDB.FindByProduct %s: %v", productID, err)
		return nil, fmt.Errorf("failed to find comments: %w", err)
	}
	return comments, nil
}

func (c *CommentDB) DeleteByID(commentID id.UUID) error {
	if err := c.db.Delete(&entity.Comment{}, "id = ?", commentID).Error; err != nil {
		c.log.Errorf("CommentDB.DeleteByID %s: %v", commentID, err)
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}