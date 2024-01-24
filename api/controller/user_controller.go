package controller

import (
	"github.com/gin-gonic/gin"
	"goback/api/response"
	"goback/domain/entity"
	"goback/mapper"
	"goback/services"
	"net/http"
	"reflect"
	"strconv"
)

type UserController struct {
	userService      services.IUserService
	workspaceService services.IWorkspaceService
	userMapper       mapper.IUserMapper
	workspaceMapper  mapper.IWorkspaceMapper
	paginateMapper   mapper.IPaginationMapper
}

func NewUserController(userService services.IUserService, workspaceService services.IWorkspaceService, userMapper mapper.IUserMapper, workspaceMapper mapper.IWorkspaceMapper, paginateMapper mapper.IPaginationMapper) *UserController {
	return &UserController{userService, workspaceService, userMapper, workspaceMapper, paginateMapper}
}

// Profile godoc
// @Summary	get user profile
// @Description	get user profile
// @Tags user
// @Produce json
// @Success 200 {object} response.ProfileResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /user/profile [get]
func (c UserController) Profile(context *gin.Context) {

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

	context.JSON(http.StatusOK, c.userMapper.ToProfileResponse(user))
}

// GetAllWorkspaces godoc
// @Summary	get all user workspace (created and invited)
// @Description	get all user workspace (created and invited)
// @Param page query int false "page number"
// @Tags user
// @Produce json
// @Success 200 {object} response.WorkspacesResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /user/get-all-workspaces [get]
func (c UserController) GetAllWorkspaces(context *gin.Context) {

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

	currentPage := 1
	currentPageAsString := context.Query("page")
	if currentPageAsString != "" {
		var castError error
		currentPage, castError = strconv.Atoi(currentPageAsString)
		if castError != nil {
			currentPage = 1
		}
	}

	workspaces, totalCount, err := c.workspaceService.GetAllByUserID(user.ID, currentPage)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	var workspacesResponse response.WorkspacesResponse

	for _, workspace := range workspaces {
		workspacesResponse.Workspaces = append(workspacesResponse.Workspaces, *c.workspaceMapper.ToWorkspaceResponse(&workspace))
	}

	workspacesResponse.Pagination = *c.paginateMapper.ToPaginationResponse(totalCount, currentPage)

	context.JSON(http.StatusOK, workspacesResponse)
}
