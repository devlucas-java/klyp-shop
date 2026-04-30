package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dseller"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbSellerSvc *gorm.DB
var sellerService *service.SellerService

func setupSellerService(t *testing.T) {
	var err error

	dbSellerSvc, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = dbSellerSvc.AutoMigrate(&entity.User{}, &entity.Seller{})
	if err != nil {
		t.Fatal(err)
	}

	log := logger.NewLogger(logger.TRACE)
	userRepo := database.NewUserDB(dbSellerSvc, log)
	sellerRepo := database.NewSellerDB(dbSellerSvc)
	sellerMapper := mapper.NewSellerMapper()
	sellerService = service.NewSellerService(log, userRepo, sellerRepo, sellerMapper)
}

func seedSellerUser(t *testing.T) *entity.User {
	user, err := entity.NewUser("Seller User", "seller@test.com", "selleruser", "password123")
	if err != nil {
		t.Fatal(err)
	}
	if err := dbSellerSvc.Create(user).Error; err != nil {
		t.Fatal(err)
	}
	return user
}

func TestSellerService_CreateSeller(t *testing.T) {
	setupSellerService(t)

	user := seedSellerUser(t)

	req := &dseller.CreateSeller{
		DisplayName: "My Store",
		Bio:         "Best store ever",
	}

	res, err := sellerService.CreateSeller(user, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, req.DisplayName, res.DisplayName)
	assert.Equal(t, req.Bio, res.Bio)

	// Verify user is now a seller
	var updated entity.User
	dbSellerSvc.First(&updated, "id = ?", user.ID)
	assert.True(t, updated.IsSeller)
}

func TestSellerService_CreateSeller_AlreadySeller(t *testing.T) {
	setupSellerService(t)

	user := seedSellerUser(t)

	// First creation should succeed
	req := &dseller.CreateSeller{
		DisplayName: "My Store",
		Bio:         "Best store ever",
	}
	_, err := sellerService.CreateSeller(user, req)
	if err != nil {
		t.Fatal(err)
	}

	// Second creation should fail
	req2 := &dseller.CreateSeller{
		DisplayName: "Another Store",
		Bio:         "Another bio",
	}
	_, err = sellerService.CreateSeller(user, req2)
	assert.Error(t, err)
}

func TestSellerService_GetSellerByID(t *testing.T) {
	setupSellerService(t)

	user := seedSellerUser(t)

	// Create a seller first
	req := &dseller.CreateSeller{
		DisplayName: "My Store",
		Bio:         "Best store ever",
	}
	sellerRes, err := sellerService.CreateSeller(user, req)
	if err != nil {
		t.Fatal(err)
	}
	uuid, err := id.Parse(sellerRes.ID)
	if err != nil {
		t.Fatalf("failed to parse product ID: %v", err)
	}

	res, err := sellerService.GetSellerByID(uuid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, sellerRes.ID, res.ID)
	assert.Equal(t, req.DisplayName, res.DisplayName)
}

func TestSellerService_GetSellerByID_NotFound(t *testing.T) {
	setupSellerService(t)

	ghostID := id.NewUUID()

	_, err := sellerService.GetSellerByID(ghostID)
	assert.Error(t, err)
}

func TestSellerService_DeleteSeller(t *testing.T) {
	setupSellerService(t)

	user := seedSellerUser(t)

	// Create a seller first
	req := &dseller.CreateSeller{
		DisplayName: "My Store",
		Bio:         "Best store ever",
	}
	_, err := sellerService.CreateSeller(user, req)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the seller
	err = sellerService.DeleteSeller(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify user is no longer a seller
	var updated entity.User
	dbSellerSvc.First(&updated, "id = ?", user.ID)
	assert.False(t, updated.IsSeller)
}

func TestSellerService_DeleteSeller_NotSeller(t *testing.T) {
	setupSellerService(t)

	user := seedSellerUser(t)

	// Try to delete a non-seller user
	err := sellerService.DeleteSeller(user)
	assert.Error(t, err)
}
