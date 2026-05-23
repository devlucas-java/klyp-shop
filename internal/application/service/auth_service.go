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
		return nil, apperrors.InvalidCredentials(nil)
	}

	match, err := user.VerifyPassword(login.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, apperrors.InvalidCredentials(nil)
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return auth.NewJWTResponse(token, a.mapper.ToResponse(user)), nil
}

func (a *AuthService) Register(dto *auth.RegisterDTO) (*auth.JWTResponse, error) {
	exists, err := a.userRepository.ExistsUserByEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.Conflict("email is already in use", nil)
	}

	exists, err = a.userRepository.ExistsUserByUserName(dto.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.Conflict("username is already in use", nil)
	}

	user, err := entity.NewUser(dto.Name, dto.Email, dto.Username, dto.Password)
	if err != nil {
		return nil, err
	}
	user.ShoppingCart = *entity.NewShoppingCart(user.ID)

	user, err = a.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return auth.NewJWTResponse(token, a.mapper.ToResponse(user)), nil
}

func (a *AuthService) VerifyPassword(req *auth.VerifyPasswordRequest, auth *entity.User) (*others.BooleanDTO, error) {
	user, err := a.userRepository.FindByID(auth.ID)
	if err != nil {
		return nil, err
	}

	match, err := user.VerifyPassword(req.Password)
	if err != nil {
		return nil, err
	}

	return &others.BooleanDTO{Result: match}, nil
}

func (a *AuthService) UpdatePassword(dto *auth.UpdatePasswordRequest, auth *entity.User) error {
	user, err := a.userRepository.FindByID(auth.ID)
	if err != nil {
		return err
	}

	if err := user.ChangePassword(dto.CurrentPassword, dto.NewPassword); err != nil {
		return err
	}

	if _, err = a.userRepository.Update(user); err != nil {
		return err
	}

	return nil
}
