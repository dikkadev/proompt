package models

import (
	"time"

	"github.com/dikkadev/proompt/server/internal/models"
)

// CreatePromptRequest represents the request body for creating a prompt
type CreatePromptRequest struct {
	Title                  string         `json:"title" validate:"required,min=1,max=255"`
	Content                string         `json:"content" validate:"required"`
	Type                   string         `json:"type" validate:"required,oneof=system user image video"`
	UseCase                string         `json:"use_case" validate:"required,min=1,max=100"`
	ModelCompatibilityTags []string       `json:"model_compatibility_tags,omitempty"`
	TemperatureSuggestion  *float64       `json:"temperature_suggestion,omitempty" validate:"omitempty,min=0,max=2"`
	OtherParameters        map[string]any `json:"other_parameters,omitempty"`
	Notes                  string         `json:"notes,omitempty"`
}

// UpdatePromptRequest represents the request body for updating a prompt
type UpdatePromptRequest struct {
	Title                  *string        `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content                *string        `json:"content,omitempty" validate:"omitempty,min=1"`
	Type                   *string        `json:"type,omitempty" validate:"omitempty,oneof=system user image video"`
	UseCase                *string        `json:"use_case,omitempty" validate:"omitempty,min=1,max=100"`
	ModelCompatibilityTags []string       `json:"model_compatibility_tags,omitempty"`
	TemperatureSuggestion  *float64       `json:"temperature_suggestion,omitempty" validate:"omitempty,min=0,max=2"`
	OtherParameters        map[string]any `json:"other_parameters,omitempty"`
	Notes                  *string        `json:"notes,omitempty"`
}

// CreateSnippetRequest represents the request body for creating a snippet
type CreateSnippetRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Content     string `json:"content" validate:"required"`
	Description string `json:"description,omitempty"`
}

// UpdateSnippetRequest represents the request body for updating a snippet
type UpdateSnippetRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content     *string `json:"content,omitempty" validate:"omitempty,min=1"`
	Description *string `json:"description,omitempty"`
}

// CreateNoteRequest represents the request body for creating a note
type CreateNoteRequest struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
	Body  string `json:"body" validate:"required"`
}

// UpdateNoteRequest represents the request body for updating a note
type UpdateNoteRequest struct {
	Title *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Body  *string `json:"body,omitempty" validate:"omitempty,min=1"`
}

// ToPrompt converts CreatePromptRequest to domain model
func (r *CreatePromptRequest) ToPrompt() *models.Prompt {
	var useCase *string
	if r.UseCase != "" {
		useCase = &r.UseCase
	}

	return &models.Prompt{
		Title:                  r.Title,
		Content:                r.Content,
		Type:                   models.PromptType(r.Type),
		UseCase:                useCase,
		ModelCompatibilityTags: models.StringSlice(r.ModelCompatibilityTags),
		TemperatureSuggestion:  r.TemperatureSuggestion,
		OtherParameters:        models.JSONMap(r.OtherParameters),
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
}

// ToSnippet converts CreateSnippetRequest to domain model
func (r *CreateSnippetRequest) ToSnippet() *models.Snippet {
	var description *string
	if r.Description != "" {
		description = &r.Description
	}

	return &models.Snippet{
		Title:       r.Title,
		Content:     r.Content,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ToNote converts CreateNoteRequest to domain model
func (r *CreateNoteRequest) ToNote(promptID string) *models.Note {
	var body *string
	if r.Body != "" {
		body = &r.Body
	}

	return &models.Note{
		PromptID:  promptID,
		Title:     r.Title,
		Body:      body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// TemplatePreviewRequest represents the request body for template preview
type TemplatePreviewRequest struct {
	Content   string            `json:"content" validate:"required"`
	Variables map[string]string `json:"variables,omitempty"`
}

// CreatePromptLinkRequest represents the request body for creating a prompt link
type CreatePromptLinkRequest struct {
	ToPromptID string `json:"to_prompt_id" validate:"required"`
	LinkType   string `json:"link_type,omitempty"`
}

// AddTagRequest represents the request body for adding a tag
type AddTagRequest struct {
	TagName string `json:"tag_name" validate:"required"`
}

// AnalyzeTemplateRequest represents the request body for template analysis
type AnalyzeTemplateRequest struct {
	Template string `json:"template" validate:"required"`
}
