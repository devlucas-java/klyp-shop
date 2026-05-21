package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/auth"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/others"
	"github.com/devlucas-java/klyp-shop/internal/domain/apperrors"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
)

const authService = "auth_service.AuthService"

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
		return nil, apperrors.Unauthorized(authService+".login: invalid credentials", err)
	}

	match, err := user.VerifyPassword(login.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, apperrors.Unauthorized(authService+".login: invalid credentials", nil)
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, apperrors.Internal(authService+".login: failed to generate token", err)
	}

	return auth.NewJWTResponse(token, a.mapper.ToResponse(user)), nil
}

func (a *AuthService) Register(dto *auth.RegisterDTO) (*auth.JWTResponse, error) {
	exists, err := a.userRepository.ExistsUserByEmail(dto.Email)
	if err != nil {
		return nil, apperrors.Database(authService+".register: failed to check email", err)
	}
	if exists {
		return nil, apperrors.Conflict(authService+".register: email already in use", nil)
	}

	exists, err = a.userRepository.ExistsUserByUserName(dto.Username)
	if err != nil {
		return nil, apperrors.Database(authService+".register: failed to check username", err)
	}
	if exists {
		return nil, apperrors.Conflict(authService+".register: username already in use", nil)
	}

	user, err := entity.NewUser(dto.Name, dto.Email, dto.Username, dto.Password)
	if err != nil {
		return nil, apperrors.Internal(authService+".register: failed to create user", err)
	}

	user, err = a.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, apperrors.Internal(authService+".register: failed to generate token", err)
	}

	return auth.NewJWTResponse(token, a.mapper.ToResponse(user)), nil
}

func (a *AuthService) VerifyPassword(req *auth.VerifyPasswordRequest, user *entity.User) (*others.BooleanDTO, error) {
	stored, err := a.userRepository.FindByID(user.ID)
	if err != nil {
		return nil, apperrors.NotFound(authService+".verify_password: user not found", err)
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
		return apperrors.NotFound(authService+".update_password: user not found", err)
	}

	if err := stored.ChangePassword(dto.CurrentPassword, dto.NewPassword); err != nil {
		return err
	}

	_, err = a.userRepository.Update(stored)
	if err != nil {
		return apperrors.Database(authService+".update_password: failed to update password", err)
	}

	return nil
}
