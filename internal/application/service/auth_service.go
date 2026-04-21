package service

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/mapper"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/request/auth_request"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/response/auth_response"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/response/others_response"
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

func (a *AuthService) Login(login *auth_request.LoginDTO) (*auth_response.JWTDTO, error) {
	user, err := a.userRepository.FindByEmailOrUsername(login.Login)
	if err != nil {
		return nil, errors.ErrInvalidCredentials(err)
	}

	if match, _ := password_encoder.Match(login.Password, user.Password); !match {
		return nil, errors.ErrInvalidCredentials(err)
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return auth_response.NewJWTDTO(token, a.mapper.UserToUserDTO(user)), nil
}

func (a *AuthService) Register(dto *auth_request.RegisterDTO) (*auth_response.JWTDTO, error) {
	user := a.mapper.RegisterDTOToUser(dto)

	pass, err := password_encoder.Encoder(dto.Password)
	if err != nil {
		return nil, errors.ErrInternal("Failed to encode password", err)
	}

	user.Password = pass
	user.Roles = []enums.Role{enums.USER}

	user, err = a.userRepository.Create(user)
	if err != nil {
		return nil, errors.ErrDatabase("Failed to create user", err)

	}

	token, _ := a.jwtService.GenerateToken(user)
	return auth_response.NewJWTDTO(token, a.mapper.UserToUserDTO(user)), nil
}

func (a *AuthService) VerifyPassword(req *auth_request.VerifyPasswordRequest, user *entity.User) (*others_response.BooleanDTO, error) {
	stored, err := a.userRepository.FindByID(user.ID)
	if err != nil {
		return nil, errors.ErrInternal("Failed to retrieve user", err)
	}

	match, err := password_encoder.Match(req.Password, stored.Password)
	if err != nil {
		return nil, errors.ErrInternal("Failed to verify password", err)
	}

	return &others_response.BooleanDTO{Result: match}, nil
}

func (a *AuthService) UpdatePassword(dto *auth_request.UpdatePasswordRequest, user *entity.User) error {
	stored, err := a.userRepository.FindByID(user.ID)
	if err != nil {
		return err
	}

	match, err := password_encoder.Match(dto.CurrentPassword, stored.Password)
	if err != nil {
		return err
	}

	if !match {
		return errors.ErrInvalidCredentials(err)
	}

	hash, err := password_encoder.Encoder(dto.NewPassword)
	if err != nil {
		return err
	}

	stored.Password = hash

	_, err = a.userRepository.Update(stored)
	return err
}
