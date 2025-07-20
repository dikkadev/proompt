package models

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standardized API error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Details map[string]string `json:"details,omitempty"`
}

// ValidationErrorResponse represents validation errors
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Fields  map[string]string `json:"fields"`
}

// WriteError writes a standardized error response
func WriteError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Code:    code,
	}

	json.NewEncoder(w).Encode(response)
}

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, fields map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := ValidationErrorResponse{
		Error:   "Validation Failed",
		Message: "The request contains invalid data",
		Code:    http.StatusBadRequest,
		Fields:  fields,
	}

	json.NewEncoder(w).Encode(response)
}

// WriteNotFound writes a 404 error response
func WriteNotFound(w http.ResponseWriter, resource string) {
	WriteError(w, http.StatusNotFound, resource+" not found")
}

// WriteInternalError writes a 500 error response
func WriteInternalError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusInternalServerError, message)
}

// WriteBadRequest writes a 400 error response
func WriteBadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, message)
}
