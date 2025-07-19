package config

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	XMLName  xml.Name `xml:"proompt"`
	Database Database `xml:"database" validate:"required,database_exclusive"`
	Storage  Storage  `xml:"storage" validate:"required"`
	Server   Server   `xml:"server" validate:"required"`
}

type Database struct {
	Local *LocalDatabase `xml:"local" validate:"omitempty"`
	Turso *TursoDatabase `xml:"turso" validate:"omitempty"`
}

type LocalDatabase struct {
	Path       string `xml:"path" validate:"required"`
	Migrations string `xml:"migrations" validate:"required"`
}

type TursoDatabase struct {
	URL   string `xml:"url" validate:"required,url"`
	Token string `xml:"token" validate:"required"`
}

type Storage struct {
	ReposDir string `xml:"repos_dir" validate:"required"`
}

type Server struct {
	Host string `xml:"host" validate:"required"`
	Port int    `xml:"port" validate:"required,min=1,max=65535"`
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

var (
	// Global validator instance with custom validators registered
	configValidator *validator.Validate
)

func init() {
	configValidator = validator.New()

	// Register custom validator for database exclusivity
	configValidator.RegisterValidation("database_exclusive", validateDatabaseExclusive)
}

// validateDatabaseExclusive ensures exactly one database type is configured
func validateDatabaseExclusive(fl validator.FieldLevel) bool {
	db := fl.Field().Interface().(Database)
	localSet := db.Local != nil
	tursoSet := db.Turso != nil

	// XOR: exactly one must be set
	return localSet != tursoSet
}

// validate checks if the configuration is valid and complete using struct tags
func (c *Config) validate() error {
	if err := configValidator.Struct(c); err != nil {
		// Convert validator errors to more user-friendly messages
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationErrors(validationErrors)
		}
		return err
	}

	// Manually validate nested structs since they're pointers
	if c.Database.Local != nil {
		if err := configValidator.Struct(c.Database.Local); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				return formatValidationErrors(validationErrors)
			}
			return err
		}
	}

	if c.Database.Turso != nil {
		if err := configValidator.Struct(c.Database.Turso); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				return formatValidationErrors(validationErrors)
			}
			return err
		}
	}

	return nil
}

// formatValidationErrors converts validator errors to user-friendly messages
func formatValidationErrors(errs validator.ValidationErrors) error {
	for _, err := range errs {
		switch err.Tag() {
		case "database_exclusive":
			return fmt.Errorf("database: must configure exactly one database type (local or turso), not both or neither")
		case "required":
			return fmt.Errorf("%s cannot be empty", getFieldPath(err))
		case "url":
			return fmt.Errorf("%s must be a valid URL", getFieldPath(err))
		case "min":
			return fmt.Errorf("%s must be at least %s", getFieldPath(err), err.Param())
		case "max":
			return fmt.Errorf("%s must be at most %s", getFieldPath(err), err.Param())
		default:
			return fmt.Errorf("%s failed validation: %s", getFieldPath(err), err.Tag())
		}
	}
	return fmt.Errorf("validation failed")
}

// getFieldPath converts struct field names to user-friendly config paths
func getFieldPath(err validator.FieldError) string {
	namespace := err.Namespace()

	// Convert struct field names to config path format
	switch namespace {
	case "Config.Database.Local.Path":
		return "database.local.path"
	case "Config.Database.Local.Migrations":
		return "database.local.migrations"
	case "Config.Database.Turso.URL":
		return "database.turso.url"
	case "Config.Database.Turso.Token":
		return "database.turso.token"
	case "Config.Storage.ReposDir":
		return "storage.repos_dir"
	case "Config.Server.Host":
		return "server.host"
	case "Config.Server.Port":
		return "server.port"
	default:
		return err.Field()
	}
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
