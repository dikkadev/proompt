package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/dikkadev/proompt/server/internal/git"
	"github.com/dikkadev/proompt/server/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// snippetRepository implements SnippetRepository interface
type snippetRepository struct {
	db         txExecutor
	gitService git.GitService
	logger     *slog.Logger
}

// newSnippetRepository creates a new snippet repository
func newSnippetRepository(db *sqlx.DB, gitService git.GitService, logger *slog.Logger) SnippetRepository {
	return &snippetRepository{
		db:         db,
		gitService: gitService,
		logger:     logger,
	}
}

// newSnippetRepositoryWithTx creates a new snippet repository with transaction
func newSnippetRepositoryWithTx(tx *sqlx.Tx, gitService git.GitService, logger *slog.Logger) SnippetRepository {
	return &snippetRepository{
		db:         tx,
		gitService: gitService,
		logger:     logger,
	}
}

// Create creates a new snippet
func (r *snippetRepository) Create(ctx context.Context, snippet *models.Snippet) error {
	if snippet.ID == "" {
		snippet.ID = uuid.New().String()
	}

	now := time.Now()
	snippet.CreatedAt = now
	snippet.UpdatedAt = now

	r.logger.Debug("Creating snippet", "id", snippet.ID, "title", snippet.Title)

	query := `
		INSERT INTO snippets (
			id, title, content, description, created_at, updated_at
		) VALUES (
			:id, :title, :content, :description, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, snippet)
	if err != nil {
		r.logger.Error("Failed to create snippet in database", "error", err, "id", snippet.ID)
		return fmt.Errorf("failed to create snippet: %w", err)
	}

	// Create git branch for versioning
	if err := r.gitService.CreateSnippetBranch(ctx, snippet, ""); err != nil {
		r.logger.Error("Failed to create git branch for snippet", "error", err, "id", snippet.ID)
		return fmt.Errorf("failed to create git branch: %w", err)
	}

	r.logger.Info("Snippet created successfully", "id", snippet.ID, "title", snippet.Title)
	return nil
}

// GetByID retrieves a snippet by ID
func (r *snippetRepository) GetByID(ctx context.Context, id string) (*models.Snippet, error) {
	r.logger.Debug("Getting snippet by ID", "id", id)

	query := `
		SELECT id, title, content, description, created_at, updated_at, git_ref
		FROM snippets 
		WHERE id = ?`

	var snippet models.Snippet
	err := r.db.GetContext(ctx, &snippet, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debug("Snippet not found", "id", id)
			return nil, fmt.Errorf("snippet not found: %s", id)
		}
		r.logger.Error("Failed to get snippet", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get snippet: %w", err)
	}

	r.logger.Debug("Snippet retrieved successfully", "id", id, "title", snippet.Title)
	return &snippet, nil
}

// Update updates an existing snippet
func (r *snippetRepository) Update(ctx context.Context, snippet *models.Snippet) error {
	snippet.UpdatedAt = time.Now()

	r.logger.Debug("Updating snippet", "id", snippet.ID, "title", snippet.Title)

	query := `
		UPDATE snippets SET
			title = :title,
			content = :content,
			description = :description,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, snippet)
	if err != nil {
		r.logger.Error("Failed to update snippet in database", "error", err, "id", snippet.ID)
		return fmt.Errorf("failed to update snippet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err, "id", snippet.ID)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Debug("Snippet not found for update", "id", snippet.ID)
		return fmt.Errorf("snippet not found: %s", snippet.ID)
	}

	// Update git branch
	if err := r.gitService.UpdateSnippetBranch(ctx, snippet, ""); err != nil {
		r.logger.Error("Failed to update git branch for snippet", "error", err, "id", snippet.ID)
		return fmt.Errorf("failed to update git branch: %w", err)
	}

	r.logger.Info("Snippet updated successfully", "id", snippet.ID, "title", snippet.Title)
	return nil
}

// Delete deletes a snippet
func (r *snippetRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug("Deleting snippet", "id", id)

	query := `DELETE FROM snippets WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete snippet from database", "error", err, "id", id)
		return fmt.Errorf("failed to delete snippet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err, "id", id)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Debug("Snippet not found for deletion", "id", id)
		return fmt.Errorf("snippet not found: %s", id)
	}

	// Delete git branch
	if err := r.gitService.DeleteSnippetBranch(ctx, id); err != nil {
		r.logger.Error("Failed to delete git branch for snippet", "error", err, "id", id)
		return fmt.Errorf("failed to delete git branch: %w", err)
	}

	r.logger.Info("Snippet deleted successfully", "id", id)
	return nil
}

// List retrieves snippets with filtering
func (r *snippetRepository) List(ctx context.Context, filters SnippetFilters) ([]*models.Snippet, error) {
	r.logger.Debug("Listing snippets with filters")

	query := `
		SELECT id, title, content, description, created_at, updated_at, git_ref
		FROM snippets`

	var conditions []string
	var args []interface{}

	if len(filters.Tags) > 0 {
		placeholders := make([]string, len(filters.Tags))
		for i, tag := range filters.Tags {
			placeholders[i] = "?"
			args = append(args, "%"+tag+"%")
		}
		conditions = append(conditions, fmt.Sprintf("tags LIKE %s", strings.Join(placeholders, " OR tags LIKE ")))
	}

	if filters.CreatedAfter != nil {
		conditions = append(conditions, "created_at > ?")
		args = append(args, *filters.CreatedAfter)
	}

	if filters.CreatedBefore != nil {
		conditions = append(conditions, "created_at < ?")
		args = append(args, *filters.CreatedBefore)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY updated_at DESC"

	if filters.Limit != nil {
		query += " LIMIT ?"
		args = append(args, *filters.Limit)
	}

	if filters.Offset != nil {
		query += " OFFSET ?"
		args = append(args, *filters.Offset)
	}

	var snippets []*models.Snippet
	err := r.db.SelectContext(ctx, &snippets, query, args...)
	if err != nil {
		r.logger.Error("Failed to list snippets", "error", err)
		return nil, fmt.Errorf("failed to list snippets: %w", err)
	}

	r.logger.Debug("Snippets listed successfully", "count", len(snippets))
	return snippets, nil
}

// Search searches snippets using full-text search
func (r *snippetRepository) Search(ctx context.Context, query string) ([]*models.Snippet, error) {
	r.logger.Debug("Searching snippets", "query", query)

	searchQuery := `
		SELECT s.id, s.title, s.content, s.description, s.created_at, s.updated_at, s.git_ref
		FROM snippets s
		JOIN snippets_fts fts ON s.id = fts.rowid
		WHERE snippets_fts MATCH ?
		ORDER BY rank`

	var snippets []*models.Snippet
	err := r.db.SelectContext(ctx, &snippets, searchQuery, query)
	if err != nil {
		r.logger.Error("Failed to search snippets", "error", err, "query", query)
		return nil, fmt.Errorf("failed to search snippets: %w", err)
	}

	r.logger.Debug("Snippets search completed", "query", query, "count", len(snippets))
	return snippets, nil
}
