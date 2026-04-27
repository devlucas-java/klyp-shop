package database_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbReview *gorm.DB
var reviewRepo *database.ReviewDB

func setupReviewDB(t *testing.T) {
	var err error

	dbReview, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = dbReview.AutoMigrate(&entity.User{}, &entity.Seller{}, &entity.Product{}, &entity.Review{})
	if err != nil {
		t.Fatal(err)
	}

	reviewRepo = database.NewReviewDB(dbReview).(*database.ReviewDB)
}

func createReviewProduct(t *testing.T) (*entity.User, *entity.Product) {
	user := &entity.User{
		ID:       id.NewUUID(),
		Name:     "reviewer",
		Email:    "reviewer@test.com",
		Username: "reviewer",
		Password: "hash",
	}
	dbReview.Create(user)

	seller := entity.NewSeller(user.ID, "Review Shop", "Bio")
	dbReview.Create(seller)

	product := entity.NewProduct("Reviewed Product", "Desc", 0.01, 5, seller.ID, []string{})
	dbReview.Create(product)

	return user, product
}

func TestCreateReview(t *testing.T) {
	setupReviewDB(t)

	user, product := createReviewProduct(t)

	review := entity.NewReview(user.ID, product.ID, 5, "Excellent!")

	res, err := reviewRepo.Create(review)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, review.Rating, res.Rating)
	assert.Equal(t, review.Comment, res.Comment)
	assert.Equal(t, review.UserID, res.UserID)
	assert.Equal(t, review.ProductID, res.ProductID)
}

func TestUpdateReview(t *testing.T) {
	setupReviewDB(t)

	user, product := createReviewProduct(t)

	review := entity.NewReview(user.ID, product.ID, 3, "Average")
	dbReview.Create(review)

	review.Rating = 4
	review.Comment = "Pretty good"

	res, err := reviewRepo.Update(review)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 4, res.Rating)
	assert.Equal(t, "Pretty good", res.Comment)
}

func TestFindReviewsByProductID(t *testing.T) {
	setupReviewDB(t)

	user, product := createReviewProduct(t)

	dbReview.Create(entity.NewReview(user.ID, product.ID, 5, "Great"))
	dbReview.Create(entity.NewReview(user.ID, product.ID, 4, "Good"))

	res, err := reviewRepo.FindByProductID(product.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 2, len(res))
}

func TestFindReviewsByProductID_Empty(t *testing.T) {
	setupReviewDB(t)

	res, err := reviewRepo.FindByProductID(id.NewUUID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 0, len(res))
}

func TestDeleteReview(t *testing.T) {
	setupReviewDB(t)

	user, product := createReviewProduct(t)

	review := entity.NewReview(user.ID, product.ID, 2, "Not great")
	dbReview.Create(review)

	err := reviewRepo.DeleteByID(review.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	dbReview.Model(&entity.Review{}).Where("id = ?", review.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
