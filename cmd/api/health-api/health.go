package health_api

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the structure of the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HealthHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Create a response object
	response := HealthResponse{
		Status:  "OK",
		Message: "The server is up",
	}

	// Set headers and send JSON response
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	err := json.NewEncoder(responseWriter).Encode(response)
	if err != nil {
		return
	}
}
