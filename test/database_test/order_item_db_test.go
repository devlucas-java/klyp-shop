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

var dbOrderItem *gorm.DB
var orderItemRepo *database.OrderItemDB
var logOrderItem *logger.Logger

func setupOrderItemDB(t *testing.T) {
	t.Helper()
	var err error

	dbOrderItem, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = dbOrderItem.AutoMigrate(&entity.Order{}, &entity.Product{}, &entity.OrderItem{})
	require.NoError(t, err)

	orderItemRepo = database.NewOrderItemDB(dbOrderItem).(*database.OrderItemDB)
}

func createOrderItemOrder(t *testing.T) *entity.Order {
	t.Helper()
	order := &entity.Order{
		ID:        id.NewUUID(),
		UserID:    id.NewUUID(),
		AddressID: id.NewUUID(),
		Status:    entity.OrderStatusPending,
		TotalBTC:  05,
	}
	require.NoError(t, dbOrderItem.Create(order).Error)
	return order
}

func createOrderItemProduct(t *testing.T) *entity.Product {
	t.Helper()
	product := &entity.Product{
		ID:          id.NewUUID(),
		Name:        "Test Product",
		Description: "Test Description",
		PriceBTC:    01,
		Stock:       10,
	}
	require.NoError(t, dbOrderItem.Create(product).Error)
	return product
}

func TestCreateOrderItem(t *testing.T) {
	setupOrderItemDB(t)

	order := createOrderItemOrder(t)
	product := createOrderItemProduct(t)

	orderItem := &entity.OrderItem{
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  2,
		PriceBTC:  02,
	}

	res, err := orderItemRepo.Create(orderItem)
	require.NoError(t, err)

	assert.Equal(t, orderItem.Quantity, res.Quantity)
	assert.Equal(t, orderItem.OrderID, res.OrderID)
	assert.Equal(t, orderItem.ProductID, res.ProductID)
}

func TestGetOrderItemByOrder(t *testing.T) {
	setupOrderItemDB(t)

	order := createOrderItemOrder(t)
	product := createOrderItemProduct(t)

	require.NoError(t, dbOrderItem.Create(&entity.OrderItem{
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  1,
		PriceBTC:  01,
	}).Error)

	res, err := orderItemRepo.FindByOrder(order.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestUpdateOrderItem(t *testing.T) {
	setupOrderItemDB(t)

	order := createOrderItemOrder(t)
	product := createOrderItemProduct(t)

	orderItem := &entity.OrderItem{
		ID:        id.NewUUID(),
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  1,
		PriceBTC:  01,
	}
	require.NoError(t, dbOrderItem.Create(orderItem).Error)

	orderItem.Quantity = 3

	res, err := orderItemRepo.Update(orderItem)
	require.NoError(t, err)
	assert.Equal(t, 3, res.Quantity)
}

func TestDeleteOrderItem(t *testing.T) {
	setupOrderItemDB(t)

	order := createOrderItemOrder(t)
	product := createOrderItemProduct(t)

	orderItem := &entity.OrderItem{
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  1,
		PriceBTC:  01,
	}
	require.NoError(t, dbOrderItem.Create(orderItem).Error)

	err := orderItemRepo.DeleteByID(orderItem.ID)
	require.NoError(t, err)

	var count int64
	dbOrderItem.Model(&entity.OrderItem{}).Where("id = ?", orderItem.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
