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

type CommentDB struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCommentDB(db *gorm.DB, log *logger.Logger) repository.CommentRepository {
	return &CommentDB{db: db, log: log}
}

func (c *CommentDB) Create(comment *entity.Comment) (*entity.Comment, error) {
	if err := c.db.WithContext(context.Background()).Create(comment).Error; err != nil {
		c.log.Errorf("CommentDB.Create: %v", err)
		return nil, domainErr.ErrDatabase("failed to create comment", err)
	}
	return comment, nil
}

func (c *CommentDB) Save(comment *entity.Comment) (*entity.Comment, error) {
	if err := c.db.WithContext(context.Background()).Where("id = ?", comment.ID).Save(comment).Error; err != nil {
		c.log.Errorf("CommentDB.Save %s: %v", comment.ID, err)
		return nil, domainErr.ErrDatabase("failed to save comment", err)
	}
	return comment, nil
}

func (c *CommentDB) Updates(comment *entity.Comment) (*entity.Comment, error) {
	if err := c.db.WithContext(context.Background()).Model(comment).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		c.log.Errorf("CommentDB.Updates %s: %v", comment.ID, err)
		return nil, domainErr.ErrDatabase("failed to update comment", err)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErr.ErrNotFound("Comment", err)
		}
		c.log.Errorf("CommentDB.FindByID %s: %v", commentID, err)
		return nil, domainErr.ErrDatabase("failed to find comment", err)
	}
	return &comment, nil
}

func (c *CommentDB) FindByUser(userID id.UUID) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	if err := c.db.WithContext(context.Background()).Where("user_id = ?", userID).Find(&comments).Error; err != nil {
		c.log.Errorf("CommentDB.FindByUser %s: %v", userID, err)
		return nil, domainErr.ErrDatabase("failed to find comments", err)
	}
	return comments, nil
}

func (c *CommentDB) FindByProduct(productID id.UUID) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	if err := c.db.WithContext(context.Background()).Where("product_id = ?", productID).Find(&comments).Error; err != nil {
		c.log.Errorf("CommentDB.FindByProduct %s: %v", productID, err)
		return nil, domainErr.ErrDatabase("failed to find comments", err)
	}
	return comments, nil
}

func (c *CommentDB) DeleteByID(commentID id.UUID) error {
	if err := c.db.WithContext(context.Background()).Delete(&entity.Comment{}, "id = ?", commentID).Error; err != nil {
		c.log.Errorf("CommentDB.DeleteByID %s: %v", commentID, err)
		return domainErr.ErrDatabase("failed to delete comment", err)
	}
	return nil
}
