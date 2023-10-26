package response

type WorkspacesResponse struct {
	Workspaces []WorkspaceResponse `json:"data"`
	Pagination PaginationResponse  `json:"pagination"`
}
