package logging

import (
	"log/slog"
	"strings"
)

// LoggingConfig represents the logging configuration from config package
// This is a copy to avoid circular imports
type LoggingConfig struct {
	Level     string
	Source    bool
	Timestamp bool
	Outputs   OutputsConfig
}

type OutputsConfig struct {
	Stdout StdoutOutputConfig
	File   FileOutputConfig
}

type StdoutOutputConfig struct {
	Enabled bool
	Colors  bool
}

type FileOutputConfig struct {
	Enabled  bool
	Path     string
	MaxSize  string
	MaxFiles int
}

// parseLogLevel converts string log level to slog.Level
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
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

// ConvertConfigLogging converts config.Logging to LoggingConfig
func ConvertConfigLogging(cfg interface{}) LoggingConfig {
	// This is a temporary solution to avoid circular imports
	// In a real implementation, we'd use interfaces or move types to a common package

	// For now, we'll use reflection or type assertion
	// This is a placeholder - we'll implement proper conversion
	return LoggingConfig{
		Level:     "debug",
		Source:    false,
		Timestamp: true,
		Outputs: OutputsConfig{
			Stdout: StdoutOutputConfig{
				Enabled: true,
				Colors:  true,
			},
			File: FileOutputConfig{
				Enabled: false,
			},
		},
	}
}
