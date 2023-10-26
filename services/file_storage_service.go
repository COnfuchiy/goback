package services

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type IFileStorageService interface {
	SaveFileToStorage(fileID uuid.UUID, fileObject *multipart.FileHeader)
}

type FileStorageService struct {
}

func NewFileStorageService() IFileStorageService {
	return &FileStorageService{}
}

func (s FileStorageService) SaveFileToStorage(fileID uuid.UUID, fileObject *multipart.FileHeader) {
}
