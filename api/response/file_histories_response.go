package response

type FileHistoriesResponse struct {
	FileHistories []FileHistoryResponse `json:"file_histories"`
	Pagination    PaginationResponse    `json:"pagination"`
}
