package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dikkadev/proompt/server/internal/api/models"
	"github.com/dikkadev/proompt/server/internal/repository"
	"github.com/google/uuid"
)

// NoteHandlers contains handlers for note operations
type NoteHandlers struct {
	repo repository.Repository
}

// NewNoteHandlers creates a new note handlers instance
func NewNoteHandlers(repo repository.Repository) *NoteHandlers {
	return &NoteHandlers{repo: repo}
}

// CreateNote godoc
// @Summary Create a note for a prompt
// @Description Create a new note associated with a specific prompt
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Prompt ID" format(uuid)
// @Param request body models.CreateNoteRequest true "Note creation data"
// @Success 201 {object} models.NoteResponse "Successfully created note"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 404 {object} models.ErrorResponse "Prompt not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id}/notes [post]
func (h *NoteHandlers) CreateNote(w http.ResponseWriter, r *http.Request) {
	promptID := r.PathValue("id")
	if promptID == "" {
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	// Verify prompt exists
	_, err := h.repo.Prompts().GetByID(r.Context(), promptID)
	if err != nil {
		models.WriteNotFound(w, "Prompt")
		return
	}

	var req models.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	// Basic validation
	if req.Title == "" {
		models.WriteBadRequest(w, "Title is required")
		return
	}
	if req.Body == "" {
		models.WriteBadRequest(w, "Body is required")
		return
	}

	// Convert to domain model
	note := req.ToNote(promptID)
	note.ID = uuid.New().String()

	// Create note
	if err := h.repo.Notes().Create(r.Context(), note); err != nil {
		models.WriteInternalError(w, "Failed to create note")
		return
	}

	// Return created note
	response := models.FromNote(note)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetNote godoc
// @Summary Get a note by ID
// @Description Retrieve a specific note by its unique identifier
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID" format(uuid)
// @Success 200 {object} models.NoteResponse "Note details"
// @Failure 400 {object} models.ErrorResponse "Invalid note ID"
// @Failure 404 {object} models.ErrorResponse "Note not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /notes/{id} [get]
func (h *NoteHandlers) GetNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		models.WriteBadRequest(w, "Note ID is required")
		return
	}

	note, err := h.repo.Notes().GetByID(r.Context(), id)
	if err != nil {
		models.WriteNotFound(w, "Note")
		return
	}

	response := models.FromNote(note)
	json.NewEncoder(w).Encode(response)
}

// UpdateNote godoc
// @Summary Update a note
// @Description Update an existing note with new data
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID" format(uuid)
// @Param request body models.UpdateNoteRequest true "Note update data"
// @Success 200 {object} models.NoteResponse "Updated note"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 404 {object} models.ErrorResponse "Note not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /notes/{id} [put]
func (h *NoteHandlers) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		models.WriteBadRequest(w, "Note ID is required")
		return
	}

	// Get existing note
	existing, err := h.repo.Notes().GetByID(r.Context(), id)
	if err != nil {
		models.WriteNotFound(w, "Note")
		return
	}

	var req models.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	// Apply updates
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Body != nil {
		existing.Body = req.Body
	}

	// Update note
	if err := h.repo.Notes().Update(r.Context(), existing); err != nil {
		models.WriteInternalError(w, "Failed to update note")
		return
	}

	// Return updated note
	response := models.FromNote(existing)
	json.NewEncoder(w).Encode(response)
}

// DeleteNote godoc
// @Summary Delete a note
// @Description Delete a note by its ID
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID" format(uuid)
// @Success 204 "Note successfully deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid note ID"
// @Failure 404 {object} models.ErrorResponse "Note not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /notes/{id} [delete]
func (h *NoteHandlers) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		models.WriteBadRequest(w, "Note ID is required")
		return
	}

	if err := h.repo.Notes().Delete(r.Context(), id); err != nil {
		models.WriteInternalError(w, "Failed to delete note")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListNotesForPrompt godoc
// @Summary List notes for a prompt
// @Description Get all notes associated with a specific prompt
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Prompt ID" format(uuid)
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Items per page (default: 20, max: 100)" minimum(1) maximum(100)
// @Success 200 {object} models.NoteListResponse "List of notes for the prompt"
// @Failure 400 {object} models.ErrorResponse "Invalid parameters"
// @Failure 404 {object} models.ErrorResponse "Prompt not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /prompts/{id}/notes [get]
func (h *NoteHandlers) ListNotesForPrompt(w http.ResponseWriter, r *http.Request) {
	promptID := r.PathValue("id")
	if promptID == "" {
		models.WriteBadRequest(w, "Prompt ID is required")
		return
	}

	// Verify prompt exists
	_, err := h.repo.Prompts().GetByID(r.Context(), promptID)
	if err != nil {
		models.WriteNotFound(w, "Prompt")
		return
	}

	notes, err := h.repo.Notes().ListByPromptID(r.Context(), promptID)
	if err != nil {
		models.WriteInternalError(w, "Failed to list notes")
		return
	}

	responses := models.FromNotes(notes)

	// Create list response
	listResponse := models.ListResponse[*models.NoteResponse]{
		Data:       responses,
		Total:      len(responses),
		Page:       1,
		PageSize:   len(responses),
		TotalPages: 1,
	}

	json.NewEncoder(w).Encode(listResponse)
}
