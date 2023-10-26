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

func (c UserController) Profile(context *gin.Context) {

	userObject, isUserExist := context.Get("user")
	if !isUserExist {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User not exist"})
		return
	}

	user, isUser := userObject.(*entity.User)
	if !isUser {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User is not type of " + reflect.TypeOf(entity.User{}).String()})
		return
	}

	context.JSON(http.StatusOK, c.userMapper.ToProfileResponse(user))
}

func (c UserController) GetAllWorkspaces(context *gin.Context) {

	userObject, isUserExist := context.Get("user")
	if !isUserExist {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User not exist"})
		return
	}

	user, isUser := userObject.(*entity.User)
	if !isUser {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User is not type of " + reflect.TypeOf(entity.User{}).String()})
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

	workspaces, totalCount, err := c.workspaceService.GetAllByUserID(user.ID, currentPage)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	var workspacesResponce response.WorkspacesResponse

	for _, workspace := range workspaces {
		workspacesResponce.Workspaces = append(workspacesResponce.Workspaces, *c.workspaceMapper.ToWorkspaceResponse(&workspace))
	}

	workspacesResponce.Pagination = *c.paginateMapper.ToPaginationResponse(totalCount, currentPage)

	context.JSON(http.StatusOK, workspacesResponce)
}
