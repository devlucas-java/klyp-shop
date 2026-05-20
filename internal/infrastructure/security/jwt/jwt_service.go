package jwt

import (
	"errors"
	"time"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/go-chi/jwtauth"
)

type JWTService struct {
	tokenAuth *jwtauth.JWTAuth
	expireIn  time.Duration // duração do access token
}

func NewJWTService(secret string, expireInMinutes int) *JWTService {
	return &JWTService{
		tokenAuth: jwtauth.New("HS256", []byte(secret), nil),
		expireIn:  time.Duration(expireInMinutes) * time.Minute,
	}
}

func (s *JWTService) GenerateToken(user *entity.User) (string, error) {
	now := time.Now()

	_, token, err := s.tokenAuth.Encode(map[string]interface{}{
		"jti":     id.NewUUID().String(),
		"iat":     now.Unix(),
		"exp":     now.Add(s.expireIn).Unix(),
		"user_id": user.ID.String(),
		"email":   user.Email,
		"roles":   user.Roles,
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *JWTService) Validate(tokenString string) (map[string]interface{}, error) {
	token, err := jwtauth.VerifyToken(s.tokenAuth, tokenString)
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, errors.New("invalid token")
	}

	return token.PrivateClaims(), nil
}
