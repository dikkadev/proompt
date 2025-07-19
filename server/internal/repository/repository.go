package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/dikkadev/proompt/server/internal/db"
	"github.com/dikkadev/proompt/server/internal/git"
	"github.com/dikkadev/proompt/server/internal/logging"
	"github.com/jmoiron/sqlx"
)

// repository implements the Repository interface
type repository struct {
	db         *db.DB
	gitService git.GitService
	logger     *slog.Logger

	prompts  PromptRepository
	snippets SnippetRepository
	notes    NoteRepository
}

// New creates a new repository instance
func New(database *db.DB, gitService git.GitService) Repository {
	logger := logging.NewLogger("repository")

	repo := &repository{
		db:         database,
		gitService: gitService,
		logger:     logger,
	}

	repo.prompts = newPromptRepository(database.DB, gitService, logger.WithGroup("prompts"))
	repo.snippets = newSnippetRepository(database.DB, gitService, logger.WithGroup("snippets"))
	repo.notes = newNoteRepository(database.DB, logger.WithGroup("notes"))

	return repo
}

// Prompts returns the prompt repository
func (r *repository) Prompts() PromptRepository {
	return r.prompts
}

// Snippets returns the snippet repository
func (r *repository) Snippets() SnippetRepository {
	return r.snippets
}

// Notes returns the note repository
func (r *repository) Notes() NoteRepository {
	return r.notes
}

// WithTx executes a function within a database transaction
func (r *repository) WithTx(ctx context.Context, fn func(Repository) error) error {
	r.logger.Debug("Starting database transaction")

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error("Failed to begin transaction", "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := &repository{
		db:         r.db,
		gitService: r.gitService,
		logger:     r.logger,
	}

	txRepo.prompts = newPromptRepositoryWithTx(tx, r.gitService, r.logger.WithGroup("prompts"))
	txRepo.snippets = newSnippetRepositoryWithTx(tx, r.gitService, r.logger.WithGroup("snippets"))
	txRepo.notes = newNoteRepositoryWithTx(tx, r.logger.WithGroup("notes"))

	defer func() {
		if p := recover(); p != nil {
			r.logger.Error("Transaction panic, rolling back", "panic", p)
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(txRepo); err != nil {
		r.logger.Debug("Transaction function failed, rolling back", "error", err)
		if rbErr := tx.Rollback(); rbErr != nil {
			r.logger.Error("Failed to rollback transaction", "error", rbErr, "original_error", err)
			return fmt.Errorf("transaction failed: %w (rollback error: %v)", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("Failed to commit transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Debug("Transaction committed successfully")
	return nil
}

// Close closes the repository and releases resources
func (r *repository) Close() error {
	r.logger.Debug("Closing repository")
	return r.db.Close()
}

// txExecutor interface for both *sqlx.DB and *sqlx.Tx
type txExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
