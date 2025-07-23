package config

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/dikkadev/prettyslog"
	"github.com/dikkadev/proompt/server/internal/logging"
	"github.com/go-playground/validator/v10"
)

// RawConfig represents the raw XML structure before environment processing
type RawConfig struct {
	XMLName   xml.Name      `xml:"proompt"`
	Databases []RawDatabase `xml:"database"`
	Storages  []RawStorage  `xml:"storage"`
	Servers   []RawServer   `xml:"server"`
	Loggings  []RawLogging  `xml:"logging"`
}

// Config represents the processed configuration for a specific environment
type Config struct {
	Database Database `validate:"required,database_exclusive"`
	Storage  Storage  `validate:"required"`
	Server   Server   `validate:"required"`
	Logging  Logging
}

type RawDatabase struct {
	Environment string         `xml:"environment,attr"`
	Local       *LocalDatabase `xml:"local"`
	Turso       *TursoDatabase `xml:"turso"`
}

type Database struct {
	Local *LocalDatabase
	Turso *TursoDatabase
}

type LocalDatabase struct {
	Path       string `xml:"path,attr" validate:"required"`
	Migrations string `xml:"migrations,attr" validate:"required"`
}

type TursoDatabase struct {
	URL   string `xml:"url,attr" validate:"required,url"`
	Token string `xml:"token,attr" validate:"required"`
}

type RawStorage struct {
	Environment string `xml:"environment,attr"`
	ReposDir    string `xml:"repos_dir,attr" validate:"required"`
}

type Storage struct {
	ReposDir string `validate:"required"`
}

type RawServer struct {
	Environment string `xml:"environment,attr"`
	Host        string `xml:"host,attr" validate:"required"`
	Port        int    `xml:"port,attr" validate:"required,min=1,max=65535"`
}

type Server struct {
	Host string `validate:"required"`
	Port int    `validate:"required,min=1,max=65535"`
}

type RawLogging struct {
	Environment string     `xml:"environment,attr"`
	Level       string     `xml:"level,attr"`
	Source      bool       `xml:"source,attr"`
	Timestamp   bool       `xml:"timestamp,attr"`
	Outputs     RawOutputs `xml:"outputs"`
}

type Logging struct {
	Level     string
	Source    bool
	Timestamp bool
	Outputs   Outputs
}

type RawOutputs struct {
	Stdout []RawStdoutOutput `xml:"stdout"`
	File   []RawFileOutput   `xml:"file"`
}

type Outputs struct {
	Stdout StdoutOutput
	File   FileOutput
}

type RawStdoutOutput struct {
	Environment string `xml:"environment,attr"`
	Enabled     bool   `xml:"enabled,attr"`
	Colors      bool   `xml:"colors,attr"`
}

type StdoutOutput struct {
	Enabled bool
	Colors  bool
}

type RawFileOutput struct {
	Environment string `xml:"environment,attr"`
	Enabled     bool   `xml:"enabled,attr"`
	Path        string `xml:"path,attr"`
	MaxSize     string `xml:"max_size,attr"`
	MaxFiles    int    `xml:"max_files,attr"`
}

type FileOutput struct {
	Enabled  bool
	Path     string
	MaxSize  string
	MaxFiles int
}

// Load loads configuration from XML file with fallback locations and processes it for the given environment
func Load(configPath string, environment string) (*Config, error) {
	var path string

	if configPath != "" {
		path = configPath
		fmt.Printf("DEBUG: Using provided config path: %s\n", path)
	} else {
		// Try default locations in order
		candidates := []string{
			"./proompt.xml",
			"/etc/proompt/proompt.xml",
		}

		fmt.Printf("DEBUG: Searching for config in default locations: %v\n", candidates)
		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				path = candidate
				fmt.Printf("DEBUG: Found config file at: %s\n", path)
				break
			} else {
				fmt.Printf("DEBUG: Config not found at: %s (%v)\n", candidate, err)
			}
		}

		if path == "" {
			return nil, fmt.Errorf("no configuration file found in default locations: %v", candidates)
		}
	}

	fmt.Printf("DEBUG: Reading config file: %s\n", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	fmt.Printf("DEBUG: Config file size: %d bytes\n", len(data))
	var rawConfig RawConfig
	if err := xml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	fmt.Printf("DEBUG: Parsed raw config for environment: %s\n", environment)

	// Process raw config for the specific environment
	config, err := processEnvironmentConfig(&rawConfig, environment)
	if err != nil {
		return nil, fmt.Errorf("failed to process environment config: %w", err)
	}

	fmt.Printf("DEBUG: Processed config - Database type: %s, Server: %s:%d\n",
		config.DatabaseType(), config.Server.Host, config.Server.Port)

	// Validate configuration
	fmt.Printf("DEBUG: Validating configuration\n")
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration in %s: %w", path, err)
	}

	fmt.Printf("DEBUG: Configuration loaded and validated successfully\n")
	return config, nil
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

// processEnvironmentConfig processes raw config for a specific environment
func processEnvironmentConfig(raw *RawConfig, environment string) (*Config, error) {
	config := &Config{}

	// Process Database
	database, err := selectDatabase(raw.Databases, environment)
	if err != nil {
		return nil, fmt.Errorf("failed to select database config: %w", err)
	}
	config.Database = *database

	// Process Storage
	storage, err := selectStorage(raw.Storages, environment)
	if err != nil {
		return nil, fmt.Errorf("failed to select storage config: %w", err)
	}
	config.Storage = *storage

	// Process Server
	server, err := selectServer(raw.Servers, environment)
	if err != nil {
		return nil, fmt.Errorf("failed to select server config: %w", err)
	}
	config.Server = *server

	// Process Logging (with defaults if not specified)
	logging := selectLogging(raw.Loggings, environment)
	config.Logging = *logging

	return config, nil
}

// selectDatabase selects the appropriate database config for the environment
func selectDatabase(databases []RawDatabase, environment string) (*Database, error) {
	var selected *RawDatabase

	// First, look for environment-specific config
	for _, db := range databases {
		if db.Environment == environment {
			selected = &db
			break
		}
	}

	// If not found, look for config without environment attribute (default)
	if selected == nil {
		for _, db := range databases {
			if db.Environment == "" {
				selected = &db
				break
			}
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("no database configuration found for environment: %s", environment)
	}

	return &Database{
		Local: selected.Local,
		Turso: selected.Turso,
	}, nil
}

// selectStorage selects the appropriate storage config for the environment
func selectStorage(storages []RawStorage, environment string) (*Storage, error) {
	var selected *RawStorage

	// First, look for environment-specific config
	for _, storage := range storages {
		if storage.Environment == environment {
			selected = &storage
			break
		}
	}

	// If not found, look for config without environment attribute (default)
	if selected == nil {
		for _, storage := range storages {
			if storage.Environment == "" {
				selected = &storage
				break
			}
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("no storage configuration found for environment: %s", environment)
	}

	return &Storage{
		ReposDir: selected.ReposDir,
	}, nil
}

// selectServer selects the appropriate server config for the environment
func selectServer(servers []RawServer, environment string) (*Server, error) {
	var selected *RawServer

	// First, look for environment-specific config
	for _, server := range servers {
		if server.Environment == environment {
			selected = &server
			break
		}
	}

	// If not found, look for config without environment attribute (default)
	if selected == nil {
		for _, server := range servers {
			if server.Environment == "" {
				selected = &server
				break
			}
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("no server configuration found for environment: %s", environment)
	}

	return &Server{
		Host: selected.Host,
		Port: selected.Port,
	}, nil
}

// selectLogging selects the appropriate logging config for the environment
func selectLogging(loggings []RawLogging, environment string) *Logging {
	var selected *RawLogging

	// First, look for environment-specific config
	for _, logging := range loggings {
		if logging.Environment == environment {
			selected = &logging
			break
		}
	}

	// If not found, look for config without environment attribute (default)
	if selected == nil {
		for _, logging := range loggings {
			if logging.Environment == "" {
				selected = &logging
				break
			}
		}
	}

	// If still not found, use defaults
	if selected == nil {
		return &Logging{
			Level:     "debug",
			Source:    false,
			Timestamp: true,
			Outputs: Outputs{
				Stdout: StdoutOutput{
					Enabled: true,
					Colors:  true,
				},
				File: FileOutput{
					Enabled: false,
				},
			},
		}
	}

	// Process outputs
	stdout := selectStdoutOutput(selected.Outputs.Stdout, environment)
	file := selectFileOutput(selected.Outputs.File, environment)

	return &Logging{
		Level:     getStringOrDefault(selected.Level, "debug"),
		Source:    selected.Source,
		Timestamp: getBoolOrDefault(selected.Timestamp, true),
		Outputs: Outputs{
			Stdout: *stdout,
			File:   *file,
		},
	}
}

// selectStdoutOutput selects the appropriate stdout output config
func selectStdoutOutput(outputs []RawStdoutOutput, environment string) *StdoutOutput {
	var selected *RawStdoutOutput

	// First, look for environment-specific config
	for _, output := range outputs {
		if output.Environment == environment {
			selected = &output
			break
		}
	}

	// If not found, look for config without environment attribute (default)
	if selected == nil {
		for _, output := range outputs {
			if output.Environment == "" {
				selected = &output
				break
			}
		}
	}

	// If still not found, use defaults
	if selected == nil {
		return &StdoutOutput{
			Enabled: true,
			Colors:  true,
		}
	}

	return &StdoutOutput{
		Enabled: selected.Enabled,
		Colors:  selected.Colors,
	}
}

// selectFileOutput selects the appropriate file output config
func selectFileOutput(outputs []RawFileOutput, environment string) *FileOutput {
	var selected *RawFileOutput

	// First, look for environment-specific config
	for _, output := range outputs {
		if output.Environment == environment {
			selected = &output
			break
		}
	}

	// If not found, look for config without environment attribute (default)
	if selected == nil {
		for _, output := range outputs {
			if output.Environment == "" {
				selected = &output
				break
			}
		}
	}

	// If still not found, use defaults
	if selected == nil {
		return &FileOutput{
			Enabled: false,
		}
	}

	return &FileOutput{
		Enabled:  selected.Enabled,
		Path:     selected.Path,
		MaxSize:  getStringOrDefault(selected.MaxSize, "100MB"),
		MaxFiles: getIntOrDefault(selected.MaxFiles, 5),
	}
}

// Helper functions
func getStringOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func getBoolOrDefault(value bool, defaultValue bool) bool {
	// For bool, we can't distinguish between false and unset, so we use the provided value
	return value
}

func getIntOrDefault(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

// CreateLogger creates a logger from the logging configuration
func (l *Logging) CreateLogger(group string) (*slog.Logger, error) {
	handler, err := l.CreateHandler(group)
	if err != nil {
		return nil, err
	}
	return slog.New(handler), nil
}

// CreateHandler creates a slog handler from the logging configuration
func (l *Logging) CreateHandler(group string) (slog.Handler, error) {
	var writers []io.Writer

	// Add stdout writer if enabled
	if l.Outputs.Stdout.Enabled {
		writers = append(writers, os.Stdout)
	}

	// Add file writer if enabled
	if l.Outputs.File.Enabled {
		fileWriter, err := l.createFileWriter()
		if err != nil {
			return nil, fmt.Errorf("failed to create file writer: %w", err)
		}
		writers = append(writers, fileWriter)
	}

	if len(writers) == 0 {
		return nil, fmt.Errorf("no log outputs enabled")
	}

	var writer io.Writer
	if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = io.MultiWriter(writers...)
	}

	// Parse log level
	level := l.parseLogLevel()

	// Create handler - only enable colors if stdout is enabled and colors are requested
	colors := l.Outputs.Stdout.Enabled && l.Outputs.Stdout.Colors

	return prettyslog.NewPrettyslogHandler(group,
		prettyslog.WithWriter(writer),
		prettyslog.WithLevel(level),
		prettyslog.WithColors(colors),
		prettyslog.WithSource(l.Source),
		prettyslog.WithTimestamp(l.Timestamp),
	), nil
}

// createFileWriter creates a file writer with proper directory creation
func (l *Logging) createFileWriter() (io.Writer, error) {
	path := l.Outputs.File.Path
	if path == "" {
		// Generate automatic filename
		path = logging.GenerateLogFileName()
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}

	// Open file for appending
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", path, err)
	}

	return file, nil
}

// parseLogLevel converts string log level to slog.Level
func (l *Logging) parseLogLevel() slog.Level {
	switch strings.ToLower(l.Level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug // default to debug
	}
}
