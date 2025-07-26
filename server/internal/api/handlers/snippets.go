package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dikkadev/proompt/server/internal/api/models"
	"github.com/dikkadev/proompt/server/internal/logging"
	"github.com/dikkadev/proompt/server/internal/repository"
	"github.com/google/uuid"
)

// SnippetHandlers contains handlers for snippet operations
type SnippetHandlers struct {
	repo   repository.Repository
	logger *slog.Logger
}

// NewSnippetHandlers creates a new snippet handlers instance
func NewSnippetHandlers(repo repository.Repository) *SnippetHandlers {
	return &SnippetHandlers{
		repo:   repo,
		logger: logging.NewLogger("handlers.snippets"),
	}
}

// CreateSnippet godoc
// @Summary Create a new snippet
// @Description Create a new code snippet or text block
// @Tags snippets
// @Accept json
// @Produce json
// @Param request body models.CreateSnippetRequest true "Snippet creation data"
// @Success 201 {object} models.SnippetResponse "Successfully created snippet"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /snippets [post]
func (h *SnippetHandlers) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("CreateSnippet handler started")

	var req models.CreateSnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Debug("Failed to decode request body", "error", err)
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	h.logger.Debug("Decoded create snippet request",
		"title", req.Title,
		"description", req.Description,
		"content_length", len(req.Content))

	// Basic validation
	if req.Title == "" {
		h.logger.Debug("Validation failed: title is required")
		models.WriteBadRequest(w, "Title is required")
		return
	}
	if req.Content == "" {
		h.logger.Debug("Validation failed: content is required")
		models.WriteBadRequest(w, "Content is required")
		return
	}

	// Convert to domain model
	snippet := req.ToSnippet()
	snippet.ID = uuid.New().String()

	h.logger.Debug("Generated snippet ID and converted to domain model",
		"snippet_id", snippet.ID,
		"has_description", snippet.Description != nil)
	// Create snippet
	if err := h.repo.Snippets().Create(r.Context(), snippet); err != nil {
		h.logger.Error("Failed to create snippet in repository",
			"snippet_id", snippet.ID,
			"error", err)
		models.WriteInternalError(w, "Failed to create snippet")
		return
	}

	h.logger.Debug("Successfully created snippet", "snippet_id", snippet.ID)

	// Return created snippet
	response := models.FromSnippet(snippet)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetSnippet godoc
// @Summary Get a snippet by ID
// @Description Retrieve a specific snippet by its unique identifier
// @Tags snippets
// @Accept json
// @Produce json
// @Param id path string true "Snippet ID" format(uuid)
// @Success 200 {object} models.SnippetResponse "Snippet details"
// @Failure 400 {object} models.ErrorResponse "Invalid snippet ID"
// @Failure 404 {object} models.ErrorResponse "Snippet not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /snippets/{id} [get]
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

// UpdateSnippet godoc
// @Summary Update a snippet
// @Description Update an existing snippet with new data
// @Tags snippets
// @Accept json
// @Produce json
// @Param id path string true "Snippet ID" format(uuid)
// @Param request body models.UpdateSnippetRequest true "Snippet update data"
// @Success 200 {object} models.SnippetResponse "Updated snippet"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 404 {object} models.ErrorResponse "Snippet not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /snippets/{id} [put]
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

// DeleteSnippet godoc
// @Summary Delete a snippet
// @Description Delete a snippet by its ID
// @Tags snippets
// @Accept json
// @Produce json
// @Param id path string true "Snippet ID" format(uuid)
// @Success 204 "Snippet successfully deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid snippet ID"
// @Failure 404 {object} models.ErrorResponse "Snippet not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /snippets/{id} [delete]
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

// ListSnippets godoc
// @Summary List snippets
// @Description Get a paginated list of snippets with optional filtering
// @Tags snippets
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Items per page (default: 20, max: 100)" minimum(1) maximum(100)
// @Param search query string false "Search term for title/content"
// @Param tags query string false "Filter by tags (comma-separated)"
// @Success 200 {object} models.SnippetListResponse "List of snippets"
// @Failure 400 {object} models.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /snippets [get]
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

// AddSnippetTag godoc
// @Summary Add a tag to a snippet
// @Description Add a new tag to an existing snippet
// @Tags snippet-tags
// @Accept json
// @Produce json
// @Param id path string true "Snippet ID" format(uuid)
// @Param request body models.AddTagRequest true "Tag data"
// @Success 201 {object} models.TagResponse "Successfully added tag"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 404 {object} models.ErrorResponse "Snippet not found"
// @Failure 409 {object} models.ErrorResponse "Tag already exists"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /snippets/{id}/tags [post]
func (h *SnippetHandlers) AddSnippetTag(w http.ResponseWriter, r *http.Request) {
	snippetID := r.PathValue("id")
	if snippetID == "" {
		models.WriteBadRequest(w, "Snippet ID is required")
		return
	}

	var req models.AddTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	// Basic validation
	if req.TagName == "" {
		models.WriteBadRequest(w, "Tag name is required")
		return
	}

	// Check if snippet exists
	_, err := h.repo.Snippets().GetByID(r.Context(), snippetID)
	if err != nil {
		models.WriteNotFound(w, "Snippet")
		return
	}

	if err := h.repo.Snippets().AddTag(r.Context(), snippetID, req.TagName); err != nil {
		models.WriteInternalError(w, "Failed to add tag to snippet")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// RemoveSnippetTag godoc
// @Summary Remove a tag from a snippet
// @Description Remove an existing tag from a snippet
// @Tags snippet-tags
// @Accept json
// @Produce json
// @Param id path string true "Snippet ID" format(uuid)
// @Param tagName path string true "Tag name to remove"
// @Success 204 "Tag successfully removed"
// @Failure 400 {object} models.ErrorResponse "Invalid parameters"
// @Failure 404 {object} models.ErrorResponse "Snippet or tag not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /snippets/{id}/tags/{tagName} [delete]
func (h *SnippetHandlers) RemoveSnippetTag(w http.ResponseWriter, r *http.Request) {
	snippetID := r.PathValue("id")
	tagName := r.PathValue("tagName")

	if snippetID == "" || tagName == "" {
		models.WriteBadRequest(w, "Both snippet ID and tag name are required")
		return
	}

	if err := h.repo.Snippets().RemoveTag(r.Context(), snippetID, tagName); err != nil {
		models.WriteNotFound(w, "Tag")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetSnippetTags godoc
// @Summary Get tags for a snippet
// @Description Get all tags associated with a specific snippet
// @Tags snippet-tags
// @Accept json
// @Produce json
// @Param id path string true "Snippet ID" format(uuid)
// @Success 200 {object} models.TagListResponse "List of snippet tags"
// @Failure 400 {object} models.ErrorResponse "Invalid snippet ID"
// @Failure 404 {object} models.ErrorResponse "Snippet not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /snippets/{id}/tags [get]
func (h *SnippetHandlers) GetSnippetTags(w http.ResponseWriter, r *http.Request) {
	snippetID := r.PathValue("id")
	if snippetID == "" {
		models.WriteBadRequest(w, "Snippet ID is required")
		return
	}

	tags, err := h.repo.Snippets().GetTags(r.Context(), snippetID)
	if err != nil {
		models.WriteInternalError(w, "Failed to get snippet tags")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"tags": tags})
}

// ListAllSnippetTags godoc
// @Summary List all snippet tags
// @Description Get a list of all available snippet tags in the system
// @Tags snippet-tags
// @Accept json
// @Produce json
// @Param search query string false "Search term for tag names"
// @Success 200 {object} models.TagListResponse "List of all snippet tags"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /snippets/tags [get]
func (h *SnippetHandlers) ListAllSnippetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.repo.Snippets().ListAllTags(r.Context())
	if err != nil {
		models.WriteInternalError(w, "Failed to list all snippet tags")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"tags": tags})
}
