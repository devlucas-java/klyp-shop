package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/request/user_request"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/response/user_response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type UserService struct {
	userRepository repository.UserRepository
	log            *logger.Logger
	userMapper     *mapper.UserMapper
}

func NewUserService(userRepository repository.UserRepository, log *logger.Logger, userMapper *mapper.UserMapper) *UserService {
	return &UserService{
		userRepository: userRepository,
		log:            log,
		userMapper:     userMapper,
	}
}

func (s *UserService) GetMe(auth *entity.User) (*user_response.UserDTO, error) {
	s.log.Infof("Getting user by ID %s", auth.ID)
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", auth.ID, err)
		return nil, errors.ErrNotFound("User", err)
	}

	return s.userMapper.UserToUserDTO(user), nil
}

func (s *UserService) UpdateMe(auth *entity.User, req *user_request.UpdateUserRequest) (*user_response.UserDTO, error) {
	s.log.Infof("Updating user %s", auth.ID)
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", auth.ID, err)
		return nil, errors.ErrNotFound("User", err)
	}

	if req.Name != "" && req.Name != user.Name {
		user.Name = req.Name
	}

	if req.Email != "" && req.Email != user.Email {
		userWithEmail, _ := s.userRepository.FindByEmailOrUsername(req.Email)
		if userWithEmail != nil {
			return nil, errors.ErrConflict("User", nil)
		}
		user.Email = req.Email
	}

	if req.Username != "" && req.Username != user.Username {
		userWithUsername, _ := s.userRepository.FindByEmailOrUsername(req.Username)
		if userWithUsername != nil {
			return nil, errors.ErrConflict("User", nil)
		}
		user.Username = req.Username
	}

	_, err = s.userRepository.Update(user)
	if err != nil {
		s.log.Errorf("Failed to update user %s: %v", auth.ID, err)
		return nil, errors.ErrNotFound("User", err)
	}

	return s.userMapper.UserToUserDTO(user), nil
}

func (s *UserService) DeleteMe(auth *entity.User) error {
	s.log.Infof("Deleting user %s", auth.ID)
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", auth.ID, err)
		return errors.ErrNotFound("User", err)
	}

	err = s.userRepository.Delete(user.ID)
	if err != nil {
		s.log.Errorf("Failed to delete user %s: %v", auth.ID, err)
		return errors.ErrNotFound("User", err)
	}

	return nil
}

func (s *UserService) PromoteToAdmin(id id.UUID) error {
	s.log.Infof("Attempting to promote user %s to admin", id)
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", id, err)
		return errors.ErrNotFound("User", err)
	}

	if user.IsSeller {
		s.log.Warnf("Cannot promote seller %s to admin", id)
		return errors.ErrInvalidRole("seller cannot be promoted to admin", err)
	}

	if user.HasRole(enums.ADMIN) {
		s.log.Warnf("User %s is already an admin", id)
		return errors.ErrInvalidRole("user is already an admin", err)
	}

	user.Roles = []enums.Role{enums.ADMIN}

	_, err = s.userRepository.Update(user)
	if err != nil {
		s.log.Errorf("Failed to update user %s to admin role: %v", id, err)
		return errors.Wrap("UPDATE_ERROR", "failed to update user roles", 500, err)
	}

	s.log.Infof("Successfully promoted user %s to admin", id)
	return nil
}

func (s *UserService) DemoteToUser(id id.UUID) error {
	s.log.Infof("Attempting to demote user %s to regular user", id)
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", id, err)
		return errors.ErrNotFound("User", err)
	}

	if user.HasRole(enums.USER) || user.IsSeller {
		s.log.Warnf("User %s cannot be demoted to user (isSeller: %v, hasUserRole: %v)", id, user.IsSeller, user.HasRole(enums.USER))
		msg := ""
		if user.IsSeller {
			msg = "seller cannot be demoted to user"
		} else if user.HasRole(enums.USER) {
			msg = "user is already a user"
		} else {
			msg = "user cannot be demoted to user"
		}
		return errors.ErrInvalidRole(msg, err)
	}

	user.IsSeller = false
	user.Roles = []enums.Role{enums.USER}

	_, err = s.userRepository.Update(user)
	if err != nil {
		s.log.Errorf("Failed to update user %s to user role: %v", id, err)
		return errors.Wrap("UPDATE_ERROR", "failed to update user roles", 500, err)
	}

	s.log.Infof("Successfully demoted user %s to regular user", id)
	return nil
}
