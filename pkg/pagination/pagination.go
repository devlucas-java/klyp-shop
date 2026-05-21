package pagination

import (
	"net/http"
	"strconv"
)

func ParsePagination(r *http.Request) InputPagination {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	search := r.URL.Query().Get("search")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	return InputPagination{
		Page:   page,
		Size:   size,
		Search: search,
	}
}
