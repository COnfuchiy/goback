package services

import (
	"github.com/google/uuid"
	"goback/domain/entity"
	"goback/repository"
)

type IUserService interface {
	Create(user *entity.User) error
	GetByID(id uuid.UUID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetUserFromContext(userIDAsString string) (*entity.User, error)
}

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &UserService{userRepository}
}

func (s UserService) Create(user *entity.User) error {
	return s.userRepository.Create(user)
}

func (s UserService) GetByID(id uuid.UUID) (*entity.User, error) {
	return s.userRepository.FindByID(id)
}

func (s UserService) GetByEmail(email string) (*entity.User, error) {
	return s.userRepository.FindByEmail(email)
}

func (s UserService) GetUserFromContext(userIDAsString string) (*entity.User, error) {
	userID, err := uuid.Parse(userIDAsString)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
