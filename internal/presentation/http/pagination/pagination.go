package pagination

import (
	"net/http"
	"strconv"
)

// Params holds pagination parameters
type Params struct {
	Limit  int
	Offset int
}

// Info holds pagination metadata for responses
type Info struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
	Page   int `json:"page"`
}

// DefaultLimit is the default number of items per page
const DefaultLimit = 20

// MaxLimit is the maximum number of items per page
const MaxLimit = 100

// Parse extracts pagination parameters from HTTP request query string
// Defaults: limit=20, offset=0
func Parse(r *http.Request) Params {
	query := r.URL.Query()

	limit := DefaultLimit
	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			if parsed > MaxLimit {
				limit = MaxLimit
			} else {
				limit = parsed
			}
		}
	}

	offset := 0
	if o := query.Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Support page parameter as alternative to offset
	if p := query.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			offset = (parsed - 1) * limit
		}
	}

	return Params{Limit: limit, Offset: offset}
}

// NewInfo creates pagination info for response
func NewInfo(limit, offset, total int) Info {
	if limit <= 0 {
		limit = DefaultLimit
	}
	page := 1
	if offset > 0 {
		page = (offset / limit) + 1
	}
	return Info{
		Limit:  limit,
		Offset: offset,
		Total:  total,
		Page:   page,
	}
}
