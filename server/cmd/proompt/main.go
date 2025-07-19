package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dikka/proompt/server/internal/config"
	"github.com/dikka/proompt/server/internal/db"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Ensure necessary directories exist
	if err := cfg.EnsureDirectories(); err != nil {
		log.Fatal("Failed to create directories:", err)
	}

	// Connect to database based on configuration
	var database *db.DB
	switch cfg.DatabaseType() {
	case "local":
		database, err = db.NewLocal(cfg.Database.Local.Path)
		if err != nil {
			log.Fatal("Failed to connect to local database:", err)
		}

		// Run migrations for local database
		if err := database.RunMigrations(cfg.Database.Local.Migrations); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}

	case "turso":
		database, err = db.NewTurso(cfg.Database.Turso.URL, cfg.Database.Turso.Token)
		if err != nil {
			log.Fatal("Failed to connect to Turso database:", err)
		}

	default:
		log.Fatal("Unknown database type:", cfg.DatabaseType())
	}
	defer database.Close()

	fmt.Println("Proompt server initialized successfully")
	fmt.Printf("Database type: %s\n", cfg.DatabaseType())
	if cfg.Database.Local != nil {
		fmt.Printf("Database path: %s\n", cfg.Database.Local.Path)
	}
	fmt.Printf("Repos: %s\n", cfg.Storage.ReposDir)
	fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
}
