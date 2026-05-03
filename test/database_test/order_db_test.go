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

var dbOrder *gorm.DB
var orderRepo *database.OrderDB
var logOrder *logger.Logger

func setupOrderDB() {
	var err error

	dbOrder, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dbOrder.AutoMigrate(&entity.User{}, &entity.Address{}, &entity.Order{}, &entity.OrderItem{}, &entity.BitcoinPayment{})
	if err != nil {
		panic(err)
	}

	logOrder = logger.NewLogger(logger.TRACE)
	orderRepo = database.NewOrderDB(dbOrder, logOrder).(*database.OrderDB)
}

func createOrderUser() *entity.User {
	user := &entity.User{
		ID:    id.NewUUID(),
		Name:  "test",
		Email: "test@test.com",
	}

	dbOrder.Create(user)
	return user
}

func createOrderAddress(userID id.UUID) *entity.Address {
	address := &entity.Address{
		ID:       id.NewUUID(),
		UserID:   userID,
		Street:   "Main St",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}

	dbOrder.Create(address)
	return address
}

func TestCreateOrder(t *testing.T) {
	setupOrderDB()

	user := createOrderUser()
	address := createOrderAddress(user.ID)

	order := &entity.Order{
		UserID:    user.ID,
		AddressID: address.ID,
		Status:    entity.OrderStatusPending,
		TotalBTC:  0.5,
	}

	res, err := orderRepo.Create(order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, order.TotalBTC, res.TotalBTC)
	assert.Equal(t, order.UserID, res.UserID)
	assert.Equal(t, order.AddressID, res.AddressID)
}

func TestGetOrderByUser(t *testing.T) {
	setupOrderDB()

	user := createOrderUser()
	address := createOrderAddress(user.ID)

	dbOrder.Create(&entity.Order{
		UserID:    user.ID,
		AddressID: address.ID,
		Status:    entity.OrderStatusPending,
		TotalBTC:  0.5,
	})

	res, err := orderRepo.FindByUser(user.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) == 0 {
		t.Fatal("expected orders")
	}
}

func TestUpdateOrder(t *testing.T) {
	setupOrderDB()

	user := createOrderUser()
	address := createOrderAddress(user.ID)

	order := &entity.Order{
		ID:        id.NewUUID(),
		UserID:    user.ID,
		AddressID: address.ID,
		Status:    entity.OrderStatusPending,
		TotalBTC:  0.5,
	}

	dbOrder.Create(order)

	order.Status = entity.OrderStatusPaid

	res, err := orderRepo.Update(order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Status != entity.OrderStatusPaid {
		t.Fatal("update failed")
	}
}

func TestDeleteOrder(t *testing.T) {
	setupOrderDB()

	user := createOrderUser()
	address := createOrderAddress(user.ID)

	order := &entity.Order{
		UserID:    user.ID,
		AddressID: address.ID,
		Status:    entity.OrderStatusPending,
		TotalBTC:  0.5,
	}

	dbOrder.Create(order)

	err := orderRepo.DeleteByID(order.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	dbOrder.Model(&entity.Order{}).Where("id = ?", order.ID).Count(&count)

	if count != 0 {
		t.Fatal("delete failed")
	}
}
