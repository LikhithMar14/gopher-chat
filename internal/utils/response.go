package utils

import (
	"encoding/json"
	"net/http"
)

type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Status  int         `json:"status"`
}

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

func HandleValidationError(w http.ResponseWriter, err error) error {
	return WriteErrorResponse(w, http.StatusBadRequest, err.Error())
}


func HandleInternalError(w http.ResponseWriter, err error) error {
	return WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
}
