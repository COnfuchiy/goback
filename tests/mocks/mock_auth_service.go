package mocks

import (
	"github.com/stretchr/testify/mock"
	"goback/domain/entity"
)

type MockAuthService struct {
	mock.Mock
}

func (mock *MockAuthService) CreateAccessAndRefreshTokens(user *entity.User) (string, string, error) {
	args := mock.Called(user)

	var userAccessToken string
	if args.Get(0) != nil {
		userAccessToken = args.Get(0).(string)
	}

	var userRefreshToken string
	if args.Get(1) != nil {
		userRefreshToken = args.Get(1).(string)
	}

	var err error
	if args.Get(2) != nil {
		err = args.Get(2).(error)
	}
	return userAccessToken, userRefreshToken, err
}

func (mock *MockAuthService) ComparePassword(hashPassword string, password string) error {
	args := mock.Called(hashPassword, password)

	var err error
	if args.Get(0) != nil {
		err = args.Get(0).(error)
	}
	return err
}

func (mock *MockAuthService) HashPassword(password string) (string, error) {
	args := mock.Called(password)

	var hashedPassword string
	if args.Get(0) != nil {
		hashedPassword = args.Get(0).(string)
	}

	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return hashedPassword, err
}

func (mock *MockAuthService) ExtractIDFromAccessToken(requestToken string) (string, error) {
	args := mock.Called(requestToken)

	var userIDasString string
	if args.Get(0) != nil {
		userIDasString = args.Get(0).(string)
	}

	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return userIDasString, err
}

func (mock *MockAuthService) ExtractIDFromRefreshToken(requestToken string) (string, error) {
	args := mock.Called(requestToken)

	var userIDasString string
	if args.Get(0) != nil {
		userIDasString = args.Get(0).(string)
	}

	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return userIDasString, err

}
