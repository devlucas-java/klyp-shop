package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/duser"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/database"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbUserSvc *gorm.DB
var userService *service.UserService

func setupUserService(t *testing.T) {
	var err error

	dbUserSvc, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = dbUserSvc.AutoMigrate(&entity.User{})
	if err != nil {
		t.Fatal(err)
	}

	log := logger.NewLogger(logger.TRACE)
	userRepo := database.NewUserDB(dbUserSvc, log)
	userMapper := mapper.NewUserMapper()
	userService = service.NewUserService(userRepo, log, userMapper)
}

func seedUser(t *testing.T) *entity.User {
	user, err := entity.NewUser("John", "john@test.com", "john123", "password123")
	if err != nil {
		t.Fatal(err)
	}
	if err := dbUserSvc.Create(user).Error; err != nil {
		t.Fatal(err)
	}
	return user
}

func TestGetMe(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)

	res, err := userService.GetMe(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, user.Name, res.Name)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.Username, res.Username)
}

func TestGetMe_NotFound(t *testing.T) {
	setupUserService(t)

	ghost := &entity.User{ID: id.NewUUID()}

	_, err := userService.GetMe(ghost)
	assert.Error(t, err)
}

func TestUpdateMe_Name(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)

	req := &duser.UpdateUserRequest{Name: "Updated Name"}

	res, err := userService.UpdateMe(user, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, "Updated Name", res.Name)
}

func TestUpdateMe_Email(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)

	req := &duser.UpdateUserRequest{Email: "newemail@test.com"}

	res, err := userService.UpdateMe(user, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, "newemail@test.com", res.Email)
}

func TestUpdateMe_EmailConflict(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)

	// Create another user with the target email
	other, err := entity.NewUser("Other", "taken@test.com", "other123", "pass")
	if err != nil {
		t.Fatal(err)
	}
	dbUserSvc.Create(other)

	req := &duser.UpdateUserRequest{Email: "taken@test.com"}

	_, err = userService.UpdateMe(user, req)
	assert.Error(t, err)
}

func TestUpdateMe_UsernameConflict(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)

	other, err := entity.NewUser("Other", "other2@test.com", "takenuser", "pass")
	if err != nil {
		t.Fatal(err)
	}
	dbUserSvc.Create(other)

	req := &duser.UpdateUserRequest{Username: "takenuser"}

	_, err = userService.UpdateMe(user, req)
	assert.Error(t, err)
}

func TestDeleteMe(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)

	err := userService.DeleteMe(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	dbUserSvc.Model(&entity.User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteMe_NotFound(t *testing.T) {
	setupUserService(t)

	ghost := &entity.User{ID: id.NewUUID()}

	err := userService.DeleteMe(ghost)
	assert.Error(t, err)
}

func TestPromoteToAdmin(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)

	err := userService.PromoteToAdmin(user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var updated entity.User
	dbUserSvc.First(&updated, "id = ?", user.ID)
	assert.True(t, updated.HasRole(enums.ADMIN))
}

func TestPromoteToAdmin_AlreadyAdmin(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)
	user.Roles = []enums.Role{enums.ADMIN}
	dbUserSvc.Save(user)

	err := userService.PromoteToAdmin(user.ID)
	assert.Error(t, err)
}

func TestPromoteToAdmin_IsSeller(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)
	user.IsSeller = true
	dbUserSvc.Save(user)

	err := userService.PromoteToAdmin(user.ID)
	assert.Error(t, err)
}

func TestDemoteToUser(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)
	user.Roles = []enums.Role{enums.ADMIN}
	dbUserSvc.Save(user)

	err := userService.DemoteToUser(user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var updated entity.User
	dbUserSvc.First(&updated, "id = ?", user.ID)
	assert.True(t, updated.HasRole(enums.USER))
}

func TestDemoteToUser_AlreadyUser(t *testing.T) {
	setupUserService(t)

	user := seedUser(t)
	// user already has USER role by default

	err := userService.DemoteToUser(user.ID)
	assert.Error(t, err)
}
