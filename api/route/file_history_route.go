package route

import (
	"github.com/gin-gonic/gin"
	"goback/api/controller"
	"goback/mapper"
	"goback/services"
)

func InitFileHistoryRoute(router *gin.RouterGroup, fileHistoryService services.IFileHistoryService, fileHistoryMapper mapper.IFileHistoryMapper) {
	fileHistoryController := controller.NewFileHistoryController(fileHistoryService, fileHistoryMapper)
	router.GET("/:file_history_id", fileHistoryController.GetFileHistory)
}
