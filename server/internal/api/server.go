package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dikkadev/proompt/server/internal/api/handlers"
	"github.com/dikkadev/proompt/server/internal/config"
	"github.com/dikkadev/proompt/server/internal/repository"
)

// Server represents the HTTP server
type Server struct {
	server *http.Server
	logger *slog.Logger
	repo   repository.Repository
}

// New creates a new HTTP server
func New(cfg *config.Config, repo repository.Repository, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("GET /api/health", handlers.Health)

	// Create handlers
	promptHandlers := handlers.NewPromptHandlers(repo)
	snippetHandlers := handlers.NewSnippetHandlers(repo)
	noteHandlers := handlers.NewNoteHandlers(repo)

	// Prompts endpoints
	mux.HandleFunc("GET /api/prompts", promptHandlers.ListPrompts)
	mux.HandleFunc("POST /api/prompts", promptHandlers.CreatePrompt)
	mux.HandleFunc("GET /api/prompts/{id}", promptHandlers.GetPrompt)
	mux.HandleFunc("PUT /api/prompts/{id}", promptHandlers.UpdatePrompt)
	mux.HandleFunc("DELETE /api/prompts/{id}", promptHandlers.DeletePrompt)

	// Snippets endpoints
	mux.HandleFunc("GET /api/snippets", snippetHandlers.ListSnippets)
	mux.HandleFunc("POST /api/snippets", snippetHandlers.CreateSnippet)
	mux.HandleFunc("GET /api/snippets/{id}", snippetHandlers.GetSnippet)
	mux.HandleFunc("PUT /api/snippets/{id}", snippetHandlers.UpdateSnippet)
	mux.HandleFunc("DELETE /api/snippets/{id}", snippetHandlers.DeleteSnippet)

	// Notes endpoints
	mux.HandleFunc("GET /api/prompts/{id}/notes", noteHandlers.ListNotesForPrompt)
	mux.HandleFunc("POST /api/prompts/{id}/notes", noteHandlers.CreateNote)
	mux.HandleFunc("GET /api/notes/{id}", noteHandlers.GetNote)
	mux.HandleFunc("PUT /api/notes/{id}", noteHandlers.UpdateNote)
	mux.HandleFunc("DELETE /api/notes/{id}", noteHandlers.DeleteNote)

	// Create middleware stack
	stack := CreateStack(
		LoggingMiddleware(logger),
		RecoveryMiddleware(logger),
		CORSMiddleware(),
		ContentTypeMiddleware(),
	)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      stack(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		server: server,
		logger: logger,
		repo:   repo,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}
