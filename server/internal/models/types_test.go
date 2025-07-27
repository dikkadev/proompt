package models

import (
	"encoding/json"
	"testing"
)

func TestStringSlice_Scan(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    StringSlice
		expectError bool
	}{
		{
			name:        "nil value",
			input:       nil,
			expected:    nil,
			expectError: false,
		},
		{
			name:        "empty JSON array as []byte",
			input:       []byte("[]"),
			expected:    StringSlice{},
			expectError: false,
		},
		{
			name:        "empty JSON array as string",
			input:       "[]",
			expected:    StringSlice{},
			expectError: false,
		},
		{
			name:        "JSON array with strings as []byte",
			input:       []byte(`["gpt-4", "claude-3", "gemini-pro"]`),
			expected:    StringSlice{"gpt-4", "claude-3", "gemini-pro"},
			expectError: false,
		},
		{
			name:        "JSON array with strings as string",
			input:       `["gpt-4", "claude-3", "gemini-pro"]`,
			expected:    StringSlice{"gpt-4", "claude-3", "gemini-pro"},
			expectError: false,
		},
		{
			name:        "JSON array with strings as *string",
			input:       stringPtr(`["api", "database", "security"]`),
			expected:    StringSlice{"api", "database", "security"},
			expectError: false,
		},
		{
			name:        "nil *string",
			input:       (*string)(nil),
			expected:    nil,
			expectError: false,
		},
		{
			name:        "single string in array",
			input:       `["single"]`,
			expected:    StringSlice{"single"},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			input:       `["invalid json`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "unsupported type",
			input:       123,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "non-string array",
			input:       `[1, 2, 3]`,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s StringSlice
			err := s.Scan(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !equalStringSlices(s, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, s)
			}
		})
	}
}

func TestStringSlice_Value(t *testing.T) {
	tests := []struct {
		name     string
		input    StringSlice
		expected string
	}{
		{
			name:     "nil slice",
			input:    nil,
			expected: "",
		},
		{
			name:     "empty slice",
			input:    StringSlice{},
			expected: "[]",
		},
		{
			name:     "single item",
			input:    StringSlice{"gpt-4"},
			expected: `["gpt-4"]`,
		},
		{
			name:     "multiple items",
			input:    StringSlice{"gpt-4", "claude-3", "gemini-pro"},
			expected: `["gpt-4","claude-3","gemini-pro"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.input.Value()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.input == nil {
				if value != nil {
					t.Errorf("expected nil value for nil slice, got %v", value)
				}
				return
			}

			jsonStr, ok := value.([]byte)
			if !ok {
				t.Errorf("expected []byte, got %T", value)
				return
			}

			if string(jsonStr) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(jsonStr))
			}
		})
	}
}

func TestJSONMap_Scan(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    JSONMap
		expectError bool
	}{
		{
			name:        "nil value",
			input:       nil,
			expected:    nil,
			expectError: false,
		},
		{
			name:        "empty JSON object as []byte",
			input:       []byte("{}"),
			expected:    JSONMap{},
			expectError: false,
		},
		{
			name:        "empty JSON object as string",
			input:       "{}",
			expected:    JSONMap{},
			expectError: false,
		},
		{
			name:        "JSON object as []byte",
			input:       []byte(`{"max_tokens": 2000, "temperature": 0.7}`),
			expected:    JSONMap{"max_tokens": float64(2000), "temperature": 0.7},
			expectError: false,
		},
		{
			name:        "JSON object as string",
			input:       `{"max_tokens": 2000, "temperature": 0.7}`,
			expected:    JSONMap{"max_tokens": float64(2000), "temperature": 0.7},
			expectError: false,
		},
		{
			name:        "JSON object as *string",
			input:       stringPtr(`{"api_version": "v1", "timeout": 30}`),
			expected:    JSONMap{"api_version": "v1", "timeout": float64(30)},
			expectError: false,
		},
		{
			name:        "nil *string",
			input:       (*string)(nil),
			expected:    nil,
			expectError: false,
		},
		{
			name:        "complex nested object",
			input:       `{"config": {"retries": 3, "endpoints": ["api1", "api2"]}}`,
			expected:    JSONMap{"config": map[string]interface{}{"retries": float64(3), "endpoints": []interface{}{"api1", "api2"}}},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			input:       `{"invalid": json}`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "unsupported type",
			input:       123,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m JSONMap
			err := m.Scan(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !equalJSONMaps(m, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, m)
			}
		})
	}
}

func TestJSONMap_Value(t *testing.T) {
	tests := []struct {
		name     string
		input    JSONMap
		expected string
	}{
		{
			name:     "nil map",
			input:    nil,
			expected: "",
		},
		{
			name:     "empty map",
			input:    JSONMap{},
			expected: "{}",
		},
		{
			name:     "simple map",
			input:    JSONMap{"max_tokens": 2000},
			expected: `{"max_tokens":2000}`,
		},
		{
			name:     "complex map",
			input:    JSONMap{"max_tokens": 2000, "temperature": 0.7, "model": "gpt-4"},
			expected: "", // We'll check this is valid JSON instead of exact match due to key ordering
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.input.Value()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.input == nil {
				if value != nil {
					t.Errorf("expected nil value for nil map, got %v", value)
				}
				return
			}

			jsonBytes, ok := value.([]byte)
			if !ok {
				t.Errorf("expected []byte, got %T", value)
				return
			}

			if tt.expected != "" {
				if string(jsonBytes) != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, string(jsonBytes))
				}
			} else {
				// For complex maps, just verify it's valid JSON
				var result map[string]interface{}
				if err := json.Unmarshal(jsonBytes, &result); err != nil {
					t.Errorf("result is not valid JSON: %v", err)
				}
			}
		})
	}
}

// Test the actual database scanning scenario that was failing
func TestStringSlice_DatabaseScenario(t *testing.T) {
	// This simulates what SQLite returns for JSON columns
	testCases := []struct {
		name     string
		dbValue  interface{}
		expected StringSlice
	}{
		{
			name:     "SQLite JSON as string (typical case)",
			dbValue:  `["gpt-4", "claude-3", "gemini-pro"]`,
			expected: StringSlice{"gpt-4", "claude-3", "gemini-pro"},
		},
		{
			name:     "SQLite JSON as []byte (alternative driver)",
			dbValue:  []byte(`["api", "database", "security"]`),
			expected: StringSlice{"api", "database", "security"},
		},
		{
			name:     "Empty array from database",
			dbValue:  "[]",
			expected: StringSlice{},
		},
		{
			name:     "NULL from database",
			dbValue:  nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var s StringSlice
			err := s.Scan(tc.dbValue)
			if err != nil {
				t.Fatalf("Scan failed: %v", err)
			}

			if !equalStringSlices(s, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, s)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func equalStringSlices(a, b StringSlice) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalJSONMaps(a, b JSONMap) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}
