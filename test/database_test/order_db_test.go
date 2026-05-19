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

var dbOrder *gorm.DB
var orderRepo *database.OrderDB
var logOrder *logger.Logger

func setupOrderDB(t *testing.T) {
	t.Helper()
	var err error

	dbOrder, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = dbOrder.AutoMigrate(
		&entity.User{},
		&entity.Address{},
		&entity.Order{},
		&entity.OrderItem{},
		&entity.BitcoinPayment{},
	)
	require.NoError(t, err)

	logOrder = logger.NewLogger(logger.TRACE)
	orderRepo = database.NewOrderDB(dbOrder, logOrder).(*database.OrderDB)
}

func createOrderUser(t *testing.T) *entity.User {
	t.Helper()
	user := &entity.User{
		ID:    id.NewUUID(),
		Name:  "test",
		Email: "test@test.com",
	}
	require.NoError(t, dbOrder.Create(user).Error)
	return user
}

func createOrderAddress(t *testing.T, userID id.UUID) *entity.Address {
	t.Helper()
	address := &entity.Address{
		ID:       id.NewUUID(),
		UserID:   userID,
		Street:   "Main St",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}
	require.NoError(t, dbOrder.Create(address).Error)
	return address
}

func TestCreateOrder(t *testing.T) {
	setupOrderDB(t)

	user := createOrderUser(t)
	address := createOrderAddress(t, user.ID)

	order := &entity.Order{
		UserID:    user.ID,
		AddressID: address.ID,
		Status:    entity.OrderStatusPending,
		TotalBTC:  0.5,
	}

	res, err := orderRepo.Create(order)
	require.NoError(t, err)

	assert.Equal(t, order.TotalBTC, res.TotalBTC)
	assert.Equal(t, order.UserID, res.UserID)
	assert.Equal(t, order.AddressID, res.AddressID)
}

func TestGetOrderByUser(t *testing.T) {
	setupOrderDB(t)

	user := createOrderUser(t)
	address := createOrderAddress(t, user.ID)

	require.NoError(t, dbOrder.Create(&entity.Order{
		UserID:    user.ID,
		AddressID: address.ID,
		Status:    entity.OrderStatusPending,
		TotalBTC:  0.5,
	}).Error)

	res, err := orderRepo.FindByUser(user.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestUpdateOrder(t *testing.T) {
	setupOrderDB(t)

	user := createOrderUser(t)
	address := createOrderAddress(t, user.ID)

	order := &entity.Order{
		ID:        id.NewUUID(),
		UserID:    user.ID,
		AddressID: address.ID,
		Status:    entity.OrderStatusPending,
		TotalBTC:  0.5,
	}
	require.NoError(t, dbOrder.Create(order).Error)

	order.Status = entity.OrderStatusPaid

	res, err := orderRepo.Update(order)
	require.NoError(t, err)
	assert.Equal(t, entity.OrderStatusPaid, res.Status)
}

func TestDeleteOrder(t *testing.T) {
	setupOrderDB(t)

	user := createOrderUser(t)
	address := createOrderAddress(t, user.ID)

	order := &entity.Order{
		UserID:    user.ID,
		AddressID: address.ID,
		Status:    entity.OrderStatusPending,
		TotalBTC:  0.5,
	}
	require.NoError(t, dbOrder.Create(order).Error)

	err := orderRepo.DeleteByID(order.ID)
	require.NoError(t, err)

	var count int64
	dbOrder.Model(&entity.Order{}).Where("id = ?", order.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
