package helpers

import (
	"net/http"
	"strconv"
)

type paginationParams struct {
	Page  int
	Limit int
}

func ParsePaginationParams(r *http.Request) paginationParams {
	queryParams := r.URL.Query()

	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 1 {
		page = 1 // default
	}

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit < 1 {
		limit = 10 // default
	}

	return paginationParams{
		Page:  page,
		Limit: limit,
	}
}
