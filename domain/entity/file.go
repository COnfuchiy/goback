package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type File struct {
	ID            uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	FileHistoryID uuid.UUID
	FileHistory   FileHistory
	Filename      string `gorm:"column:filename;not null"`
	Tag           string `gorm:"column:tag"`
	Size          int64  `gorm:"column:size;not null"`
	DownloadURL   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UserId        uuid.UUID
	User          User
}

func (f File) AfterFind(context *gorm.DB) error {
	f.DownloadURL = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	return nil
}
