package pagination

type OutPutPagination struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type InputPagination struct {
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Search string `json:"search"`
}
