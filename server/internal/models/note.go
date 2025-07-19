package models

import (
	"time"
)

type Note struct {
	ID        string    `json:"id" db:"id"`
	PromptID  string    `json:"prompt_id" db:"prompt_id"`
	Title     string    `json:"title" db:"title"`
	Body      *string   `json:"body" db:"body"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
