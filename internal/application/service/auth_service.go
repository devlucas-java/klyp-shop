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
}

func NewAuthService(userRepository repository.UserRepository, jwtService *jwt.JWTService) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (a *AuthService) Login(login *auth_request.LoginDTO) (*auth_response.JWTDTO, error) {
	user, err := a.userRepository.FindByEmailOrUsername(login.Login)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	match, _ := password_encoder.Match(login.Password, user.Password)
	if !match {
		return nil, errors.ErrInvalidCredentials
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return auth_response.NewJWTDTO(token, mapper.UserToUserDTO(user)), nil
}

func (a *AuthService) Register(dto *auth_request.RegisterDTO) (*auth_response.JWTDTO, error) {
	user := mapper.RegisterDTOToUser(dto)

	pass, err := password_encoder.Encoder(dto.Password)
	if err != nil {
		return nil, err
	}

	user.Password = pass
	user.Roles = []enums.Role{enums.USER}

	user, err = a.userRepository.Create(user)
	if err != nil {
		return nil, errors.Wrap(
			errors.ErrUserAlreadyExists.Code,
			errors.ErrUserAlreadyExists.Message,
			errors.ErrUserAlreadyExists.Status,
			err, // causa real fica no log, não vaza pro cliente
		)
	}

	token, _ := a.jwtService.GenerateToken(user)

	return auth_response.NewJWTDTO(token, mapper.UserToUserDTO(user)), nil
}

func (a *AuthService) VerifyPassword(req *auth_request.VerifyPasswordRequest, user *entity.User) (*others_response.BooleanDTO, error) {
	stored, err := a.userRepository.FindByID(user.ID)
	if err != nil {
		return nil, err
	}

	match, err := password_encoder.Match(req.Password, stored.Password)
	if err != nil {
		return nil, err
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
		return errors.ErrInvalidCredentials
	}

	hash, err := password_encoder.Encoder(dto.NewPassword)
	if err != nil {
		return err
	}

	stored.Password = hash

	_, err = a.userRepository.Update(stored)
	return err
}
