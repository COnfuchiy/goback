package request

type CreateFileRequest struct {
	Filename string `form:"filename" binding:"required"`
	Size     int64  `form:"size" binding:"required"`
}
