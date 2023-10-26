package mapper

import (
	"goback/api/response"
	"goback/domain/entity"
)

type IFileHistoryMapper interface {
	ToFileHistoryResponse(fileHistory *entity.FileHistory) *response.FileHistoryResponse
}

type FileHistoryMapper struct {
	fileMapper IFileMapper
}

func NewFileHistoryMapper(fileMapper IFileMapper) IFileHistoryMapper {
	return &FileHistoryMapper{fileMapper: fileMapper}
}

func (m FileHistoryMapper) ToFileHistoryResponse(fileHistory *entity.FileHistory) *response.FileHistoryResponse {
	var filesResponses []response.FileResponse
	for _, file := range fileHistory.Files {
		filesResponses = append(filesResponses, *m.fileMapper.ToFileResponse(&file))
	}
	return &response.FileHistoryResponse{
		ID:    fileHistory.ID.String(),
		Files: filesResponses,
	}
}
