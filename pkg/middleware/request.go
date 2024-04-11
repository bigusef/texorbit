package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

var limitKey contextKey = "limit"
var offsetKey contextKey = "offset"

func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, ok := r.URL.Query()["limit"]
		if !ok || len(limit[0]) < 1 {
			limit = []string{"10"}
		}
		offset, ok := r.URL.Query()["offset"]
		if !ok || len(offset[0]) < 1 {
			offset = []string{"0"}
		}

		limitInt, err := strconv.Atoi(limit[0])
		if err != nil {
			http.Error(w, "parameter 'limit' must be a number", http.StatusBadRequest)
			return
		}

		offsetInt, err := strconv.Atoi(offset[0])
		if err != nil {
			http.Error(w, "parameter 'offset' must be a number", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), limitKey, limitInt)
		ctx = context.WithValue(ctx, offsetKey, offsetInt)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
