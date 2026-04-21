package mapper

import (
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/request/auth_request"
	"github.com/devlucas-java/klyp-shop/internal/delivery/http/dto/response/user_response"
	"github.com/devlucas-java/klyp-shop/internal/domain/entity"
)

type UserMapper struct {
}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) UserToUserDTO(user *entity.User) *user_response.UserDTO {
	return &user_response.UserDTO{
		ID:       user.ID.String(),
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		IsSeller: user.IsSeller,
		Roles:    user.Roles,
	}
}

func (m *UserMapper) RegisterDTOToUser(dto *auth_request.RegisterDTO) *entity.User {

	return &entity.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Username: dto.Username,
	}
}
