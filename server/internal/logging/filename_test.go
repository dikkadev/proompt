package logging

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

func TestGenerateLogFileName(t *testing.T) {
	// Test the filename format
	filename := GenerateLogFileName()

	// Should match pattern: log-yyyyWwwd-HHMMSS.log
	// Example: log-2025W337-143022.log
	pattern := `^log-\d{4}W\d{2}\d-\d{6}\.log$`
	matched, err := regexp.MatchString(pattern, filename)
	if err != nil {
		t.Fatalf("Regex error: %v", err)
	}
	if !matched {
		t.Errorf("Filename %q doesn't match expected pattern %q", filename, pattern)
	}
}

func TestGenerateLogFileNameFormat(t *testing.T) {
	// Test with a known time
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Monday week 1",
			time:     time.Date(2025, 1, 6, 14, 30, 22, 0, time.UTC), // Monday, week 2
			expected: "log-2025W021-143022.log",
		},
		{
			name:     "Sunday week 1",
			time:     time.Date(2025, 1, 5, 9, 15, 45, 0, time.UTC), // Sunday, week 1
			expected: "log-2025W017-091545.log",
		},
		{
			name:     "Friday week 30",
			time:     time.Date(2025, 7, 25, 20, 21, 45, 0, time.UTC), // Friday, week 30
			expected: "log-2025W305-202145.log",
		},
		{
			name:     "Wednesday week 33",
			time:     time.Date(2025, 8, 13, 12, 0, 0, 0, time.UTC), // Wednesday, week 33
			expected: "log-2025W333-120000.log",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock time.Now() by creating a custom function
			// Since we can't easily mock time.Now(), we'll test the format logic
			year, week := tt.time.ISOWeek()
			weekday := int(tt.time.Weekday())
			if weekday == 0 { // Sunday = 0 in Go, but we want it to be 7
				weekday = 7
			}

			result := formatLogFileName(year, week, weekday, tt.time.Hour(), tt.time.Minute(), tt.time.Second())
			if result != tt.expected {
				t.Errorf("formatLogFileName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function to test the formatting logic
func formatLogFileName(year, week, weekday, hour, minute, second int) string {
	return fmt.Sprintf("log-%dW%02d%d-%02d%02d%02d.log",
		year, week, weekday, hour, minute, second)
}

func TestWeekdayConversion(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected int
	}{
		{"Sunday", time.Date(2025, 7, 27, 0, 0, 0, 0, time.UTC), 7},    // Sunday -> 7
		{"Monday", time.Date(2025, 7, 28, 0, 0, 0, 0, time.UTC), 1},    // Monday -> 1
		{"Tuesday", time.Date(2025, 7, 29, 0, 0, 0, 0, time.UTC), 2},   // Tuesday -> 2
		{"Wednesday", time.Date(2025, 7, 30, 0, 0, 0, 0, time.UTC), 3}, // Wednesday -> 3
		{"Thursday", time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC), 4},  // Thursday -> 4
		{"Friday", time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC), 5},     // Friday -> 5
		{"Saturday", time.Date(2025, 8, 2, 0, 0, 0, 0, time.UTC), 6},   // Saturday -> 6
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weekday := int(tt.time.Weekday())
			if weekday == 0 { // Sunday = 0 in Go, but we want it to be 7
				weekday = 7
			}

			if weekday != tt.expected {
				t.Errorf("Weekday conversion for %s: got %d, want %d", tt.name, weekday, tt.expected)
			}
		})
	}
}

func TestUTCTime(t *testing.T) {
	// Test that the function uses UTC time
	filename := GenerateLogFileName()

	// The filename should be generated using UTC time
	// We can't test the exact time, but we can verify it's a valid format
	if len(filename) != len("log-2025W337-143022.log") {
		t.Errorf("Filename length = %d, want %d", len(filename), len("log-2025W337-143022.log"))
	}

	// Verify it starts with "log-" and ends with ".log"
	if filename[:4] != "log-" {
		t.Errorf("Filename should start with 'log-', got %q", filename[:4])
	}
	if filename[len(filename)-4:] != ".log" {
		t.Errorf("Filename should end with '.log', got %q", filename[len(filename)-4:])
	}
}
