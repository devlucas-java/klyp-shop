package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dauth"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbAuthSvc *gorm.DB
var authService *service.AuthService

func setupAuthService(t *testing.T) {
	var err error

	dbAuthSvc, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = dbAuthSvc.AutoMigrate(&entity.User{})
	if err != nil {
		t.Fatal(err)
	}

	log := logger.NewLogger(logger.TRACE)
	userRepo := database.NewUserDB(dbAuthSvc, log)
	jwtSvc := jwt.NewJWTService("test-secret-key")
	userMapper := mapper.NewUserMapper()
	authService = service.NewAuthService(userRepo, jwtSvc, userMapper)
}

func TestAuthService_Register(t *testing.T) {
	setupAuthService(t)

	dto := &dauth.RegisterDTO{
		Name:     "Alice",
		Email:    "alice@test.com",
		Username: "alice",
		Password: "securepass",
	}

	res, err := authService.Register(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "Bearer", res.TypeToken)
	assert.Equal(t, "alice@test.com", res.User.Email)
	assert.Equal(t, "alice", res.User.Username)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	setupAuthService(t)

	dto := &dauth.RegisterDTO{
		Name:     "Bob",
		Email:    "bob@test.com",
		Username: "bob",
		Password: "pass",
	}

	_, err := authService.Register(dto)
	if err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	dto2 := &dauth.RegisterDTO{
		Name:     "Bob2",
		Email:    "bob@test.com",
		Username: "bob2",
		Password: "pass",
	}

	_, err = authService.Register(dto2)
	assert.Error(t, err)
}

func TestAuthService_Login_ByEmail(t *testing.T) {
	setupAuthService(t)

	dto := &dauth.RegisterDTO{
		Name:     "Carol",
		Email:    "carol@test.com",
		Username: "carol",
		Password: "mypassword",
	}
	_, err := authService.Register(dto)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	req := &dauth.LoginRequest{
		Login:    "carol@test.com",
		Password: "mypassword",
	}

	res, err := authService.Login(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "carol@test.com", res.User.Email)
}

func TestAuthService_Login_ByUsername(t *testing.T) {
	setupAuthService(t)

	dto := &dauth.RegisterDTO{
		Name:     "Dave",
		Email:    "dave@test.com",
		Username: "dave",
		Password: "davepass",
	}
	_, err := authService.Register(dto)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	req := &dauth.LoginRequest{
		Login:    "dave",
		Password: "davepass",
	}

	res, err := authService.Login(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "dave", res.User.Username)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	setupAuthService(t)

	dto := &dauth.RegisterDTO{
		Name:     "Eve",
		Email:    "eve@test.com",
		Username: "eve",
		Password: "correctpass",
	}
	_, err := authService.Register(dto)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	req := &dauth.LoginRequest{
		Login:    "eve@test.com",
		Password: "wrongpass",
	}

	_, err = authService.Login(req)
	assert.Error(t, err)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	setupAuthService(t)

	req := &dauth.LoginRequest{
		Login:    "ghost@test.com",
		Password: "pass",
	}

	_, err := authService.Login(req)
	assert.Error(t, err)
}

func TestAuthService_VerifyPassword_Correct(t *testing.T) {
	setupAuthService(t)

	dto := &dauth.RegisterDTO{
		Name:     "Frank",
		Email:    "frank@test.com",
		Username: "frank",
		Password: "frankpass",
	}
	registered, err := authService.Register(dto)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user := &entity.User{}
	dbAuthSvc.First(user, "email = ?", registered.User.Email)

	req := &dauth.VerifyPasswordRequest{Password: "frankpass"}

	res, err := authService.VerifyPassword(req, user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.True(t, res.Result)
}

func TestAuthService_VerifyPassword_Wrong(t *testing.T) {
	setupAuthService(t)

	dto := &dauth.RegisterDTO{
		Name:     "Grace",
		Email:    "grace@test.com",
		Username: "grace",
		Password: "gracepass",
	}
	registered, err := authService.Register(dto)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user := &entity.User{}
	dbAuthSvc.First(user, "email = ?", registered.User.Email)

	req := &dauth.VerifyPasswordRequest{Password: "wrongpass"}

	res, err := authService.VerifyPassword(req, user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.False(t, res.Result)
}

func TestAuthService_UpdatePassword(t *testing.T) {
	setupAuthService(t)

	dto := &dauth.RegisterDTO{
		Name:     "Hank",
		Email:    "hank@test.com",
		Username: "hank",
		Password: "oldpass",
	}
	registered, err := authService.Register(dto)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user := &entity.User{}
	dbAuthSvc.First(user, "email = ?", registered.User.Email)

	req := &dauth.UpdatePasswordRequest{
		CurrentPassword: "oldpass",
		NewPassword:     "newpass123",
	}

	err = authService.UpdatePassword(req, user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify new password works
	loginReq := &dauth.LoginRequest{
		Login:    "hank@test.com",
		Password: "newpass123",
	}
	res, err := authService.Login(loginReq)
	if err != nil {
		t.Fatalf("login with new password failed: %v", err)
	}
	assert.NotEmpty(t, res.Token)
}

func TestAuthService_UpdatePassword_WrongCurrent(t *testing.T) {
	setupAuthService(t)

	dto := &dauth.RegisterDTO{
		Name:     "Iris",
		Email:    "iris@test.com",
		Username: "iris",
		Password: "irispass",
	}
	registered, err := authService.Register(dto)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user := &entity.User{}
	dbAuthSvc.First(user, "email = ?", registered.User.Email)

	req := &dauth.UpdatePasswordRequest{
		CurrentPassword: "wrongcurrent",
		NewPassword:     "newpass",
	}

	err = authService.UpdatePassword(req, user)
	assert.Error(t, err)
}
