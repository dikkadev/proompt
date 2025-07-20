package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dikkadev/proompt/server/internal/api/models"
)

// Health handles the health check endpoint
func Health(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0", // TODO: Get from build info
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
