package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// GetFileHistory godoc
// @Summary	get file history
// @Description	get all file versions
// @Tags file
// @Param workspace_id path string true "workspace id"
// @Param file_history_id path string true "workspace id"
// @Produce json
// @Success 200 {object} response.FileHistoryResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /workspace/{workspace_id}/file-history/{file_history_id} [get]
// @Security Bearer
func (c FileHistoryController) GetFileHistory(context *gin.Context) {

	fileHistoryID := context.Param("file_history_id")
	if fileHistoryID != "" {
		if _, err := uuid.Parse(fileHistoryID); err != nil {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "File history ID is not uuid"})
			return
		}
	} else {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "File history ID is not specified"})
		return
	}

	fileHistory, err := c.fileHistoryService.GetFromContext(fileHistoryID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, c.fileHistoryMapper.ToFileHistoryResponse(fileHistory))
}
