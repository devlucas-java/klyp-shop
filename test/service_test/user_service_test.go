package service_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/internal/application/service"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	userDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/user"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
	"github.com/devlucas-java/klyp-shop/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newUserService(userRepo *mocks.UserRepositoryMock) *service.UserService {
	return service.NewUserService(userRepo, logger.NewLogger(logger.TRACE), mapper.NewUserMapper())
}

func newUser() *entity.User {
	return &entity.User{
		ID:       id.NewUUID(),
		Name:     "John",
		Email:    "john@test.com",
		Username: "john123",
		Password: "hash",
		Roles:    []enums.Role{enums.USER},
	}
}

// ── GetMe ─────────────────────────────────────────────────────────────────────

func TestGetMe(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	userRepo.On("FindByID", user.ID).Return(user, nil)

	res, err := svc.GetMe(user)

	assert.NoError(t, err)
	assert.Equal(t, user.Name, res.Name)
	assert.Equal(t, user.Email, res.Email)
	userRepo.AssertExpectations(t)
}

func TestGetMe_NotFound(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	ghost := &entity.User{ID: id.NewUUID()}
	userRepo.On("FindByID", ghost.ID).Return(nil, apperrors.NotFound("User", nil))

	_, err := svc.GetMe(ghost)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

// ── UpdateMe ──────────────────────────────────────────────────────────────────

func TestUpdateMe_Name(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	updated := *user
	updated.Name = "Updated Name"

	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(&updated, nil)

	res, err := svc.UpdateMe(user, &userDTO.UpdateUserRequest{Name: "Updated Name"})

	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", res.Name)
	userRepo.AssertExpectations(t)
}

func TestUpdateMe_Email_NewAvailable(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	updated := *user
	updated.Email = "new@test.com"

	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("ExistsUserByEmail", "new@test.com").Return(false, nil)
	userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(&updated, nil)

	res, err := svc.UpdateMe(user, &userDTO.UpdateUserRequest{Email: "new@test.com"})

	assert.NoError(t, err)
	assert.Equal(t, "new@test.com", res.Email)
	userRepo.AssertExpectations(t)
}

func TestUpdateMe_Email_SameAsCurrent_SkipsConflictCheck(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	updated := *user

	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(&updated, nil)

	// sending the same email should not trigger ExistsUserByEmail
	res, err := svc.UpdateMe(user, &userDTO.UpdateUserRequest{Email: user.Email})

	assert.NoError(t, err)
	assert.Equal(t, user.Email, res.Email)
	userRepo.AssertNotCalled(t, "ExistsUserByEmail")
	userRepo.AssertExpectations(t)
}

func TestUpdateMe_EmailConflict(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()

	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("ExistsUserByEmail", "taken@test.com").Return(true, nil)

	_, err := svc.UpdateMe(user, &userDTO.UpdateUserRequest{Email: "taken@test.com"})

	assert.Error(t, err)
	userRepo.AssertNotCalled(t, "Update")
	userRepo.AssertExpectations(t)
}

func TestUpdateMe_Username_NewAvailable(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	updated := *user
	updated.Username = "newuser"

	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("ExistsUserByUserName", "newuser").Return(false, nil)
	userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(&updated, nil)

	res, err := svc.UpdateMe(user, &userDTO.UpdateUserRequest{Username: "newuser"})

	assert.NoError(t, err)
	assert.Equal(t, "newuser", res.Username)
	userRepo.AssertExpectations(t)
}

func TestUpdateMe_Username_SameAsCurrent_SkipsConflictCheck(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	updated := *user

	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(&updated, nil)

	res, err := svc.UpdateMe(user, &userDTO.UpdateUserRequest{Username: user.Username})

	assert.NoError(t, err)
	assert.Equal(t, user.Username, res.Username)
	userRepo.AssertNotCalled(t, "ExistsUserByUserName")
	userRepo.AssertExpectations(t)
}

func TestUpdateMe_UsernameConflict(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()

	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("ExistsUserByUserName", "takenuser").Return(true, nil)

	_, err := svc.UpdateMe(user, &userDTO.UpdateUserRequest{Username: "takenuser"})

	assert.Error(t, err)
	userRepo.AssertNotCalled(t, "Update")
	userRepo.AssertExpectations(t)
}

func TestUpdateMe_UserNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	ghost := &entity.User{ID: id.NewUUID()}
	userRepo.On("FindByID", ghost.ID).Return(nil, apperrors.NotFound("User", nil))

	_, err := svc.UpdateMe(ghost, &userDTO.UpdateUserRequest{Name: "X"})

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

// ── DeleteMe ──────────────────────────────────────────────────────────────────

func TestDeleteMe(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("DeleteByID", user.ID).Return(nil)

	err := svc.DeleteMe(user)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestDeleteMe_NotFound(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	ghost := &entity.User{ID: id.NewUUID()}
	userRepo.On("FindByID", ghost.ID).Return(nil, apperrors.NotFound("User", nil))

	err := svc.DeleteMe(ghost)

	assert.Error(t, err)
	userRepo.AssertNotCalled(t, "DeleteByID")
	userRepo.AssertExpectations(t)
}

// ── PromoteToAdmin ────────────────────────────────────────────────────────────

func TestPromoteToAdmin(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	promoted := *user
	promoted.Roles = []enums.Role{enums.ADMIN}

	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(&promoted, nil)

	err := svc.PromoteToAdmin(user.ID)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestPromoteToAdmin_AlreadyAdmin(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	user.Roles = []enums.Role{enums.ADMIN}
	userRepo.On("FindByID", user.ID).Return(user, nil)

	err := svc.PromoteToAdmin(user.ID)

	assert.Error(t, err)
	userRepo.AssertNotCalled(t, "Update")
	userRepo.AssertExpectations(t)
}

func TestPromoteToAdmin_IsSeller(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	user.IsSeller = true
	userRepo.On("FindByID", user.ID).Return(user, nil)

	err := svc.PromoteToAdmin(user.ID)

	assert.Error(t, err)
	userRepo.AssertNotCalled(t, "Update")
	userRepo.AssertExpectations(t)
}

func TestPromoteToAdmin_UserNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	ghostID := id.NewUUID()
	userRepo.On("FindByID", ghostID).Return(nil, apperrors.NotFound("User", nil))

	err := svc.PromoteToAdmin(ghostID)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

// ── DemoteToUser ──────────────────────────────────────────────────────────────

func TestDemoteToUser(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	user.Roles = []enums.Role{enums.ADMIN}
	demoted := *user
	demoted.Roles = []enums.Role{enums.USER}

	userRepo.On("FindByID", user.ID).Return(user, nil)
	userRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(&demoted, nil)

	err := svc.DemoteToUser(user.ID)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestDemoteToUser_AlreadyUser(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser() // role is USER by default
	userRepo.On("FindByID", user.ID).Return(user, nil)

	err := svc.DemoteToUser(user.ID)

	assert.Error(t, err)
	userRepo.AssertNotCalled(t, "Update")
	userRepo.AssertExpectations(t)
}

func TestDemoteToUser_IsSeller(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	user := newUser()
	user.IsSeller = true
	user.Roles = []enums.Role{enums.SELLER}
	userRepo.On("FindByID", user.ID).Return(user, nil)

	err := svc.DemoteToUser(user.ID)

	assert.Error(t, err)
	userRepo.AssertNotCalled(t, "Update")
	userRepo.AssertExpectations(t)
}

func TestDemoteToUser_UserNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	svc := newUserService(userRepo)

	ghostID := id.NewUUID()
	userRepo.On("FindByID", ghostID).Return(nil, apperrors.NotFound("User", nil))

	err := svc.DemoteToUser(ghostID)

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}
