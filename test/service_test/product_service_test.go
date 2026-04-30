package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dproduct"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dseller"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbProductSvc *gorm.DB
var productService *service.ProductService
var sellerServiceForTest *service.SellerService

func setupProductService(t *testing.T) {
	var err error

	dbProductSvc, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = dbProductSvc.AutoMigrate(&entity.User{}, &entity.Seller{}, &entity.Product{})
	if err != nil {
		t.Fatal(err)
	}

	log := logger.NewLogger(logger.TRACE)
	userRepo := database.NewUserDB(dbProductSvc, log)
	sellerRepo := database.NewSellerDB(dbProductSvc)
	productRepo := database.NewProductDB(dbProductSvc)
	productMapper := mapper.NewProductMapper()
	sellerMapper := mapper.NewSellerMapper()

	sellerServiceForTest = service.NewSellerService(log, userRepo, sellerRepo, sellerMapper)
	productService = service.NewProductService(log, productRepo, userRepo, sellerRepo, productMapper)

	// Store userRepo for later use in seedProductSeller
	_productUserRepo = userRepo
}

var _productUserRepo repository.UserRepository

func seedProductSeller(t *testing.T) *entity.User {
	user, err := entity.NewUser("Product Seller", "productseller@test.com", "productseller", "password123")
	if err != nil {
		t.Fatal(err)
	}
	if err := dbProductSvc.Create(user).Error; err != nil {
		t.Fatal(err)
	}

	// Convert user to seller
	sellerReq := &dseller.CreateSeller{
		DisplayName: "Product Store",
		Bio:         "We sell products",
	}
	_, err = sellerServiceForTest.CreateSeller(user, sellerReq)
	if err != nil {
		t.Fatal(err)
	}

	// Reload user with seller
	user, err = _productUserRepo.FindByIDWithSeller(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func TestProductService_CreateProduct(t *testing.T) {
	setupProductService(t)

	user := seedProductSeller(t)

	req := &dproduct.CreateProduct{
		Name:        "Test Product",
		Description: "A test product",
		PriceBTC:    0.01,
		Stock:       100,
		Categories:  []string{"electronics", "gadgets"},
	}

	res, err := productService.CreateProduct(user, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, req.Description, res.Description)
	assert.Equal(t, req.PriceBTC, res.PriceBTC)
	assert.Equal(t, req.Stock, res.Stock)
	assert.Equal(t, user.Seller.ID, res.SellerID)
}

func TestProductService_CreateProduct_NotSeller(t *testing.T) {
	setupProductService(t)

	user, err := entity.NewUser("Not Seller", "notseller@test.com", "notseller", "password123")
	if err != nil {
		t.Fatal(err)
	}
	dbProductSvc.Create(user)

	req := &dproduct.CreateProduct{
		Name:        "Test Product",
		Description: "A test product",
		PriceBTC:    0.01,
		Stock:       100,
	}

	_, err = productService.CreateProduct(user, req)
	assert.Error(t, err)
}

func TestProductService_GetProductByID(t *testing.T) {
	setupProductService(t)

	user := seedProductSeller(t)

	// Create a product first
	req := &dproduct.CreateProduct{
		Name:        "Test Product",
		Description: "A test product",
		PriceBTC:    0.01,
		Stock:       100,
	}
	productRes, err := productService.CreateProduct(user, req)
	if err != nil {
		t.Fatal(err)
	}

	uuid, err := id.Parse(productRes.ID)
	if err != nil {
		t.Fatalf("failed to parse product ID: %v", err)
	}

	res, err := productService.GetProductByID(uuid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, productRes.ID, res.ID)
	assert.Equal(t, req.Name, res.Name)
}

func TestProductService_GetProductByID_NotFound(t *testing.T) {
	setupProductService(t)

	ghostID := id.NewUUID()

	_, err := productService.GetProductByID(ghostID)
	assert.Error(t, err)
}

func TestProductService_UpdateProduct(t *testing.T) {
	setupProductService(t)

	user := seedProductSeller(t)

	// Create a product first
	req := &dproduct.CreateProduct{
		Name:        "Original Name",
		Description: "Original description",
		PriceBTC:    0.01,
		Stock:       100,
	}
	productRes, err := productService.CreateProduct(user, req)
	if err != nil {
		t.Fatal(err)
	}

	// Update the product
	updateReq := &dproduct.UpdateProduct{
		Name:        "Updated Name",
		Description: "Updated description",
		PriceBTC:    0.02,
		Stock:       50,
	}
	uuid, err := id.Parse(productRes.ID)
	if err != nil {
		t.Fatalf("failed to parse product ID: %v", err)
	}

	res, err := productService.UpdateProduct(user, updateReq, uuid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, "Updated Name", res.Name)
	assert.Equal(t, "Updated description", res.Description)
	assert.Equal(t, 0.02, res.PriceBTC)
	assert.Equal(t, 50, res.Stock)
}

func TestProductService_UpdateProduct_NotOwner(t *testing.T) {
	setupProductService(t)

	user1 := seedProductSeller(t)

	// Create a product as user1
	req := &dproduct.CreateProduct{
		Name:        "Test Product",
		Description: "A test product",
		PriceBTC:    0.01,
		Stock:       100,
	}
	productRes, err := productService.CreateProduct(user1, req)
	if err != nil {
		t.Fatal(err)
	}

	// Create another seller
	user2, err := entity.NewUser("Seller 2", "seller2@test.com", "seller2", "password123")
	if err != nil {
		t.Fatal(err)
	}
	dbProductSvc.Create(user2)
	sellerReq := &dseller.CreateSeller{
		DisplayName: "Store 2",
		Bio:         "Store 2",
	}
	_, err = sellerServiceForTest.CreateSeller(user2, sellerReq)
	if err != nil {
		t.Fatal(err)
	}
	user2, _ = database.NewUserDB(dbProductSvc, logger.NewLogger(logger.TRACE)).FindByIDWithSeller(user2.ID)

	// Try to update as user2
	updateReq := &dproduct.UpdateProduct{
		Name: "Hacked Name",
	}

	uuid, err := id.Parse(productRes.ID)
	if err != nil {
		t.Fatalf("failed to parse product ID: %v", err)
	}

	_, err = productService.UpdateProduct(user2, updateReq, uuid)
	assert.Error(t, err)
}

func TestProductService_DeleteProduct(t *testing.T) {
	setupProductService(t)

	user := seedProductSeller(t)

	// Create a product first
	req := &dproduct.CreateProduct{
		Name:        "Test Product",
		Description: "A test product",
		PriceBTC:    0.01,
		Stock:       100,
	}
	productRes, err := productService.CreateProduct(user, req)
	if err != nil {
		t.Fatal(err)
	}

	uuid, err := id.Parse(productRes.ID)
	if err != nil {
		t.Fatalf("failed to parse product ID: %v", err)
	}

	// Delete the product
	err = productService.DeleteProduct(user, uuid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = productService.GetProductByID(uuid)
	assert.Error(t, err)
}

func TestProductService_DeleteProduct_NotOwner(t *testing.T) {
	setupProductService(t)

	user1 := seedProductSeller(t)

	// Create a product as user1
	req := &dproduct.CreateProduct{
		Name:        "Test Product",
		Description: "A test product",
		PriceBTC:    0.01,
		Stock:       100,
	}
	productRes, err := productService.CreateProduct(user1, req)
	if err != nil {
		t.Fatal(err)
	}

	uuid, err := id.Parse(productRes.ID)
	if err != nil {
		t.Fatalf("failed to parse product ID: %v", err)
	}

	user2, err := entity.NewUser("Seller 2", "seller2@test.com", "seller2", "password123")
	if err != nil {
		t.Fatal(err)
	}
	dbProductSvc.Create(user2)
	sellerReq := &dseller.CreateSeller{
		DisplayName: "Store 2",
		Bio:         "Store 2",
	}

	_, err = sellerServiceForTest.CreateSeller(user2, sellerReq)
	if err != nil {
		t.Fatal(err)
	}
	user2, _ = database.NewUserDB(dbProductSvc, logger.NewLogger(logger.TRACE)).FindByIDWithSeller(user2.ID)

	err = productService.DeleteProduct(user2, uuid)
	assert.Error(t, err)
}
