package services

import (
	"github.com/google/uuid"
	"goback/domain/entity"
	"goback/repository"
)

type IWorkspaceService interface {
	Create(workspace *entity.Workspace) error
	GetAllByUserID(userID uuid.UUID, currentPage int) ([]entity.Workspace, int64, error)
	GetWorkspaceFromContext(contextWorkspaceID string) (*entity.Workspace, error)
	CheckUserAccess(workspace *entity.Workspace, userID uuid.UUID) (bool, error)
}

type WorkspaceService struct {
	workspaceRepository  repository.IWorkspaceRepository
	paginationRepository repository.IPaginateRepository
}

func NewWorkspaceService(workspaceRepository repository.IWorkspaceRepository, paginationRepository repository.IPaginateRepository) IWorkspaceService {
	return &WorkspaceService{workspaceRepository, paginationRepository}
}

func (s WorkspaceService) Create(workspace *entity.Workspace) error {
	return s.workspaceRepository.Create(workspace)
}

func (s WorkspaceService) GetAllByUserID(userID uuid.UUID, currentPage int) ([]entity.Workspace, int64, error) {
	offset, limit := s.paginationRepository.GetOffsetAndLimitFromPage(currentPage)
	return s.workspaceRepository.FindAllByUserID(userID, offset, limit)
}

func (s WorkspaceService) GetWorkspaceFromContext(contextWorkspaceID string) (*entity.Workspace, error) {
	workspaceID, err := uuid.Parse(contextWorkspaceID)
	if err != nil {
		return nil, err
	}
	workspace, err := s.workspaceRepository.FindByID(workspaceID)
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (s WorkspaceService) CheckUserAccess(workspace *entity.Workspace, userID uuid.UUID) (bool, error) {
	if workspace.CreatorID == userID {
		return true, nil
	}
	for _, workspaceUser := range workspace.Users {
		if workspaceUser.ID == userID {
			return true, nil
		}
	}
	return false, nil
}
