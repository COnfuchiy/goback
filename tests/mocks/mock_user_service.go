package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"goback/domain/entity"
)

type MockUserService struct {
	mock.Mock
}

func (mock *MockUserService) Create(user *entity.User) error {
	args := mock.Called(user)

	var err error
	if args.Get(0) != nil {
		err = args.Get(0).(error)
	}
	return err
}

func (mock *MockUserService) GetByID(id uuid.UUID) (*entity.User, error) {
	return nil, nil
}

func (mock *MockUserService) GetByEmail(email string) (*entity.User, error) {
	args := mock.Called(email)

	var user *entity.User
	if args.Get(0) != nil {
		user = args.Get(0).(*entity.User)
	}

	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return user, err
}

func (mock *MockUserService) GetUserFromContext(userIDAsString string) (*entity.User, error) {
	args := mock.Called(userIDAsString)

	var user *entity.User
	if args.Get(0) != nil {
		user = args.Get(0).(*entity.User)
	}

	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return user, err
}
