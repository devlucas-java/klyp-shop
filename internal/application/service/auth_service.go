package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/auth"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/others"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
)

type AuthService struct {
	userRepository repository.UserRepository
	jwtService     *jwt.JWTService
	mapper         *mapper.UserMapper
}

func NewAuthService(userRepository repository.UserRepository, jwtService *jwt.JWTService, mapper *mapper.UserMapper) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
		mapper:         mapper,
	}
}

func (a *AuthService) Login(login *auth.LoginRequest) (*auth.JWTResponse, error) {
	user, err := a.userRepository.FindByEmailOrUsername(login.Login)
	if err != nil {
		return nil, errors.ErrInvalidCredentials(err)
	}

	match, err := user.VerifyPassword(login.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, errors.ErrInvalidCredentials(nil)
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate token", err)
	}

	return auth.NewJWTResponse(token, a.mapper.ToResponse(user)), nil
}

func (a *AuthService) Register(dto *auth.RegisterDTO) (*auth.JWTResponse, error) {

	exists, err := a.userRepository.ExistsUserByEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrConflict("email", err)
	}

	exists, err = a.userRepository.ExistsUserByUserName(dto.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrConflict("username", err)
	}

	user, err := entity.NewUser(dto.Name, dto.Email, dto.Username, dto.Password)
	if err != nil {
		return nil, errors.ErrInternal("failed to create user", err)
	}

	user, err = a.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate token", err)
	}

	return auth.NewJWTResponse(token, a.mapper.ToResponse(user)), nil
}

func (a *AuthService) VerifyPassword(req *auth.VerifyPasswordRequest, user *entity.User) (*others.BooleanDTO, error) {
	stored, err := a.userRepository.FindByID(user.ID)
	if err != nil {
		return nil, errors.ErrInternal("failed to retrieve user", err)
	}

	match, err := stored.VerifyPassword(req.Password)
	if err != nil {
		return nil, err
	}

	return &others.BooleanDTO{Result: match}, nil
}

func (a *AuthService) UpdatePassword(dto *auth.UpdatePasswordRequest, user *entity.User) error {
	stored, err := a.userRepository.FindByID(user.ID)
	if err != nil {
		return errors.ErrInternal("failed to retrieve user", err)
	}

	if err := stored.ChangePassword(dto.CurrentPassword, dto.NewPassword); err != nil {
		return err
	}

	_, err = a.userRepository.Update(stored)
	if err != nil {
		return errors.ErrDatabase("failed to update password", err)
	}

	return nil
}
