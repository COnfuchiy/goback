package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Username   *string   `gorm:"column:username;unique;not null"`
	Email      *string   `gorm:"column:email;unique;not null"`
	Password   *string   `gorm:"column:password;not null" json:"-"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Workspaces []Workspace `gorm:"many2many:user_workspaces"`
}
