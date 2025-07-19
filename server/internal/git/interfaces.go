package git

import (
	"context"
	"time"

	"github.com/dikkadev/proompt/server/internal/models"
)

// GitService handles git operations for prompt and snippet versioning
type GitService interface {
	// Repository initialization
	InitializeRepo(ctx context.Context) error

	// Prompt operations
	CreatePromptBranch(ctx context.Context, prompt *models.Prompt, userNote string) error
	UpdatePromptBranch(ctx context.Context, prompt *models.Prompt, userNote string) error
	DeletePromptBranch(ctx context.Context, promptID string) error

	// Snippet operations
	CreateSnippetBranch(ctx context.Context, snippet *models.Snippet, userNote string) error
	UpdateSnippetBranch(ctx context.Context, snippet *models.Snippet, userNote string) error
	DeleteSnippetBranch(ctx context.Context, snippetID string) error

	// History and versioning
	GetPromptHistory(ctx context.Context, promptID string) ([]GitCommit, error)
	GetSnippetHistory(ctx context.Context, snippetID string) ([]GitCommit, error)
	GetPromptVersion(ctx context.Context, promptID string, commitHash string) (*models.Prompt, error)
	GetSnippetVersion(ctx context.Context, snippetID string, commitHash string) (*models.Snippet, error)

	// Repository health
	ValidateRepo(ctx context.Context) error
}

// GitCommit represents a git commit with metadata
type GitCommit struct {
	Hash      string    `json:"hash"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
	Body      string    `json:"body,omitempty"`
}

// PromptContent represents the content stored in git for a prompt
type PromptContent struct {
	ID                 string             `json:"id"`
	Title              string             `json:"title"`
	Content            string             `json:"content"`
	Type               string             `json:"type"`
	UseCase            string             `json:"use_case"`
	ModelCompatibility models.StringSlice `json:"model_compatibility"`
	Parameters         models.JSONMap     `json:"parameters"`
	Variables          models.StringSlice `json:"variables"`
	Tags               models.StringSlice `json:"tags"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

// SnippetContent represents the content stored in git for a snippet
type SnippetContent struct {
	ID        string             `json:"id"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	Variables models.StringSlice `json:"variables"`
	Tags      models.StringSlice `json:"tags"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
