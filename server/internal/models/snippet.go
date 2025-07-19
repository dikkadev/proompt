package models

import (
	"time"
)

type Snippet struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	GitRef      *string   `json:"git_ref" db:"git_ref"`
}

type SnippetTag struct {
	SnippetID string `json:"snippet_id" db:"snippet_id"`
	TagName   string `json:"tag_name" db:"tag_name"`
}
