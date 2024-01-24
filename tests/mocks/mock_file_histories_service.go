package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"goback/domain/entity"
)

type MockFileHistoriesService struct {
	mock.Mock
}

func (mock *MockFileHistoriesService) Create(fileHistory *entity.FileHistory) error {
	args := mock.Called(fileHistory.WorkspaceID.String())

	var err error
	if args.Get(0) != nil {
		err = args.Get(0).(error)
	}
	return err
}

func (mock *MockFileHistoriesService) GetAllByWorkspaceID(workspaceID uuid.UUID, currentPage int) ([]entity.FileHistory, int64, error) {
	args := mock.Called(workspaceID, currentPage)

	var fileHistories []entity.FileHistory
	if args.Get(0) != nil {
		fileHistories = args.Get(0).([]entity.FileHistory)
	}

	var totalCount int64
	if args.Get(1) != nil {
		totalCount = int64(args.Get(1).(int))
	}

	var err error
	if args.Get(2) != nil {
		err = args.Get(2).(error)
	}
	return fileHistories, totalCount, err
}

func (mock *MockFileHistoriesService) GetFromContext(contextFileHistoryID string) (*entity.FileHistory, error) {
	args := mock.Called(contextFileHistoryID)

	var fileHistory *entity.FileHistory
	if args.Get(0) != nil {
		fileHistory = args.Get(0).(*entity.FileHistory)
	}

	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return fileHistory, err
}
