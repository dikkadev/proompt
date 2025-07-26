package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dikkadev/proompt/server/internal/api/models"
	"github.com/dikkadev/proompt/server/internal/repository"
	"github.com/dikkadev/proompt/server/internal/template"
)

// TemplateHandler handles template-related HTTP requests
type TemplateHandler struct {
	repo repository.Repository
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(repo repository.Repository) *TemplateHandler {
	return &TemplateHandler{
		repo: repo,
	}
}

// PreviewTemplate godoc
// @Summary Preview template rendering
// @Description Preview how a template will render with provided data
// @Tags templates
// @Accept json
// @Produce json
// @Param request body models.TemplatePreviewRequest true "Template preview data"
// @Success 200 {object} models.TemplatePreviewResponse "Template preview result"
// @Failure 400 {object} models.ErrorResponse "Invalid request data or template syntax"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /template/preview [post]
func (h *TemplateHandler) PreviewTemplate(w http.ResponseWriter, r *http.Request) {
	var req models.TemplatePreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	// Basic validation
	if req.Content == "" {
		models.WriteBadRequest(w, "Content is required")
		return
	}

	// Get all snippets for resolution
	snippets, err := h.repo.Snippets().List(r.Context(), repository.SnippetFilters{})
	if err != nil {
		models.WriteInternalError(w, "Failed to fetch snippets")
		return
	}

	// Create snippet resolver
	snippetResolver := template.NewSnippetResolver(snippets, req.Variables)

	// Resolve template with snippets and variables
	result := snippetResolver.ResolveWithSnippets(req.Content)

	// Get variable status
	allVariables := snippetResolver.GetAllVariables(req.Content)
	variableStatus := snippetResolver.GetVariableStatusWithSnippets(req.Content)

	// Convert to response format
	var responseVars []models.TemplateVariable
	for _, v := range allVariables {
		status := variableStatus[v.Name]
		responseVars = append(responseVars, models.TemplateVariable{
			Name:         v.Name,
			DefaultValue: v.DefaultValue,
			HasDefault:   v.HasDefault,
			Status:       status,
		})
	}

	response := models.TemplatePreviewResponse{
		ResolvedContent: result.Content,
		Variables:       responseVars,
		Warnings:        result.Warnings,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AnalyzeTemplate godoc
// @Summary Analyze template structure
// @Description Analyze a template to extract variables, functions, and structure information
// @Tags templates
// @Accept json
// @Produce json
// @Param request body models.TemplatePreviewRequest true "Template analysis data"
// @Success 200 {object} models.TemplatePreviewResponse "Template analysis result"
// @Failure 400 {object} models.ErrorResponse "Invalid request data or template syntax"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /template/analyze [post]
func (h *TemplateHandler) AnalyzeTemplate(w http.ResponseWriter, r *http.Request) {
	var req models.TemplatePreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteBadRequest(w, "Invalid JSON body")
		return
	}

	// Basic validation
	if req.Content == "" {
		models.WriteBadRequest(w, "Content is required")
		return
	}

	// Get all snippets for analysis
	snippets, err := h.repo.Snippets().List(r.Context(), repository.SnippetFilters{})
	if err != nil {
		models.WriteInternalError(w, "Failed to fetch snippets")
		return
	}

	// Create snippet resolver
	snippetResolver := template.NewSnippetResolver(snippets, req.Variables)

	// Analyze without resolving
	allVariables := snippetResolver.GetAllVariables(req.Content)
	variableStatus := snippetResolver.GetVariableStatusWithSnippets(req.Content)

	// Get snippet insertion result for warnings
	snippetResult := snippetResolver.InsertSnippets(req.Content)

	// Convert to response format
	var responseVars []models.TemplateVariable
	for _, v := range allVariables {
		status := variableStatus[v.Name]
		responseVars = append(responseVars, models.TemplateVariable{
			Name:         v.Name,
			DefaultValue: v.DefaultValue,
			HasDefault:   v.HasDefault,
			Status:       status,
		})
	}

	response := models.TemplatePreviewResponse{
		ResolvedContent: snippetResult.Content, // Content with snippets inserted but variables not resolved
		Variables:       responseVars,
		Warnings:        snippetResult.Warnings,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
