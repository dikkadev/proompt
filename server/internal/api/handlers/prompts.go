package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dikkadev/proompt/server/internal/api/models"
	domainModels "github.com/dikkadev/proompt/server/internal/models"
	"github.com/dikkadev/proompt/server/internal/repository"
	"github.com/google/uuid"
)

// PromptHandlers contains handlers for prompt operations
type PromptHandlers struct {
	repo repository.Repository
}

// NewPromptHandlers creates a new prompt handlers instance
func NewPromptHandlers(repo repository.Repository) *PromptHandlers {
	return &PromptHandlers{repo: repo}
}

// CreatePrompt handles POST /api/prompts
func (h *PromptHandlers) CreatePrompt(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePromptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	// TODO: Add validation middleware
	// For now, basic validation
	if req.Title == "" {
		models.WriteBadRequest(w, "Title is required")
		return
	}
	if req.Content == "" {
		models.WriteBadRequest(w, "Content is required")
		return
	}

	// Convert to domain model
	prompt := req.ToPrompt()
	prompt.ID = uuid.New().String()

	// Create prompt
	if err := h.repo.Prompts().Create(r.Context(), prompt); err != nil {
		models.WriteInternalError(w, "Failed to create prompt")
		return
	}

	// Return created prompt
	response := models.FromPrompt(prompt)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetPrompt handles GET /api/prompts/{id}
func (h *PromptHandlers) GetPrompt(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	prompt, err := h.repo.Prompts().GetByID(r.Context(), id)
	if err != nil {
		models.WriteNotFound(w, "Prompt")
		return
	}

	response := models.FromPrompt(prompt)
	json.NewEncoder(w).Encode(response)
}

// UpdatePrompt handles PUT /api/prompts/{id}
func (h *PromptHandlers) UpdatePrompt(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	// Get existing prompt
	existing, err := h.repo.Prompts().GetByID(r.Context(), id)
	if err != nil {
		models.WriteNotFound(w, "Prompt")
		return
	}

	var req models.UpdatePromptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	// Apply updates
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Content != nil {
		existing.Content = *req.Content
	}
	if req.Type != nil {
		existing.Type = domainModels.PromptType(*req.Type)
	}
	if req.UseCase != nil {
		existing.UseCase = req.UseCase
	}
	if req.ModelCompatibilityTags != nil {
		existing.ModelCompatibilityTags = domainModels.StringSlice(req.ModelCompatibilityTags)
	}
	if req.TemperatureSuggestion != nil {
		existing.TemperatureSuggestion = req.TemperatureSuggestion
	}
	if req.OtherParameters != nil {
		existing.OtherParameters = domainModels.JSONMap(req.OtherParameters)
	}

	// Update prompt
	if err := h.repo.Prompts().Update(r.Context(), existing); err != nil {
		models.WriteInternalError(w, "Failed to update prompt")
		return
	}

	// Return updated prompt
	response := models.FromPrompt(existing)
	json.NewEncoder(w).Encode(response)
}

// DeletePrompt handles DELETE /api/prompts/{id}
func (h *PromptHandlers) DeletePrompt(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	if err := h.repo.Prompts().Delete(r.Context(), id); err != nil {
		models.WriteInternalError(w, "Failed to delete prompt")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListPrompts handles GET /api/prompts
func (h *PromptHandlers) ListPrompts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filters := repository.PromptFilters{}

	if typeParam := r.URL.Query().Get("type"); typeParam != "" {
		filters.Type = &typeParam
	}
	if useCaseParam := r.URL.Query().Get("use_case"); useCaseParam != "" {
		filters.UseCase = &useCaseParam
	}
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if limit, err := strconv.Atoi(limitParam); err == nil {
			filters.Limit = &limit
		}
	}
	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if offset, err := strconv.Atoi(offsetParam); err == nil {
			filters.Offset = &offset
		}
	}

	prompts, err := h.repo.Prompts().List(r.Context(), filters)
	if err != nil {
		models.WriteInternalError(w, "Failed to list prompts")
		return
	}

	responses := models.FromPrompts(prompts)

	// Create list response
	listResponse := models.ListResponse[*models.PromptResponse]{
		Data:       responses,
		Total:      len(responses), // TODO: Get actual total count
		Page:       1,              // TODO: Calculate from offset/limit
		PageSize:   len(responses),
		TotalPages: 1, // TODO: Calculate from total/page_size
	}

	json.NewEncoder(w).Encode(listResponse)
}
