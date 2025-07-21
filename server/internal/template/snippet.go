package template

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dikkadev/proompt/server/internal/models"
)

// SnippetResolver handles snippet insertion and variable resolution
type SnippetResolver struct {
	snippets  map[string]*models.Snippet
	variables map[string]string
}

// NewSnippetResolver creates a new snippet resolver
func NewSnippetResolver(snippets []*models.Snippet, variables map[string]string) *SnippetResolver {
	snippetMap := make(map[string]*models.Snippet)
	for _, snippet := range snippets {
		snippetMap[snippet.Title] = snippet
	}

	if variables == nil {
		variables = make(map[string]string)
	}

	return &SnippetResolver{
		snippets:  snippetMap,
		variables: variables,
	}
}

// snippetRegex matches @snippet_name or @{snippet name with spaces}
var snippetRegex = regexp.MustCompile(`@(?:\{([^}]+)\}|([a-zA-Z_][a-zA-Z0-9_]*))`)

// SnippetInsertResult contains the result of snippet insertion
type SnippetInsertResult struct {
	Content   string
	Warnings  []string
	Variables []Variable
}

// InsertSnippets replaces snippet references with their content and resolves variables
func (sr *SnippetResolver) InsertSnippets(content string) SnippetInsertResult {
	var warnings []string
	var allVariables []Variable

	// Track processed snippets to prevent infinite recursion
	processed := make(map[string]bool)

	result := sr.insertSnippetsRecursive(content, processed, &warnings, &allVariables)

	return SnippetInsertResult{
		Content:   result,
		Warnings:  warnings,
		Variables: allVariables,
	}
}

func (sr *SnippetResolver) insertSnippetsRecursive(content string, processed map[string]bool, warnings *[]string, allVariables *[]Variable) string {
	return snippetRegex.ReplaceAllStringFunc(content, func(match string) string {
		submatch := snippetRegex.FindStringSubmatch(match)
		if len(submatch) < 3 {
			return match
		}

		// Extract snippet name (either from {name} or direct name)
		var snippetName string
		if submatch[1] != "" {
			snippetName = strings.TrimSpace(submatch[1])
		} else if submatch[2] != "" {
			snippetName = submatch[2]
		} else {
			return match
		}

		// Check for recursion
		if processed[snippetName] {
			*warnings = append(*warnings, fmt.Sprintf("Circular reference detected for snippet '%s'", snippetName))
			return match
		}

		// Find the snippet
		snippet, exists := sr.snippets[snippetName]
		if !exists {
			*warnings = append(*warnings, fmt.Sprintf("Snippet '%s' not found", snippetName))
			return match
		}

		// Mark as processed
		processed[snippetName] = true

		// Extract variables from snippet content
		snippetVars := ExtractVariables(snippet.Content)
		*allVariables = append(*allVariables, snippetVars...)

		// Recursively process the snippet content (in case it contains other snippets)
		processedContent := sr.insertSnippetsRecursive(snippet.Content, processed, warnings, allVariables)

		// Unmark to allow reuse in different contexts
		delete(processed, snippetName)

		return processedContent
	})
}

// ResolveWithSnippets performs both snippet insertion and variable resolution
func (sr *SnippetResolver) ResolveWithSnippets(content string) ResolveResult {
	// First, insert snippets
	snippetResult := sr.InsertSnippets(content)

	// Then resolve variables
	resolver := NewVariableResolver(sr.variables)
	variableResult := resolver.Resolve(snippetResult.Content)

	// Combine warnings
	allWarnings := append(snippetResult.Warnings, variableResult.Warnings...)

	return ResolveResult{
		Content:  variableResult.Content,
		Warnings: allWarnings,
	}
}

// GetAllVariables returns all variables from content and any referenced snippets
func (sr *SnippetResolver) GetAllVariables(content string) []Variable {
	snippetResult := sr.InsertSnippets(content)
	contentVars := ExtractVariables(content)

	// Combine and deduplicate variables
	varMap := make(map[string]Variable)

	for _, v := range contentVars {
		varMap[v.Name] = v
	}

	for _, v := range snippetResult.Variables {
		if existing, exists := varMap[v.Name]; !exists || (!existing.HasDefault && v.HasDefault) {
			varMap[v.Name] = v
		}
	}

	// Convert back to slice
	var result []Variable
	for _, v := range varMap {
		result = append(result, v)
	}

	return result
}

// GetVariableStatusWithSnippets returns variable status considering snippet variables
func (sr *SnippetResolver) GetVariableStatusWithSnippets(content string) map[string]string {
	allVars := sr.GetAllVariables(content)
	status := make(map[string]string)

	for _, variable := range allVars {
		if _, exists := sr.variables[variable.Name]; exists {
			status[variable.Name] = "provided"
		} else if variable.HasDefault {
			status[variable.Name] = "default"
		} else {
			status[variable.Name] = "missing"
		}
	}

	return status
}
