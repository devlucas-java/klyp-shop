package service

import (
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type UserService struct {
	userRepository repository.UserRepository
	log            *logger.Logger
}

func NewUserService(userRepository repository.UserRepository, log *logger.Logger) *UserService {
	return &UserService{
		userRepository: userRepository,
		log:            log,
	}
}

func (s *UserService) PromoteToAdmin(userID id.UUID) error {
	s.log.Infof("Attempting to promote user %s to admin", userID)
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", userID, err)
		return errors.Wrap(errors.ErrUserNotFound.Code, "failed to find user", errors.ErrUserNotFound.Status, err)
	}

	if user.IsSeller {
		s.log.Warnf("Cannot promote seller %s to admin", userID)
		return errors.New(
			errors.ErrInvalidRole.Code,
			"seller cannot be promoted to admin",
			errors.ErrInvalidRole.Status,
			nil,
		)
	}

	if user.HasRole(enums.ADMIN) {
		s.log.Warnf("User %s is already an admin", userID)
		return errors.New(
			errors.ErrInvalidRole.Code,
			"user is already an admin",
			errors.ErrInvalidRole.Status,
			nil,
		)
	}

	user.Roles = []enums.Role{enums.ADMIN}

	_, err = s.userRepository.Update(user)
	if err != nil {
		s.log.Errorf("Failed to update user %s to admin role: %v", userID, err)
		return errors.Wrap("UPDATE_ERROR", "failed to update user roles", 500, err)
	}

	s.log.Infof("Successfully promoted user %s to admin", userID)
	return nil
}

func (s *UserService) DemoteToUser(userID id.UUID) error {
	s.log.Infof("Attempting to demote user %s to regular user", userID)
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", userID, err)
		return errors.Wrap(errors.ErrUserNotFound.Code, "failed to find user", errors.ErrUserNotFound.Status, err)
	}

	if user.HasRole(enums.USER) || user.IsSeller {
		s.log.Warnf("User %s cannot be demoted to user (isSeller: %v, hasUserRole: %v)", userID, user.IsSeller, user.HasRole(enums.USER))
		return errors.New(
			errors.ErrInvalidRole.Code,
			"seller or existing user not authorized to be demoted to user",
			errors.ErrInvalidRole.Status,
			nil,
		)
	}

	user.IsSeller = false
	user.Roles = []enums.Role{enums.USER}

	_, err = s.userRepository.Update(user)
	if err != nil {
		s.log.Errorf("Failed to update user %s to user role: %v", userID, err)
		return errors.Wrap("UPDATE_ERROR", "failed to update user roles", 500, err)
	}

	s.log.Infof("Successfully demoted user %s to regular user", userID)
	return nil
}

func (s *UserService) BecomeSeller(userID id.UUID) error {
	s.log.Infof("Attempting to make user %s a seller", userID)
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		s.log.Errorf("Failed to find user by ID %s: %v", userID, err)
		return errors.Wrap(errors.ErrUserNotFound.Code, "failed to find user", errors.ErrUserNotFound.Status, err)
	}

	if user.HasRole(enums.ADMIN) {
		s.log.Warnf("Admin %s cannot become a seller", userID)
		return errors.New(
			errors.ErrInvalidRole.Code,
			"admin cannot become a seller",
			errors.ErrInvalidRole.Status,
			nil,
		)
	}

	if user.IsSeller {
		s.log.Warnf("User %s is already a seller", userID)
		return errors.ErrUserIsAlreadySeller
	}

	user.IsSeller = true
	user.Roles = []enums.Role{enums.SELLER}

	_, err = s.userRepository.Update(user)
	if err != nil {
		s.log.Errorf("Failed to update user %s to seller role: %v", userID, err)
		return errors.Wrap("UPDATE_ERROR", "failed to update user to seller", 500, err)
	}

	s.log.Infof("Successfully updated user %s to seller", userID)
	return nil
}
