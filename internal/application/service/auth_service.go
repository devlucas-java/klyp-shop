package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dauth"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dothers"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
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

func (a *AuthService) Login(login *dauth.LoginRequest) (*dauth.JWTResponse, error) {
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

	return dauth.NewJWTResponse(token, a.mapper.ToResponse(user)), nil
}

func (a *AuthService) Register(dto *dauth.RegisterDTO) (*dauth.JWTResponse, error) {
	user, err := entity.NewUser(dto.Name, dto.Email, dto.Username, dto.Password)
	if err != nil {
		return nil, errors.ErrInternal("failed to create user", err)
	}

	user, err = a.userRepository.Create(user)
	if err != nil {
		return nil, errors.ErrDatabase("failed to create user", err)
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate token", err)
	}

	return dauth.NewJWTResponse(token, a.mapper.ToResponse(user)), nil
}

func (a *AuthService) VerifyPassword(req *dauth.VerifyPasswordRequest, user *entity.User) (*dothers.BooleanDTO, error) {
	stored, err := a.userRepository.FindByID(user.ID)
	if err != nil {
		return nil, errors.ErrInternal("failed to retrieve user", err)
	}

	match, err := stored.VerifyPassword(req.Password)
	if err != nil {
		return nil, err
	}

	return &dothers.BooleanDTO{Result: match}, nil
}

func (a *AuthService) UpdatePassword(dto *dauth.UpdatePasswordRequest, user *entity.User) error {
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
