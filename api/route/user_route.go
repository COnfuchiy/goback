package route

import (
	"github.com/gin-gonic/gin"
	"goback/api/controller"
	"goback/mapper"
	"goback/services"
)

func InitUserRoute(router *gin.RouterGroup, userService services.IUserService, workspaceService services.IWorkspaceService, userMapper mapper.IUserMapper, workspaceMapper mapper.IWorkspaceMapper, paginateMapper mapper.IPaginationMapper) {
	userController := controller.NewUserController(userService, workspaceService, userMapper, workspaceMapper, paginateMapper)

	router.GET("/profile", userController.Profile)
	router.GET("/workspaces", userController.GetAllWorkspaces)
}
