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

var dbProduct *gorm.DB
var productRepo *database.ProductDB

func setupProductDB(t *testing.T) {
	var err error

	dbProduct, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = dbProduct.AutoMigrate(&entity.User{}, &entity.Seller{}, &entity.Product{}, &entity.Review{})
	if err != nil {
		t.Fatal(err)
	}

	productRepo = database.NewProductDB(dbProduct).(*database.ProductDB)
}

func createProductSeller(t *testing.T) *entity.Seller {
	user := &entity.User{
		ID:       id.NewUUID(),
		Name:     "prod-user",
		Email:    "produser@test.com",
		Username: "produser",
		Password: "hash",
	}
	dbProduct.Create(user)

	seller := entity.NewSeller(user.ID, "Prod Shop", "Bio")
	dbProduct.Create(seller)
	return seller
}

func TestCreateProduct(t *testing.T) {
	setupProductDB(t)

	seller := createProductSeller(t)

	product := entity.NewProduct("Laptop", "A great laptop", 0.05, 10, seller.ID, []string{"electronics"})

	res, err := productRepo.Create(product)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, product.Name, res.Name)
	assert.Equal(t, product.PriceBTC, res.PriceBTC)
	assert.Equal(t, product.Stock, res.Stock)
	assert.Equal(t, product.SellerID, res.SellerID)
}

func TestFindProductByID(t *testing.T) {
	setupProductDB(t)

	seller := createProductSeller(t)

	product := entity.NewProduct("Phone", "A phone", 0.02, 5, seller.ID, []string{"electronics"})
	dbProduct.Create(product)

	found, err := productRepo.FindByID(product.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, product.ID, found.ID)
	assert.Equal(t, product.Name, found.Name)
}

func TestFindProductByID_NotFound(t *testing.T) {
	setupProductDB(t)

	_, err := productRepo.FindByID(id.NewUUID())
	assert.Error(t, err)
}

func TestUpdateProduct(t *testing.T) {
	setupProductDB(t)

	seller := createProductSeller(t)

	product := entity.NewProduct("Old Product", "Desc", 0.01, 3, seller.ID, []string{"misc"})
	dbProduct.Create(product)

	product.Name = "Updated Product"
	product.Stock = 99

	res, err := productRepo.Update(product)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, "Updated Product", res.Name)
	assert.Equal(t, 99, res.Stock)
}

func TestDeleteProduct(t *testing.T) {
	setupProductDB(t)

	seller := createProductSeller(t)

	product := entity.NewProduct("To Delete", "Desc", 0.01, 1, seller.ID, []string{})
	dbProduct.Create(product)

	err := productRepo.DeleteByID(product.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	dbProduct.Model(&entity.Product{}).Where("id = ?", product.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestFindProductsBySellerID(t *testing.T) {
	setupProductDB(t)

	seller := createProductSeller(t)

	for i := 0; i < 3; i++ {
		dbProduct.Create(entity.NewProduct(
			"Product",
			"Desc",
			0.01,
			1,
			seller.ID,
			[]string{},
		))
	}

	res, err := productRepo.FindBySellerID(seller.ID, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 3, len(res))
}
