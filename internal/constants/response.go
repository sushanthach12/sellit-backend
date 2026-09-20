package constants

// Pagination holds paging metadata
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// Data wraps results + pagination for any result type T
type Data[T any] struct {
	Results    []T        `json:"results"`
	Pagination Pagination `json:"pagination"`
}

// Response is the outer envelope for any result type T
type Response[T any] struct {
	Data     Data[T] `json:"data"`
	Metadata any     `json:"metadata,omitempty"`
}

// New builds a Response, defaulting nil results to an empty slice
// so the JSON output is `[]` instead of `null`.
func NewPaginatedResponse[T any](results []T, pagination Pagination) Response[T] {
	if results == nil {
		results = []T{}
	}
	return Response[T]{
		Data: Data[T]{
			Results:    results,
			Pagination: pagination,
		},
	}
}

// ---- New, separate type — do NOT reuse Response[T] here ----

type SimpleResponse[T any] struct {
	Data     T   `json:"data"`
	Metadata any `json:"metadata,omitempty"`
}

func NewResponse[T any](data T, metadata any) SimpleResponse[T] {
	return SimpleResponse[T]{
		Data:     data, // ✅ T into T, fine
		Metadata: metadata,
	}
}
