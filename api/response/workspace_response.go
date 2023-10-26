package response

type WorkspaceResponse struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Creator ProfileResponse   `json:"creator"`
	Users   []ProfileResponse `json:"users"`
	// FilesHistories []FileHistoryResponse `json:"files"`
}
