package services

import (
	"github.com/google/uuid"
	"goback/domain/entity"
	"goback/repository"
)

type IFileService interface {
	Create(file *entity.File) error
	GetFromContext(contextFileID string) (*entity.File, error)
	CheckExisting(filename string) (bool, error)
}

type FileService struct {
	fileRepository repository.IFileRepository
}

func NewFileService(fileRepository repository.IFileRepository) IFileService {
	return &FileService{fileRepository: fileRepository}
}

func (s FileService) Create(file *entity.File) error {
	return s.fileRepository.Create(file)
}

func (s FileService) GetFromContext(contextFileID string) (*entity.File, error) {
	fileID, err := uuid.Parse(contextFileID)
	if err != nil {
		return nil, err
	}
	file, err := s.fileRepository.FindByID(fileID)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s FileService) CheckExisting(filename string) (bool, error) {
	return s.fileRepository.CheckExisting(filename)
}
