package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dikkadev/proompt/server/internal/api/models"
	"github.com/dikkadev/proompt/server/internal/logging"
	domainModels "github.com/dikkadev/proompt/server/internal/models"
	"github.com/dikkadev/proompt/server/internal/repository"
	"github.com/google/uuid"
)

// PromptHandlers contains handlers for prompt operations
type PromptHandlers struct {
	repo   repository.Repository
	logger *slog.Logger
}

// NewPromptHandlers creates a new prompt handlers instance
func NewPromptHandlers(repo repository.Repository) *PromptHandlers {
	return &PromptHandlers{
		repo:   repo,
		logger: logging.NewLogger("handlers.prompts"),
	}
}

// CreatePrompt handles POST /api/prompts
func (h *PromptHandlers) CreatePrompt(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("CreatePrompt handler started")

	var req models.CreatePromptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Debug("Failed to decode request body", "error", err)
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	h.logger.Debug("Decoded create prompt request",
		"title", req.Title,
		"type", req.Type,
		"content_length", len(req.Content))

	// TODO: Add validation middleware
	// For now, basic validation
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
	prompt := req.ToPrompt()
	prompt.ID = uuid.New().String()

	h.logger.Debug("Generated prompt ID and converted to domain model",
		"prompt_id", prompt.ID,
		"prompt_type", prompt.Type)

	// Create prompt
	if err := h.repo.Prompts().Create(r.Context(), prompt); err != nil {
		h.logger.Error("Failed to create prompt in repository",
			"prompt_id", prompt.ID,
			"error", err)
		models.WriteInternalError(w, "Failed to create prompt")
		return
	}

	h.logger.Debug("Successfully created prompt", "prompt_id", prompt.ID)

	// Return created prompt
	response := models.FromPrompt(prompt)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	h.logger.Debug("CreatePrompt handler completed successfully", "prompt_id", prompt.ID)
}

// GetPrompt handles GET /api/prompts/{id}
func (h *PromptHandlers) GetPrompt(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h.logger.Debug("GetPrompt handler started", "prompt_id", id)

	if id == "" {
		h.logger.Debug("Validation failed: prompt ID is required")
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	prompt, err := h.repo.Prompts().GetByID(r.Context(), id)
	if err != nil {
		h.logger.Debug("Prompt not found in repository", "prompt_id", id, "error", err)
		models.WriteNotFound(w, "Prompt")
		return
	}

	h.logger.Debug("Successfully retrieved prompt",
		"prompt_id", id,
		"title", prompt.Title,
		"type", prompt.Type)

	response := models.FromPrompt(prompt)
	json.NewEncoder(w).Encode(response)

	h.logger.Debug("GetPrompt handler completed successfully", "prompt_id", id)
}

// UpdatePrompt handles PUT /api/prompts/{id}
func (h *PromptHandlers) UpdatePrompt(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h.logger.Debug("UpdatePrompt handler started", "prompt_id", id)

	if id == "" {
		h.logger.Debug("Validation failed: prompt ID is required")
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	// Get existing prompt
	existing, err := h.repo.Prompts().GetByID(r.Context(), id)
	if err != nil {
		h.logger.Debug("Prompt not found for update", "prompt_id", id, "error", err)
		models.WriteNotFound(w, "Prompt")
		return
	}

	h.logger.Debug("Retrieved existing prompt for update",
		"prompt_id", id,
		"current_title", existing.Title)

	var req models.UpdatePromptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Debug("Failed to decode update request body", "prompt_id", id, "error", err)
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	// Track what fields are being updated
	var updatedFields []string

	// Apply updates
	if req.Title != nil {
		h.logger.Debug("Updating title", "prompt_id", id, "old_title", existing.Title, "new_title", *req.Title)
		existing.Title = *req.Title
		updatedFields = append(updatedFields, "title")
	}
	if req.Content != nil {
		h.logger.Debug("Updating content", "prompt_id", id, "new_content_length", len(*req.Content))
		existing.Content = *req.Content
		updatedFields = append(updatedFields, "content")
	}
	if req.Type != nil {
		h.logger.Debug("Updating type", "prompt_id", id, "old_type", existing.Type, "new_type", *req.Type)
		existing.Type = domainModels.PromptType(*req.Type)
		updatedFields = append(updatedFields, "type")
	}
	if req.UseCase != nil {
		existing.UseCase = req.UseCase
		updatedFields = append(updatedFields, "use_case")
	}
	if req.ModelCompatibilityTags != nil {
		existing.ModelCompatibilityTags = domainModels.StringSlice(req.ModelCompatibilityTags)
		updatedFields = append(updatedFields, "model_compatibility_tags")
	}
	if req.TemperatureSuggestion != nil {
		existing.TemperatureSuggestion = req.TemperatureSuggestion
		updatedFields = append(updatedFields, "temperature_suggestion")
	}
	if req.OtherParameters != nil {
		existing.OtherParameters = domainModels.JSONMap(req.OtherParameters)
		updatedFields = append(updatedFields, "other_parameters")
	}

	h.logger.Debug("Applying updates to prompt",
		"prompt_id", id,
		"updated_fields", updatedFields)

	// Update prompt
	if err := h.repo.Prompts().Update(r.Context(), existing); err != nil {
		h.logger.Error("Failed to update prompt in repository",
			"prompt_id", id,
			"updated_fields", updatedFields,
			"error", err)
		models.WriteInternalError(w, "Failed to update prompt")
		return
	}

	h.logger.Debug("Successfully updated prompt",
		"prompt_id", id,
		"updated_fields", updatedFields)

	// Return updated prompt
	response := models.FromPrompt(existing)
	json.NewEncoder(w).Encode(response)

	h.logger.Debug("UpdatePrompt handler completed successfully", "prompt_id", id)
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
	h.logger.Debug("ListPrompts handler started")

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
		} else {
			h.logger.Debug("Invalid limit parameter", "limit", limitParam, "error", err)
		}
	}
	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if offset, err := strconv.Atoi(offsetParam); err == nil {
			filters.Offset = &offset
		} else {
			h.logger.Debug("Invalid offset parameter", "offset", offsetParam, "error", err)
		}
	}

	h.logger.Debug("Parsed query filters",
		"type", filters.Type,
		"use_case", filters.UseCase,
		"limit", filters.Limit,
		"offset", filters.Offset)

	prompts, err := h.repo.Prompts().List(r.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to list prompts from repository",
			"filters", filters,
			"error", err)
		models.WriteInternalError(w, "Failed to list prompts")
		return
	}

	h.logger.Debug("Successfully retrieved prompts from repository",
		"count", len(prompts))

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

	h.logger.Debug("ListPrompts handler completed successfully",
		"returned_count", len(responses))
}

// CreatePromptLink handles POST /api/prompts/{id}/links
func (h *PromptHandlers) CreatePromptLink(w http.ResponseWriter, r *http.Request) {
	fromPromptID := r.PathValue("id")
	if fromPromptID == "" {
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	var req models.CreatePromptLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	// Basic validation
	if req.ToPromptID == "" {
		models.WriteBadRequest(w, "To prompt ID is required")
		return
	}

	if fromPromptID == req.ToPromptID {
		models.WriteBadRequest(w, "Cannot link prompt to itself")
		return
	}

	// Check if both prompts exist
	_, err := h.repo.Prompts().GetByID(r.Context(), fromPromptID)
	if err != nil {
		models.WriteNotFound(w, "From prompt")
		return
	}

	_, err = h.repo.Prompts().GetByID(r.Context(), req.ToPromptID)
	if err != nil {
		models.WriteNotFound(w, "To prompt")
		return
	}

	// Create the link
	link := &domainModels.PromptLink{
		FromPromptID: fromPromptID,
		ToPromptID:   req.ToPromptID,
		LinkType:     req.LinkType,
	}

	if err := h.repo.Prompts().CreateLink(r.Context(), link); err != nil {
		models.WriteInternalError(w, "Failed to create prompt link")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.FromPromptLink(link))
}

// AddPromptTag handles POST /api/prompts/{id}/tags
func (h *PromptHandlers) AddPromptTag(w http.ResponseWriter, r *http.Request) {
	promptID := r.PathValue("id")
	if promptID == "" {
		models.WriteBadRequest(w, "Prompt ID is required")
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

	// Check if prompt exists
	_, err := h.repo.Prompts().GetByID(r.Context(), promptID)
	if err != nil {
		models.WriteNotFound(w, "Prompt")
		return
	}

	if err := h.repo.Prompts().AddTag(r.Context(), promptID, req.TagName); err != nil {
		models.WriteInternalError(w, "Failed to add tag to prompt")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// RemovePromptTag handles DELETE /api/prompts/{id}/tags/{tagName}
func (h *PromptHandlers) RemovePromptTag(w http.ResponseWriter, r *http.Request) {
	promptID := r.PathValue("id")
	tagName := r.PathValue("tagName")

	if promptID == "" || tagName == "" {
		models.WriteBadRequest(w, "Both prompt ID and tag name are required")
		return
	}

	if err := h.repo.Prompts().RemoveTag(r.Context(), promptID, tagName); err != nil {
		models.WriteNotFound(w, "Tag")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPromptTags handles GET /api/prompts/{id}/tags
func (h *PromptHandlers) GetPromptTags(w http.ResponseWriter, r *http.Request) {
	promptID := r.PathValue("id")
	if promptID == "" {
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	tags, err := h.repo.Prompts().GetTags(r.Context(), promptID)
	if err != nil {
		models.WriteInternalError(w, "Failed to get prompt tags")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"tags": tags})
}

// ListAllPromptTags handles GET /api/prompts/tags
func (h *PromptHandlers) ListAllPromptTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.repo.Prompts().ListAllTags(r.Context())
	if err != nil {
		models.WriteInternalError(w, "Failed to list all prompt tags")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"tags": tags})
}

// DeletePromptLink handles DELETE /api/prompts/{id}/links/{toId}
func (h *PromptHandlers) DeletePromptLink(w http.ResponseWriter, r *http.Request) {
	fromPromptID := r.PathValue("id")
	toPromptID := r.PathValue("toId")

	if fromPromptID == "" || toPromptID == "" {
		models.WriteBadRequest(w, "Both prompt IDs are required")
		return
	}

	if err := h.repo.Prompts().DeleteLink(r.Context(), fromPromptID, toPromptID); err != nil {
		models.WriteNotFound(w, "Prompt link")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPromptLinksFrom handles GET /api/prompts/{id}/links
func (h *PromptHandlers) GetPromptLinksFrom(w http.ResponseWriter, r *http.Request) {
	promptID := r.PathValue("id")
	if promptID == "" {
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	links, err := h.repo.Prompts().GetLinksFrom(r.Context(), promptID)
	if err != nil {
		models.WriteInternalError(w, "Failed to get prompt links")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.FromPromptLinks(links))
}

// GetPromptLinksTo handles GET /api/prompts/{id}/backlinks
func (h *PromptHandlers) GetPromptLinksTo(w http.ResponseWriter, r *http.Request) {
	promptID := r.PathValue("id")
	if promptID == "" {
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	links, err := h.repo.Prompts().GetLinksTo(r.Context(), promptID)
	if err != nil {
		models.WriteInternalError(w, "Failed to get prompt backlinks")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.FromPromptLinks(links))
}
