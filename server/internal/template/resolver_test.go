package template

import (
	"reflect"
	"testing"
)

func TestExtractVariables(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []Variable
	}{
		{
			name:     "no variables",
			content:  "Hello world",
			expected: []Variable{},
		},
		{
			name:    "single variable without default",
			content: "Hello {{name}}",
			expected: []Variable{
				{Name: "name", DefaultValue: "", HasDefault: false},
			},
		},
		{
			name:    "single variable with default",
			content: "Hello {{name:World}}",
			expected: []Variable{
				{Name: "name", DefaultValue: "World", HasDefault: true},
			},
		},
		{
			name:    "multiple variables mixed",
			content: "Hello {{name:World}}, you are {{age}} years old",
			expected: []Variable{
				{Name: "name", DefaultValue: "World", HasDefault: true},
				{Name: "age", DefaultValue: "", HasDefault: false},
			},
		},
		{
			name:    "duplicate variables",
			content: "{{name}} and {{name}} again",
			expected: []Variable{
				{Name: "name", DefaultValue: "", HasDefault: false},
			},
		},
		{
			name:    "variable with spaces",
			content: "{{ name : default value }}",
			expected: []Variable{
				{Name: "name", DefaultValue: " default value ", HasDefault: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractVariables(tt.content)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractVariables() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVariableResolver_Resolve(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		variables        map[string]string
		expectedContent  string
		expectedWarnings int
	}{
		{
			name:             "no variables",
			content:          "Hello world",
			variables:        map[string]string{},
			expectedContent:  "Hello world",
			expectedWarnings: 0,
		},
		{
			name:             "variable with value",
			content:          "Hello {{name}}",
			variables:        map[string]string{"name": "Alice"},
			expectedContent:  "Hello Alice",
			expectedWarnings: 0,
		},
		{
			name:             "variable with default, no value provided",
			content:          "Hello {{name:World}}",
			variables:        map[string]string{},
			expectedContent:  "Hello World",
			expectedWarnings: 0,
		},
		{
			name:             "variable with default, value provided",
			content:          "Hello {{name:World}}",
			variables:        map[string]string{"name": "Alice"},
			expectedContent:  "Hello Alice",
			expectedWarnings: 0,
		},
		{
			name:             "missing variable without default",
			content:          "Hello {{name}}",
			variables:        map[string]string{},
			expectedContent:  "Hello {{name}}",
			expectedWarnings: 1,
		},
		{
			name:    "mixed variables",
			content: "Hello {{name:World}}, you are {{age}} years old and live in {{city:Unknown}}",
			variables: map[string]string{
				"name": "Alice",
				"age":  "25",
			},
			expectedContent:  "Hello Alice, you are 25 years old and live in Unknown",
			expectedWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewVariableResolver(tt.variables)
			result := resolver.Resolve(tt.content)

			if result.Content != tt.expectedContent {
				t.Errorf("Resolve() content = %v, want %v", result.Content, tt.expectedContent)
			}

			if len(result.Warnings) != tt.expectedWarnings {
				t.Errorf("Resolve() warnings count = %v, want %v", len(result.Warnings), tt.expectedWarnings)
			}
		})
	}
}

func TestVariableResolver_GetMissingVariables(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		variables map[string]string
		expected  []string
	}{
		{
			name:      "no missing variables",
			content:   "Hello {{name}}",
			variables: map[string]string{"name": "Alice"},
			expected:  []string{},
		},
		{
			name:      "one missing variable",
			content:   "Hello {{name}}",
			variables: map[string]string{},
			expected:  []string{"name"},
		},
		{
			name:      "variable with default not missing",
			content:   "Hello {{name:World}}",
			variables: map[string]string{},
			expected:  []string{},
		},
		{
			name:      "mixed missing and provided",
			content:   "Hello {{name}}, you are {{age:unknown}} years old",
			variables: map[string]string{"name": "Alice"},
			expected:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewVariableResolver(tt.variables)
			result := resolver.GetMissingVariables(tt.content)

			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetMissingVariables() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestVariableResolver_GetVariableStatus(t *testing.T) {
	resolver := NewVariableResolver(map[string]string{
		"provided_var": "value",
	})

	content := "{{provided_var}} {{default_var:default}} {{missing_var}}"
	status := resolver.GetVariableStatus(content)

	expected := map[string]string{
		"provided_var": "provided",
		"default_var":  "default",
		"missing_var":  "missing",
	}

	if !reflect.DeepEqual(status, expected) {
		t.Errorf("GetVariableStatus() = %v, want %v", status, expected)
	}
}
