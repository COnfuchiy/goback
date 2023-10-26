package controller

import (
	"github.com/gin-gonic/gin"
	"goback/api/response"
	"goback/mapper"
	"goback/services"
	"net/http"
)

type FileHistoryController struct {
	fileHistoryService services.IFileHistoryService
	fileHistoryMapper  mapper.IFileHistoryMapper
}

func NewFileHistoryController(fileHistoryService services.IFileHistoryService, fileHistoryMapper mapper.IFileHistoryMapper) *FileHistoryController {
	return &FileHistoryController{fileHistoryService, fileHistoryMapper}
}

func (c FileHistoryController) GetFileHistory(context *gin.Context) {

	fileHistoryID := context.Param("file_history_id")
	if fileHistoryID == "" {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "File ID is not specified"})
		return
	}

	fileHistory, err := c.fileHistoryService.GetFileHistoryFromContext(fileHistoryID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, c.fileHistoryMapper.ToFileHistoryResponse(fileHistory))
}
