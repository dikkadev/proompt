package repository

import (
	"context"
	"testing"

	"github.com/dikkadev/proompt/server/internal/config"
	"github.com/dikkadev/proompt/server/internal/db"
	"github.com/dikkadev/proompt/server/internal/git"
	"github.com/dikkadev/proompt/server/internal/models"
)

func setupTestRepo(t *testing.T) (Repository, func()) {
	// Create in-memory database
	database, err := db.NewLocal(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations("../db/migrations"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create test config
	cfg := &config.Config{
		Storage: config.Storage{
			ReposDir: "/tmp/test-repos",
		},
	}

	// Create git service
	gitService, err := git.NewGitService(cfg)
	if err != nil {
		t.Fatalf("Failed to create git service: %v", err)
	}

	// Create repository
	repo := New(database, gitService)

	// Cleanup function
	cleanup := func() {
		repo.Close()
	}

	return repo, cleanup
}

func TestPromptCRUD(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// Test Create
	prompt := &models.Prompt{
		Title:   "Test Prompt",
		Content: "This is a test prompt content",
		Type:    models.PromptTypeSystem,
	}

	err := repo.Prompts().Create(ctx, prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt: %v", err)
	}

	if prompt.ID == "" {
		t.Error("Expected prompt ID to be generated")
	}

	// Test GetByID
	retrieved, err := repo.Prompts().GetByID(ctx, prompt.ID)
	if err != nil {
		t.Fatalf("Failed to get prompt: %v", err)
	}

	if retrieved.Title != prompt.Title {
		t.Errorf("Expected title %s, got %s", prompt.Title, retrieved.Title)
	}

	// Test Update
	retrieved.Title = "Updated Test Prompt"
	err = repo.Prompts().Update(ctx, retrieved)
	if err != nil {
		t.Fatalf("Failed to update prompt: %v", err)
	}

	// Verify update
	updated, err := repo.Prompts().GetByID(ctx, prompt.ID)
	if err != nil {
		t.Fatalf("Failed to get updated prompt: %v", err)
	}

	if updated.Title != "Updated Test Prompt" {
		t.Errorf("Expected updated title, got %s", updated.Title)
	}

	// Test List
	prompts, err := repo.Prompts().List(ctx, PromptFilters{})
	if err != nil {
		t.Fatalf("Failed to list prompts: %v", err)
	}

	if len(prompts) != 1 {
		t.Errorf("Expected 1 prompt, got %d", len(prompts))
	}

	// Test Delete
	err = repo.Prompts().Delete(ctx, prompt.ID)
	if err != nil {
		t.Fatalf("Failed to delete prompt: %v", err)
	}

	// Verify deletion
	_, err = repo.Prompts().GetByID(ctx, prompt.ID)
	if err == nil {
		t.Error("Expected error when getting deleted prompt")
	}
}

func TestSnippetCRUD(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// Test Create
	snippet := &models.Snippet{
		Title:   "Test Snippet",
		Content: "This is a test snippet content",
	}

	err := repo.Snippets().Create(ctx, snippet)
	if err != nil {
		t.Fatalf("Failed to create snippet: %v", err)
	}

	if snippet.ID == "" {
		t.Error("Expected snippet ID to be generated")
	}

	// Test GetByID
	retrieved, err := repo.Snippets().GetByID(ctx, snippet.ID)
	if err != nil {
		t.Fatalf("Failed to get snippet: %v", err)
	}

	if retrieved.Title != snippet.Title {
		t.Errorf("Expected title %s, got %s", snippet.Title, retrieved.Title)
	}

	// Test Update
	retrieved.Title = "Updated Test Snippet"
	err = repo.Snippets().Update(ctx, retrieved)
	if err != nil {
		t.Fatalf("Failed to update snippet: %v", err)
	}

	// Test Delete
	err = repo.Snippets().Delete(ctx, snippet.ID)
	if err != nil {
		t.Fatalf("Failed to delete snippet: %v", err)
	}
}

func TestTransactions(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()

	// Test successful transaction
	var promptID string
	err := repo.WithTx(ctx, func(txRepo Repository) error {
		prompt := &models.Prompt{
			Title:   "Transaction Test Prompt",
			Content: "This prompt should be created",
			Type:    models.PromptTypeSystem,
		}

		if err := txRepo.Prompts().Create(ctx, prompt); err != nil {
			return err
		}

		promptID = prompt.ID
		return nil
	})

	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}

	// Verify prompt was created
	_, err = repo.Prompts().GetByID(ctx, promptID)
	if err != nil {
		t.Fatalf("Prompt not found after successful transaction: %v", err)
	}

	// Test failed transaction (rollback)
	err = repo.WithTx(ctx, func(txRepo Repository) error {
		prompt := &models.Prompt{
			Title:   "Failed Transaction Prompt",
			Content: "This prompt should not be created",
			Type:    models.PromptTypeSystem,
		}

		if err := txRepo.Prompts().Create(ctx, prompt); err != nil {
			return err
		}

		// Force an error to trigger rollback
		return context.Canceled
	})

	if err == nil {
		t.Error("Expected transaction to fail")
	}

	// Verify no extra prompts were created
	prompts, err := repo.Prompts().List(ctx, PromptFilters{})
	if err != nil {
		t.Fatalf("Failed to list prompts: %v", err)
	}

	if len(prompts) != 1 {
		t.Errorf("Expected 1 prompt after rollback, got %d", len(prompts))
	}
}
