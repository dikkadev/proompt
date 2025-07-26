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

// CreatePrompt godoc
// @Summary Create a new prompt
// @Description Create a new prompt with the provided details
// @Tags prompts
// @Accept json
// @Produce json
// @Param request body models.CreatePromptRequest true "Prompt creation data"
// @Success 201 {object} models.PromptResponse "Successfully created prompt"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts [post]
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

// GetPrompt godoc
// @Summary Get a prompt by ID
// @Description Retrieve a specific prompt by its unique identifier
// @Tags prompts
// @Accept json
// @Produce json
// @Param id path string true "Prompt ID" format(uuid)
// @Success 200 {object} models.PromptResponse "Prompt details"
// @Failure 400 {object} models.ErrorResponse "Invalid prompt ID"
// @Failure 404 {object} models.ErrorResponse "Prompt not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id} [get]
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

// UpdatePrompt godoc
// @Summary Update a prompt
// @Description Update an existing prompt with new data
// @Tags prompts
// @Accept json
// @Produce json
// @Param id path string true "Prompt ID" format(uuid)
// @Param request body models.UpdatePromptRequest true "Prompt update data"
// @Success 200 {object} models.PromptResponse "Updated prompt"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 404 {object} models.ErrorResponse "Prompt not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id} [put]
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

// DeletePrompt godoc
// @Summary Delete a prompt
// @Description Delete a prompt by its ID
// @Tags prompts
// @Accept json
// @Produce json
// @Param id path string true "Prompt ID" format(uuid)
// @Success 204 "Prompt successfully deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid prompt ID"
// @Failure 404 {object} models.ErrorResponse "Prompt not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id} [delete]
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

// ListPrompts godoc
// @Summary List prompts
// @Description Get a paginated list of prompts with optional filtering
// @Tags prompts
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Items per page (default: 20, max: 100)" minimum(1) maximum(100)
// @Param search query string false "Search term for title/content"
// @Param type query string false "Filter by prompt type" Enums(system,user,image,video)
// @Param use_case query string false "Filter by use case"
// @Param tags query string false "Filter by tags (comma-separated)"
// @Success 200 {object} models.PromptListResponse "List of prompts"
// @Failure 400 {object} models.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts [get]
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

// CreatePromptLink godoc
// @Summary Create a link between prompts
// @Description Create a relationship link from one prompt to another
// @Tags prompt-links
// @Accept json
// @Produce json
// @Param id path string true "Source Prompt ID" format(uuid)
// @Param request body models.CreatePromptLinkRequest true "Link creation data"
// @Success 201 {object} models.PromptLinkResponse "Successfully created link"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 404 {object} models.ErrorResponse "Prompt not found"
// @Failure 409 {object} models.ErrorResponse "Link already exists"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id}/links [post]
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

// AddPromptTag godoc
// @Summary Add a tag to a prompt
// @Description Add a new tag to an existing prompt
// @Tags prompt-tags
// @Accept json
// @Produce json
// @Param id path string true "Prompt ID" format(uuid)
// @Param request body models.AddTagRequest true "Tag data"
// @Success 201 {object} models.TagResponse "Successfully added tag"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 404 {object} models.ErrorResponse "Prompt not found"
// @Failure 409 {object} models.ErrorResponse "Tag already exists"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id}/tags [post]
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

// RemovePromptTag godoc
// @Summary Remove a tag from a prompt
// @Description Remove an existing tag from a prompt
// @Tags prompt-tags
// @Accept json
// @Produce json
// @Param id path string true "Prompt ID" format(uuid)
// @Param tagName path string true "Tag name to remove"
// @Success 204 "Tag successfully removed"
// @Failure 400 {object} models.ErrorResponse "Invalid parameters"
// @Failure 404 {object} models.ErrorResponse "Prompt or tag not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id}/tags/{tagName} [delete]
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

// GetPromptTags godoc
// @Summary Get tags for a prompt
// @Description Get all tags associated with a specific prompt
// @Tags prompt-tags
// @Accept json
// @Produce json
// @Param id path string true "Prompt ID" format(uuid)
// @Success 200 {object} models.TagListResponse "List of prompt tags"
// @Failure 400 {object} models.ErrorResponse "Invalid prompt ID"
// @Failure 404 {object} models.ErrorResponse "Prompt not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id}/tags [get]
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

// ListAllPromptTags godoc
// @Summary List all prompt tags
// @Description Get a list of all available prompt tags in the system
// @Tags prompt-tags
// @Accept json
// @Produce json
// @Param search query string false "Search term for tag names"
// @Success 200 {object} models.TagListResponse "List of all prompt tags"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/tags [get]
func (h *PromptHandlers) ListAllPromptTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.repo.Prompts().ListAllTags(r.Context())
	if err != nil {
		models.WriteInternalError(w, "Failed to list all prompt tags")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"tags": tags})
}

// DeletePromptLink godoc
// @Summary Delete a link between prompts
// @Description Remove a relationship link between two prompts
// @Tags prompt-links
// @Accept json
// @Produce json
// @Param id path string true "Source Prompt ID" format(uuid)
// @Param toId path string true "Target Prompt ID" format(uuid)
// @Success 204 "Link successfully deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid prompt IDs"
// @Failure 404 {object} models.ErrorResponse "Link or prompt not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id}/links/{toId} [delete]
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

// GetPromptLinksFrom godoc
// @Summary Get outgoing links from a prompt
// @Description Get all prompts that this prompt links to
// @Tags prompt-links
// @Accept json
// @Produce json
// @Param id path string true "Prompt ID" format(uuid)
// @Success 200 {object} models.PromptLinkListResponse "List of outgoing links"
// @Failure 400 {object} models.ErrorResponse "Invalid prompt ID"
// @Failure 404 {object} models.ErrorResponse "Prompt not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id}/links [get]
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

// GetPromptLinksTo godoc
// @Summary Get incoming links to a prompt
// @Description Get all prompts that link to this prompt (backlinks)
// @Tags prompt-links
// @Accept json
// @Produce json
// @Param id path string true "Prompt ID" format(uuid)
// @Success 200 {object} models.PromptLinkListResponse "List of incoming links"
// @Failure 400 {object} models.ErrorResponse "Invalid prompt ID"
// @Failure 404 {object} models.ErrorResponse "Prompt not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id}/backlinks [get]
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
