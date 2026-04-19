package auth_response

import "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/response/user_response"

type JWTDTO struct {
	Token     string                 `json:"token"`
	TypeToken string                 `json:"typeToken"`
	User      *user_response.UserDTO `json:"user_response"`
}

func NewJWTDTO(token string, user *user_response.UserDTO) *JWTDTO {
	return &JWTDTO{
		Token:     token,
		TypeToken: "Bearer",
		User:      user,
	}
}
