package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dseller"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newSellerService(userRepo *mocks.UserRepositoryMock, sellerRepo *mocks.SellerRepositoryMock) *service.SellerService {
	return service.NewSellerService(logger.NewLogger(logger.TRACE), userRepo, sellerRepo, mapper.NewSellerMapper())
}

func newSellerUser() *entity.User {
	return &entity.User{
		ID:       id.NewUUID(),
		Name:     "Seller User",
		Email:    "seller@test.com",
		Username: "selleruser",
		IsSeller: false,
	}
}

func newSeller(userID id.UUID) *entity.Seller {
	return &entity.Seller{
		ID:          id.NewUUID(),
		UserID:      userID,
		DisplayName: "My Store",
		Bio:         "Best store",
	}
}

func TestSellerService_CreateSeller(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newSellerService(userRepo, sellerRepo)

	user := newSellerUser()
	seller := newSeller(user.ID)

	userRepo.On("FindByID", user.ID).Return(user, nil)
	sellerRepo.On("Create", mock.AnythingOfType("*entity.Seller")).Return(seller, nil)
	userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(user, nil)

	res, err := svc.CreateSeller(user, &dseller.CreateSeller{DisplayName: "My Store", Bio: "Best store"})

	assert.NoError(t, err)
	assert.Equal(t, seller.DisplayName, res.DisplayName)
	userRepo.AssertExpectations(t)
	sellerRepo.AssertExpectations(t)
}

func TestSellerService_CreateSeller_AlreadySeller(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newSellerService(userRepo, sellerRepo)

	user := newSellerUser()
	user.IsSeller = true
	userRepo.On("FindByID", user.ID).Return(user, nil)

	_, err := svc.CreateSeller(user, &dseller.CreateSeller{DisplayName: "Store", Bio: "Bio"})

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestSellerService_GetSellerByID(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newSellerService(userRepo, sellerRepo)

	user := newSellerUser()
	seller := newSeller(user.ID)
	sellerRepo.On("FindByID", seller.ID).Return(seller, nil)

	res, err := svc.GetSellerByID(seller.ID)

	assert.NoError(t, err)
	assert.Equal(t, seller.ID.String(), res.ID)
	sellerRepo.AssertExpectations(t)
}

func TestSellerService_GetSellerByID_NotFound(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newSellerService(userRepo, sellerRepo)

	ghostID := id.NewUUID()
	sellerRepo.On("FindByID", ghostID).Return(nil, domainErr.ErrNotFound("Seller", nil))

	_, err := svc.GetSellerByID(ghostID)

	assert.Error(t, err)
	sellerRepo.AssertExpectations(t)
}

func TestSellerService_UpdateSeller(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newSellerService(userRepo, sellerRepo)

	user := newSellerUser()
	user.IsSeller = true
	seller := newSeller(user.ID)
	user.Seller = seller

	updated := *seller
	updated.DisplayName = "New Name"

	userRepo.On("FindByIDWithSeller", user.ID).Return(user, nil)
	sellerRepo.On("Updates", mock.AnythingOfType("*entity.Seller")).Return(&updated, nil)

	res, err := svc.UpdateSeller(user, &dseller.UpdateSeller{DisplayName: "New Name"})

	assert.NoError(t, err)
	assert.Equal(t, "New Name", res.DisplayName)
	userRepo.AssertExpectations(t)
	sellerRepo.AssertExpectations(t)
}

func TestSellerService_UpdateSeller_NotSeller(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newSellerService(userRepo, sellerRepo)

	user := newSellerUser()
	userRepo.On("FindByIDWithSeller", user.ID).Return(user, nil)

	_, err := svc.UpdateSeller(user, &dseller.UpdateSeller{DisplayName: "Name"})

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestSellerService_DeleteSeller(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newSellerService(userRepo, sellerRepo)

	user := newSellerUser()
	user.IsSeller = true
	seller := newSeller(user.ID)
	user.Seller = seller

	userRepo.On("FindByIDWithSeller", user.ID).Return(user, nil)
	sellerRepo.On("DeleteByID", seller.ID).Return(nil)
	userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(user, nil)

	err := svc.DeleteSeller(user)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	sellerRepo.AssertExpectations(t)
}

func TestSellerService_DeleteSeller_NotSeller(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	sellerRepo := new(mocks.SellerRepositoryMock)
	svc := newSellerService(userRepo, sellerRepo)

	user := newSellerUser()
	userRepo.On("FindByIDWithSeller", user.ID).Return(user, nil)

	err := svc.DeleteSeller(user)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}
