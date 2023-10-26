package entity

import (
	"github.com/google/uuid"
	"time"
)

type FileHistory struct {
	ID          uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	WorkspaceID uuid.UUID
	Workspace   Workspace
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Files       []File
}
