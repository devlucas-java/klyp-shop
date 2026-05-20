package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/duser"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type UserService struct {
	userRepository repository.UserRepository
	log            *logger.Logger
	userMapper     *mapper.UserMapper
	userPolicy     *policy.UserPolicy
}

func NewUserService(
	userRepository repository.UserRepository,
	log *logger.Logger,
	userMapper *mapper.UserMapper,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		log:            log,
		userMapper:     userMapper,
		userPolicy:     policy.NewUserPolicy(),
	}
}

func (s *UserService) GetMe(auth *entity.User) (*duser.UserResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", auth.ID, err)
		return nil, errors.ErrNotFound("User", err)
	}
	return s.userMapper.ToResponse(user), nil
}

func (s *UserService) UpdateMe(auth *entity.User, req *duser.UpdateUserRequest) (*duser.UserResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", auth.ID, err)
		return nil, errors.ErrNotFound("User", err)
	}

	if req.Email != "" && req.Email != user.Email {
		existing, _ := s.userRepository.FindByEmailOrUsername(req.Email)
		if existing != nil {
			return nil, errors.ErrConflict("Email", nil)
		}
		user.ChangeEmail(req.Email)
	}

	if req.Username != "" && req.Username != user.Username {
		existing, _ := s.userRepository.FindByEmailOrUsername(req.Username)
		if existing != nil {
			return nil, errors.ErrConflict("Username", nil)
		}
		user.ChangeUsername(req.Username)
	}

	user.ChangeName(req.Name)

	if _, err = s.userRepository.Update(user); err != nil {
		s.log.Errorf("Failed to update user %s: %v", auth.ID, err)
		return nil, errors.ErrDatabase("failed to update user", err)
	}

	return s.userMapper.ToResponse(user), nil
}

func (s *UserService) DeleteMe(auth *entity.User) error {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", auth.ID, err)
		return errors.ErrNotFound("User", err)
	}

	if err := s.userRepository.DeleteByID(user.ID); err != nil {
		s.log.Errorf("Failed to delete user %s: %v", auth.ID, err)
		return errors.ErrDatabase("failed to delete user", err)
	}

	return nil
}

func (s *UserService) PromoteToAdmin(userID id.UUID) error {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", userID, err)
		return errors.ErrNotFound("User", err)
	}

	if err := s.userPolicy.CanPromoteToAdmin(user); err != nil {
		s.log.Warnf("Cannot promote user %s to admin: %v", userID, err)
		return err
	}

	user.ChangerToAdmin()

	if _, err := s.userRepository.Update(user); err != nil {
		s.log.Errorf("Failed to update user %s roles: %v", userID, err)
		return errors.ErrDatabase("failed to update user roles", err)
	}

	s.log.Infof("User %s promoted to admin", userID)
	return nil
}

func (s *UserService) DemoteToUser(userID id.UUID) error {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", userID, err)
		return errors.ErrNotFound("User", err)
	}

	if err := s.userPolicy.CanDemoteToUser(user); err != nil {
		s.log.Warnf("Cannot demote user %s: %v", userID, err)
		return err
	}

	user.ChangerToUser()

	if _, err := s.userRepository.Update(user); err != nil {
		s.log.Errorf("Failed to update user %s roles: %v", userID, err)
		return errors.ErrDatabase("failed to update user roles", err)
	}

	s.log.Infof("User %s demoted to regular user", userID)
	return nil
}
