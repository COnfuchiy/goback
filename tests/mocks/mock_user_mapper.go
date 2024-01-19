package mocks

import (
	"github.com/stretchr/testify/mock"
	"goback/api/request"
	"goback/api/response"
	"goback/domain/entity"
)

type MockUserMapper struct {
	mock.Mock
}

func (mock *MockUserMapper) FromRegisterRequest(registerRequest *request.RegisterRequest) *entity.User {
	args := mock.Called(registerRequest)

	var user *entity.User
	if args.Get(0) != nil {
		user = args.Get(0).(*entity.User)
	}
	return user
}

func (mock *MockUserMapper) ToProfileResponse(user *entity.User) *response.ProfileResponse {
	return nil
}
