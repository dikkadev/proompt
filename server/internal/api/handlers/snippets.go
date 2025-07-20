package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dikkadev/proompt/server/internal/api/models"
	"github.com/dikkadev/proompt/server/internal/repository"
	"github.com/google/uuid"
)

// SnippetHandlers contains handlers for snippet operations
type SnippetHandlers struct {
	repo repository.Repository
}

// NewSnippetHandlers creates a new snippet handlers instance
func NewSnippetHandlers(repo repository.Repository) *SnippetHandlers {
	return &SnippetHandlers{repo: repo}
}

// CreateSnippet handles POST /api/snippets
func (h *SnippetHandlers) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	// Basic validation
	if req.Title == "" {
		models.WriteBadRequest(w, "Title is required")
		return
	}
	if req.Content == "" {
		models.WriteBadRequest(w, "Content is required")
		return
	}

	// Convert to domain model
	snippet := req.ToSnippet()
	snippet.ID = uuid.New().String()

	// Create snippet
	if err := h.repo.Snippets().Create(r.Context(), snippet); err != nil {
		models.WriteInternalError(w, "Failed to create snippet")
		return
	}

	// Return created snippet
	response := models.FromSnippet(snippet)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetSnippet handles GET /api/snippets/{id}
func (h *SnippetHandlers) GetSnippet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		models.WriteBadRequest(w, "Snippet ID is required")
		return
	}

	snippet, err := h.repo.Snippets().GetByID(r.Context(), id)
	if err != nil {
		models.WriteNotFound(w, "Snippet")
		return
	}

	response := models.FromSnippet(snippet)
	json.NewEncoder(w).Encode(response)
}

// UpdateSnippet handles PUT /api/snippets/{id}
func (h *SnippetHandlers) UpdateSnippet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		models.WriteBadRequest(w, "Snippet ID is required")
		return
	}

	// Get existing snippet
	existing, err := h.repo.Snippets().GetByID(r.Context(), id)
	if err != nil {
		models.WriteNotFound(w, "Snippet")
		return
	}

	var req models.UpdateSnippetRequest
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
	if req.Description != nil {
		existing.Description = req.Description
	}

	// Update snippet
	if err := h.repo.Snippets().Update(r.Context(), existing); err != nil {
		models.WriteInternalError(w, "Failed to update snippet")
		return
	}

	// Return updated snippet
	response := models.FromSnippet(existing)
	json.NewEncoder(w).Encode(response)
}

// DeleteSnippet handles DELETE /api/snippets/{id}
func (h *SnippetHandlers) DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		models.WriteBadRequest(w, "Snippet ID is required")
		return
	}

	if err := h.repo.Snippets().Delete(r.Context(), id); err != nil {
		models.WriteInternalError(w, "Failed to delete snippet")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListSnippets handles GET /api/snippets
func (h *SnippetHandlers) ListSnippets(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filters := repository.SnippetFilters{}

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

	snippets, err := h.repo.Snippets().List(r.Context(), filters)
	if err != nil {
		models.WriteInternalError(w, "Failed to list snippets")
		return
	}

	responses := models.FromSnippets(snippets)

	// Create list response
	listResponse := models.ListResponse[*models.SnippetResponse]{
		Data:       responses,
		Total:      len(responses), // TODO: Get actual total count
		Page:       1,              // TODO: Calculate from offset/limit
		PageSize:   len(responses),
		TotalPages: 1, // TODO: Calculate from total/page_size
	}

	json.NewEncoder(w).Encode(listResponse)
}
