package service_test

// import (
// 	"testing"

// 	"github.com/devlucas-java/klyp-shop/internal/application/service"
// 	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/daddress"
// 	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
// 	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
// )

// func createTestUser(t *testing.T) *entity.User {
// 	user := &entity.User{
// 		Name:     "test",
// 		Email:    "test@test.com",
// 		Username: "testuser",
// 	}

// 	err := testDB.Create(user).Error
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	return user
// }

// func TestCreateAddress(t *testing.T) {

// 	user := createTestUser(t)

// 	mapper := &mapper.AddressMapper{}
// 	service := service.NewAddressService(addressRepo, log, mapper, userRepo)

// 	req := &daddress.CreateAddressRequest{
// 		Street:   "Main St",
// 		City:     "City",
// 		State:    "State",
// 		Country:  "Country",
// 		Postcode: "12345",
// 		Number:   10,
// 	}

// 	res, err := service.CreateAddress(user, req)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if res == nil {
// 		t.Fatal("expected response, got nil")
// 	}
// }

// func TestCreateAddress_MaxLimit(t *testing.T) {

// 	user := createTestUser(t)

// 	for i := 0; i < 3; i++ {
// 		testDB.Create(&entity.Address{
// 			UserID:   user.ID,
// 			Street:   "Street",
// 			City:     "City",
// 			State:    "State",
// 			Country:  "Country",
// 			Postcode: "12345",
// 			Number:   int32(i),
// 		})
// 	}

// 	mapper := &mapper.AddressMapper{}
// 	service := service.NewAddressService(addressRepo, log, mapper, userRepo)

// 	req := &daddress.CreateAddressRequest{
// 		Street: "New",
// 		City:   "City",
// 	}

// 	_, err := service.CreateAddress(user, req)

// 	if err == nil {
// 		t.Fatal("expected error but got nil")
// 	}
// }

// func TestGetAddresses(t *testing.T) {

// 	user := createTestUser(t)

// 	testDB.Create(&entity.Address{
// 		UserID:   user.ID,
// 		Street:   "Street 1",
// 		City:     "City",
// 		State:    "State",
// 		Country:  "Country",
// 		Postcode: "12345",
// 	})

// 	mapper := &mapper.AddressMapper{}
// 	service := service.NewAddressService(addressRepo, log, mapper, userRepo)

// 	res, err := service.GetAddresses(user)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if len(res) == 0 {
// 		t.Fatal("expected addresses")
// 	}
// }

// func TestUpdateAddress(t *testing.T) {

// 	user := createTestUser(t)

// 	addr := &entity.Address{
// 		UserID:   user.ID,
// 		Street:   "Old Street",
// 		City:     "City",
// 		State:    "State",
// 		Country:  "Country",
// 		Postcode: "12345",
// 	}

// 	testDB.Create(addr)

// 	mapper := &mapper.AddressMapper{}
// 	service := service.NewAddressService(addressRepo, log, mapper, userRepo)

// 	req := &daddress.UpdateAddressRequest{
// 		Street: "Updated Street",
// 	}

// 	res, err := service.UpdateAddress(user, req, addr.ID)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if res == nil {
// 		t.Fatal("expected response")
// 	}
// }

// func TestDeleteAddress(t *testing.T) {

// 	user := createTestUser(t)

// 	addr := &entity.Address{
// 		UserID:   user.ID,
// 		Street:   "Street",
// 		City:     "City",
// 		State:    "State",
// 		Country:  "Country",
// 		Postcode: "12345",
// 	}

// 	testDB.Create(addr)

// 	mapper := &mapper.AddressMapper{}
// 	service := service.NewAddressService(addressRepo, log, mapper, userRepo)

// 	err := service.DeleteAddress(user, addr.ID)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// }
