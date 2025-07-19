package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/dikkadev/proompt/server/internal/config"
	"github.com/dikkadev/proompt/server/internal/db"
	"github.com/dikkadev/proompt/server/internal/logging"
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

	slog.Info("Proompt server initialized successfully",
		"database_type", cfg.DatabaseType(),
		"repos_dir", cfg.Storage.ReposDir,
		"server_host", cfg.Server.Host,
		"server_port", cfg.Server.Port,
	)

	if cfg.Database.Local != nil {
		slog.Info("Local database configuration", "path", cfg.Database.Local.Path)
	}
}
