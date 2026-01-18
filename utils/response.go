package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse for consistent error handling
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondJSON writes a JSON response with the given data and status code
func RespondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// RespondError writes an error response with the given message and status code
func RespondError(w http.ResponseWriter, message string, statusCode int) {
	RespondJSON(w, ErrorResponse{Error: message}, statusCode)
}
