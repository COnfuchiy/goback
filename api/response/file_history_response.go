package response

type FileHistoryResponse struct {
	ID    string         `json:"id"`
	Files []FileResponse `json:"files"`
}
