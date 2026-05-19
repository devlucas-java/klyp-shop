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

var dbReview *gorm.DB
var reviewRepo *database.ReviewDB
var logReview *logger.Logger

func setupReviewDB(t *testing.T) {
	t.Helper()
	var err error

	dbReview, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = dbReview.AutoMigrate(&entity.User{}, &entity.Seller{}, &entity.Product{}, &entity.Review{})
	require.NoError(t, err)

	logReview = logger.NewLogger(logger.TRACE)
	reviewRepo = database.NewReviewDB(dbReview, logReview).(*database.ReviewDB)
}

func createReviewProduct(t *testing.T) (*entity.User, *entity.Product) {
	t.Helper()
	user := &entity.User{
		ID:       id.NewUUID(),
		Name:     "reviewer",
		Email:    "reviewer@test.com",
		Username: "reviewer",
		Password: "hash",
	}
	require.NoError(t, dbReview.Create(user).Error)

	seller := entity.NewSeller(user.ID, "Review Shop", "Bio")
	require.NoError(t, dbReview.Create(seller).Error)

	product, err := entity.NewProduct("Reviewed Product", "Desc", 0.01, 5, []string{})
	require.NoError(t, err)
	product.SellerID = seller.ID
	require.NoError(t, dbReview.Create(product).Error)

	return user, product
}

func TestCreateReview(t *testing.T) {
	setupReviewDB(t)

	user, product := createReviewProduct(t)

	review := entity.NewReview(user.ID, product.ID, 5, "Excellent!")

	res, err := reviewRepo.Create(review)
	require.NoError(t, err)

	assert.Equal(t, review.Rating, res.Rating)
	assert.Equal(t, review.Comment, res.Comment)
	assert.Equal(t, review.UserID, res.UserID)
	assert.Equal(t, review.ProductID, res.ProductID)
}

func TestUpdateReview(t *testing.T) {
	setupReviewDB(t)

	user, product := createReviewProduct(t)

	review := entity.NewReview(user.ID, product.ID, 3, "Average")
	require.NoError(t, dbReview.Create(review).Error)

	review.Rating = 4
	review.Comment = "Pretty good"

	res, err := reviewRepo.Update(review)
	require.NoError(t, err)

	assert.Equal(t, 4, res.Rating)
	assert.Equal(t, "Pretty good", res.Comment)
}

func TestFindReviewsByProductID(t *testing.T) {
	setupReviewDB(t)

	user, product := createReviewProduct(t)

	require.NoError(t, dbReview.Create(entity.NewReview(user.ID, product.ID, 5, "Great")).Error)
	require.NoError(t, dbReview.Create(entity.NewReview(user.ID, product.ID, 4, "Good")).Error)

	res, err := reviewRepo.FindByProductID(product.ID)
	require.NoError(t, err)

	assert.Equal(t, 2, len(res))
}

func TestFindReviewsByProductID_Empty(t *testing.T) {
	setupReviewDB(t)

	res, err := reviewRepo.FindByProductID(id.NewUUID())
	require.NoError(t, err)

	assert.Equal(t, 0, len(res))
}

func TestDeleteReview(t *testing.T) {
	setupReviewDB(t)

	user, product := createReviewProduct(t)

	review := entity.NewReview(user.ID, product.ID, 2, "Not great")
	require.NoError(t, dbReview.Create(review).Error)

	err := reviewRepo.DeleteByID(review.ID)
	require.NoError(t, err)

	var count int64
	dbReview.Model(&entity.Review{}).Where("id = ?", review.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
