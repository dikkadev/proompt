package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dikka/proompt/server/internal/db"
)

func main() {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory:", err)
	}

	// Create proompt directory if it doesn't exist
	proomptDir := filepath.Join(homeDir, ".proompt")
	if err := os.MkdirAll(proomptDir, 0755); err != nil {
		log.Fatal("Failed to create proompt directory:", err)
	}

	// Database path
	dbPath := filepath.Join(proomptDir, "database.db")

	// Connect to database
	database, err := db.New(dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Run migrations
	migrationsPath := "internal/db/migrations"
	if err := database.RunMigrations(migrationsPath); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("Proompt server initialized successfully")
	fmt.Printf("Database: %s\n", dbPath)
}
