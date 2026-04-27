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

var dbSeller *gorm.DB
var sellerRepo *database.SellerDB

func setupSellerDB(t *testing.T) {
	var err error

	dbSeller, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = dbSeller.AutoMigrate(&entity.User{}, &entity.Seller{})
	if err != nil {
		t.Fatal(err)
	}

	sellerRepo = database.NewSellerDB(dbSeller).(*database.SellerDB)
}

func createSellerUser(t *testing.T) *entity.User {
	user := &entity.User{
		ID:       id.NewUUID(),
		Name:     "seller-user",
		Email:    "seller@test.com",
		Username: "selleruser",
		Password: "hash",
	}
	if err := dbSeller.Create(user).Error; err != nil {
		t.Fatal(err)
	}
	return user
}

func TestCreateSeller(t *testing.T) {
	setupSellerDB(t)

	user := createSellerUser(t)

	seller := entity.NewSeller(user.ID, "My Shop", "Best shop around")

	res, err := sellerRepo.Create(seller)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, seller.DisplayName, res.DisplayName)
	assert.Equal(t, seller.Bio, res.Bio)
	assert.Equal(t, seller.UserID, res.UserID)
}

func TestFindSellerByID(t *testing.T) {
	setupSellerDB(t)

	user := createSellerUser(t)

	seller := entity.NewSeller(user.ID, "Find Shop", "Bio")
	dbSeller.Create(seller)

	found, err := sellerRepo.FindByID(seller.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, seller.ID, found.ID)
	assert.Equal(t, seller.DisplayName, found.DisplayName)
}

func TestFindSellerByID_NotFound(t *testing.T) {
	setupSellerDB(t)

	_, err := sellerRepo.FindByID(id.NewUUID())
	assert.Error(t, err)
}

func TestUpdateSeller(t *testing.T) {
	setupSellerDB(t)

	user := createSellerUser(t)

	seller := entity.NewSeller(user.ID, "Old Name", "Old Bio")
	dbSeller.Create(seller)

	seller.DisplayName = "New Name"
	seller.Bio = "New Bio"

	res, err := sellerRepo.Update(seller)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, "New Name", res.DisplayName)
	assert.Equal(t, "New Bio", res.Bio)
}

func TestDeleteSeller(t *testing.T) {
	setupSellerDB(t)

	user := createSellerUser(t)

	seller := entity.NewSeller(user.ID, "To Delete", "Bio")
	dbSeller.Create(seller)

	err := sellerRepo.DeleteByID(seller.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	dbSeller.Model(&entity.Seller{}).Where("id = ?", seller.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestFindSellers_Pagination(t *testing.T) {
	setupSellerDB(t)

	for i := 0; i < 3; i++ {
		user := &entity.User{
			ID:       id.NewUUID(),
			Name:     "user",
			Email:    "user" + string(rune('a'+i)) + "@test.com",
			Username: "user" + string(rune('a'+i)),
			Password: "hash",
		}
		dbSeller.Create(user)
		dbSeller.Create(entity.NewSeller(user.ID, "Shop "+string(rune('A'+i)), "Bio"))
	}

	res, err := sellerRepo.Find(1, 2, "desc", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 2, len(res))
}

func TestFindSellers_Search(t *testing.T) {
	setupSellerDB(t)

	user := createSellerUser(t)
	dbSeller.Create(entity.NewSeller(user.ID, "UniqueShopName", "Bio"))

	res, err := sellerRepo.Find(1, 10, "desc", "UniqueShop")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 1, len(res))
	assert.Equal(t, "UniqueShopName", res[0].DisplayName)
}
