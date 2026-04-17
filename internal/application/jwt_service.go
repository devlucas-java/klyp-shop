package application

import (
	"errors"

	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
	"github.com/go-chi/jwtauth"
)

type JWTService struct {
	tokenAuth *jwtauth.JWTAuth
}

func NewJWTService(secret string) *JWTService {
	j := jwtauth.New("HS256", []byte(secret), nil)
	return &JWTService{tokenAuth: j}
}

func (s *JWTService) GenerateToken(user *entity.User) (string, error) {

	roles, authorities := extractRoles(user.Roles)

	_, token, err := s.tokenAuth.Encode(map[string]interface{}{
		"user_id":     user.ID.String(),
		"email":       user.Email,
		"roles":       roles,
		"authorities": authorities,
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

	claims := token.PrivateClaims()

	return claims, nil
}

func extractRoles(roles []entity.Role) ([]string, []string) {

	roleSet := make(map[string]bool)
	authSet := make(map[string]bool)

	for _, role := range roles {
		roleSet[role.Name] = true

		for _, auth := range role.Authorities {
			authSet[auth.Name] = true
		}
	}

	var roleList []string
	for r := range roleSet {
		roleList = append(roleList, r)
	}

	var authList []string
	for a := range authSet {
		authList = append(authList, a)
	}

	return roleList, authList
}
