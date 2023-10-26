package route

import (
	"github.com/gin-gonic/gin"
	"goback/api/controller"
	"goback/mapper"
	"goback/services"
)

func InitFileRoute(router *gin.RouterGroup, userService services.IUserService, fileHistoryService services.IFileHistoryService, fileService services.IFileService, fileStorageService services.IFileStorageService, workspaceService services.IWorkspaceService, fileMapper mapper.IFileMapper) {
	fileController := controller.NewFileController(userService, fileHistoryService, fileService, fileStorageService, workspaceService, fileMapper)
	router.POST("/create", fileController.Create)
	router.GET("/check-filename-existing", fileController.CheckFilenameExisting)
	router.GET("/:file_id/get-file-download-link", fileController.GetFileDownloadLink)
}
