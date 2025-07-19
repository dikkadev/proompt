package config

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	XMLName  xml.Name `xml:"proompt"`
	Database Database `xml:"database"`
	Storage  Storage  `xml:"storage"`
	Server   Server   `xml:"server"`
}

type Database struct {
	Local *LocalDatabase `xml:"local"`
	Turso *TursoDatabase `xml:"turso"`
}

type LocalDatabase struct {
	Path       string `xml:"path"`
	Migrations string `xml:"migrations"`
}

type TursoDatabase struct {
	URL   string `xml:"url"`
	Token string `xml:"token"`
}

type Storage struct {
	ReposDir string `xml:"repos_dir"`
}

type Server struct {
	Host string `xml:"host"`
	Port int    `xml:"port"`
}

// Load loads configuration from XML file with fallback locations
func Load(configPath string) (*Config, error) {
	var path string

	if configPath != "" {
		path = configPath
	} else {
		// Try default locations in order
		candidates := []string{
			"./proompt.xml",
			"/etc/proompt/proompt.xml",
		}

		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				path = candidate
				break
			}
		}

		if path == "" {
			return nil, fmt.Errorf("no configuration file found in default locations: %v", candidates)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config Config
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Validate configuration (no defaults!)
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration in %s: %w", path, err)
	}

	return &config, nil
}

// validate checks if the configuration is valid and complete
func (c *Config) validate() error {
	// Database: exactly one must be configured
	if c.Database.Local != nil && c.Database.Turso != nil {
		return fmt.Errorf("database: cannot configure both local and turso - choose exactly one")
	}

	if c.Database.Local == nil && c.Database.Turso == nil {
		return fmt.Errorf("database: must configure either local or turso database")
	}

	// Validate local database config
	if c.Database.Local != nil {
		if c.Database.Local.Path == "" {
			return fmt.Errorf("database.local.path cannot be empty")
		}
		if c.Database.Local.Migrations == "" {
			return fmt.Errorf("database.local.migrations cannot be empty")
		}
	}

	// Validate turso database config
	if c.Database.Turso != nil {
		if c.Database.Turso.URL == "" {
			return fmt.Errorf("database.turso.url cannot be empty")
		}
		if c.Database.Turso.Token == "" {
			return fmt.Errorf("database.turso.token cannot be empty")
		}
	}

	// Storage validation
	if c.Storage.ReposDir == "" {
		return fmt.Errorf("storage.repos_dir cannot be empty")
	}

	// Server validation
	if c.Server.Host == "" {
		return fmt.Errorf("server.host cannot be empty")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535, got %d", c.Server.Port)
	}

	return nil
}

// DatabaseType returns the configured database type
func (c *Config) DatabaseType() string {
	if c.Database.Local != nil {
		return "local"
	}
	if c.Database.Turso != nil {
		return "turso"
	}
	return "unknown"
}

// EnsureDirectories creates necessary directories based on configuration
func (c *Config) EnsureDirectories() error {
	// Create database directory (only for local)
	if c.Database.Local != nil {
		dbDir := filepath.Dir(c.Database.Local.Path)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory %s: %w", dbDir, err)
		}
	}

	// Create repos directory
	if err := os.MkdirAll(c.Storage.ReposDir, 0755); err != nil {
		return fmt.Errorf("failed to create repos directory %s: %w", c.Storage.ReposDir, err)
	}

	return nil
}
