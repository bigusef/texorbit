package util

import (
	"encoding/json"
	"net/http"
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
		Result interface{} `json:"result"`
		Count  int64       `json:"count"`
	}{payload, count}

	JsonResponseWriter(w, code, response)
}
