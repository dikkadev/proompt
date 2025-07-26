package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dikkadev/proompt/server/internal/api/models"
)

// Health godoc
// @Summary Health check endpoint
// @Description Returns the health status of the API server
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse "Health status"
// @Router /health [get]
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
