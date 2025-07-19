package repository

import (
	"context"

	"github.com/dikkadev/proompt/server/internal/models"
)

// Repository provides access to all data repositories with transaction support
type Repository interface {
	Prompts() PromptRepository
	Snippets() SnippetRepository
	Notes() NoteRepository

	// WithTx executes a function within a database transaction
	// If the function returns an error, the transaction is rolled back
	WithTx(ctx context.Context, fn func(Repository) error) error

	// Close closes the repository and releases resources
	Close() error
}

// PromptRepository handles CRUD operations for prompts
type PromptRepository interface {
	Create(ctx context.Context, prompt *models.Prompt) error
	GetByID(ctx context.Context, id string) (*models.Prompt, error)
	Update(ctx context.Context, prompt *models.Prompt) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters PromptFilters) ([]*models.Prompt, error)
	Search(ctx context.Context, query string) ([]*models.Prompt, error)
}

// SnippetRepository handles CRUD operations for snippets
type SnippetRepository interface {
	Create(ctx context.Context, snippet *models.Snippet) error
	GetByID(ctx context.Context, id string) (*models.Snippet, error)
	Update(ctx context.Context, snippet *models.Snippet) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters SnippetFilters) ([]*models.Snippet, error)
	Search(ctx context.Context, query string) ([]*models.Snippet, error)
}

// NoteRepository handles CRUD operations for notes
type NoteRepository interface {
	Create(ctx context.Context, note *models.Note) error
	GetByID(ctx context.Context, id string) (*models.Note, error)
	Update(ctx context.Context, note *models.Note) error
	Delete(ctx context.Context, id string) error
	ListByPromptID(ctx context.Context, promptID string) ([]*models.Note, error)
	Search(ctx context.Context, query string) ([]*models.Note, error)
}

// PromptFilters defines filtering options for prompt queries
type PromptFilters struct {
	Type          *string
	UseCase       *string
	Tags          []string
	HasVariables  *bool
	CreatedAfter  *string
	CreatedBefore *string
	UpdatedAfter  *string
	UpdatedBefore *string
	Limit         *int
	Offset        *int
}

// SnippetFilters defines filtering options for snippet queries
type SnippetFilters struct {
	Tags          []string
	HasVariables  *bool
	CreatedAfter  *string
	CreatedBefore *string
	UpdatedAfter  *string
	UpdatedBefore *string
	Limit         *int
	Offset        *int
}
