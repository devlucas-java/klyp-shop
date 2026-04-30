package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/daddress"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newAddressService(addrRepo *mocks.AddressRepositoryMock, userRepo *mocks.UserRepositoryMock) *service.AddressService {
	return service.NewAddressService(addrRepo, logger.NewLogger(logger.TRACE), mapper.NewAddressMapper(), userRepo)
}

func newAddressUser() *entity.User {
	return &entity.User{
		ID:       id.NewUUID(),
		Name:     "Addr User",
		Email:    "addr@test.com",
		Username: "addruser",
	}
}

func newAddress(userID id.UUID) *entity.Address {
	return &entity.Address{
		ID:       id.NewUUID(),
		UserID:   userID,
		Street:   "Main St",
		City:     "City",
		State:    "State",
		Country:  "Country",
		Postcode: "12345",
		Number:   10,
	}
}

func TestAddressService_CreateAddress(t *testing.T) {
	addrRepo := new(mocks.AddressRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAddressService(addrRepo, userRepo)

	user := newAddressUser()
	req := &daddress.CreateAddressRequest{
		Street: "Main St", City: "City", State: "State",
		Country: "Country", PostCode: "12345", Number: 10,
	}

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addrRepo.On("FindByUser", user.ID).Return([]*entity.Address{}, nil)
	addrRepo.On("Create", mock.AnythingOfType("*entity.Address")).Return(newAddress(user.ID), nil)

	res, err := svc.CreateAddress(user, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Street, res.Street)
	userRepo.AssertExpectations(t)
	addrRepo.AssertExpectations(t)
}

func TestAddressService_CreateAddress_MaxLimit(t *testing.T) {
	addrRepo := new(mocks.AddressRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAddressService(addrRepo, userRepo)

	user := newAddressUser()
	existing := []*entity.Address{
		newAddress(user.ID), newAddress(user.ID), newAddress(user.ID),
	}

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addrRepo.On("FindByUser", user.ID).Return(existing, nil)

	req := &daddress.CreateAddressRequest{
		Street: "Extra", City: "City", State: "State",
		Country: "Country", PostCode: "00000", Number: 1,
	}

	_, err := svc.CreateAddress(user, req)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
	addrRepo.AssertExpectations(t)
}

func TestAddressService_GetAddresses(t *testing.T) {
	addrRepo := new(mocks.AddressRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAddressService(addrRepo, userRepo)

	user := newAddressUser()
	addrs := []*entity.Address{newAddress(user.ID), newAddress(user.ID)}

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addrRepo.On("FindByUser", user.ID).Return(addrs, nil)

	res, err := svc.GetAddresses(user)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	userRepo.AssertExpectations(t)
	addrRepo.AssertExpectations(t)
}

func TestAddressService_GetAddresses_Empty(t *testing.T) {
	addrRepo := new(mocks.AddressRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAddressService(addrRepo, userRepo)

	user := newAddressUser()

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addrRepo.On("FindByUser", user.ID).Return([]*entity.Address{}, nil)

	res, err := svc.GetAddresses(user)

	assert.NoError(t, err)
	assert.Len(t, res, 0)
	userRepo.AssertExpectations(t)
	addrRepo.AssertExpectations(t)
}

func TestAddressService_UpdateAddress(t *testing.T) {
	addrRepo := new(mocks.AddressRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAddressService(addrRepo, userRepo)

	user := newAddressUser()
	addr := newAddress(user.ID)
	updated := *addr
	updated.Street = "New Street"

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addrRepo.On("FindByID", addr.ID).Return(addr, nil)
	addrRepo.On("Update", mock.AnythingOfType("*entity.Address")).Return(&updated, nil)

	req := &daddress.UpdateAddressRequest{
		Street: "New Street", City: "City", State: "State",
		Country: "Country", PostCode: "12345", Number: 10,
	}

	res, err := svc.UpdateAddress(user, req, addr.ID)

	assert.NoError(t, err)
	assert.Equal(t, "New Street", res.Street)
	userRepo.AssertExpectations(t)
	addrRepo.AssertExpectations(t)
}

func TestAddressService_UpdateAddress_WrongOwner(t *testing.T) {
	addrRepo := new(mocks.AddressRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAddressService(addrRepo, userRepo)

	owner := newAddressUser()
	other := &entity.User{ID: id.NewUUID()}
	addr := newAddress(owner.ID)

	userRepo.On("FindByID", other.ID).Return(other, nil)
	addrRepo.On("FindByID", addr.ID).Return(addr, nil)

	req := &daddress.UpdateAddressRequest{Street: "Hacked"}

	_, err := svc.UpdateAddress(other, req, addr.ID)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
	addrRepo.AssertExpectations(t)
}

func TestAddressService_DeleteAddress(t *testing.T) {
	addrRepo := new(mocks.AddressRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAddressService(addrRepo, userRepo)

	user := newAddressUser()
	addr := newAddress(user.ID)

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addrRepo.On("FindByID", addr.ID).Return(addr, nil)
	addrRepo.On("DeleteByID", addr.ID).Return(nil)

	err := svc.DeleteAddress(user, addr.ID)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	addrRepo.AssertExpectations(t)
}

func TestAddressService_DeleteAddress_WrongOwner(t *testing.T) {
	addrRepo := new(mocks.AddressRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAddressService(addrRepo, userRepo)

	owner := newAddressUser()
	other := &entity.User{ID: id.NewUUID()}
	addr := newAddress(owner.ID)

	userRepo.On("FindByID", other.ID).Return(other, nil)
	addrRepo.On("FindByID", addr.ID).Return(addr, nil)

	err := svc.DeleteAddress(other, addr.ID)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
	addrRepo.AssertExpectations(t)
}

func TestAddressService_DeleteAddress_NotFound(t *testing.T) {
	addrRepo := new(mocks.AddressRepositoryMock)
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAddressService(addrRepo, userRepo)

	user := newAddressUser()
	ghostID := id.NewUUID()

	userRepo.On("FindByID", user.ID).Return(user, nil)
	addrRepo.On("FindByID", ghostID).Return(nil, domainErr.ErrNotFound("Address", nil))

	err := svc.DeleteAddress(user, ghostID)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
	addrRepo.AssertExpectations(t)
}
