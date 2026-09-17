package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"lab-4/backend/internal/repository"
	"lab-4/backend/internal/service"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, err error) {
	status, code, message := http.StatusInternalServerError, "internal_error", "Internal server error"
	switch {
	case errors.Is(err, service.ErrValidation):
		status, code, message = http.StatusBadRequest, "validation_error", "Invalid request data"
	case errors.Is(err, repository.ErrNotFound):
		status, code, message = http.StatusNotFound, "not_found", "Resource not found"
	case errors.Is(err, repository.ErrConflict):
		status, code, message = http.StatusConflict, "conflict", "Resource already exists"
	}
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
