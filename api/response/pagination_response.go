package response

type PaginationResponse struct {
	TotalEntities int64  `json:"total_entities"`
	TotalPages    int64  `json:"total_pages"`
	CurrentPage   int64  `json:"current_page"`
	NextPage      *int64 `json:"next_page"`
	PrevPage      *int64 `json:"prev_page"`
}
