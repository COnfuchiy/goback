package repository

import (
	"github.com/google/uuid"
	"goback/domain/entity"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *entity.User) error
	FindByID(id uuid.UUID) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) Create(user *entity.User) error {
	return r.db.Model(&entity.User{}).Create(&user).Error
}

func (r UserRepository) FindByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	result := r.db.Model(&entity.User{}).Where("id = ?", id.String()).First(&user)
	return &user, result.Error
}

func (r UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := r.db.Model(&entity.User{}).Where("email = ?", email).First(&user)
	return &user, result.Error
}
