package request

type FileHistoryRequest struct {
	FileHistoryID string `form:"file_history_id" binding:"required"`
}
