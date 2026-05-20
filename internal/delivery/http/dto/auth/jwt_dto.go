package auth

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/user"
)

type JWTResponse struct {
	Token     string             `json:"token"`
	TypeToken string             `json:"typeToken"`
	User      *user.UserResponse `json:"user_response"`
}

func NewJWTResponse(token string, user *user.UserResponse) *JWTResponse {
	return &JWTResponse{
		Token:     token,
		TypeToken: "Bearer",
		User:      user,
	}
}
