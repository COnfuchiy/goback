package request

import "mime/multipart"

type CreateFileRequest struct {
	Filename string                `form:"filename" binding:"required"`
	Size     int64                 `form:"size" binding:"required"`
	File     *multipart.FileHeader `form:"file" binding:"required"`
}
