package controller

import (
	"github.com/gin-gonic/gin"
	"goback/api/request"
	"goback/api/response"
	"goback/domain/entity"
	"goback/mapper"
	"goback/services"
	"net/http"
	"reflect"
)

type FileController struct {
	userService        services.IUserService
	fileHistoryService services.IFileHistoryService
	fileService        services.IFileService
	fileStorageService services.IFileStorageService
	workspaceService   services.IWorkspaceService
	fileMapper         mapper.IFileMapper
}

func NewFileController(userService services.IUserService, fileHistoryService services.IFileHistoryService, fileService services.IFileService, fileStorageService services.IFileStorageService, workspaceService services.IWorkspaceService, fileMapper mapper.IFileMapper) *FileController {
	return &FileController{userService, fileHistoryService, fileService, fileStorageService, workspaceService, fileMapper}
}

func (c FileController) Create(context *gin.Context) {
	var req request.CreateFileRequest

	err := context.ShouldBind(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	userObject, isUserExist := context.Get("user")
	if !isUserExist {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User not exist"})
		context.Abort()
	}

	user, isUser := userObject.(*entity.User)
	if !isUser {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User is not type of " + reflect.TypeOf(entity.User{}).String()})
		context.Abort()
	}

	workspaceObject, isWorkspaceExist := context.Get("workspace")
	if !isWorkspaceExist {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Workspace not exist"})
		return
	}

	workspace, isWorkspace := workspaceObject.(*entity.Workspace)
	if !isWorkspace {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Workspace is not type of " + reflect.TypeOf(entity.Workspace{}).String()})
		return
	}

	fileObject, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	newFileHistory := entity.FileHistory{
		WorkspaceID: workspace.ID,
	}

	err = c.fileHistoryService.Create(&newFileHistory)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	file := c.fileMapper.FromCreateFileResponce(req)
	file.FileHistoryID = newFileHistory.ID
	file.FileHistory = newFileHistory
	file.UserId = user.ID
	file.User = *user

	err = c.fileService.Create(&file)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	go c.fileStorageService.SaveFileToStorage(file.ID, fileObject)

	context.JSON(http.StatusOK, c.fileMapper.ToFileResponse(&file))
}

func (c FileController) GetFileDownloadLink(context *gin.Context) {

	fileID := context.Param("file_id")
	if fileID == "" {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "File ID is not specified"})
		return
	}

	file, err := c.fileService.GetFromContext(fileID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, c.fileMapper.ToDownloadFileLinkResponse(file))
}

func (c FileController) CheckFilenameExisting(context *gin.Context) {
	var req request.CheckFileRequest

	err := context.ShouldBind(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	isFileExist, err := c.fileService.CheckExisting(req.Filename)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response.CheckFileResponse{IsFileExist: isFileExist})
}
