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

// promptRepository implements PromptRepository interface
type promptRepository struct {
	db         txExecutor
	gitService git.GitService
	logger     *slog.Logger
}

// newPromptRepository creates a new prompt repository
func newPromptRepository(db *sqlx.DB, gitService git.GitService, logger *slog.Logger) PromptRepository {
	return &promptRepository{
		db:         db,
		gitService: gitService,
		logger:     logger,
	}
}

// newPromptRepositoryWithTx creates a new prompt repository with transaction
func newPromptRepositoryWithTx(tx *sqlx.Tx, gitService git.GitService, logger *slog.Logger) PromptRepository {
	return &promptRepository{
		db:         tx,
		gitService: gitService,
		logger:     logger,
	}
}

// Create creates a new prompt
func (r *promptRepository) Create(ctx context.Context, prompt *models.Prompt) error {
	if prompt.ID == "" {
		prompt.ID = uuid.New().String()
	}

	now := time.Now()
	prompt.CreatedAt = now
	prompt.UpdatedAt = now

	r.logger.Debug("Creating prompt", "id", prompt.ID, "title", prompt.Title)

	query := `
		INSERT INTO prompts (
			id, title, content, type, use_case, model_compatibility_tags, 
			temperature_suggestion, other_parameters, created_at, updated_at
		) VALUES (
			:id, :title, :content, :type, :use_case, :model_compatibility_tags,
			:temperature_suggestion, :other_parameters, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, prompt)
	if err != nil {
		r.logger.Error("Failed to create prompt in database", "error", err, "id", prompt.ID)
		return fmt.Errorf("failed to create prompt: %w", err)
	}

	// Create git branch for versioning
	if err := r.gitService.CreatePromptBranch(ctx, prompt, ""); err != nil {
		r.logger.Error("Failed to create git branch for prompt", "error", err, "id", prompt.ID)
		return fmt.Errorf("failed to create git branch: %w", err)
	}

	r.logger.Info("Prompt created successfully", "id", prompt.ID, "title", prompt.Title)
	return nil
}

// GetByID retrieves a prompt by ID
func (r *promptRepository) GetByID(ctx context.Context, id string) (*models.Prompt, error) {
	r.logger.Debug("Getting prompt by ID", "id", id)

	query := `
		SELECT id, title, content, type, use_case, model_compatibility_tags,
		       temperature_suggestion, other_parameters, created_at, updated_at, git_ref
		FROM prompts 
		WHERE id = ?`

	var prompt models.Prompt
	err := r.db.GetContext(ctx, &prompt, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debug("Prompt not found", "id", id)
			return nil, fmt.Errorf("prompt not found: %s", id)
		}
		r.logger.Error("Failed to get prompt", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	r.logger.Debug("Prompt retrieved successfully", "id", id, "title", prompt.Title)
	return &prompt, nil
}

// Update updates an existing prompt
func (r *promptRepository) Update(ctx context.Context, prompt *models.Prompt) error {
	prompt.UpdatedAt = time.Now()

	r.logger.Debug("Updating prompt", "id", prompt.ID, "title", prompt.Title)

	query := `
		UPDATE prompts SET
			title = :title,
			content = :content,
			type = :type,
			use_case = :use_case,
			model_compatibility_tags = :model_compatibility_tags,
			temperature_suggestion = :temperature_suggestion,
			other_parameters = :other_parameters,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, prompt)
	if err != nil {
		r.logger.Error("Failed to update prompt in database", "error", err, "id", prompt.ID)
		return fmt.Errorf("failed to update prompt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err, "id", prompt.ID)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Debug("Prompt not found for update", "id", prompt.ID)
		return fmt.Errorf("prompt not found: %s", prompt.ID)
	}

	// Update git branch
	if err := r.gitService.UpdatePromptBranch(ctx, prompt, ""); err != nil {
		r.logger.Error("Failed to update git branch for prompt", "error", err, "id", prompt.ID)
		return fmt.Errorf("failed to update git branch: %w", err)
	}

	r.logger.Info("Prompt updated successfully", "id", prompt.ID, "title", prompt.Title)
	return nil
}

// Delete deletes a prompt
func (r *promptRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug("Deleting prompt", "id", id)

	query := `DELETE FROM prompts WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete prompt from database", "error", err, "id", id)
		return fmt.Errorf("failed to delete prompt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err, "id", id)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Debug("Prompt not found for deletion", "id", id)
		return fmt.Errorf("prompt not found: %s", id)
	}

	// Delete git branch
	if err := r.gitService.DeletePromptBranch(ctx, id); err != nil {
		r.logger.Error("Failed to delete git branch for prompt", "error", err, "id", id)
		return fmt.Errorf("failed to delete git branch: %w", err)
	}

	r.logger.Info("Prompt deleted successfully", "id", id)
	return nil
}

// List retrieves prompts with filtering
func (r *promptRepository) List(ctx context.Context, filters PromptFilters) ([]*models.Prompt, error) {
	r.logger.Debug("Listing prompts with filters")

	query := `
		SELECT id, title, content, type, use_case, model_compatibility_tags,
		       temperature_suggestion, other_parameters, created_at, updated_at, git_ref
		FROM prompts`

	var conditions []string
	var args []interface{}

	if filters.Type != nil {
		conditions = append(conditions, "type = ?")
		args = append(args, *filters.Type)
	}

	if filters.UseCase != nil {
		conditions = append(conditions, "use_case = ?")
		args = append(args, *filters.UseCase)
	}

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

	var prompts []*models.Prompt
	err := r.db.SelectContext(ctx, &prompts, query, args...)
	if err != nil {
		r.logger.Error("Failed to list prompts", "error", err)
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	r.logger.Debug("Prompts listed successfully", "count", len(prompts))
	return prompts, nil
}

// Search searches prompts using full-text search
func (r *promptRepository) Search(ctx context.Context, query string) ([]*models.Prompt, error) {
	r.logger.Debug("Searching prompts", "query", query)

	searchQuery := `
		SELECT p.id, p.title, p.content, p.type, p.use_case, p.model_compatibility_tags,
		       p.temperature_suggestion, p.other_parameters, p.created_at, p.updated_at, p.git_ref
		FROM prompts p
		JOIN prompts_fts fts ON p.id = fts.rowid
		WHERE prompts_fts MATCH ?
		ORDER BY rank`

	var prompts []*models.Prompt
	err := r.db.SelectContext(ctx, &prompts, searchQuery, query)
	if err != nil {
		r.logger.Error("Failed to search prompts", "error", err, "query", query)
		return nil, fmt.Errorf("failed to search prompts: %w", err)
	}

	r.logger.Debug("Prompts search completed", "query", query, "count", len(prompts))
	return prompts, nil
}
