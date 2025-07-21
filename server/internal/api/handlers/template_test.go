package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dikkadev/proompt/server/internal/api/models"
	domainModels "github.com/dikkadev/proompt/server/internal/models"
	"github.com/dikkadev/proompt/server/internal/repository"
)

// mockSnippetRepository implements SnippetRepository for testing
type mockSnippetRepository struct {
	snippets map[string]*domainModels.Snippet
}

func newMockSnippetRepository() *mockSnippetRepository {
	return &mockSnippetRepository{
		snippets: make(map[string]*domainModels.Snippet),
	}
}

func (m *mockSnippetRepository) Create(ctx context.Context, snippet *domainModels.Snippet) error {
	m.snippets[snippet.ID] = snippet
	return nil
}

func (m *mockSnippetRepository) GetByID(ctx context.Context, id string) (*domainModels.Snippet, error) {
	snippet, exists := m.snippets[id]
	if !exists {
		return nil, ErrNotFound
	}
	return snippet, nil
}

func (m *mockSnippetRepository) Update(ctx context.Context, snippet *domainModels.Snippet) error {
	if _, exists := m.snippets[snippet.ID]; !exists {
		return ErrNotFound
	}
	m.snippets[snippet.ID] = snippet
	return nil
}

func (m *mockSnippetRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.snippets[id]; !exists {
		return ErrNotFound
	}
	delete(m.snippets, id)
	return nil
}

func (m *mockSnippetRepository) List(ctx context.Context, filters repository.SnippetFilters) ([]*domainModels.Snippet, error) {
	var result []*domainModels.Snippet
	for _, snippet := range m.snippets {
		result = append(result, snippet)
	}
	return result, nil
}

func (m *mockSnippetRepository) Search(ctx context.Context, query string) ([]*domainModels.Snippet, error) {
	return nil, nil // Not implemented for tests
}

func (m *mockSnippetRepository) AddTag(ctx context.Context, snippetID, tagName string) error {
	return nil // Not implemented for tests
}

func (m *mockSnippetRepository) RemoveTag(ctx context.Context, snippetID, tagName string) error {
	return nil // Not implemented for tests
}

func (m *mockSnippetRepository) GetTags(ctx context.Context, snippetID string) ([]string, error) {
	return nil, nil // Not implemented for tests
}

func (m *mockSnippetRepository) ListAllTags(ctx context.Context) ([]string, error) {
	return nil, nil // Not implemented for tests
}

// mockTemplateRepository implements Repository for template testing
type mockTemplateRepository struct {
	snippets *mockSnippetRepository
}

func newMockTemplateRepository() *mockTemplateRepository {
	return &mockTemplateRepository{
		snippets: newMockSnippetRepository(),
	}
}

func (m *mockTemplateRepository) Prompts() repository.PromptRepository {
	return nil // Not needed for template tests
}

func (m *mockTemplateRepository) Snippets() repository.SnippetRepository {
	return m.snippets
}

func (m *mockTemplateRepository) Notes() repository.NoteRepository {
	return nil // Not needed for template tests
}

func (m *mockTemplateRepository) WithTx(ctx context.Context, fn func(repository.Repository) error) error {
	return fn(m) // Simple implementation for tests
}

func (m *mockTemplateRepository) Close() error {
	return nil
}

func TestTemplatePreview(t *testing.T) {
	repo := newMockTemplateRepository()

	// Add test snippets
	greeting := &domainModels.Snippet{
		ID:      "greeting",
		Title:   "greeting",
		Content: "Hello {{name:World}}!",
	}
	repo.snippets.snippets["greeting"] = greeting

	signature := &domainModels.Snippet{
		ID:      "signature",
		Title:   "signature",
		Content: "Best regards,\n{{author}}",
	}
	repo.snippets.snippets["signature"] = signature

	handler := NewTemplateHandler(repo)

	tests := []struct {
		name           string
		requestBody    models.TemplatePreviewRequest
		expectedStatus int
		checkResponse  func(t *testing.T, response models.TemplatePreviewResponse)
	}{
		{
			name: "simple variable resolution",
			requestBody: models.TemplatePreviewRequest{
				Content: "Hello {{name}}!",
				Variables: map[string]string{
					"name": "Alice",
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response models.TemplatePreviewResponse) {
				if response.ResolvedContent != "Hello Alice!" {
					t.Errorf("Expected 'Hello Alice!', got %s", response.ResolvedContent)
				}
				if len(response.Warnings) != 0 {
					t.Errorf("Expected no warnings, got %v", response.Warnings)
				}
			},
		},
		{
			name: "variable with default",
			requestBody: models.TemplatePreviewRequest{
				Content:   "Hello {{name:World}}!",
				Variables: map[string]string{},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response models.TemplatePreviewResponse) {
				if response.ResolvedContent != "Hello World!" {
					t.Errorf("Expected 'Hello World!', got %s", response.ResolvedContent)
				}
			},
		},
		{
			name: "missing variable warning",
			requestBody: models.TemplatePreviewRequest{
				Content:   "Hello {{name}}!",
				Variables: map[string]string{},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response models.TemplatePreviewResponse) {
				if response.ResolvedContent != "Hello {{name}}!" {
					t.Errorf("Expected 'Hello {{name}}!', got %s", response.ResolvedContent)
				}
				if len(response.Warnings) != 1 {
					t.Errorf("Expected 1 warning, got %d", len(response.Warnings))
				}
			},
		},
		{
			name: "snippet insertion",
			requestBody: models.TemplatePreviewRequest{
				Content: "@greeting How are you?",
				Variables: map[string]string{
					"name": "Bob",
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response models.TemplatePreviewResponse) {
				if response.ResolvedContent != "Hello Bob! How are you?" {
					t.Errorf("Expected 'Hello Bob! How are you?', got %s", response.ResolvedContent)
				}
			},
		},
		{
			name: "snippet with missing variable",
			requestBody: models.TemplatePreviewRequest{
				Content:   "@signature",
				Variables: map[string]string{},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, response models.TemplatePreviewResponse) {
				if response.ResolvedContent != "Best regards,\n{{author}}" {
					t.Errorf("Expected 'Best regards,\\n{{author}}', got %s", response.ResolvedContent)
				}
				if len(response.Warnings) != 1 {
					t.Errorf("Expected 1 warning, got %d", len(response.Warnings))
				}
			},
		},
		{
			name: "invalid JSON",
			requestBody: models.TemplatePreviewRequest{
				Content: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/template/preview", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.PreviewTemplate(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response body for successful requests
			if tt.expectedStatus == http.StatusOK && tt.checkResponse != nil {
				var response models.TemplatePreviewResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestTemplateAnalyze(t *testing.T) {
	repo := newMockTemplateRepository()

	// Add test snippet
	greeting := &domainModels.Snippet{
		ID:      "greeting",
		Title:   "greeting",
		Content: "Hello {{name:World}}!",
	}
	repo.snippets.snippets["greeting"] = greeting

	handler := NewTemplateHandler(repo)

	requestBody := models.TemplatePreviewRequest{
		Content: "@greeting How are you, {{user}}?",
		Variables: map[string]string{
			"name": "Alice",
		},
	}

	// Create request
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/template/analyze", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.AnalyzeTemplate(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check response
	var response models.TemplatePreviewResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Should have snippet inserted but variables not resolved
	expectedContent := "Hello {{name:World}}! How are you, {{user}}?"
	if response.ResolvedContent != expectedContent {
		t.Errorf("Expected %s, got %s", expectedContent, response.ResolvedContent)
	}

	// Should have variable information
	if len(response.Variables) == 0 {
		t.Error("Expected variables in response")
	}

	// Check variable status
	foundName := false
	foundUser := false
	for _, v := range response.Variables {
		if v.Name == "name" {
			foundName = true
			if v.Status != "provided" {
				t.Errorf("Expected name variable status 'provided', got %s", v.Status)
			}
		}
		if v.Name == "user" {
			foundUser = true
			if v.Status != "missing" {
				t.Errorf("Expected user variable status 'missing', got %s", v.Status)
			}
		}
	}

	if !foundName {
		t.Error("Expected to find 'name' variable")
	}
	if !foundUser {
		t.Error("Expected to find 'user' variable")
	}
}
