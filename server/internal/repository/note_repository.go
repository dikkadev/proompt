package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/dikkadev/proompt/server/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// noteRepository implements NoteRepository interface
type noteRepository struct {
	db     txExecutor
	logger *slog.Logger
}

// newNoteRepository creates a new note repository
func newNoteRepository(db *sqlx.DB, logger *slog.Logger) NoteRepository {
	return &noteRepository{
		db:     db,
		logger: logger,
	}
}

// newNoteRepositoryWithTx creates a new note repository with transaction
func newNoteRepositoryWithTx(tx *sqlx.Tx, logger *slog.Logger) NoteRepository {
	return &noteRepository{
		db:     tx,
		logger: logger,
	}
}

// Create creates a new note
func (r *noteRepository) Create(ctx context.Context, note *models.Note) error {
	if note.ID == "" {
		note.ID = uuid.New().String()
	}

	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now

	r.logger.Debug("Creating note", "id", note.ID, "title", note.Title, "prompt_id", note.PromptID)

	query := `
		INSERT INTO notes (
			id, prompt_id, title, body, created_at, updated_at
		) VALUES (
			:id, :prompt_id, :title, :body, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, note)
	if err != nil {
		r.logger.Error("Failed to create note in database", "error", err, "id", note.ID)
		return fmt.Errorf("failed to create note: %w", err)
	}

	r.logger.Info("Note created successfully", "id", note.ID, "title", note.Title)
	return nil
}

// GetByID retrieves a note by ID
func (r *noteRepository) GetByID(ctx context.Context, id string) (*models.Note, error) {
	r.logger.Debug("Getting note by ID", "id", id)

	query := `
		SELECT id, prompt_id, title, body, created_at, updated_at
		FROM notes 
		WHERE id = ?`

	var note models.Note
	err := r.db.GetContext(ctx, &note, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debug("Note not found", "id", id)
			return nil, fmt.Errorf("note not found: %s", id)
		}
		r.logger.Error("Failed to get note", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	r.logger.Debug("Note retrieved successfully", "id", id, "title", note.Title)
	return &note, nil
}

// Update updates an existing note
func (r *noteRepository) Update(ctx context.Context, note *models.Note) error {
	note.UpdatedAt = time.Now()

	r.logger.Debug("Updating note", "id", note.ID, "title", note.Title)

	query := `
		UPDATE notes SET
			title = :title,
			body = :body,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, note)
	if err != nil {
		r.logger.Error("Failed to update note in database", "error", err, "id", note.ID)
		return fmt.Errorf("failed to update note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err, "id", note.ID)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Debug("Note not found for update", "id", note.ID)
		return fmt.Errorf("note not found: %s", note.ID)
	}

	r.logger.Info("Note updated successfully", "id", note.ID, "title", note.Title)
	return nil
}

// Delete deletes a note
func (r *noteRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug("Deleting note", "id", id)

	query := `DELETE FROM notes WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete note from database", "error", err, "id", id)
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err, "id", id)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Debug("Note not found for deletion", "id", id)
		return fmt.Errorf("note not found: %s", id)
	}

	r.logger.Info("Note deleted successfully", "id", id)
	return nil
}

// ListByPromptID retrieves all notes for a specific prompt
func (r *noteRepository) ListByPromptID(ctx context.Context, promptID string) ([]*models.Note, error) {
	r.logger.Debug("Listing notes by prompt ID", "prompt_id", promptID)

	query := `
		SELECT id, prompt_id, title, body, created_at, updated_at
		FROM notes
		WHERE prompt_id = ?
		ORDER BY created_at DESC`

	var notes []*models.Note
	err := r.db.SelectContext(ctx, &notes, query, promptID)
	if err != nil {
		r.logger.Error("Failed to list notes by prompt ID", "error", err, "prompt_id", promptID)
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	r.logger.Debug("Notes listed successfully", "prompt_id", promptID, "count", len(notes))
	return notes, nil
}

// Search searches notes using full-text search
func (r *noteRepository) Search(ctx context.Context, query string) ([]*models.Note, error) {
	r.logger.Debug("Searching notes", "query", query)

	searchQuery := `
		SELECT n.id, n.prompt_id, n.title, n.body, n.created_at, n.updated_at
		FROM notes n
		JOIN notes_fts fts ON n.id = fts.rowid
		WHERE notes_fts MATCH ?
		ORDER BY rank`

	var notes []*models.Note
	err := r.db.SelectContext(ctx, &notes, searchQuery, query)
	if err != nil {
		r.logger.Error("Failed to search notes", "error", err, "query", query)
		return nil, fmt.Errorf("failed to search notes: %w", err)
	}

	r.logger.Debug("Notes search completed", "query", query, "count", len(notes))
	return notes, nil
}
