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
	"strconv"
)

type WorkspaceController struct {
	workspaceService   services.IWorkspaceService
	fileHistoryService services.IFileHistoryService
	workspaceMapper    mapper.IWorkspaceMapper
	fileHistoryMapper  mapper.IFileHistoryMapper
	paginateMapper     mapper.IPaginationMapper
}

func NewWorkspaceController(workspaceService services.IWorkspaceService, fileHistoryService services.IFileHistoryService, workspaceMapper mapper.IWorkspaceMapper, fileHistoryMapper mapper.IFileHistoryMapper, paginateMapper mapper.IPaginationMapper) *WorkspaceController {
	return &WorkspaceController{workspaceService, fileHistoryService, workspaceMapper, fileHistoryMapper, paginateMapper}
}

// CreateWorkspace godoc
// @Summary	create workspace
// @Description	create workspace
// @Param name formData string true "workspace name"
// @Tags workspace
// @Produce json
// @Success 200 {object} response.WorkspaceResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /workspace/create [post]
// @Security Bearer
func (c WorkspaceController) CreateWorkspace(context *gin.Context) {
	var req request.CreateWorkspaceRequest

	err := context.ShouldBind(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	userObject, isUserExist := context.Get("user")
	if !isUserExist {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User not exist"})
		return
	}

	user, _ := userObject.(*entity.User)
	if reflect.ValueOf(user).IsNil() {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User is not type of " + reflect.TypeOf(entity.User{}).String()})
		return
	}

	newWorkspace := c.workspaceMapper.FromCreateRequest(&req)
	newWorkspace.CreatorID = user.ID
	newWorkspace.Creator = *user

	err = c.workspaceService.Create(newWorkspace)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, c.workspaceMapper.ToWorkspaceResponse(newWorkspace))
}

// GetWorkspace godoc
// @Summary	get workspace
// @Description	get workspace
// @Tags workspace
// @Param workspace_id path string true "workspace id"
// @Produce json
// @Success 200 {object} response.WorkspaceResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /workspace/{workspace_id} [get]
// @Security Bearer
func (c WorkspaceController) GetWorkspace(context *gin.Context) {

	workspaceObject, isWorkspaceExist := context.Get("workspace")
	if !isWorkspaceExist {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Workspace not exist"})
		return
	}

	workspace, _ := workspaceObject.(*entity.Workspace)
	if reflect.ValueOf(workspace).IsNil() {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Workspace is not type of " + reflect.TypeOf(entity.Workspace{}).String()})
		return
	}

	context.JSON(http.StatusOK, c.workspaceMapper.ToWorkspaceResponse(workspace))
}

// GetAllFilesHistories godoc
// @Summary	get all files histories
// @Description	get all files histories
// @Tags workspace
// @Param workspace_id path string true "workspace id"
// @Produce json
// @Success 200 {object} response.FileHistoriesResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /workspace/{workspace_id}/file-histories [get]
// @Security Bearer
func (c WorkspaceController) GetAllFilesHistories(context *gin.Context) {

	workspaceObject, isWorkspaceExist := context.Get("workspace")
	if !isWorkspaceExist {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Workspace ID is not specified"})
		return
	}

	workspace, _ := workspaceObject.(*entity.Workspace)
	if reflect.ValueOf(workspace).IsNil() {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Workspace is not type of " + reflect.TypeOf(entity.Workspace{}).String()})
		return
	}

	currentPage := 1
	currentPageAsString := context.Param("page")

	if currentPageAsString != "" {
		var castError error
		currentPage, castError = strconv.Atoi(currentPageAsString)
		if castError != nil {
			currentPage = 1
		}
	}

	fileHistories, totalCount, err := c.fileHistoryService.GetAllByWorkspaceID(workspace.ID, currentPage)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	var fileHistoriesResponse response.FileHistoriesResponse

	for _, history := range fileHistories {
		fileHistoriesResponse.FileHistories = append(fileHistoriesResponse.FileHistories, *c.fileHistoryMapper.ToFileHistoryResponse(&history))
	}

	fileHistoriesResponse.Pagination = *c.paginateMapper.ToPaginationResponse(totalCount, currentPage)

	context.JSON(http.StatusOK, fileHistoriesResponse)
}
