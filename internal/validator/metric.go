package validator

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ValidateHandleNewRequest(rw http.ResponseWriter, r *http.Request) bool {
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if value == "" {
		http.Error(rw, "No metric value specified", http.StatusBadRequest)
		return false
	}

	if name == "" {
		http.Error(rw, "No metric name specified", http.StatusNotFound)
		return false
	}
	return true
}

func ValidateHandleItemRequest(rw http.ResponseWriter, r *http.Request) bool {
	name := chi.URLParam(r, "name")

	if name == "" {
		http.Error(rw, "No metric name specified", http.StatusNotFound)
		return false
	}
	return true
}
