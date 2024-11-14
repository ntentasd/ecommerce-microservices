package utils

import (
	"encoding/json"
	"net/http"
	"slices"
)

func MethodHandler(next http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(methods, r.Method) {
			next(w, r)
			return
		}

		FormatResponse(w, http.StatusMethodNotAllowed, Envelope{"error": "Method not allowed"})
	}
}

func ReadJSON(w http.ResponseWriter, r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
