package models

import (
	"time"

	"github.com/dikkadev/proompt/server/internal/models"
)

// PromptResponse represents a prompt in API responses
type PromptResponse struct {
	ID                     string         `json:"id"`
	Title                  string         `json:"title"`
	Content                string         `json:"content"`
	Type                   string         `json:"type"`
	UseCase                *string        `json:"use_case"`
	ModelCompatibilityTags []string       `json:"model_compatibility_tags"`
	TemperatureSuggestion  *float64       `json:"temperature_suggestion"`
	OtherParameters        map[string]any `json:"other_parameters"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	GitRef                 *string        `json:"git_ref"`
}

// SnippetResponse represents a snippet in API responses
type SnippetResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GitRef      *string   `json:"git_ref"`
}

// NoteResponse represents a note in API responses
type NoteResponse struct {
	ID        string    `json:"id"`
	PromptID  string    `json:"prompt_id"`
	Title     string    `json:"title"`
	Body      *string   `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListResponse represents a paginated list response
type ListResponse[T any] struct {
	Data       []T `json:"data"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// FromPrompt converts domain model to API response
func FromPrompt(p *models.Prompt) *PromptResponse {
	var modelTags []string
	if p.ModelCompatibilityTags != nil {
		modelTags = []string(p.ModelCompatibilityTags)
	}

	var otherParams map[string]any
	if p.OtherParameters != nil {
		otherParams = map[string]any(p.OtherParameters)
	}

	return &PromptResponse{
		ID:                     p.ID,
		Title:                  p.Title,
		Content:                p.Content,
		Type:                   string(p.Type),
		UseCase:                p.UseCase,
		ModelCompatibilityTags: modelTags,
		TemperatureSuggestion:  p.TemperatureSuggestion,
		OtherParameters:        otherParams,
		CreatedAt:              p.CreatedAt,
		UpdatedAt:              p.UpdatedAt,
		GitRef:                 p.GitRef,
	}
}

// FromSnippet converts domain model to API response
func FromSnippet(s *models.Snippet) *SnippetResponse {
	return &SnippetResponse{
		ID:          s.ID,
		Title:       s.Title,
		Content:     s.Content,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		GitRef:      s.GitRef,
	}
}

// FromNote converts domain model to API response
func FromNote(n *models.Note) *NoteResponse {
	return &NoteResponse{
		ID:        n.ID,
		PromptID:  n.PromptID,
		Title:     n.Title,
		Body:      n.Body,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

// FromPrompts converts slice of domain models to API responses
func FromPrompts(prompts []*models.Prompt) []*PromptResponse {
	responses := make([]*PromptResponse, len(prompts))
	for i, p := range prompts {
		responses[i] = FromPrompt(p)
	}
	return responses
}

// FromSnippets converts slice of domain models to API responses
func FromSnippets(snippets []*models.Snippet) []*SnippetResponse {
	responses := make([]*SnippetResponse, len(snippets))
	for i, s := range snippets {
		responses[i] = FromSnippet(s)
	}
	return responses
}

// FromNotes converts slice of domain models to API responses
func FromNotes(notes []*models.Note) []*NoteResponse {
	responses := make([]*NoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = FromNote(n)
	}
	return responses
}
