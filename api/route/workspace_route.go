package route

import (
	"github.com/gin-gonic/gin"
	"goback/api/controller"
	"goback/mapper"
	"goback/services"
)

func InitWorkspaceRoute(router *gin.RouterGroup, workspaceService services.IWorkspaceService, fileHistoryService services.IFileHistoryService, workspaceMapper mapper.IWorkspaceMapper, fileHistoryMapper mapper.IFileHistoryMapper, paginateMapper mapper.IPaginationMapper) {
	workspaceController := controller.NewWorkspaceController(workspaceService, fileHistoryService, workspaceMapper, fileHistoryMapper, paginateMapper)
	router.POST("/create", workspaceController.CreateWorkspace)
	router.GET("/:workspace_id", workspaceController.GetWorkspace)
	router.GET("/:workspace_id/file-histories", workspaceController.GetAllFilesHistories)
}
