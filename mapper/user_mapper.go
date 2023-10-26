package mapper

import (
	"goback/api/request"
	"goback/api/response"
	"goback/domain/entity"
)

type IUserMapper interface {
	FromRegisterRequest(registerRequest *request.RegisterRequest) *entity.User
	ToProfileResponse(user *entity.User) *response.ProfileResponse
}

type UserMapper struct {
}

func NewUserMapper() IUserMapper {
	return &UserMapper{}
}

func (m UserMapper) FromRegisterRequest(registerRequest *request.RegisterRequest) *entity.User {
	return &entity.User{
		Username: &registerRequest.Username,
		Email:    &registerRequest.Email,
		Password: &registerRequest.Password,
	}
}

func (m UserMapper) ToProfileResponse(user *entity.User) *response.ProfileResponse {
	return &response.ProfileResponse{
		Username: *user.Username,
		Email:    *user.Email,
	}
}
