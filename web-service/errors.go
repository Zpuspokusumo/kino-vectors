package main

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

func handleError(w http.ResponseWriter, r *http.Request, statusCode int, message string, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResponse := ErrorResponse{
		Message: message,
		Code:    statusCode, // Using HTTP status code as the error code
		Details: details,
	}

	json.NewEncoder(w).Encode(errResponse)
}
