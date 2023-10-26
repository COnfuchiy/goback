package request

type CheckFileRequest struct {
	Filename string `form:"filename" binding:"required"`
}
