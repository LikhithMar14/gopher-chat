package utils

import (
	"encoding/json"
	"net/http"
)

// StandardResponse represents the standard API response format
type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Status  int         `json:"status"`
}

// WriteSuccessResponse writes a standardized success response
func WriteSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	response := StandardResponse{
		Success: true,
		Data:    data,
		Status:  statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(response)
}

// WriteErrorResponse writes a standardized error response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) error {
	response := StandardResponse{
		Success: false,
		Error:   message,
		Status:  statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(response)
}

// HandleValidationError handles validation errors consistently
func HandleValidationError(w http.ResponseWriter, err error) error {
	return WriteErrorResponse(w, http.StatusBadRequest, err.Error())
}

// HandleInternalError handles internal server errors consistently
func HandleInternalError(w http.ResponseWriter, err error) error {
	return WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
}
