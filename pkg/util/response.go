package util

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func JsonResponseWriter(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func JsonListResponseWriter(w http.ResponseWriter, code int, payload interface{}, count int64) {
	response := struct {
		Items interface{} `json:"items"`
		Count int64       `json:"count"`
	}{payload, count}

	JsonResponseWriter(w, code, response)
}

func ErrorResponseWriter(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	JsonResponseWriter(w, code, errorResponse{msg})
}

func GetPaginationParams(r *http.Request) (int64, int64, error) {
	var l, o int64 = 10, 0

	if param, ok := r.URL.Query()["limit"]; ok && len(param[0]) > 0 {
		limitInt, err := strconv.ParseInt(param[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		l = limitInt
	}

	if param, ok := r.URL.Query()["offset"]; ok && len(param[0]) > 0 {
		offsetInt, err := strconv.ParseInt(param[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		o = offsetInt
	}
	return l, o, nil
}
