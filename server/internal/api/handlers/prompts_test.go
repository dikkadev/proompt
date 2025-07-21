package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dikkadev/proompt/server/internal/api/models"
	domainModels "github.com/dikkadev/proompt/server/internal/models"
	"github.com/dikkadev/proompt/server/internal/repository"
)

var ErrNotFound = errors.New("not found")

// mockPromptRepository implements PromptRepository for testing
type mockPromptRepository struct {
	prompts map[string]*domainModels.Prompt
}

func newMockPromptRepository() *mockPromptRepository {
	return &mockPromptRepository{
		prompts: make(map[string]*domainModels.Prompt),
	}
}

func (m *mockPromptRepository) Create(ctx context.Context, prompt *domainModels.Prompt) error {
	m.prompts[prompt.ID] = prompt
	return nil
}

func (m *mockPromptRepository) GetByID(ctx context.Context, id string) (*domainModels.Prompt, error) {
	prompt, exists := m.prompts[id]
	if !exists {
		return nil, ErrNotFound
	}
	return prompt, nil
}

func (m *mockPromptRepository) Update(ctx context.Context, prompt *domainModels.Prompt) error {
	if _, exists := m.prompts[prompt.ID]; !exists {
		return ErrNotFound
	}
	m.prompts[prompt.ID] = prompt
	return nil
}

func (m *mockPromptRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.prompts[id]; !exists {
		return ErrNotFound
	}
	delete(m.prompts, id)
	return nil
}

func (m *mockPromptRepository) List(ctx context.Context, filters repository.PromptFilters) ([]*domainModels.Prompt, error) {
	var result []*domainModels.Prompt
	for _, prompt := range m.prompts {
		result = append(result, prompt)
	}
	return result, nil
}

func (m *mockPromptRepository) Search(ctx context.Context, query string) ([]*domainModels.Prompt, error) {
	return nil, nil // Not implemented for tests
}

func (m *mockPromptRepository) CreateLink(ctx context.Context, link *domainModels.PromptLink) error {
	return nil // Not implemented for tests
}

func (m *mockPromptRepository) DeleteLink(ctx context.Context, fromPromptID, toPromptID string) error {
	return nil // Not implemented for tests
}

func (m *mockPromptRepository) GetLinksFrom(ctx context.Context, promptID string) ([]*domainModels.PromptLink, error) {
	return nil, nil // Not implemented for tests
}

func (m *mockPromptRepository) GetLinksTo(ctx context.Context, promptID string) ([]*domainModels.PromptLink, error) {
	return nil, nil // Not implemented for tests
}

func (m *mockPromptRepository) AddTag(ctx context.Context, promptID, tagName string) error {
	return nil // Not implemented for tests
}

func (m *mockPromptRepository) RemoveTag(ctx context.Context, promptID, tagName string) error {
	return nil // Not implemented for tests
}

func (m *mockPromptRepository) GetTags(ctx context.Context, promptID string) ([]string, error) {
	return nil, nil // Not implemented for tests
}

func (m *mockPromptRepository) ListAllTags(ctx context.Context) ([]string, error) {
	return nil, nil // Not implemented for tests
}

// mockRepository implements Repository for testing
type mockRepository struct {
	prompts *mockPromptRepository
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		prompts: newMockPromptRepository(),
	}
}

func (m *mockRepository) Prompts() repository.PromptRepository {
	return m.prompts
}

func (m *mockRepository) Snippets() repository.SnippetRepository {
	return nil // Not needed for prompt tests
}

func (m *mockRepository) Notes() repository.NoteRepository {
	return nil // Not needed for prompt tests
}

func (m *mockRepository) WithTx(ctx context.Context, fn func(repository.Repository) error) error {
	return fn(m) // Simple implementation for tests
}

func (m *mockRepository) Close() error {
	return nil
}

func TestCreatePrompt(t *testing.T) {
	repo := newMockRepository()
	handlers := NewPromptHandlers(repo)

	// Test valid request
	reqBody := models.CreatePromptRequest{
		Title:   "Test Prompt",
		Content: "This is a test prompt",
		Type:    "user",
		UseCase: "testing",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/prompts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlers.CreatePrompt(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response models.PromptResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Title != reqBody.Title {
		t.Errorf("Expected title %s, got %s", reqBody.Title, response.Title)
	}
	if response.Content != reqBody.Content {
		t.Errorf("Expected content %s, got %s", reqBody.Content, response.Content)
	}
}

func TestCreatePromptInvalidJSON(t *testing.T) {
	repo := newMockRepository()
	handlers := NewPromptHandlers(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/prompts", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlers.CreatePrompt(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreatePromptMissingTitle(t *testing.T) {
	repo := newMockRepository()
	handlers := NewPromptHandlers(repo)

	reqBody := models.CreatePromptRequest{
		Content: "This is a test prompt",
		Type:    "user",
		UseCase: "testing",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/prompts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlers.CreatePrompt(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetPrompt(t *testing.T) {
	repo := newMockRepository()
	handlers := NewPromptHandlers(repo)

	// Create a test prompt
	prompt := &domainModels.Prompt{
		ID:      "test-id",
		Title:   "Test Prompt",
		Content: "Test content",
		Type:    domainModels.PromptTypeUser,
	}
	repo.prompts.Create(context.Background(), prompt)

	req := httptest.NewRequest(http.MethodGet, "/api/prompts/test-id", nil)
	req.SetPathValue("id", "test-id")
	w := httptest.NewRecorder()

	handlers.GetPrompt(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.PromptResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != prompt.ID {
		t.Errorf("Expected ID %s, got %s", prompt.ID, response.ID)
	}
}

func TestGetPromptNotFound(t *testing.T) {
	repo := newMockRepository()
	handlers := NewPromptHandlers(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/prompts/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	handlers.GetPrompt(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestListPrompts(t *testing.T) {
	repo := newMockRepository()
	handlers := NewPromptHandlers(repo)

	// Create test prompts
	prompt1 := &domainModels.Prompt{
		ID:      "test-id-1",
		Title:   "Test Prompt 1",
		Content: "Test content 1",
		Type:    domainModels.PromptTypeUser,
	}
	prompt2 := &domainModels.Prompt{
		ID:      "test-id-2",
		Title:   "Test Prompt 2",
		Content: "Test content 2",
		Type:    domainModels.PromptTypeSystem,
	}
	repo.prompts.Create(context.Background(), prompt1)
	repo.prompts.Create(context.Background(), prompt2)

	req := httptest.NewRequest(http.MethodGet, "/api/prompts", nil)
	w := httptest.NewRecorder()

	handlers.ListPrompts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.ListResponse[*models.PromptResponse]
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 prompts, got %d", len(response.Data))
	}
}
