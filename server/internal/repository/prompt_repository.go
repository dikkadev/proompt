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
	r.logger.Debug("Listing prompts with filters",
		"type", filters.Type,
		"use_case", filters.UseCase,
		"tags", filters.Tags,
		"limit", filters.Limit,
		"offset", filters.Offset)

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

	// For now, use simple LIKE search since we removed FTS
	searchQuery := `
		SELECT id, title, content, type, use_case, model_compatibility_tags,
		       temperature_suggestion, other_parameters, created_at, updated_at, git_ref
		FROM prompts
		WHERE title LIKE ? OR content LIKE ?
		ORDER BY updated_at DESC`

	searchTerm := "%" + query + "%"
	var prompts []*models.Prompt
	err := r.db.SelectContext(ctx, &prompts, searchQuery, searchTerm, searchTerm)
	if err != nil {
		r.logger.Error("Failed to search prompts", "error", err, "query", query)
		return nil, fmt.Errorf("failed to search prompts: %w", err)
	}

	r.logger.Debug("Prompts search completed", "query", query, "count", len(prompts))
	return prompts, nil
}

// CreateLink creates a bidirectional link between two prompts
func (r *promptRepository) CreateLink(ctx context.Context, link *models.PromptLink) error {
	r.logger.Debug("Creating prompt link", "from", link.FromPromptID, "to", link.ToPromptID, "type", link.LinkType)

	// Set default link type if not provided
	if link.LinkType == "" {
		link.LinkType = "followup"
	}

	link.CreatedAt = time.Now()

	query := `
		INSERT INTO prompt_links (from_prompt_id, to_prompt_id, link_type, created_at)
		VALUES (?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, link.FromPromptID, link.ToPromptID, link.LinkType, link.CreatedAt)
	if err != nil {
		r.logger.Error("Failed to create prompt link", "error", err, "from", link.FromPromptID, "to", link.ToPromptID)
		return fmt.Errorf("failed to create prompt link: %w", err)
	}

	r.logger.Info("Prompt link created successfully", "from", link.FromPromptID, "to", link.ToPromptID, "type", link.LinkType)
	return nil
}

// DeleteLink deletes a link between two prompts
func (r *promptRepository) DeleteLink(ctx context.Context, fromPromptID, toPromptID string) error {
	r.logger.Debug("Deleting prompt link", "from", fromPromptID, "to", toPromptID)

	query := `DELETE FROM prompt_links WHERE from_prompt_id = ? AND to_prompt_id = ?`
	result, err := r.db.ExecContext(ctx, query, fromPromptID, toPromptID)
	if err != nil {
		r.logger.Error("Failed to delete prompt link", "error", err, "from", fromPromptID, "to", toPromptID)
		return fmt.Errorf("failed to delete prompt link: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err, "from", fromPromptID, "to", toPromptID)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Debug("Prompt link not found for deletion", "from", fromPromptID, "to", toPromptID)
		return fmt.Errorf("prompt link not found")
	}

	r.logger.Info("Prompt link deleted successfully", "from", fromPromptID, "to", toPromptID)
	return nil
}

// GetLinksFrom retrieves all links from a specific prompt
func (r *promptRepository) GetLinksFrom(ctx context.Context, promptID string) ([]*models.PromptLink, error) {
	r.logger.Debug("Getting links from prompt", "id", promptID)

	query := `
		SELECT from_prompt_id, to_prompt_id, link_type, created_at
		FROM prompt_links
		WHERE from_prompt_id = ?
		ORDER BY created_at DESC`

	var links []*models.PromptLink
	err := r.db.SelectContext(ctx, &links, query, promptID)
	if err != nil {
		r.logger.Error("Failed to get links from prompt", "error", err, "id", promptID)
		return nil, fmt.Errorf("failed to get links from prompt: %w", err)
	}

	r.logger.Debug("Links from prompt retrieved successfully", "id", promptID, "count", len(links))
	return links, nil
}

// GetLinksTo retrieves all links to a specific prompt
func (r *promptRepository) GetLinksTo(ctx context.Context, promptID string) ([]*models.PromptLink, error) {
	r.logger.Debug("Getting links to prompt", "id", promptID)

	query := `
		SELECT from_prompt_id, to_prompt_id, link_type, created_at
		FROM prompt_links
		WHERE to_prompt_id = ?
		ORDER BY created_at DESC`

	var links []*models.PromptLink
	err := r.db.SelectContext(ctx, &links, query, promptID)
	if err != nil {
		r.logger.Error("Failed to get links to prompt", "error", err, "id", promptID)
		return nil, fmt.Errorf("failed to get links to prompt: %w", err)
	}

	r.logger.Debug("Links to prompt retrieved successfully", "id", promptID, "count", len(links))
	return links, nil
}

// AddTag adds a tag to a prompt
func (r *promptRepository) AddTag(ctx context.Context, promptID, tagName string) error {
	r.logger.Debug("Adding tag to prompt", "id", promptID, "tag", tagName)

	query := `INSERT OR IGNORE INTO prompt_tags (prompt_id, tag_name) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, promptID, tagName)
	if err != nil {
		r.logger.Error("Failed to add tag to prompt", "error", err, "id", promptID, "tag", tagName)
		return fmt.Errorf("failed to add tag to prompt: %w", err)
	}

	r.logger.Info("Tag added to prompt successfully", "id", promptID, "tag", tagName)
	return nil
}

// RemoveTag removes a tag from a prompt
func (r *promptRepository) RemoveTag(ctx context.Context, promptID, tagName string) error {
	r.logger.Debug("Removing tag from prompt", "id", promptID, "tag", tagName)

	query := `DELETE FROM prompt_tags WHERE prompt_id = ? AND tag_name = ?`
	result, err := r.db.ExecContext(ctx, query, promptID, tagName)
	if err != nil {
		r.logger.Error("Failed to remove tag from prompt", "error", err, "id", promptID, "tag", tagName)
		return fmt.Errorf("failed to remove tag from prompt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "error", err, "id", promptID, "tag", tagName)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.Debug("Tag not found for removal", "id", promptID, "tag", tagName)
		return fmt.Errorf("tag not found on prompt")
	}

	r.logger.Info("Tag removed from prompt successfully", "id", promptID, "tag", tagName)
	return nil
}

// GetTags retrieves all tags for a prompt
func (r *promptRepository) GetTags(ctx context.Context, promptID string) ([]string, error) {
	r.logger.Debug("Getting tags for prompt", "id", promptID)

	query := `SELECT tag_name FROM prompt_tags WHERE prompt_id = ? ORDER BY tag_name`
	var tags []string
	err := r.db.SelectContext(ctx, &tags, query, promptID)
	if err != nil {
		r.logger.Error("Failed to get tags for prompt", "error", err, "id", promptID)
		return nil, fmt.Errorf("failed to get tags for prompt: %w", err)
	}

	r.logger.Debug("Tags for prompt retrieved successfully", "id", promptID, "count", len(tags))
	return tags, nil
}

// ListAllTags retrieves all unique tags used by prompts
func (r *promptRepository) ListAllTags(ctx context.Context) ([]string, error) {
	r.logger.Debug("Listing all prompt tags")

	query := `SELECT DISTINCT tag_name FROM prompt_tags ORDER BY tag_name`
	var tags []string
	err := r.db.SelectContext(ctx, &tags, query)
	if err != nil {
		r.logger.Error("Failed to list all prompt tags", "error", err)
		return nil, fmt.Errorf("failed to list all prompt tags: %w", err)
	}

	r.logger.Debug("All prompt tags listed successfully", "count", len(tags))
	return tags, nil
}
