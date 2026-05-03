package database_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbComment *gorm.DB
var commentRepo *database.CommentDB
var logComment *logger.Logger

func setupCommentDB() {
	var err error

	dbComment, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dbComment.AutoMigrate(&entity.User{}, &entity.Product{}, &entity.Comment{})
	if err != nil {
		panic(err)
	}

	logComment = logger.NewLogger(logger.TRACE)
	commentRepo = database.NewCommentDB(dbComment, logComment).(*database.CommentDB)
}

func createCommentUser() *entity.User {
	user := &entity.User{
		ID:    id.NewUUID(),
		Name:  "test",
		Email: "test@test.com",
	}

	dbComment.Create(user)
	return user
}

func createCommentProduct() *entity.Product {
	product := &entity.Product{
		ID:          id.NewUUID(),
		Name:        "Test Product",
		Description: "Test Description",
		PriceBTC:    0.1,
		Stock:       10,
	}

	dbComment.Create(product)
	return product
}

func TestCreateComment(t *testing.T) {
	setupCommentDB()

	user := createCommentUser()
	product := createCommentProduct()

	comment := &entity.Comment{
		UserID:    user.ID,
		ProductID: product.ID,
		Content:   "This is a test comment",
	}

	res, err := commentRepo.Create(comment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, comment.Content, res.Content)
	assert.Equal(t, comment.UserID, res.UserID)
	assert.Equal(t, comment.ProductID, res.ProductID)
}

func TestGetCommentByUser(t *testing.T) {
	setupCommentDB()

	user := createCommentUser()
	product := createCommentProduct()

	dbComment.Create(&entity.Comment{
		UserID:    user.ID,
		ProductID: product.ID,
		Content:   "Comment 1",
	})

	res, err := commentRepo.FindByUser(user.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) == 0 {
		t.Fatal("expected comments")
	}
}

func TestGetCommentByProduct(t *testing.T) {
	setupCommentDB()

	user := createCommentUser()
	product := createCommentProduct()

	dbComment.Create(&entity.Comment{
		UserID:    user.ID,
		ProductID: product.ID,
		Content:   "Comment 1",
	})

	res, err := commentRepo.FindByProduct(product.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) == 0 {
		t.Fatal("expected comments")
	}
}

func TestUpdateComment(t *testing.T) {
	setupCommentDB()

	user := createCommentUser()
	product := createCommentProduct()

	comment := &entity.Comment{
		ID:        id.NewUUID(),
		UserID:    user.ID,
		ProductID: product.ID,
		Content:   "Original comment",
	}

	dbComment.Create(comment)

	comment.Content = "Updated comment"

	res, err := commentRepo.Update(comment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Content != "Updated comment" {
		t.Fatal("update failed")
	}
}

func TestDeleteComment(t *testing.T) {
	setupCommentDB()

	user := createCommentUser()
	product := createCommentProduct()

	comment := &entity.Comment{
		UserID:    user.ID,
		ProductID: product.ID,
		Content:   "Comment to delete",
	}

	dbComment.Create(comment)

	err := commentRepo.DeleteByID(comment.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	dbComment.Model(&entity.Comment{}).Where("id = ?", comment.ID).Count(&count)

	if count != 0 {
		t.Fatal("delete failed")
	}
}