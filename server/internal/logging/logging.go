package logging

import (
	"log/slog"

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
