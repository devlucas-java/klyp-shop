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

var dbAddress *gorm.DB
var addressRepo *database.AddressDB
var logAddress *logger.Logger

func setupAddressDB(t *testing.T) {
	t.Helper()
	var err error

	dbAddress, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = dbAddress.AutoMigrate(&entity.User{}, &entity.Address{})
	require.NoError(t, err)

	logAddress = logger.NewLogger(logger.TRACE)
	addressRepo = database.NewAddressDB(dbAddress, logAddress).(*database.AddressDB)
}

func createAddressUser(t *testing.T) *entity.User {
	t.Helper()
	user := &entity.User{
		ID:    id.NewUUID(),
		Name:  "test",
		Email: "test@test.com",
	}
	require.NoError(t, dbAddress.Create(user).Error)
	return user
}

func TestCreateAddress(t *testing.T) {
	setupAddressDB(t)

	user := createAddressUser(t)

	addr := &entity.Address{
		UserID:   user.ID,
		Street:   "Main",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}

	res, err := addressRepo.Create(addr)
	require.NoError(t, err)

	assert.Equal(t, addr.Street, res.Street)
	assert.Equal(t, addr.City, res.City)
	assert.Equal(t, addr.State, res.State)
	assert.Equal(t, addr.Country, res.Country)
	assert.Equal(t, addr.Postcode, res.Postcode)
	assert.Equal(t, addr.UserID, res.UserID)
}

func TestGetAddress(t *testing.T) {
	setupAddressDB(t)

	user := createAddressUser(t)

	require.NoError(t, dbAddress.Create(&entity.Address{
		UserID:   user.ID,
		Street:   "Street 1",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}).Error)

	res, err := addressRepo.FindByUser(user.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestUpdateAddress(t *testing.T) {
	setupAddressDB(t)

	user := createAddressUser(t)

	addr := &entity.Address{
		ID:       id.NewUUID(),
		UserID:   user.ID,
		Street:   "Main",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}
	require.NoError(t, dbAddress.Create(addr).Error)

	addr.Street = "New Street"

	res, err := addressRepo.Update(addr)
	require.NoError(t, err)
	assert.Equal(t, "New Street", res.Street)
}

func TestDeleteAddress(t *testing.T) {
	setupAddressDB(t)

	user := createAddressUser(t)

	addr := &entity.Address{
		UserID:   user.ID,
		Street:   "Street",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}
	require.NoError(t, dbAddress.Create(addr).Error)

	err := addressRepo.DeleteByID(addr.ID)
	require.NoError(t, err)

	var count int64
	dbAddress.Model(&entity.Address{}).Where("id = ?", addr.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
