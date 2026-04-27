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

var dbAddress *gorm.DB
var addressRepo *database.AddressDB
var logAddress *logger.Logger

func setupAddressDB() {
	var err error

	dbAddress, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dbAddress.AutoMigrate(&entity.User{}, &entity.Address{})
	if err != nil {
		panic(err)
	}

	logAddress = logger.NewLogger(logger.TRACE)
	addressRepo = database.NewAddressDB(dbAddress, logAddress).(*database.AddressDB)
}

func createUser() *entity.User {
	user := &entity.User{
		ID:    id.NewUUID(),
		Name:  "test",
		Email: "test@test.com",
	}

	dbAddress.Create(user)
	return user
}

func TestCreateAddress(t *testing.T) {
	setupAddressDB()

	user := createUser()

	addr := &entity.Address{
		UserID:   user.ID,
		Street:   "Main",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}

	res, err := addressRepo.Create(addr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, addr.Street, res.Street)
	assert.Equal(t, addr.City, res.City)
	assert.Equal(t, addr.State, res.State)
	assert.Equal(t, addr.Country, res.Country)
	assert.Equal(t, addr.Postcode, res.Postcode)
	assert.Equal(t, addr.UserID, res.UserID)
}

func TestGetAddress(t *testing.T) {
	setupAddressDB()

	user := createUser()

	dbAddress.Create(&entity.Address{
		UserID:   user.ID,
		Street:   "Street 1",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	})

	res, err := addressRepo.FindByUser(user.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) == 0 {
		t.Fatal("expected addresses")
	}
}

func TestUpdateAddress(t *testing.T) {
	setupAddressDB()

	user := createUser()

	addr := &entity.Address{
		ID:       id.NewUUID(),
		UserID:   user.ID,
		Street:   "Main",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}

	dbAddress.Create(addr)

	addr.Street = "New Street"

	res, err := addressRepo.Update(addr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Street != "New Street" {
		t.Fatal("update failed")
	}
}

func TestDeleteAddress(t *testing.T) {
	setupAddressDB()

	user := createUser()

	addr := &entity.Address{
		UserID:   user.ID,
		Street:   "Street",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}

	dbAddress.Create(addr)

	err := addressRepo.DeleteByID(addr.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	dbAddress.Model(&entity.Address{}).Where("id = ?", addr.ID).Count(&count)

	if count != 0 {
		t.Fatal("expected address to be deleted")
	}
}
