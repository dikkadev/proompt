package template

import (
	"reflect"
	"testing"

	"github.com/dikkadev/proompt/server/internal/models"
)

func TestSnippetResolver_InsertSnippets(t *testing.T) {
	snippets := []*models.Snippet{
		{
			Title:   "greeting",
			Content: "Hello {{name:World}}!",
		},
		{
			Title:   "signature",
			Content: "Best regards,\n{{author}}",
		},
		{
			Title:   "nested",
			Content: "@greeting How are you?",
		},
	}

	tests := []struct {
		name             string
		content          string
		variables        map[string]string
		expectedContent  string
		expectedWarnings int
	}{
		{
			name:             "no snippets",
			content:          "Hello world",
			variables:        map[string]string{},
			expectedContent:  "Hello world",
			expectedWarnings: 0,
		},
		{
			name:             "simple snippet insertion",
			content:          "@greeting",
			variables:        map[string]string{},
			expectedContent:  "Hello {{name:World}}!",
			expectedWarnings: 0,
		},
		{
			name:             "snippet with spaces",
			content:          "@{greeting}",
			variables:        map[string]string{},
			expectedContent:  "Hello {{name:World}}!",
			expectedWarnings: 0,
		},
		{
			name:             "snippet not found",
			content:          "@unknown",
			variables:        map[string]string{},
			expectedContent:  "@unknown",
			expectedWarnings: 1,
		},
		{
			name:             "multiple snippets",
			content:          "@greeting\n\n@signature",
			variables:        map[string]string{},
			expectedContent:  "Hello {{name:World}}!\n\nBest regards,\n{{author}}",
			expectedWarnings: 0,
		},
		{
			name:             "nested snippet",
			content:          "@nested",
			variables:        map[string]string{},
			expectedContent:  "Hello {{name:World}}! How are you?",
			expectedWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewSnippetResolver(snippets, tt.variables)
			result := resolver.InsertSnippets(tt.content)

			if result.Content != tt.expectedContent {
				t.Errorf("InsertSnippets() content = %v, want %v", result.Content, tt.expectedContent)
			}

			if len(result.Warnings) != tt.expectedWarnings {
				t.Errorf("InsertSnippets() warnings count = %v, want %v", len(result.Warnings), tt.expectedWarnings)
			}
		})
	}
}

func TestSnippetResolver_ResolveWithSnippets(t *testing.T) {
	snippets := []*models.Snippet{
		{
			Title:   "greeting",
			Content: "Hello {{name:World}}!",
		},
	}

	tests := []struct {
		name             string
		content          string
		variables        map[string]string
		expectedContent  string
		expectedWarnings int
	}{
		{
			name:    "snippet with variable resolution",
			content: "@greeting",
			variables: map[string]string{
				"name": "Alice",
			},
			expectedContent:  "Hello Alice!",
			expectedWarnings: 0,
		},
		{
			name:             "snippet with default variable",
			content:          "@greeting",
			variables:        map[string]string{},
			expectedContent:  "Hello World!",
			expectedWarnings: 0,
		},
		{
			name:    "mixed content with snippets and variables",
			content: "Hi there! @greeting\n\nSigned by {{author}}",
			variables: map[string]string{
				"name":   "Bob",
				"author": "System",
			},
			expectedContent:  "Hi there! Hello Bob!\n\nSigned by System",
			expectedWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewSnippetResolver(snippets, tt.variables)
			result := resolver.ResolveWithSnippets(tt.content)

			if result.Content != tt.expectedContent {
				t.Errorf("ResolveWithSnippets() content = %v, want %v", result.Content, tt.expectedContent)
			}

			if len(result.Warnings) != tt.expectedWarnings {
				t.Errorf("ResolveWithSnippets() warnings count = %v, want %v", len(result.Warnings), tt.expectedWarnings)
			}
		})
	}
}

func TestSnippetResolver_GetAllVariables(t *testing.T) {
	snippets := []*models.Snippet{
		{
			Title:   "greeting",
			Content: "Hello {{name:World}}!",
		},
		{
			Title:   "signature",
			Content: "Best regards,\n{{author}}",
		},
	}

	resolver := NewSnippetResolver(snippets, map[string]string{})

	tests := []struct {
		name     string
		content  string
		expected []string // Just variable names for simplicity
	}{
		{
			name:     "no variables",
			content:  "Hello world",
			expected: []string{},
		},
		{
			name:     "content variables only",
			content:  "Hello {{user}}",
			expected: []string{"user"},
		},
		{
			name:     "snippet variables only",
			content:  "@greeting",
			expected: []string{"name"},
		},
		{
			name:     "mixed variables",
			content:  "{{intro}} @greeting @signature",
			expected: []string{"intro", "name", "author"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.GetAllVariables(tt.content)

			var resultNames []string
			for _, v := range result {
				resultNames = append(resultNames, v.Name)
			}

			// Sort both slices for comparison
			if len(resultNames) != len(tt.expected) {
				t.Errorf("GetAllVariables() count = %v, want %v", len(resultNames), len(tt.expected))
				return
			}

			// Check if all expected variables are present
			expectedMap := make(map[string]bool)
			for _, name := range tt.expected {
				expectedMap[name] = true
			}

			for _, name := range resultNames {
				if !expectedMap[name] {
					t.Errorf("GetAllVariables() unexpected variable: %v", name)
				}
				delete(expectedMap, name)
			}

			if len(expectedMap) > 0 {
				t.Errorf("GetAllVariables() missing variables: %v", expectedMap)
			}
		})
	}
}

func TestSnippetResolver_GetVariableStatusWithSnippets(t *testing.T) {
	snippets := []*models.Snippet{
		{
			Title:   "greeting",
			Content: "Hello {{name:World}}!",
		},
		{
			Title:   "signature",
			Content: "Best regards,\n{{author}}",
		},
	}

	resolver := NewSnippetResolver(snippets, map[string]string{
		"name": "Alice",
	})

	content := "{{intro}} @greeting @signature"
	status := resolver.GetVariableStatusWithSnippets(content)

	expected := map[string]string{
		"intro":  "missing",
		"name":   "provided",
		"author": "missing",
	}

	if !reflect.DeepEqual(status, expected) {
		t.Errorf("GetVariableStatusWithSnippets() = %v, want %v", status, expected)
	}
}
