package user_response

import "github.com/devlucas-java/klyp-shop/internal/domain/enums"

type UserDTO struct {
	Name     string       `json:"name"`
	Email    string       `json:"email"`
	Username string       `json:"username"`
	IsSeller bool         `json:"is_seller"`
	Roles    []enums.Role `json:"roles"`
}
