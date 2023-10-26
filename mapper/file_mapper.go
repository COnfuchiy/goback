package mapper

import (
	"goback/api/request"
	"goback/api/response"
	"goback/domain/entity"
)

type IFileMapper interface {
	ToFileResponse(file *entity.File) *response.FileResponse
	FromCreateFileResponce(fileRequest request.CreateFileRequest) entity.File
	ToDownloadFileLinkResponse(file *entity.File) *response.DownloadFileLinkResponse
}

type FileMapper struct {
	userMapper IUserMapper
}

func NewFileMapper(userMapper IUserMapper) IFileMapper {
	return &FileMapper{userMapper: userMapper}
}

func (m FileMapper) ToFileResponse(file *entity.File) *response.FileResponse {
	return &response.FileResponse{
		Filename:    file.Filename,
		Tag:         file.Tag,
		Size:        file.Size,
		DownloadURL: file.DownloadURL,
		CreatedAt:   file.CreatedAt.Format("YYYY-MM-DD hh:mm:ss"),
		UpdatedAt:   file.UpdatedAt.Format("YYYY-MM-DD hh:mm:ss"),
		User:        *m.userMapper.ToProfileResponse(&file.User),
	}
}

func (m FileMapper) ToDownloadFileLinkResponse(file *entity.File) *response.DownloadFileLinkResponse {
	return &response.DownloadFileLinkResponse{DownloadLink: file.DownloadURL}
}

func (m FileMapper) FromCreateFileResponce(fileRequest request.CreateFileRequest) entity.File {
	return entity.File{
		Filename: fileRequest.Filename,
		Size:     fileRequest.Size,
	}
}
