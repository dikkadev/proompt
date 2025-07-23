package logging

import (
	"testing"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"debug", "DEBUG"},
		{"DEBUG", "DEBUG"},
		{"info", "INFO"},
		{"INFO", "INFO"},
		{"warn", "WARN"},
		{"WARN", "WARN"},
		{"warning", "WARN"},
		{"WARNING", "WARN"},
		{"error", "ERROR"},
		{"ERROR", "ERROR"},
		{"invalid", "DEBUG"}, // defaults to debug
		{"", "DEBUG"},        // defaults to debug
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			level := parseLogLevel(tt.input)
			if level.String() != tt.expected {
				t.Errorf("parseLogLevel(%q) = %v, want %v", tt.input, level.String(), tt.expected)
			}
		})
	}
}

func TestLoggingConfigDefaults(t *testing.T) {
	config := ConvertConfigLogging(nil)

	// Test defaults
	if config.Level != "debug" {
		t.Errorf("Default level = %v, want debug", config.Level)
	}
	if config.Source != false {
		t.Errorf("Default source = %v, want false", config.Source)
	}
	if config.Timestamp != true {
		t.Errorf("Default timestamp = %v, want true", config.Timestamp)
	}
	if !config.Outputs.Stdout.Enabled {
		t.Errorf("Default stdout enabled = %v, want true", config.Outputs.Stdout.Enabled)
	}
	if !config.Outputs.Stdout.Colors {
		t.Errorf("Default stdout colors = %v, want true", config.Outputs.Stdout.Colors)
	}
	if config.Outputs.File.Enabled {
		t.Errorf("Default file enabled = %v, want false", config.Outputs.File.Enabled)
	}
}
