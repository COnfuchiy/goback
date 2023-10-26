package response

type FileResponse struct {
	ID          string          `json:"id"`
	Filename    string          `json:"filename"`
	Tag         string          `json:"tag"`
	Size        int64           `json:"size"`
	DownloadURL string          `json:"download_url"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	User        ProfileResponse `json:"user"`
}
