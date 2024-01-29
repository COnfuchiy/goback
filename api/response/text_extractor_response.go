package response

type TextExtractorResponse struct {
	FileText string `json:"file text"`
	Success  bool   `json:"success"`
	Error    string `json:"error"`
}
