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
)

func main() {
	// Set up prettyslog as the default logger
	logging.SetDefault("proompt")

	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

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
	logger := slog.Default()
	server := api.New(cfg, repo, logger)

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
