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

var dbUser *gorm.DB
var userRepo *database.UserDB
var logUser *logger.Logger

func setupUserDB(t *testing.T) {
	t.Helper()
	var err error

	dbUser, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = dbUser.AutoMigrate(&entity.User{})
	require.NoError(t, err)

	logUser = logger.NewLogger(logger.TRACE)
	userRepo = database.NewUserDB(dbUser, logUser).(*database.UserDB)
}

func TestCreateUser(t *testing.T) {
	setupUserDB(t)

	user, err := entity.NewUser("John", "john@test.com", "john123", "hashed-password")
	require.NoError(t, err)

	res, err := userRepo.Create(user)
	require.NoError(t, err)

	assert.Equal(t, user.Name, res.Name)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.Username, res.Username)
}

func TestFindByID(t *testing.T) {
	setupUserDB(t)

	user := &entity.User{
		ID:       id.NewUUID(),
		Name:     "Jane",
		Email:    "jane@test.com",
		Username: "jane123",
		Password: "hash",
	}
	require.NoError(t, dbUser.Create(user).Error)

	found, err := userRepo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Email, found.Email)
}

func TestUpdateUser(t *testing.T) {
	setupUserDB(t)

	user, err := entity.NewUser("old", "old@email.com", "olduser", "hash")
	require.NoError(t, err)
	require.NoError(t, dbUser.Create(user).Error)

	user.Name = "New"

	updated, err := userRepo.Update(user)
	require.NoError(t, err)
	assert.Equal(t, "New", updated.Name)
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
	require.NoError(t, dbUser.Create(user).Error)

	err := userRepo.DeleteByID(user.ID)
	require.NoError(t, err)

	var count int64
	dbUser.Model(&entity.User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(0), count)
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
	require.NoError(t, dbUser.Create(user).Error)

	found, err := userRepo.FindByEmailOrUsername("search@test.com")
	require.NoError(t, err)
	assert.Equal(t, user.Username, found.Username)

	found2, err := userRepo.FindByEmailOrUsername("searchuser")
	require.NoError(t, err)
	assert.Equal(t, user.Email, found2.Email)
}
