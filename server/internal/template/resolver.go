package template

import (
	"fmt"
	"regexp"
	"strings"
)

// Variable represents a template variable with optional default value
type Variable struct {
	Name         string
	DefaultValue string
	HasDefault   bool
}

// ResolveResult contains the resolved content and any warnings
type ResolveResult struct {
	Content  string
	Warnings []string
}

// VariableResolver handles template variable resolution
type VariableResolver struct {
	variables map[string]string
}

// NewVariableResolver creates a new variable resolver with the given variables
func NewVariableResolver(variables map[string]string) *VariableResolver {
	if variables == nil {
		variables = make(map[string]string)
	}
	return &VariableResolver{
		variables: variables,
	}
}

// variableRegex matches {{variable_name}} or {{variable_name:default_value}}
var variableRegex = regexp.MustCompile(`\{\{([^}:]+)(?::([^}]*))?\}\}`)

// ExtractVariables extracts all variables from the given content
func ExtractVariables(content string) []Variable {
	matches := variableRegex.FindAllStringSubmatch(content, -1)
	variables := make([]Variable, 0, len(matches))

	seen := make(map[string]bool)
	for _, match := range matches {
		name := strings.TrimSpace(match[1])
		if seen[name] {
			continue
		}
		seen[name] = true

		variable := Variable{
			Name: name,
		}

		if len(match) > 2 && match[2] != "" {
			variable.DefaultValue = match[2]
			variable.HasDefault = true
		}

		variables = append(variables, variable)
	}

	return variables
}

// Resolve replaces all variables in the content with their values
func (r *VariableResolver) Resolve(content string) ResolveResult {
	var warnings []string

	resolved := variableRegex.ReplaceAllStringFunc(content, func(match string) string {
		submatch := variableRegex.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}

		name := strings.TrimSpace(submatch[1])

		// Check if we have a value for this variable
		if value, exists := r.variables[name]; exists {
			return value
		}

		// Check if there's a default value
		if len(submatch) > 2 && submatch[2] != "" {
			return submatch[2]
		}

		// No value and no default - add warning and keep original
		warnings = append(warnings, fmt.Sprintf("Variable '%s' is not defined and has no default value", name))
		return match
	})

	return ResolveResult{
		Content:  resolved,
		Warnings: warnings,
	}
}

// GetMissingVariables returns variables that are required but not provided
func (r *VariableResolver) GetMissingVariables(content string) []string {
	variables := ExtractVariables(content)
	var missing []string

	for _, variable := range variables {
		if _, exists := r.variables[variable.Name]; !exists && !variable.HasDefault {
			missing = append(missing, variable.Name)
		}
	}

	return missing
}

// GetVariableStatus returns the status of each variable (provided, default, missing)
func (r *VariableResolver) GetVariableStatus(content string) map[string]string {
	variables := ExtractVariables(content)
	status := make(map[string]string)

	for _, variable := range variables {
		if _, exists := r.variables[variable.Name]; exists {
			status[variable.Name] = "provided"
		} else if variable.HasDefault {
			status[variable.Name] = "default"
		} else {
			status[variable.Name] = "missing"
		}
	}

	return status
}
