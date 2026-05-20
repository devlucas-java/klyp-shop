package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/duser"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	domainErrors "github.com/devlucas-java/klyp-shop/internal/domain/errors"
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
		if appErr := domainErrors.As(err); appErr != nil {
			return nil, appErr
		}
		return nil, domainErrors.ErrDatabase("failed to retrieve user", err)
	}
	return s.userMapper.ToResponse(user), nil
}

func (s *UserService) UpdateMe(auth *entity.User, req *duser.UpdateUserRequest) (*duser.UserResponse, error) {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", auth.ID, err)
		if appErr := domainErrors.As(err); appErr != nil {
			return nil, appErr
		}
		return nil, domainErrors.ErrDatabase("failed to retrieve user", err)
	}

	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepository.ExistsUserByEmail(req.Email)
		if err != nil {
			s.log.Errorf("Failed to check existing email %s: %v", req.Email, err)
			if appErr := domainErrors.As(err); appErr != nil {
				return nil, appErr
			}
			return nil, domainErrors.ErrDatabase("failed to validate email", err)
		}
		if exists {
			return nil, domainErrors.ErrConflict("Email", nil)
		}

		user.ChangeEmail(req.Email)
	}

	if req.Username != "" && req.Username != user.Username {
		exists, err := s.userRepository.ExistsUserByUserName(req.Username)
		if err != nil {
			s.log.Errorf("Failed to check existing username %s: %v", req.Username, err)
			if appErr := domainErrors.As(err); appErr != nil {
				return nil, appErr
			}
			return nil, domainErrors.ErrDatabase("failed to validate username", err)
		}
		if exists {
			return nil, domainErrors.ErrConflict("Username", nil)
		}

		user.ChangeUsername(req.Username)
	}

	user.ChangeName(req.Name)

	if _, err = s.userRepository.Update(user); err != nil {
		s.log.Errorf("Failed to update user %s: %v", auth.ID, err)
		if appErr := domainErrors.As(err); appErr != nil {
			return nil, appErr
		}
		return nil, domainErrors.ErrDatabase("failed to update user", err)
	}

	return s.userMapper.ToResponse(user), nil
}

func (s *UserService) DeleteMe(auth *entity.User) error {
	user, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", auth.ID, err)
		if appErr := domainErrors.As(err); appErr != nil {
			return appErr
		}
		return domainErrors.ErrDatabase("failed to retrieve user", err)
	}

	if err := s.userRepository.DeleteByID(user.ID); err != nil {
		s.log.Errorf("Failed to delete user %s: %v", auth.ID, err)
		if appErr := domainErrors.As(err); appErr != nil {
			return appErr
		}
		return domainErrors.ErrDatabase("failed to delete user", err)
	}

	return nil
}

func (s *UserService) PromoteToAdmin(userID id.UUID) error {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", userID, err)
		if appErr := domainErrors.As(err); appErr != nil {
			return appErr
		}
		return domainErrors.ErrDatabase("failed to retrieve user", err)
	}

	if err := s.userPolicy.CanPromoteToAdmin(user); err != nil {
		s.log.Warnf("Cannot promote user %s to admin: %v", userID, err)
		return err
	}

	user.ChangerToAdmin()

	if _, err := s.userRepository.Update(user); err != nil {
		s.log.Errorf("Failed to update user %s roles: %v", userID, err)
		if appErr := domainErrors.As(err); appErr != nil {
			return appErr
		}
		return domainErrors.ErrDatabase("failed to update user roles", err)
	}

	s.log.Infof("User %s promoted to admin", userID)
	return nil
}

func (s *UserService) DemoteToUser(userID id.UUID) error {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", userID, err)
		if appErr := domainErrors.As(err); appErr != nil {
			return appErr
		}
		return domainErrors.ErrDatabase("failed to retrieve user", err)
	}

	if err := s.userPolicy.CanDemoteToUser(user); err != nil {
		s.log.Warnf("Cannot demote user %s: %v", userID, err)
		return err
	}

	user.ChangerToUser()

	if _, err := s.userRepository.Update(user); err != nil {
		s.log.Errorf("Failed to update user %s roles: %v", userID, err)
		if appErr := domainErrors.As(err); appErr != nil {
			return appErr
		}
		return domainErrors.ErrDatabase("failed to update user roles", err)
	}

	s.log.Infof("User %s demoted to regular user", userID)
	return nil
}
