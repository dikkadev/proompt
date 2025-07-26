// Package main provides the entry point for the Proompt server application.
//
// @title Proompt API
// @version 1.0
// @description A comprehensive API for managing prompts, snippets, notes, and templates for AI interactions.
// @description
// @description ## Features
// @description - **Prompts**: Create, manage, and organize prompts for various AI models
// @description - **Snippets**: Store and manage reusable code snippets and text blocks
// @description - **Notes**: Add contextual notes to prompts for better organization
// @description - **Templates**: Preview and analyze template content
// @description - **Tagging**: Organize content with flexible tagging system
// @description - **Linking**: Create relationships between prompts
//
// @contact.name Proompt API Support
// @contact.url https://github.com/dikkadev/proompt
// @contact.email support@proompt.dev
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /api
// @schemes http https
//
// @tag.name health
// @tag.description Health check endpoints
//
// @tag.name prompts
// @tag.description Operations on prompts
//
// @tag.name prompt-links
// @tag.description Manage relationships between prompts
//
// @tag.name prompt-tags
// @tag.description Manage tags for prompts
//
// @tag.name snippets
// @tag.description Operations on code snippets and text blocks
//
// @tag.name snippet-tags
// @tag.description Manage tags for snippets
//
// @tag.name notes
// @tag.description Manage notes associated with prompts
//
// @tag.name templates
// @tag.description Template analysis and preview operations
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description API key for authentication (if implemented)
package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dikkadev/proompt/server/internal/api"
	"github.com/dikkadev/proompt/server/internal/config"
	"github.com/dikkadev/proompt/server/internal/db"
	"github.com/dikkadev/proompt/server/internal/git"
	"github.com/dikkadev/proompt/server/internal/logging"
	"github.com/dikkadev/proompt/server/internal/repository"

	// Import for swagger docs generation
	_ "github.com/dikkadev/proompt/server/docs"
)

// determineEnvironment determines the environment from CLI flag or env var
func determineEnvironment(cliEnv string) string {
	if cliEnv != "" {
		return cliEnv
	}
	if envVar := os.Getenv("PROOMPT_ENV"); envVar != "" {
		return envVar
	}
	return "dev" // default
}

func main() {
	// Set up basic logger for startup (will be reconfigured after config load)
	logging.SetDefault("proompt")

	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	environment := flag.String("env", "", "Environment (dev/prod), overrides PROOMPT_ENV")
	silent := flag.Bool("s", false, "Silent mode - disable stdout logging")
	logLevel := flag.String("l", "", "Log level (debug/info/warn/error), overrides config")
	flag.Parse()

	// Determine environment
	env := determineEnvironment(*environment)

	// Load configuration
	cfg, err := config.Load(*configPath, env)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Apply CLI overrides
	if *silent {
		cfg.Logging.Outputs.Stdout.Enabled = false
	}
	if *logLevel != "" {
		cfg.Logging.Level = *logLevel
	}

	// Set up the configured logger
	configuredLogger, err := cfg.Logging.CreateLogger("proompt")
	if err != nil {
		slog.Error("Failed to create logger", "error", err)
		os.Exit(1)
	}
	slog.SetDefault(configuredLogger)

	// Ensure necessary directories exist
	if err := cfg.EnsureDirectories(); err != nil {
		slog.Error("Failed to create directories", "error", err)
		os.Exit(1)
	}

	// Connect to database based on configuration
	var database *db.DB
	switch cfg.DatabaseType() {
	case "local":
		database, err = db.NewLocal(cfg.Database.Local.Path)
		if err != nil {
			slog.Error("Failed to connect to local database", "error", err)
			os.Exit(1)
		}

		// Run migrations for local database
		if err := database.RunMigrations(cfg.Database.Local.Migrations); err != nil {
			slog.Error("Failed to run migrations", "error", err)
			os.Exit(1)
		}

	case "turso":
		database, err = db.NewTurso(cfg.Database.Turso.URL, cfg.Database.Turso.Token)
		if err != nil {
			slog.Error("Failed to connect to Turso database", "error", err)
			os.Exit(1)
		}

	default:
		slog.Error("Unknown database type", "type", cfg.DatabaseType())
		os.Exit(1)
	}
	defer database.Close()

	// Initialize git service
	gitService, err := git.NewGitService(cfg)
	if err != nil {
		slog.Error("Failed to initialize git service", "error", err)
		os.Exit(1)
	}

	// Initialize repository
	repo := repository.New(database, gitService)
	defer repo.Close()

	// Create API server
	server := api.New(cfg, repo, slog.Default())

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		slog.Info("Received shutdown signal")

		// Give the server 30 seconds to shut down gracefully
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown error", "error", err)
		}
		cancel()
	}()

	slog.Info("Proompt server initialized successfully",
		"database_type", cfg.DatabaseType(),
		"repos_dir", cfg.Storage.ReposDir,
		"server_host", cfg.Server.Host,
		"server_port", cfg.Server.Port,
	)

	// Start the server
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}

	// Wait for shutdown
	<-ctx.Done()
	slog.Info("Server shutdown complete")
}
