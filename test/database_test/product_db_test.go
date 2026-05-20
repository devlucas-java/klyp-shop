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

var dbProduct *gorm.DB
var productRepo *database.ProductDB

func setupProductDB(t *testing.T) {
	t.Helper()
	var err error

	dbProduct, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = dbProduct.AutoMigrate(&entity.User{}, &entity.Seller{}, &entity.Product{}, &entity.Review{})
	require.NoError(t, err)

	log := logger.NewLogger(logger.TRACE)

	productRepo = database.NewProductDB(dbProduct, log).(*database.ProductDB)
}

func createProductSeller(t *testing.T) *entity.Seller {
	t.Helper()
	user := &entity.User{
		ID:       id.NewUUID(),
		Name:     "prod-user",
		Email:    "produser@test.com",
		Username: "produser",
		Password: "hash",
	}
	require.NoError(t, dbProduct.Create(user).Error)

	seller := entity.NewSeller(user.ID, "Prod Shop", "Bio")
	require.NoError(t, dbProduct.Create(seller).Error)
	return seller
}

func newTestProduct(t *testing.T, name, desc string, price float64, stock int, cats []string) *entity.Product {
	t.Helper()
	p, err := entity.NewProduct(name, desc, price, stock, cats)
	require.NoError(t, err)
	return p
}

func TestCreateProduct(t *testing.T) {
	setupProductDB(t)
	seller := createProductSeller(t)

	product := newTestProduct(t, "Laptop", "A great laptop", 0.05, 10, []string{"electronics"})
	product.SellerID = seller.ID

	res, err := productRepo.Create(product)
	require.NoError(t, err)

	assert.Equal(t, product.Name, res.Name)
	assert.Equal(t, product.PriceBTC, res.PriceBTC)
	assert.Equal(t, product.Stock, res.Stock)
	assert.Equal(t, product.SellerID, res.SellerID)
}

func TestFindProductByID(t *testing.T) {
	setupProductDB(t)
	seller := createProductSeller(t)

	product := newTestProduct(t, "Phone", "A phone", 0.02, 5, []string{"electronics"})
	product.SellerID = seller.ID
	require.NoError(t, dbProduct.Create(product).Error)

	found, err := productRepo.FindByID(product.ID)
	require.NoError(t, err)

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

	product := newTestProduct(t, "Old Product", "Desc", 0.01, 3, []string{"misc"})
	product.SellerID = seller.ID
	require.NoError(t, dbProduct.Create(product).Error)

	product.Name = "Updated Product"
	product.Stock = 99

	res, err := productRepo.Updates(product)
	require.NoError(t, err)

	assert.Equal(t, "Updated Product", res.Name)
	assert.Equal(t, 99, res.Stock)
}

func TestDeleteProduct(t *testing.T) {
	setupProductDB(t)
	seller := createProductSeller(t)

	product := newTestProduct(t, "To Delete", "Desc", 0.01, 1, []string{})
	product.SellerID = seller.ID
	require.NoError(t, dbProduct.Create(product).Error)

	err := productRepo.DeleteByID(product.ID)
	require.NoError(t, err)

	var count int64
	dbProduct.Model(&entity.Product{}).Where("id = ?", product.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestFindProductsBySellerID(t *testing.T) {
	setupProductDB(t)
	seller := createProductSeller(t)

	for i := 0; i < 3; i++ {
		p := newTestProduct(t, "Product", "Desc", 0.01, 1, []string{})
		p.SellerID = seller.ID
		require.NoError(t, dbProduct.Create(p).Error)
	}

	res, err := productRepo.FindBySellerID(seller.ID, 1, 10)
	require.NoError(t, err)

	assert.Equal(t, 3, len(res))
}
