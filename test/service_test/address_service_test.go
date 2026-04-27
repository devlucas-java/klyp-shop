package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/daddress"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbAddrSvc *gorm.DB
var addressService *service.AddressService

func setupAddressService(t *testing.T) {
	var err error

	dbAddrSvc, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = dbAddrSvc.AutoMigrate(&entity.User{}, &entity.Address{})
	if err != nil {
		t.Fatal(err)
	}

	log := logger.NewLogger(logger.TRACE)
	addrRepo := database.NewAddressDB(dbAddrSvc, log)
	userRepo := database.NewUserDB(dbAddrSvc, log)
	addrMapper := mapper.NewAddressMapper()
	addressService = service.NewAddressService(addrRepo, log, addrMapper, userRepo)
}

func seedAddressUser(t *testing.T) *entity.User {
	user, err := entity.NewUser("Addr User", "addruser@test.com", "addruser", "password123")
	if err != nil {
		t.Fatal(err)
	}
	if err := dbAddrSvc.Create(user).Error; err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAddressService_CreateAddress(t *testing.T) {
	setupAddressService(t)

	user := seedAddressUser(t)

	req := &daddress.CreateAddressRequest{
		Street:   "Main St",
		City:     "Springfield",
		State:    "IL",
		Country:  "US",
		PostCode: "62701",
		Number:   42,
	}

	res, err := addressService.CreateAddress(user, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, req.Street, res.Street)
	assert.Equal(t, req.City, res.City)
	assert.Equal(t, req.PostCode, res.PostCode)
	assert.Equal(t, req.Number, res.Number)
}

func TestAddressService_CreateAddress_MaxLimit(t *testing.T) {
	setupAddressService(t)

	user := seedAddressUser(t)

	for i := 0; i < 3; i++ {
		dbAddrSvc.Create(&entity.Address{
			ID:       id.NewUUID(),
			UserID:   user.ID,
			Street:   "Street",
			City:     "City",
			State:    "State",
			Country:  "Country",
			Postcode: "12345",
			Number:   int32(i + 1),
		})
	}

	req := &daddress.CreateAddressRequest{
		Street:   "Extra St",
		City:     "City",
		State:    "State",
		Country:  "US",
		PostCode: "00000",
		Number:   1,
	}

	_, err := addressService.CreateAddress(user, req)
	assert.Error(t, err)
}

func TestAddressService_GetAddresses(t *testing.T) {
	setupAddressService(t)

	user := seedAddressUser(t)

	dbAddrSvc.Create(&entity.Address{
		ID:       id.NewUUID(),
		UserID:   user.ID,
		Street:   "Street 1",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	})

	res, err := addressService.GetAddresses(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 1, len(res))
}

func TestAddressService_GetAddresses_Empty(t *testing.T) {
	setupAddressService(t)

	user := seedAddressUser(t)

	res, err := addressService.GetAddresses(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, 0, len(res))
}

func TestAddressService_UpdateAddress(t *testing.T) {
	setupAddressService(t)

	user := seedAddressUser(t)

	addr := &entity.Address{
		ID:       id.NewUUID(),
		UserID:   user.ID,
		Street:   "Old Street",
		City:     "Old City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}
	dbAddrSvc.Create(addr)

	req := &daddress.UpdateAddressRequest{
		Street:   "New Street",
		City:     "New City",
		State:    "State",
		Country:  "Country",
		PostCode: "99999",
		Number:   10,
	}

	res, err := addressService.UpdateAddress(user, req, addr.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, "New Street", res.Street)
	assert.Equal(t, "New City", res.City)
}

func TestAddressService_UpdateAddress_WrongOwner(t *testing.T) {
	setupAddressService(t)

	owner := seedAddressUser(t)

	// Create a second user
	other, err := entity.NewUser("Other", "other@test.com", "otheruser", "pass")
	if err != nil {
		t.Fatal(err)
	}
	dbAddrSvc.Create(other)

	addr := &entity.Address{
		ID:       id.NewUUID(),
		UserID:   owner.ID,
		Street:   "Owner Street",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}
	dbAddrSvc.Create(addr)

	req := &daddress.UpdateAddressRequest{Street: "Hacked Street"}

	_, err = addressService.UpdateAddress(other, req, addr.ID)
	assert.Error(t, err)
}

func TestAddressService_DeleteAddress(t *testing.T) {
	setupAddressService(t)

	user := seedAddressUser(t)

	addr := &entity.Address{
		ID:       id.NewUUID(),
		UserID:   user.ID,
		Street:   "Delete Me",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}
	dbAddrSvc.Create(addr)

	err := addressService.DeleteAddress(user, addr.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	dbAddrSvc.Model(&entity.Address{}).Where("id = ?", addr.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestAddressService_DeleteAddress_WrongOwner(t *testing.T) {
	setupAddressService(t)

	owner := seedAddressUser(t)

	other, err := entity.NewUser("Other", "other3@test.com", "otheruser3", "pass")
	if err != nil {
		t.Fatal(err)
	}
	dbAddrSvc.Create(other)

	addr := &entity.Address{
		ID:       id.NewUUID(),
		UserID:   owner.ID,
		Street:   "Owner Street",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
	}
	dbAddrSvc.Create(addr)

	err = addressService.DeleteAddress(other, addr.ID)
	assert.Error(t, err)
}
