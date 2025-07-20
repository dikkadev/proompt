package git

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/dikkadev/proompt/server/internal/config"
	"github.com/dikkadev/proompt/server/internal/logging"
	"github.com/dikkadev/proompt/server/internal/models"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/afero"
)

// gitService implements GitService interface
type gitService struct {
	config   *config.Config
	fs       afero.Fs
	logger   *slog.Logger
	repo     *git.Repository
	repoPath string
}

// NewGitService creates a new git service instance
func NewGitService(cfg *config.Config) (GitService, error) {
	logger := logging.NewLogger("git")

	service := &gitService{
		config:   cfg,
		fs:       cfg.GetFilesystem(),
		logger:   logger,
		repoPath: filepath.Join(cfg.Storage.ReposDir, "git-repo"),
	}

	if err := service.InitializeRepo(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize git repository: %w", err)
	}

	return service, nil
}

// InitializeRepo initializes the main git repository
func (s *gitService) InitializeRepo(ctx context.Context) error {
	s.logger.Debug("Initializing git repository", "path", s.repoPath)

	// Ensure directory exists
	if err := s.fs.MkdirAll(s.repoPath, 0755); err != nil {
		return fmt.Errorf("failed to create repo directory: %w", err)
	}

	// Try to open existing repository
	repo, err := git.PlainOpen(s.repoPath)
	if err != nil {
		// Repository doesn't exist, create it
		s.logger.Info("Creating new git repository", "path", s.repoPath)
		repo, err = git.PlainInit(s.repoPath, false)
		if err != nil {
			return fmt.Errorf("failed to initialize git repository: %w", err)
		}

		// Create initial commit
		if err := s.createInitialCommit(repo); err != nil {
			return fmt.Errorf("failed to create initial commit: %w", err)
		}
	}

	s.repo = repo
	s.logger.Info("Git repository initialized successfully", "path", s.repoPath)
	return nil
}

// createInitialCommit creates the initial commit in the repository
func (s *gitService) createInitialCommit(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Create a README file
	readmePath := filepath.Join(s.repoPath, "README.md")
	readmeContent := "# Proompt Git Repository\n\nThis repository contains versioned prompts and snippets.\n"

	if err := afero.WriteFile(s.fs, readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README: %w", err)
	}

	// Add and commit
	if _, err := worktree.Add("README.md"); err != nil {
		return fmt.Errorf("failed to add README: %w", err)
	}

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Proompt",
			Email: "proompt@local",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	return nil
}

// CreatePromptBranch creates a new orphan branch for a prompt
func (s *gitService) CreatePromptBranch(ctx context.Context, prompt *models.Prompt, userNote string) error {
	branchName := fmt.Sprintf("prompts/%s", prompt.ID)
	s.logger.Debug("Creating prompt branch", "branch", branchName, "title", prompt.Title)

	// Create orphan branch and commit content
	content := &PromptContent{
		ID:                 prompt.ID,
		Title:              prompt.Title,
		Content:            prompt.Content,
		Type:               string(prompt.Type),
		UseCase:            getStringValue(prompt.UseCase),
		ModelCompatibility: prompt.ModelCompatibilityTags,
		Parameters:         prompt.OtherParameters,
		Variables:          models.StringSlice{}, // TODO: Extract from content
		Tags:               models.StringSlice{}, // TODO: Get from tags table
		CreatedAt:          prompt.CreatedAt,
		UpdatedAt:          prompt.UpdatedAt,
	}

	commitMessage := fmt.Sprintf("Create: %s", prompt.Title)
	if userNote != "" {
		commitMessage += "\n\n" + userNote
	}

	if err := s.createOrphanBranchWithContent(branchName, "content.json", content, commitMessage); err != nil {
		return fmt.Errorf("failed to create orphan branch with content: %w", err)
	}

	s.logger.Info("Prompt branch created successfully", "branch", branchName, "title", prompt.Title)
	return nil
}

// UpdatePromptBranch updates an existing prompt branch
func (s *gitService) UpdatePromptBranch(ctx context.Context, prompt *models.Prompt, userNote string) error {
	branchName := fmt.Sprintf("prompts/%s", prompt.ID)
	s.logger.Debug("Updating prompt branch", "branch", branchName, "title", prompt.Title)

	// Update branch with new content
	content := &PromptContent{
		ID:                 prompt.ID,
		Title:              prompt.Title,
		Content:            prompt.Content,
		Type:               string(prompt.Type),
		UseCase:            getStringValue(prompt.UseCase),
		ModelCompatibility: prompt.ModelCompatibilityTags,
		Parameters:         prompt.OtherParameters,
		Variables:          models.StringSlice{}, // TODO: Extract from content
		Tags:               models.StringSlice{}, // TODO: Get from tags table
		CreatedAt:          prompt.CreatedAt,
		UpdatedAt:          prompt.UpdatedAt,
	}

	commitMessage := fmt.Sprintf("Update: %s", prompt.Title)
	if userNote != "" {
		commitMessage += "\n\n" + userNote
	}

	if err := s.updateBranchWithContent(branchName, "content.json", content, commitMessage); err != nil {
		return fmt.Errorf("failed to update branch with content: %w", err)
	}

	s.logger.Info("Prompt branch updated successfully", "branch", branchName, "title", prompt.Title)
	return nil
}

// DeletePromptBranch deletes a prompt branch
func (s *gitService) DeletePromptBranch(ctx context.Context, promptID string) error {
	branchName := fmt.Sprintf("prompts/%s", promptID)
	s.logger.Debug("Deleting prompt branch", "branch", branchName)

	// Delete the branch
	if err := s.repo.Storer.RemoveReference(plumbing.NewBranchReferenceName(branchName)); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	s.logger.Info("Prompt branch deleted successfully", "branch", branchName)
	return nil
}

// CreateSnippetBranch creates a new orphan branch for a snippet
func (s *gitService) CreateSnippetBranch(ctx context.Context, snippet *models.Snippet, userNote string) error {
	branchName := fmt.Sprintf("snippets/%s", snippet.ID)
	s.logger.Debug("Creating snippet branch", "branch", branchName, "title", snippet.Title)

	// Create orphan branch and commit content
	content := &SnippetContent{
		ID:        snippet.ID,
		Title:     snippet.Title,
		Content:   snippet.Content,
		Variables: models.StringSlice{}, // TODO: Extract from content
		Tags:      models.StringSlice{}, // TODO: Get from tags table
		CreatedAt: snippet.CreatedAt,
		UpdatedAt: snippet.UpdatedAt,
	}

	commitMessage := fmt.Sprintf("Create: %s", snippet.Title)
	if userNote != "" {
		commitMessage += "\n\n" + userNote
	}

	if err := s.createOrphanBranchWithContent(branchName, "content.json", content, commitMessage); err != nil {
		return fmt.Errorf("failed to create orphan branch with content: %w", err)
	}

	s.logger.Info("Snippet branch created successfully", "branch", branchName, "title", snippet.Title)
	return nil
}

// UpdateSnippetBranch updates an existing snippet branch
func (s *gitService) UpdateSnippetBranch(ctx context.Context, snippet *models.Snippet, userNote string) error {
	branchName := fmt.Sprintf("snippets/%s", snippet.ID)
	s.logger.Debug("Updating snippet branch", "branch", branchName, "title", snippet.Title)

	// Update branch with new content
	content := &SnippetContent{
		ID:        snippet.ID,
		Title:     snippet.Title,
		Content:   snippet.Content,
		Variables: models.StringSlice{}, // TODO: Extract from content
		Tags:      models.StringSlice{}, // TODO: Get from tags table
		CreatedAt: snippet.CreatedAt,
		UpdatedAt: snippet.UpdatedAt,
	}

	commitMessage := fmt.Sprintf("Update: %s", snippet.Title)
	if userNote != "" {
		commitMessage += "\n\n" + userNote
	}

	if err := s.updateBranchWithContent(branchName, "content.json", content, commitMessage); err != nil {
		return fmt.Errorf("failed to update branch with content: %w", err)
	}

	s.logger.Info("Snippet branch updated successfully", "branch", branchName, "title", snippet.Title)
	return nil
}

// DeleteSnippetBranch deletes a snippet branch
func (s *gitService) DeleteSnippetBranch(ctx context.Context, snippetID string) error {
	branchName := fmt.Sprintf("snippets/%s", snippetID)
	s.logger.Debug("Deleting snippet branch", "branch", branchName)

	// Delete the branch
	if err := s.repo.Storer.RemoveReference(plumbing.NewBranchReferenceName(branchName)); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	s.logger.Info("Snippet branch deleted successfully", "branch", branchName)
	return nil
}

// GetPromptHistory retrieves commit history for a prompt
func (s *gitService) GetPromptHistory(ctx context.Context, promptID string) ([]GitCommit, error) {
	branchName := fmt.Sprintf("prompts/%s", promptID)
	return s.getBranchHistory(branchName)
}

// GetSnippetHistory retrieves commit history for a snippet
func (s *gitService) GetSnippetHistory(ctx context.Context, snippetID string) ([]GitCommit, error) {
	branchName := fmt.Sprintf("snippets/%s", snippetID)
	return s.getBranchHistory(branchName)
}

// GetPromptVersion retrieves a specific version of a prompt
func (s *gitService) GetPromptVersion(ctx context.Context, promptID string, commitHash string) (*models.Prompt, error) {
	// Get the commit
	hash := plumbing.NewHash(commitHash)
	commit, err := s.repo.CommitObject(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	// Get the tree
	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	// Get content.json file
	file, err := tree.File("content.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get content.json: %w", err)
	}

	content, err := file.Contents()
	if err != nil {
		return nil, fmt.Errorf("failed to read file contents: %w", err)
	}

	// Parse JSON content
	var promptContent PromptContent
	if err := json.Unmarshal([]byte(content), &promptContent); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to models.Prompt
	prompt := &models.Prompt{
		ID:                     promptContent.ID,
		Title:                  promptContent.Title,
		Content:                promptContent.Content,
		Type:                   models.PromptType(promptContent.Type),
		UseCase:                getStringPointer(promptContent.UseCase),
		ModelCompatibilityTags: promptContent.ModelCompatibility,
		OtherParameters:        promptContent.Parameters,
		CreatedAt:              promptContent.CreatedAt,
		UpdatedAt:              promptContent.UpdatedAt,
	}

	return prompt, nil
}

// GetSnippetVersion retrieves a specific version of a snippet
func (s *gitService) GetSnippetVersion(ctx context.Context, snippetID string, commitHash string) (*models.Snippet, error) {
	// Get the commit
	hash := plumbing.NewHash(commitHash)
	commit, err := s.repo.CommitObject(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	// Get the tree
	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	// Get content.json file
	file, err := tree.File("content.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get content.json: %w", err)
	}

	content, err := file.Contents()
	if err != nil {
		return nil, fmt.Errorf("failed to read file contents: %w", err)
	}

	// Parse JSON content
	var snippetContent SnippetContent
	if err := json.Unmarshal([]byte(content), &snippetContent); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to models.Snippet
	snippet := &models.Snippet{
		ID:        snippetContent.ID,
		Title:     snippetContent.Title,
		Content:   snippetContent.Content,
		CreatedAt: snippetContent.CreatedAt,
		UpdatedAt: snippetContent.UpdatedAt,
	}

	return snippet, nil
}

// ValidateRepo validates the git repository health
func (s *gitService) ValidateRepo(ctx context.Context) error {
	if s.repo == nil {
		return fmt.Errorf("git repository not initialized")
	}

	// Check if repository is accessible
	_, err := s.repo.Head()
	if err != nil {
		return fmt.Errorf("failed to access repository head: %w", err)
	}

	return nil
}

// Helper methods

// createOrphanBranchWithContent creates a new orphan branch with content
func (s *gitService) createOrphanBranchWithContent(branchName, filename string, content interface{}, commitMessage string) error {
	s.logger.Debug("Creating orphan branch with content", "branch", branchName, "file", filename)

	// Store current HEAD to restore later (if it exists)
	var originalHead *plumbing.Reference
	if head, err := s.repo.Head(); err == nil {
		originalHead = head
	}

	// Get worktree
	worktree, err := s.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Create orphan branch by setting HEAD to point to new branch (this is the key insight from the research!)
	branchRef := plumbing.NewBranchReferenceName(branchName)
	if err := s.repo.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, branchRef)); err != nil {
		return fmt.Errorf("failed to create orphan branch: %w", err)
	}

	// Clear the index and working tree to make it truly orphan
	if err := s.clearWorktree(worktree); err != nil {
		// Restore original HEAD on error (if it existed)
		if originalHead != nil {
			s.repo.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, originalHead.Name()))
		}
		return fmt.Errorf("failed to clear worktree: %w", err)
	}

	// Write content and commit
	if err := s.writeContentAndCommit(worktree, filename, content, commitMessage); err != nil {
		// Restore original HEAD on error (if it existed)
		if originalHead != nil {
			s.repo.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, originalHead.Name()))
		}
		return fmt.Errorf("failed to write content and commit: %w", err)
	}

	s.logger.Debug("Orphan branch created successfully", "branch", branchName)
	return nil
}

// updateBranchWithContent updates an existing branch with new content
func (s *gitService) updateBranchWithContent(branchName, filename string, content interface{}, commitMessage string) error {
	s.logger.Debug("Updating branch with content", "branch", branchName, "file", filename)

	// Store current HEAD to restore later (if it exists)
	var originalHead *plumbing.Reference
	if head, err := s.repo.Head(); err == nil {
		originalHead = head
	}

	// Get worktree
	worktree, err := s.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Switch to the branch
	branchRef := plumbing.NewBranchReferenceName(branchName)
	if err := s.repo.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, branchRef)); err != nil {
		return fmt.Errorf("failed to switch to branch: %w", err)
	}

	// Checkout the branch content
	checkoutOptions := &git.CheckoutOptions{
		Branch: branchRef,
		Force:  true,
	}
	if err := worktree.Checkout(checkoutOptions); err != nil {
		// Restore original HEAD on error (if it existed)
		if originalHead != nil {
			s.repo.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, originalHead.Name()))
		}
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	// Write updated content and commit
	if err := s.writeContentAndCommit(worktree, filename, content, commitMessage); err != nil {
		// Restore original HEAD on error (if it existed)
		if originalHead != nil {
			s.repo.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, originalHead.Name()))
		}
		return fmt.Errorf("failed to write content and commit: %w", err)
	}

	s.logger.Debug("Branch updated successfully", "branch", branchName)
	return nil
}

// clearWorktree removes all files from the worktree (for orphan branches)
func (s *gitService) clearWorktree(worktree *git.Worktree) error {
	// Get the status to see what files exist
	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("failed to get worktree status: %w", err)
	}

	// Remove each file
	for file := range status {
		if _, err := worktree.Remove(file); err != nil {
			s.logger.Warn("Failed to remove file", "file", file, "error", err)
		}
	}

	return nil
}

// writeContentAndCommit writes content to a file and commits it
func (s *gitService) writeContentAndCommit(worktree *git.Worktree, filename string, content interface{}, commitMessage string) error {
	s.logger.Debug("Writing content and committing", "file", filename, "message", commitMessage)

	// Marshal content to JSON
	jsonData, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	// Write file
	filePath := filepath.Join(worktree.Filesystem.Root(), filename)
	if err := afero.WriteFile(s.fs, filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Add file to git
	if _, err := worktree.Add(filename); err != nil {
		return fmt.Errorf("failed to add file to git: %w", err)
	}

	// Commit
	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Proompt",
			Email: "proompt@local",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	s.logger.Debug("Content written and committed successfully", "file", filename)
	return nil
}

// getBranchHistory retrieves commit history for a branch
func (s *gitService) getBranchHistory(branchName string) ([]GitCommit, error) {
	s.logger.Debug("Getting branch history", "branch", branchName)

	// Get branch reference
	ref, err := s.repo.Reference(plumbing.NewBranchReferenceName(branchName), true)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch reference: %w", err)
	}

	// Get commit iterator
	commitIter, err := s.repo.Log(&git.LogOptions{
		From: ref.Hash(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}
	defer commitIter.Close()

	var commits []GitCommit
	err = commitIter.ForEach(func(commit *object.Commit) error {
		gitCommit := GitCommit{
			Hash:      commit.Hash.String(),
			Message:   commit.Message,
			Author:    commit.Author.Name,
			Email:     commit.Author.Email,
			Timestamp: commit.Author.When,
		}
		commits = append(commits, gitCommit)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	s.logger.Debug("Branch history retrieved", "branch", branchName, "commits", len(commits))
	return commits, nil
}

// Helper utility functions

// getStringValue safely gets string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// getStringPointer creates a string pointer from value
func getStringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
