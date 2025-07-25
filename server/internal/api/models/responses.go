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

// TemplateVariable represents a variable in template responses
type TemplateVariable struct {
	Name         string `json:"name"`
	DefaultValue string `json:"default_value,omitempty"`
	HasDefault   bool   `json:"has_default"`
	Status       string `json:"status"` // "provided", "default", "missing"
}

// TemplatePreviewResponse represents the response for template preview
type TemplatePreviewResponse struct {
	ResolvedContent string             `json:"resolved_content"`
	Variables       []TemplateVariable `json:"variables"`
	Warnings        []string           `json:"warnings"`
}

// PromptLinkResponse represents a prompt link in API responses
type PromptLinkResponse struct {
	FromPromptID string    `json:"from_prompt_id"`
	ToPromptID   string    `json:"to_prompt_id"`
	LinkType     string    `json:"link_type"`
	CreatedAt    time.Time `json:"created_at"`
}

// FromPromptLink converts domain model to API response
func FromPromptLink(link *models.PromptLink) *PromptLinkResponse {
	return &PromptLinkResponse{
		FromPromptID: link.FromPromptID,
		ToPromptID:   link.ToPromptID,
		LinkType:     link.LinkType,
		CreatedAt:    link.CreatedAt,
	}
}

// FromPromptLinks converts slice of domain models to API responses
func FromPromptLinks(links []*models.PromptLink) []*PromptLinkResponse {
	responses := make([]*PromptLinkResponse, len(links))
	for i, link := range links {
		responses[i] = FromPromptLink(link)
	}
	return responses
}

// TagResponse represents a tag in API responses
type TagResponse struct {
	Name      string    `json:"name"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

// Non-generic list response types for Swagger documentation
// (older swagger generators don't support Go generics)

// PromptListResponse represents a list of prompts
type PromptListResponse struct {
	Data       []PromptResponse `json:"data"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// SnippetListResponse represents a list of snippets
type SnippetListResponse struct {
	Data       []SnippetResponse `json:"data"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// NoteListResponse represents a list of notes
type NoteListResponse struct {
	Data       []NoteResponse `json:"data"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// TagListResponse represents a list of tags
type TagListResponse struct {
	Data       []TagResponse `json:"data"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// PromptLinkListResponse represents a list of prompt links
type PromptLinkListResponse struct {
	Data       []PromptLinkResponse `json:"data"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}
