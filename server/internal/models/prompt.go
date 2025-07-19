package models

import (
	"time"
)

type Prompt struct {
	ID                     string      `json:"id" db:"id"`
	Title                  string      `json:"title" db:"title"`
	Content                string      `json:"content" db:"content"`
	Type                   PromptType  `json:"type" db:"type"`
	UseCase                *string     `json:"use_case" db:"use_case"`
	ModelCompatibilityTags StringSlice `json:"model_compatibility_tags" db:"model_compatibility_tags"`
	TemperatureSuggestion  *float64    `json:"temperature_suggestion" db:"temperature_suggestion"`
	OtherParameters        JSONMap     `json:"other_parameters" db:"other_parameters"`
	CreatedAt              time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at" db:"updated_at"`
	GitRef                 *string     `json:"git_ref" db:"git_ref"`
}

type PromptTag struct {
	PromptID string `json:"prompt_id" db:"prompt_id"`
	TagName  string `json:"tag_name" db:"tag_name"`
}

type PromptLink struct {
	FromPromptID string    `json:"from_prompt_id" db:"from_prompt_id"`
	ToPromptID   string    `json:"to_prompt_id" db:"to_prompt_id"`
	LinkType     string    `json:"link_type" db:"link_type"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
