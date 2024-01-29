package response

type TaggingResponse struct {
	FileTag string `json:"file tag"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
