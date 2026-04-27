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

var dbUser *gorm.DB
var userRepo *database.UserDB
var logUser *logger.Logger

func setupUserDB(t *testing.T) {
	var err error

	dbUser, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = dbUser.AutoMigrate(&entity.User{})
	if err != nil {
		t.Fatal(err)
	}

	logUser = logger.NewLogger(logger.TRACE)
	userRepo = database.NewUserDB(dbUser, logUser).(*database.UserDB)
}

func TestCreateUser(t *testing.T) {
	setupUserDB(t)

	user, err := entity.NewUser(
		"John",
		"john@test.com",
		"john123",
		"hashed-password",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, res.Name, user.Name)
	assert.Equal(t, res.Email, user.Email)
	assert.Equal(t, res.Username, user.Username)
	assert.Equal(t, res.Password, user.Password)
}

func TestFindByID(t *testing.T) {
	setupUserDB(t)

	user := &entity.User{
		Name:     "Jane",
		Email:    "jane@test.com",
		Username: "jane123",
		Password: "hash",
	}

	dbUser.Create(user)

	found, err := userRepo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email != user.Email {
		t.Fatal("user mismatch")
	}
}

func TestUpdateUser(t *testing.T) {
	setupUserDB(t)

	user, err := entity.NewUser(
		"old",
		"old@email.com",
		"olduser",
		"hash",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dbUser.Create(user)

	user.Name = "New"

	updated, err := userRepo.Update(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.Name != "New" {
		t.Fatal("update failed")
	}
}

func TestDeleteUser(t *testing.T) {
	setupUserDB(t)

	user := &entity.User{
		ID:       id.NewUUID(),
		Name:     "ToDelete",
		Email:    "delete@test.com",
		Username: "delete",
		Password: "hash",
	}

	dbUser.Create(user)

	err := userRepo.DeleteByID(user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int64
	dbUser.Model(&entity.User{}).Where("id = ?", user.ID).Count(&count)

	if count != 0 {
		t.Fatal("user was not deleted")
	}
}

func TestFindByEmailOrUsername(t *testing.T) {
	setupUserDB(t)

	user := &entity.User{
		ID:       id.NewUUID(),
		Name:     "Search",
		Email:    "search@test.com",
		Username: "searchuser",
		Password: "hash",
	}

	dbUser.Create(user)

	found, err := userRepo.FindByEmailOrUsername("search@test.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Username != user.Username {
		t.Fatal("email search failed")
	}

	found2, err := userRepo.FindByEmailOrUsername("searchuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found2.Email != user.Email {
		t.Fatal("username search failed")
	}
}
