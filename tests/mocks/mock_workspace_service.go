package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"goback/domain/entity"
)

type MockWorkspaceService struct {
	mock.Mock
}

func (mock *MockWorkspaceService) Create(workspace *entity.Workspace) error {
	// isn`t work?

	args := mock.Called(workspace.Name)

	var err error
	if args.Get(0) != nil {
		err = args.Get(0).(error)
	}

	return err
}

func (mock *MockWorkspaceService) GetAllByUserID(userID uuid.UUID, currentPage int) ([]entity.Workspace, int64, error) {
	args := mock.Called(userID, currentPage)

	var workspaces []entity.Workspace
	if args.Get(0) != nil {
		workspaces = args.Get(0).([]entity.Workspace)
	}

	var totalCount int64
	if args.Get(1) != nil {
		totalCount = int64(args.Get(1).(int))
	}

	var err error
	if args.Get(2) != nil {
		err = args.Get(2).(error)
	}
	return workspaces, totalCount, err

}

func (mock *MockWorkspaceService) GetWorkspaceFromContext(contextWorkspaceID string) (*entity.Workspace, error) {
	// TODO implement me
	panic("implement me")
}

func (mock *MockWorkspaceService) CheckUserAccess(workspace *entity.Workspace, userID uuid.UUID) (bool, error) {
	// TODO implement me
	panic("implement me")
}
