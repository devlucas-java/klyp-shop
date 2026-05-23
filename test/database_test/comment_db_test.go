package database_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbComment *gorm.DB
var commentRepo *database.CommentDB
var logComment *logger.Logger

func setupCommentDB(t *testing.T) {
	t.Helper()
	var err error

	dbComment, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = dbComment.AutoMigrate(&entity.User{}, &entity.Product{}, &entity.Comment{})
	require.NoError(t, err)

	commentRepo = database.NewCommentDB(dbComment).(*database.CommentDB)
}

func createCommentUser(t *testing.T) *entity.User {
	t.Helper()
	user := &entity.User{
		ID:    id.NewUUID(),
		Name:  "test",
		Email: "test@test.com",
	}
	require.NoError(t, dbComment.Create(user).Error)
	return user
}

func createCommentProduct(t *testing.T) *entity.Product {
	t.Helper()
	product := &entity.Product{
		ID:          id.NewUUID(),
		Name:        "Test Product",
		Description: "Test Description",
		PriceBTC:    01,
		Stock:       10,
	}
	require.NoError(t, dbComment.Create(product).Error)
	return product
}

func TestCreateComment(t *testing.T) {
	setupCommentDB(t)

	user := createCommentUser(t)
	product := createCommentProduct(t)

	comment := &entity.Comment{
		UserID:    user.ID,
		ProductID: product.ID,
		Content:   "This is a test comment",
	}

	res, err := commentRepo.Create(comment)
	require.NoError(t, err)

	assert.Equal(t, comment.Content, res.Content)
	assert.Equal(t, comment.UserID, res.UserID)
	assert.Equal(t, comment.ProductID, res.ProductID)
}

func TestGetCommentByUser(t *testing.T) {
	setupCommentDB(t)

	user := createCommentUser(t)
	product := createCommentProduct(t)

	require.NoError(t, dbComment.Create(&entity.Comment{
		UserID:    user.ID,
		ProductID: product.ID,
		Content:   "Comment 1",
	}).Error)

	res, err := commentRepo.FindByUser(user.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestGetCommentByProduct(t *testing.T) {
	setupCommentDB(t)

	user := createCommentUser(t)
	product := createCommentProduct(t)

	require.NoError(t, dbComment.Create(&entity.Comment{
		UserID:    user.ID,
		ProductID: product.ID,
		Content:   "Comment 1",
	}).Error)

	res, err := commentRepo.FindByProduct(product.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestUpdateComment(t *testing.T) {
	setupCommentDB(t)

	user := createCommentUser(t)
	product := createCommentProduct(t)

	comment := &entity.Comment{
		ID:        id.NewUUID(),
		UserID:    user.ID,
		ProductID: product.ID,
		Content:   "Original comment",
	}
	require.NoError(t, dbComment.Create(comment).Error)

	comment.Content = "Updated comment"

	res, err := commentRepo.Update(comment)
	require.NoError(t, err)
	assert.Equal(t, "Updated comment", res.Content)
}

func TestDeleteComment(t *testing.T) {
	setupCommentDB(t)

	user := createCommentUser(t)
	product := createCommentProduct(t)

	comment := &entity.Comment{
		UserID:    user.ID,
		ProductID: product.ID,
		Content:   "Comment to delete",
	}
	require.NoError(t, dbComment.Create(comment).Error)

	err := commentRepo.DeleteByID(comment.ID)
	require.NoError(t, err)

	var count int64
	dbComment.Model(&entity.Comment{}).Where("id = ?", comment.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
