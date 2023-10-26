package services

import (
	"github.com/google/uuid"
	"goback/domain/entity"
	"goback/repository"
)

type IFileHistoryService interface {
	Create(fileHistory *entity.FileHistory) error
	GetAllByWorkspaceID(workspaceID uuid.UUID, currentPage int) ([]entity.FileHistory, int64, error)
	GetFileHistoryFromContext(contextFileHistoryID string) (*entity.FileHistory, error)
}

type FileHistoryService struct {
	fileHistoryRepository repository.IFileHistoryRepository
	paginateRepository    repository.IPaginateRepository
}

func NewFileHistoryService(fileHistoryRepository repository.IFileHistoryRepository, paginateRepository repository.IPaginateRepository) *FileHistoryService {
	return &FileHistoryService{fileHistoryRepository, paginateRepository}
}

func (s FileHistoryService) Create(fileHistory *entity.FileHistory) error {
	return s.fileHistoryRepository.Create(fileHistory)
}

func (s FileHistoryService) GetAllByWorkspaceID(workspaceID uuid.UUID, currentPage int) ([]entity.FileHistory, int64, error) {
	offset, limit := s.paginateRepository.GetOffsetAndLimitFromPage(currentPage)
	return s.fileHistoryRepository.FindAllByWorkspaceID(workspaceID, offset, limit)
}

func (s FileHistoryService) GetFileHistoryFromContext(contextFileHistoryID string) (*entity.FileHistory, error) {
	fileHistoryID, err := uuid.Parse(contextFileHistoryID)
	if err != nil {
		return nil, err
	}
	fileHistory, err := s.fileHistoryRepository.FindByID(fileHistoryID)
	if err != nil {
		return nil, err
	}
	return fileHistory, nil
}
