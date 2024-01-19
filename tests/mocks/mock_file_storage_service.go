package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

type MockFileStorageService struct {
	mock.Mock
}

func (mock *MockFileStorageService) SaveFileToStorage(fileID uuid.UUID, fileObject *multipart.FileHeader) {
	_ = mock.Called(fileID, fileObject.Filename)
}
