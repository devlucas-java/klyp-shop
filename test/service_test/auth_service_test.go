package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/auth"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	domainErr "github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/password_encoder"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newAuthService(userRepo *mocks.UserRepositoryMock) *service.AuthService {
	return service.NewAuthService(userRepo, jwt.NewJWTService("test-secret", 5), mapper.NewUserMapper())
}

func newHashedUser(email, username, plainPass string) *entity.User {
	hash, _ := password_encoder.Encoder(plainPass)
	return &entity.User{
		ID:       id.NewUUID(),
		Name:     "Test",
		Email:    email,
		Username: username,
		Password: hash,
		Roles:    []enums.Role{enums.USER},
	}
}

func TestAuthService_Register(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAuthService(userRepo)

	dto := &auth.RegisterDTO{
		Name: "Alice", Email: "alice@test.com", Username: "alice", Password: "securepass",
	}

	userRepo.On("ExistsUserByEmail", dto.Email).Return(false, nil)
	userRepo.On("ExistsUserByUserName", dto.Username).Return(false, nil)
	userRepo.On("Create", mock.AnythingOfType("*entity.User")).
		Return(&entity.User{
			ID: id.NewUUID(), Name: dto.Name, Email: dto.Email,
			Username: dto.Username, Roles: []enums.Role{enums.USER},
		}, nil)

	res, err := svc.Register(dto)

	assert.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "Bearer", res.TypeToken)
	assert.Equal(t, dto.Email, res.User.Email)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_DBError(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAuthService(userRepo)

	dto := &auth.RegisterDTO{
		Name: "Bob", Email: "bob@test.com", Username: "bob", Password: "pass123",
	}

	userRepo.On("ExistsUserByEmail", dto.Email).Return(false, nil)
	userRepo.On("ExistsUserByUserName", dto.Username).Return(false, nil)
	userRepo.On("Create", mock.AnythingOfType("*entity.User")).
		Return(nil, domainErr.ErrDatabase("duplicate email", nil))

	_, err := svc.Register(dto)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_ByEmail(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAuthService(userRepo)

	user := newHashedUser("carol@test.com", "carol", "mypassword")
	userRepo.On("FindByEmailOrUsername", "carol@test.com").Return(user, nil)

	res, err := svc.Login(&auth.LoginRequest{Login: "carol@test.com", Password: "mypassword"})

	assert.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_ByUsername(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAuthService(userRepo)

	user := newHashedUser("dave@test.com", "dave", "davepass")
	userRepo.On("FindByEmailOrUsername", "dave").Return(user, nil)

	res, err := svc.Login(&auth.LoginRequest{Login: "dave", Password: "davepass"})

	assert.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAuthService(userRepo)

	user := newHashedUser("eve@test.com", "eve", "correctpass")
	userRepo.On("FindByEmailOrUsername", "eve@test.com").Return(user, nil)

	_, err := svc.Login(&auth.LoginRequest{Login: "eve@test.com", Password: "wrongpass"})

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAuthService(userRepo)

	userRepo.On("FindByEmailOrUsername", "ghost@test.com").
		Return(nil, domainErr.ErrNotFound("User", nil))

	_, err := svc.Login(&auth.LoginRequest{Login: "ghost@test.com", Password: "pass"})

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestAuthService_VerifyPassword_Correct(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAuthService(userRepo)

	user := newHashedUser("frank@test.com", "frank", "frankpass")
	userRepo.On("FindByID", user.ID).Return(user, nil)

	res, err := svc.VerifyPassword(&auth.VerifyPasswordRequest{Password: "frankpass"}, user)

	assert.NoError(t, err)
	assert.True(t, res.Result)
	userRepo.AssertExpectations(t)
}

func TestAuthService_VerifyPassword_Wrong(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAuthService(userRepo)

	user := newHashedUser("grace@test.com", "grace", "gracepass")
	userRepo.On("FindByID", user.ID).Return(user, nil)

	res, err := svc.VerifyPassword(&auth.VerifyPasswordRequest{Password: "wrongpass"}, user)

	assert.NoError(t, err)
	assert.False(t, res.Result)
	userRepo.AssertExpectations(t)
}

func TestAuthService_UpdatePassword(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAuthService(userRepo)

	user := newHashedUser("hank@test.com", "hank", "oldpass")
	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(user, nil)

	err := svc.UpdatePassword(&auth.UpdatePasswordRequest{
		CurrentPassword: "oldpass",
		NewPassword:     "newpass123",
	}, user)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestAuthService_UpdatePassword_WrongCurrent(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newAuthService(userRepo)

	user := newHashedUser("iris@test.com", "iris", "irispass")
	userRepo.On("FindByID", user.ID).Return(user, nil)

	err := svc.UpdatePassword(&auth.UpdatePasswordRequest{
		CurrentPassword: "wrongcurrent",
		NewPassword:     "newpass123",
	}, user)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}
