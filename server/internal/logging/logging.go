package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/dikkadev/prettyslog"
)

// DefaultConfig holds the default configuration for prettyslog handlers
var DefaultConfig = Config{
	Level:     slog.LevelDebug,
	Colors:    true,
	Source:    false,
	Timestamp: true,
}

// Config represents the logging configuration that can be shared across components
type Config struct {
	Level     slog.Level
	Colors    bool
	Source    bool
	Timestamp bool
}

// NewHandler creates a new prettyslog handler with the given group name using the default configuration
func NewHandler(group string) *prettyslog.PrettyslogHandler {
	return NewHandlerWithConfig(group, DefaultConfig)
}

// NewHandlerWithConfig creates a new prettyslog handler with the given group name and custom configuration
func NewHandlerWithConfig(group string, config Config) *prettyslog.PrettyslogHandler {
	return prettyslog.NewPrettyslogHandler(group,
		prettyslog.WithLevel(config.Level),
		prettyslog.WithColors(config.Colors),
		prettyslog.WithSource(config.Source),
		prettyslog.WithTimestamp(config.Timestamp),
	)
}

// NewLogger creates a new slog.Logger with the given group name using the default configuration
func NewLogger(group string) *slog.Logger {
	return slog.New(NewHandler(group))
}

// NewLoggerWithConfig creates a new slog.Logger with the given group name and custom configuration
func NewLoggerWithConfig(group string, config Config) *slog.Logger {
	return slog.New(NewHandlerWithConfig(group, config))
}

// SetDefault sets the default slog logger using the given group name and default configuration
func SetDefault(group string) {
	slog.SetDefault(NewLogger(group))
}

// SetDefaultWithConfig sets the default slog logger using the given group name and custom configuration
func SetDefaultWithConfig(group string, config Config) {
	slog.SetDefault(NewLoggerWithConfig(group, config))
}

// UpdateDefaultConfig updates the default configuration for future handlers
func UpdateDefaultConfig(config Config) {
	DefaultConfig = config
}

// NewLoggerFromConfig creates a logger from the new config structure
func NewLoggerFromConfig(group string, config LoggingConfig) (*slog.Logger, error) {
	handler, err := NewHandlerFromConfig(group, config)
	if err != nil {
		return nil, err
	}
	return slog.New(handler), nil
}

// NewHandlerFromConfig creates a handler from the new config structure
func NewHandlerFromConfig(group string, config LoggingConfig) (slog.Handler, error) {
	var writers []io.Writer

	// Add stdout writer if enabled
	if config.Outputs.Stdout.Enabled {
		writers = append(writers, os.Stdout)
	}

	// Add file writer if enabled
	if config.Outputs.File.Enabled {
		fileWriter, err := createFileWriter(config.Outputs.File)
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
	level := parseLogLevel(config.Level)

	// Create handler - only enable colors if stdout is enabled and colors are requested
	colors := config.Outputs.Stdout.Enabled && config.Outputs.Stdout.Colors

	return prettyslog.NewPrettyslogHandler(group,
		prettyslog.WithWriter(writer),
		prettyslog.WithLevel(level),
		prettyslog.WithColors(colors),
		prettyslog.WithSource(config.Source),
		prettyslog.WithTimestamp(config.Timestamp),
	), nil
}

// createFileWriter creates a file writer with proper directory creation
func createFileWriter(config FileOutputConfig) (io.Writer, error) {
	if config.Path == "" {
		// Generate automatic filename
		config.Path = GenerateLogFileName()
	}

	// Ensure directory exists
	dir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}

	// Open file for appending
	file, err := os.OpenFile(config.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", config.Path, err)
	}

	return file, nil
}
