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

func TestPromptLinks(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create two test prompts
	prompt1 := &models.Prompt{
		ID:      "prompt-1",
		Title:   "First Prompt",
		Content: "This is the first prompt",
		Type:    "user",
	}
	prompt2 := &models.Prompt{
		ID:      "prompt-2",
		Title:   "Second Prompt",
		Content: "This is the second prompt",
		Type:    "user",
	}

	ctx := context.Background()

	// Create prompts
	err := repo.Prompts().Create(ctx, prompt1)
	if err != nil {
		t.Fatalf("Failed to create first prompt: %v", err)
	}

	err = repo.Prompts().Create(ctx, prompt2)
	if err != nil {
		t.Fatalf("Failed to create second prompt: %v", err)
	}

	// Test creating a link
	link := &models.PromptLink{
		FromPromptID: prompt1.ID,
		ToPromptID:   prompt2.ID,
		LinkType:     "followup",
	}

	err = repo.Prompts().CreateLink(ctx, link)
	if err != nil {
		t.Fatalf("Failed to create link: %v", err)
	}

	// Test getting links from prompt1
	linksFrom, err := repo.Prompts().GetLinksFrom(ctx, prompt1.ID)
	if err != nil {
		t.Fatalf("Failed to get links from prompt1: %v", err)
	}

	if len(linksFrom) != 1 {
		t.Fatalf("Expected 1 link from prompt1, got %d", len(linksFrom))
	}

	if linksFrom[0].ToPromptID != prompt2.ID {
		t.Errorf("Expected link to prompt2, got %s", linksFrom[0].ToPromptID)
	}

	if linksFrom[0].LinkType != "followup" {
		t.Errorf("Expected link type 'followup', got %s", linksFrom[0].LinkType)
	}

	// Test getting links to prompt2
	linksTo, err := repo.Prompts().GetLinksTo(ctx, prompt2.ID)
	if err != nil {
		t.Fatalf("Failed to get links to prompt2: %v", err)
	}

	if len(linksTo) != 1 {
		t.Fatalf("Expected 1 link to prompt2, got %d", len(linksTo))
	}

	if linksTo[0].FromPromptID != prompt1.ID {
		t.Errorf("Expected link from prompt1, got %s", linksTo[0].FromPromptID)
	}

	// Test deleting the link
	err = repo.Prompts().DeleteLink(ctx, prompt1.ID, prompt2.ID)
	if err != nil {
		t.Fatalf("Failed to delete link: %v", err)
	}

	// Verify link is deleted
	linksFrom, err = repo.Prompts().GetLinksFrom(ctx, prompt1.ID)
	if err != nil {
		t.Fatalf("Failed to get links from prompt1 after deletion: %v", err)
	}

	if len(linksFrom) != 0 {
		t.Errorf("Expected 0 links after deletion, got %d", len(linksFrom))
	}

	// Test deleting non-existent link
	err = repo.Prompts().DeleteLink(ctx, prompt1.ID, prompt2.ID)
	if err == nil {
		t.Error("Expected error when deleting non-existent link")
	}
}

func TestPromptTags(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test prompt
	prompt := &models.Prompt{
		ID:      "prompt-1",
		Title:   "Test Prompt",
		Content: "This is a test prompt",
		Type:    "user",
	}

	ctx := context.Background()

	err := repo.Prompts().Create(ctx, prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt: %v", err)
	}

	// Test adding tags
	err = repo.Prompts().AddTag(ctx, prompt.ID, "tag1")
	if err != nil {
		t.Fatalf("Failed to add tag1: %v", err)
	}

	err = repo.Prompts().AddTag(ctx, prompt.ID, "tag2")
	if err != nil {
		t.Fatalf("Failed to add tag2: %v", err)
	}

	// Test adding duplicate tag (should not error)
	err = repo.Prompts().AddTag(ctx, prompt.ID, "tag1")
	if err != nil {
		t.Fatalf("Failed to add duplicate tag: %v", err)
	}

	// Test getting tags
	tags, err := repo.Prompts().GetTags(ctx, prompt.ID)
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("Expected 2 tags, got %d", len(tags))
	}

	expectedTags := map[string]bool{"tag1": true, "tag2": true}
	for _, tag := range tags {
		if !expectedTags[tag] {
			t.Errorf("Unexpected tag: %s", tag)
		}
	}

	// Test listing all tags
	allTags, err := repo.Prompts().ListAllTags(ctx)
	if err != nil {
		t.Fatalf("Failed to list all tags: %v", err)
	}

	if len(allTags) != 2 {
		t.Fatalf("Expected 2 total tags, got %d", len(allTags))
	}

	// Test removing a tag
	err = repo.Prompts().RemoveTag(ctx, prompt.ID, "tag1")
	if err != nil {
		t.Fatalf("Failed to remove tag1: %v", err)
	}

	// Verify tag is removed
	tags, err = repo.Prompts().GetTags(ctx, prompt.ID)
	if err != nil {
		t.Fatalf("Failed to get tags after removal: %v", err)
	}

	if len(tags) != 1 {
		t.Fatalf("Expected 1 tag after removal, got %d", len(tags))
	}

	if tags[0] != "tag2" {
		t.Errorf("Expected remaining tag to be 'tag2', got %s", tags[0])
	}

	// Test removing non-existent tag
	err = repo.Prompts().RemoveTag(ctx, prompt.ID, "nonexistent")
	if err == nil {
		t.Error("Expected error when removing non-existent tag")
	}
}

func TestSnippetTags(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test snippet
	snippet := &models.Snippet{
		ID:      "snippet-1",
		Title:   "Test Snippet",
		Content: "This is a test snippet",
	}

	ctx := context.Background()

	err := repo.Snippets().Create(ctx, snippet)
	if err != nil {
		t.Fatalf("Failed to create snippet: %v", err)
	}

	// Test adding tags
	err = repo.Snippets().AddTag(ctx, snippet.ID, "snippet-tag1")
	if err != nil {
		t.Fatalf("Failed to add snippet-tag1: %v", err)
	}

	err = repo.Snippets().AddTag(ctx, snippet.ID, "snippet-tag2")
	if err != nil {
		t.Fatalf("Failed to add snippet-tag2: %v", err)
	}

	// Test getting tags
	tags, err := repo.Snippets().GetTags(ctx, snippet.ID)
	if err != nil {
		t.Fatalf("Failed to get snippet tags: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("Expected 2 snippet tags, got %d", len(tags))
	}

	// Test listing all snippet tags
	allTags, err := repo.Snippets().ListAllTags(ctx)
	if err != nil {
		t.Fatalf("Failed to list all snippet tags: %v", err)
	}

	if len(allTags) != 2 {
		t.Fatalf("Expected 2 total snippet tags, got %d", len(allTags))
	}

	// Test removing a tag
	err = repo.Snippets().RemoveTag(ctx, snippet.ID, "snippet-tag1")
	if err != nil {
		t.Fatalf("Failed to remove snippet-tag1: %v", err)
	}

	// Verify tag is removed
	tags, err = repo.Snippets().GetTags(ctx, snippet.ID)
	if err != nil {
		t.Fatalf("Failed to get snippet tags after removal: %v", err)
	}

	if len(tags) != 1 {
		t.Fatalf("Expected 1 snippet tag after removal, got %d", len(tags))
	}
}
