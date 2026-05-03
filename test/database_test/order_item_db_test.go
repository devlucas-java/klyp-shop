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

var dbOrderItem *gorm.DB
var orderItemRepo *database.OrderItemDB
var logOrderItem *logger.Logger

func setupOrderItemDB() {
	var err error

	dbOrderItem, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dbOrderItem.AutoMigrate(&entity.Order{}, &entity.Product{}, &entity.OrderItem{})
	if err != nil {
		panic(err)
	}

	logOrderItem = logger.NewLogger(logger.TRACE)
	orderItemRepo = database.NewOrderItemDB(dbOrderItem, logOrderItem).(*database.OrderItemDB)
}

func createOrderItemOrder() *entity.Order {
	order := &entity.Order{
		ID:        id.NewUUID(),
		UserID:    id.NewUUID(),
		AddressID: id.NewUUID(),
		Status:    entity.OrderStatusPending,
		TotalBTC:  0.5,
	}

	dbOrderItem.Create(order)
	return order
}

func createOrderItemProduct() *entity.Product {
	product := &entity.Product{
		ID:          id.NewUUID(),
		Name:        "Test Product",
		Description: "Test Description",
		PriceBTC:    0.1,
		Stock:       10,
	}

	dbOrderItem.Create(product)
	return product
}

func TestCreateOrderItem(t *testing.T) {
	setupOrderItemDB()

	order := createOrderItemOrder()
	product := createOrderItemProduct()

	orderItem := &entity.OrderItem{
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  2,
		PriceBTC:  0.2,
	}

	res, err := orderItemRepo.Create(orderItem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, orderItem.Quantity, res.Quantity)
	assert.Equal(t, orderItem.OrderID, res.OrderID)
	assert.Equal(t, orderItem.ProductID, res.ProductID)
}

func TestGetOrderItemByOrder(t *testing.T) {
	setupOrderItemDB()

	order := createOrderItemOrder()
	product := createOrderItemProduct()

	dbOrderItem.Create(&entity.OrderItem{
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  1,
		PriceBTC:  0.1,
	})

	res, err := orderItemRepo.FindByOrder(order.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) == 0 {
		t.Fatal("expected order items")
	}
}

func TestUpdateOrderItem(t *testing.T) {
	setupOrderItemDB()

	order := createOrderItemOrder()
	product := createOrderItemProduct()

	orderItem := &entity.OrderItem{
		ID:        id.NewUUID(),
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  1,
		PriceBTC:  0.1,
	}

	dbOrderItem.Create(orderItem)

	orderItem.Quantity = 3

	res, err := orderItemRepo.Update(orderItem)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Quantity != 3 {
		t.Fatal("update failed")
	}
}

func TestDeleteOrderItem(t *testing.T) {
	setupOrderItemDB()

	order := createOrderItemOrder()
	product := createOrderItemProduct()

	orderItem := &entity.OrderItem{
		OrderID:   order.ID,
		ProductID: product.ID,
		Quantity:  1,
		PriceBTC:  0.1,
	}

	dbOrderItem.Create(orderItem)

	err := orderItemRepo.DeleteByID(orderItem.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	dbOrderItem.Model(&entity.OrderItem{}).Where("id = ?", orderItem.ID).Count(&count)

	if count != 0 {
		t.Fatal("delete failed")
	}
}
