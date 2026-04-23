package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dauth"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/dothers"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/internal/domain/enums"
	"github.com/devlucas-java/klyp-shop/internal/domain/errors"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/repository"
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/security/jwt"
	"github.com/devlucas-java/klyp-shop/pkg/password_encoder"
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

	match, err := password_encoder.Match(login.Password, user.Password)
	if err != nil {
		return nil, errors.ErrInternal("failed to verify password", err)
	}

	if !match {
		return nil, errors.ErrInvalidCredentials(nil)
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate token", err)
	}

	return dauth.NewJWTResponse(token, a.mapper.UserToUserDTO(user)), nil
}

func (a *AuthService) Register(dto *dauth.RegisterDTO) (*dauth.JWTResponse, error) {
	user := a.mapper.RegisterDTOToUser(dto)

	pass, err := password_encoder.Encoder(dto.Password)
	if err != nil {
		return nil, errors.ErrInternal("failed to encode password", err)
	}

	user.Password = pass
	user.Roles = []enums.Role{enums.USER}

	user, err = a.userRepository.Create(user)
	if err != nil {
		return nil, errors.ErrDatabase("failed to create duser", err)
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate token", err)
	}

	return dauth.NewJWTResponse(token, a.mapper.UserToUserDTO(user)), nil
}

func (a *AuthService) VerifyPassword(req *dauth.VerifyPasswordRequest, user *entity.User) (*dothers.BooleanDTO, error) {
	stored, err := a.userRepository.FindByID(user.ID)
	if err != nil {
		return nil, errors.ErrInternal("failed to retrieve duser", err)
	}

	match, err := password_encoder.Match(req.Password, stored.Password)
	if err != nil {
		return nil, errors.ErrInternal("failed to verify password", err)
	}

	return &dothers.BooleanDTO{Result: match}, nil
}

func (a *AuthService) UpdatePassword(dto *dauth.UpdatePasswordRequest, user *entity.User) error {
	stored, err := a.userRepository.FindByID(user.ID)
	if err != nil {
		return errors.ErrInternal("failed to retrieve duser", err)
	}

	match, err := password_encoder.Match(dto.CurrentPassword, stored.Password)
	if err != nil {
		return errors.ErrInternal("failed to verify password", err)
	}

	if !match {
		return errors.ErrInvalidCredentials(nil)
	}

	hash, err := password_encoder.Encoder(dto.NewPassword)
	if err != nil {
		return errors.ErrInternal("failed to encode password", err)
	}

	stored.Password = hash

	_, err = a.userRepository.Update(stored)
	if err != nil {
		return errors.ErrDatabase("failed to update password", err)
	}

	return nil
}
