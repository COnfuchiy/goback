package repository

import (
	"github.com/google/uuid"
	"goback/domain/entity"
	"gorm.io/gorm"
)

type IFileHistoryRepository interface {
	Create(fileHistory *entity.FileHistory) error
	FindByID(id uuid.UUID) (*entity.FileHistory, error)
	FindAllByWorkspaceID(workspaceID uuid.UUID, offset int, limit int) ([]entity.FileHistory, int64, error)
}

type FileHistoryRepository struct {
	db *gorm.DB
}

func NewFileHistoryRepository(db *gorm.DB) IFileHistoryRepository {
	return &FileHistoryRepository{db: db}
}

func (r FileHistoryRepository) Create(fileHistory *entity.FileHistory) error {
	return r.db.Model(&entity.FileHistory{}).Create(&fileHistory).Error
}

func (r FileHistoryRepository) FindByID(id uuid.UUID) (*entity.FileHistory, error) {
	var fileHistory entity.FileHistory

	err := r.db.Model(&entity.FileHistory{}).Where("id = ?", id).Preload("Files").First(&fileHistory).Error

	return &fileHistory, err
}

func (r FileHistoryRepository) FindAllByWorkspaceID(workspaceID uuid.UUID, offset int, limit int) ([]entity.FileHistory, int64, error) {
	var fileHistories []entity.FileHistory
	var totalCount int64

	err := r.db.Model(&entity.FileHistory{}).Where("workspace_id = ?", workspaceID).Preload("Files",
		func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Order("files.created_at DESC").Limit(1)
		}).Offset(offset).Limit(limit).Count(&totalCount).Find(&fileHistories).Error
	return fileHistories, totalCount, err
}
