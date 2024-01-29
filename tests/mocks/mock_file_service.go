package mocks

import (
	"github.com/stretchr/testify/mock"
	"goback/domain/entity"
	"mime/multipart"
)

type MockFileService struct {
	mock.Mock
}

func (mock *MockFileService) GetFileTag(file *entity.File, fileObject *multipart.FileHeader) error {
	// TODO implement me
	panic("implement me")
}

func (mock *MockFileService) Create(file *entity.File) error {
	args := mock.Called(file.Filename)

	var err error
	if args.Get(0) != nil {
		err = args.Get(0).(error)
	}
	return err
}

func (mock *MockFileService) GetFromContext(contextFileID string) (*entity.File, error) {
	args := mock.Called(contextFileID)

	var file *entity.File
	if args.Get(0) != nil {
		file = args.Get(0).(*entity.File)
	}

	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return file, err
}

func (mock *MockFileService) CheckExisting(filename string) (bool, error) {
	args := mock.Called(filename)

	var isFileExist bool
	if args.Get(0) != nil {
		isFileExist = args.Get(0).(bool)
	}

	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return isFileExist, err
}
