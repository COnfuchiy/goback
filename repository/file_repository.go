package repository

import (
	"github.com/google/uuid"
	"goback/domain/entity"
	"gorm.io/gorm"
)

type IFileRepository interface {
	Create(file *entity.File) error
	FindByID(id uuid.UUID) (*entity.File, error)
	CheckExisting(filename string) (bool, error)
}

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) IFileRepository {
	return &FileRepository{db: db}
}

func (r FileRepository) Create(file *entity.File) error {
	return r.db.Model(&entity.File{}).Create(&file).Error
}

func (r FileRepository) FindByID(id uuid.UUID) (*entity.File, error) {
	var file entity.File
	err := r.db.Model(&entity.File{}).Where("id = ?", id).First(&file).Error
	return &file, err
}

func (r FileRepository) CheckExisting(filename string) (bool, error) {
	var file entity.File
	err := r.db.Model(&entity.File{}).Where("filename = ?", filename).First(&file).Error
	if err == nil {
		return false, nil
	}
	if err != gorm.ErrRecordNotFound {
		return false, err
	} else {
		return true, nil
	}
}
