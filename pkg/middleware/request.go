package middleware

import (
	"context"
	"net/http"
	"strconv"
)

// Paginator contains limit and offset values extracted from query parameters
type Paginator struct {
	Limit  int64
	Offset int64
}

func Pagination(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10 // Default limit if invalid or missing
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0 // Default offset if invalid or missing
		}

		ctx := context.WithValue(r.Context(), "pagination", &Paginator{
			Limit:  int64(limit),
			Offset: int64(offset),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
