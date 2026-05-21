package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/user"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/policy"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

const userServiceTrace = "user_service.UserService"

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

func (s *UserService) GetMe(auth *entity.User) (*user.UserResponse, error) {
	u, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, apperrors.NotFound(userServiceTrace+".get_me: user not found", err)
	}
	return s.userMapper.ToResponse(u), nil
}

func (s *UserService) UpdateMe(auth *entity.User, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	u, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, apperrors.NotFound(userServiceTrace+".update_me: user not found", err)
	}

	if req.Email != "" && req.Email != u.Email {
		exists, err := s.userRepository.ExistsUserByEmail(req.Email)
		if err != nil {
			return nil, apperrors.Database(userServiceTrace+".update_me: failed to validate email", err)
		}
		if exists {
			return nil, apperrors.Conflict(userServiceTrace+".update_me: email already in use", nil)
		}
		u.ChangeEmail(req.Email)
	}

	if req.Username != "" && req.Username != u.Username {
		exists, err := s.userRepository.ExistsUserByUserName(req.Username)
		if err != nil {
			return nil, apperrors.Database(userServiceTrace+".update_me: failed to validate username", err)
		}
		if exists {
			return nil, apperrors.Conflict(userServiceTrace+".update_me: username already in use", nil)
		}
		u.ChangeUsername(req.Username)
	}

	u.ChangeName(req.Name)

	if _, err = s.userRepository.Update(u); err != nil {
		return nil, apperrors.Database(userServiceTrace+".update_me: failed to update user", err)
	}

	return s.userMapper.ToResponse(u), nil
}

func (s *UserService) DeleteMe(auth *entity.User) error {
	u, err := s.userRepository.FindByID(auth.ID)
	if err != nil {
		return apperrors.NotFound(userServiceTrace+".delete_me: user not found", err)
	}

	if err := s.userRepository.DeleteByID(u.ID); err != nil {
		return apperrors.Database(userServiceTrace+".delete_me: failed to delete user", err)
	}

	return nil
}

func (s *UserService) PromoteToAdmin(userID id.UUID) error {
	u, err := s.userRepository.FindByID(userID)
	if err != nil {
		return apperrors.NotFound(userServiceTrace+".promote_to_admin: user not found", err)
	}

	if err := s.userPolicy.CanPromoteToAdmin(u); err != nil {
		return err
	}

	u.ChangerToAdmin()

	if _, err := s.userRepository.Update(u); err != nil {
		return apperrors.Database(userServiceTrace+".promote_to_admin: failed to update user roles", err)
	}

	return nil
}

func (s *UserService) DemoteToUser(userID id.UUID) error {
	u, err := s.userRepository.FindByID(userID)
	if err != nil {
		return apperrors.NotFound(userServiceTrace+".demote_to_user: user not found", err)
	}

	if err := s.userPolicy.CanDemoteToUser(u); err != nil {
		return err
	}

	u.ChangerToUser()

	if _, err := s.userRepository.Update(u); err != nil {
		return apperrors.Database(userServiceTrace+".demote_to_user: failed to update user roles", err)
	}

	return nil
}
