package mapper

import (
	userDTO "github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/user"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToResponse(user *entity.User) *userDTO.UserResponse {
	return &userDTO.UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		IsSeller: user.IsSeller,
		Roles:    user.Roles,
	}
}
