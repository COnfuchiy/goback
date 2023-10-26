package middleware

import (
	"github.com/gin-gonic/gin"
	"goback/api/response"
	"goback/domain/entity"
	"goback/services"
	"net/http"
	"reflect"
	"regexp"
)

type WorkSpaceMiddleware struct {
	workspaceService services.IWorkspaceService
}

func NewWorkSpaceMiddleware(workspaceService services.IWorkspaceService) *WorkSpaceMiddleware {
	return &WorkSpaceMiddleware{workspaceService: workspaceService}
}

func (m WorkSpaceMiddleware) Handle() gin.HandlerFunc {
	return func(context *gin.Context) {
		isCreateMethod, _ := regexp.MatchString("/workspace/create", context.Request.URL.Path)

		if isCreateMethod {
			context.Next()
			return
		}

		workspaceID := context.Param("workspace_id")
		if workspaceID == "" {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Workspace ID is required"})
			context.Abort()
			return
		}

		workspace, err := m.workspaceService.GetWorkspaceFromContext(workspaceID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
			context.Abort()
			return
		}

		userObject, isUserExist := context.Get("user")
		if !isUserExist {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User not exist"})
			context.Abort()
			return
		}

		user, isUser := userObject.(*entity.User)
		if !isUser {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User is not type of " + reflect.TypeOf(entity.User{}).String()})
			context.Abort()
			return
		}

		isUserAccess, err := m.workspaceService.CheckUserAccess(workspace, user.ID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
			context.Abort()
			return
		}

		if !isUserAccess {
			context.JSON(http.StatusForbidden, response.ErrorResponse{Message: err.Error()})
			context.Abort()
			return
		}
		context.Set("workspace", workspace)
		context.Next()
		return
	}
}
