package entity

import (
	"github.com/google/uuid"
	"time"
)

type Workspace struct {
	ID             uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name           string    `gorm:"column:name"`
	CreatorID      uuid.UUID
	Creator        User
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Users          []User `gorm:"many2many:user_workspaces"`
	FilesHistories []FileHistory
}
