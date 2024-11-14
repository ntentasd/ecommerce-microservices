package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Envelope map[string]any

func FormatResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ParseInt(path string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(path, "/products/"))
}
