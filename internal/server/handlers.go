package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kon1790/rpg/internal/importer"
	"github.com/kon1790/rpg/internal/languages"
	"github.com/kon1790/rpg/internal/parser"
	"github.com/kon1790/rpg/pkg/spec"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool input types

// ListLanguagesInput is empty - no parameters needed
type ListLanguagesInput struct{}

// ListLanguagesOutput contains the list of supported languages
type ListLanguagesOutput struct {
	Languages []languages.Language `json:"languages"`
}

// ParseSpecInput contains the path to a spec file
type ParseSpecInput struct {
	SpecPath string `json:"specPath" jsonschema:"required" jsonschema_description:"Absolute path to the .spec.md file"`
}

// ParseSpecOutput contains the parsed spec
type ParseSpecOutput struct {
	Spec *spec.Spec `json:"spec"`
}

// ValidateSpecInput contains the path to a spec file
type ValidateSpecInput struct {
	SpecPath string `json:"specPath" jsonschema:"required" jsonschema_description:"Absolute path to the .spec.md file"`
}

// ValidateSpecOutput contains validation results
type ValidateSpecOutput struct {
	Valid  bool                   `json:"valid"`
	Errors []spec.ValidationError `json:"errors,omitempty"`
}

// GetGenerationContextInput contains spec path and target language
type GetGenerationContextInput struct {
	SpecPath string `json:"specPath" jsonschema:"required" jsonschema_description:"Absolute path to the .spec.md file"`
	Language string `json:"language" jsonschema:"required" jsonschema_description:"Target language ID (go, rust, java, python, typescript, csharp)"`
}

// GetGenerationContextOutput contains full generation context
type GetGenerationContextOutput struct {
	Spec           *spec.Spec         `json:"spec"`
	Language       languages.Language `json:"language"`
	PromptTemplate string             `json:"promptTemplate"`
	OutputDir      string             `json:"outputDir"`
}

// GetProjectStructureInput contains spec path and target language
type GetProjectStructureInput struct {
	SpecPath string `json:"specPath" jsonschema:"required" jsonschema_description:"Absolute path to the .spec.md file"`
	Language string `json:"language" jsonschema:"required" jsonschema_description:"Target language ID (go, rust, java, python, typescript, csharp)"`
}

// GetProjectStructureOutput contains recommended file structure
type GetProjectStructureOutput struct {
	Files     []languages.ProjectFile `json:"files"`
	OutputDir string                  `json:"outputDir"`
}

// EnsureParityInput contains projects to check for parity
type EnsureParityInput struct {
	SpecPath string          `json:"specPath" jsonschema:"required" jsonschema_description:"Path to the original spec file"`
	Projects []ProjectInfo   `json:"projects" jsonschema:"required" jsonschema_description:"List of generated projects to compare"`
}

// ProjectInfo identifies a generated project
type ProjectInfo struct {
	Language string `json:"language" jsonschema:"required" jsonschema_description:"Language ID (go, java, csharp, etc.)"`
	Path     string `json:"path" jsonschema:"required" jsonschema_description:"Path to the generated project directory"`
}

// EnsureParityOutput contains the parity analysis and fix instructions
type EnsureParityOutput struct {
	ParityScore       float64           `json:"parityScore"`
	ReferenceLanguage string            `json:"referenceLanguage"`
	FeatureMatrix     []FeatureStatus   `json:"featureMatrix"`
	Gaps              []ParityGap       `json:"gaps"`
	FixInstructions   string            `json:"fixInstructions"`
}

// FeatureStatus tracks a feature across all implementations
type FeatureStatus struct {
	ID            string                     `json:"id"`
	Name          string                     `json:"name"`
	Category      string                     `json:"category"`
	Implementations map[string]Implementation `json:"implementations"`
}

// Implementation details for a feature in a specific language
type Implementation struct {
	Present     bool   `json:"present"`
	FilePath    string `json:"filePath,omitempty"`
	LineNumber  int    `json:"lineNumber,omitempty"`
	CodeSnippet string `json:"codeSnippet,omitempty"`
}

// ParityGap represents a missing or different feature
type ParityGap struct {
	FeatureID        string `json:"featureId"`
	FeatureName      string `json:"featureName"`
	Category         string `json:"category"`
	MissingIn        string `json:"missingIn"`
	ReferenceFile    string `json:"referenceFile"`
	ReferenceCode    string `json:"referenceCode"`
	SuggestedFix     string `json:"suggestedFix"`
	TargetFile       string `json:"targetFile"`
}

// ImportSpecFromSourceInput contains the input directory path
type ImportSpecFromSourceInput struct {
	InputPath string `json:"inputPath" jsonschema:"required" jsonschema_description:"Absolute path to directory containing source code to analyze (e.g., /Users/me/projects/my-app or ~/code/my-lib)"`
	Name      string `json:"name,omitempty" jsonschema_description:"Optional name for the generated spec (defaults to the input directory name)"`
}

// ImportSpecFromSourceOutput contains the analysis prompt and file statistics for AI-powered spec generation
type ImportSpecFromSourceOutput struct {
	ProjectName      string `json:"projectName"`
	DetectedLanguage string `json:"detectedLanguage"`
	SourceFileCount  int    `json:"sourceFileCount"`
	TestFileCount    int    `json:"testFileCount"`
	APISpecCount     int    `json:"apiSpecCount"`
	ConfigFileCount  int    `json:"configFileCount"`
	DocFileCount     int    `json:"docFileCount"`
	TotalContentSize int    `json:"totalContentSize"`
	AnalysisPrompt   string `json:"analysisPrompt"`
	SpecOutputPath   string `json:"specOutputPath"`
}

// Tool handlers

func (s *Server) handleListLanguages(ctx context.Context, req *mcp.CallToolRequest, input ListLanguagesInput) (*mcp.CallToolResult, ListLanguagesOutput, error) {
	return nil, ListLanguagesOutput{
		Languages: s.registry.List(),
	}, nil
}

func (s *Server) handleParseSpec(ctx context.Context, req *mcp.CallToolRequest, input ParseSpecInput) (*mcp.CallToolResult, ParseSpecOutput, error) {
	// Read the spec file
	content, err := os.ReadFile(input.SpecPath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to read spec file: %v", err)},
			},
		}, ParseSpecOutput{}, nil
	}

	// Parse the spec
	parsedSpec, err := parser.Parse(string(content))
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to parse spec: %v", err)},
			},
		}, ParseSpecOutput{}, nil
	}

	// Normalize to ensure all slices are non-nil for valid JSON schema output
	parsedSpec.Normalize()

	return nil, ParseSpecOutput{Spec: parsedSpec}, nil
}

func (s *Server) handleValidateSpec(ctx context.Context, req *mcp.CallToolRequest, input ValidateSpecInput) (*mcp.CallToolResult, ValidateSpecOutput, error) {
	// Read the spec file
	content, err := os.ReadFile(input.SpecPath)
	if err != nil {
		return nil, ValidateSpecOutput{
			Valid: false,
			Errors: []spec.ValidationError{
				{Severity: "error", Code: "FILE_READ_ERROR", Message: fmt.Sprintf("Failed to read file: %v", err)},
			},
		}, nil
	}

	// Validate the spec
	result := parser.Validate(string(content))
	return nil, ValidateSpecOutput{
		Valid:  result.Valid,
		Errors: result.Errors,
	}, nil
}

func (s *Server) handleGetGenerationContext(ctx context.Context, req *mcp.CallToolRequest, input GetGenerationContextInput) (*mcp.CallToolResult, GetGenerationContextOutput, error) {
	// Get the language adapter
	adapter, err := s.registry.Get(input.Language)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Unsupported language: %s. Use list_languages to see supported languages.", input.Language)},
			},
		}, GetGenerationContextOutput{}, nil
	}

	// Read the spec file
	content, err := os.ReadFile(input.SpecPath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to read spec file: %v", err)},
			},
		}, GetGenerationContextOutput{}, nil
	}

	// Parse the spec
	parsedSpec, err := parser.Parse(string(content))
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to parse spec: %v", err)},
			},
		}, GetGenerationContextOutput{}, nil
	}

	// Normalize to ensure all slices are non-nil for valid JSON schema output
	parsedSpec.Normalize()

	// Build output directory path: outputDir/<project-name>/<language>/
	outputPath := filepath.Join(s.outputDir, parsedSpec.Name, input.Language)

	// Build the response
	return nil, GetGenerationContextOutput{
		Spec:           parsedSpec,
		Language:       adapter.GetLanguage(),
		PromptTemplate: adapter.GetPromptContext(),
		OutputDir:      outputPath,
	}, nil
}

func (s *Server) handleGetProjectStructure(ctx context.Context, req *mcp.CallToolRequest, input GetProjectStructureInput) (*mcp.CallToolResult, GetProjectStructureOutput, error) {
	// Get the language adapter
	adapter, err := s.registry.Get(input.Language)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Unsupported language: %s. Use list_languages to see supported languages.", input.Language)},
			},
		}, GetProjectStructureOutput{}, nil
	}

	// Read and parse the spec to determine if it has tests
	content, err := os.ReadFile(input.SpecPath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to read spec file: %v", err)},
			},
		}, GetProjectStructureOutput{}, nil
	}

	parsedSpec, err := parser.Parse(string(content))
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to parse spec: %v", err)},
			},
		}, GetProjectStructureOutput{}, nil
	}

	// Normalize to ensure all slices are non-nil for valid JSON schema output
	parsedSpec.Normalize()

	hasTests := len(parsedSpec.Tests) > 0
	files := adapter.GetProjectStructure(parsedSpec.Name, hasTests)

	// Build output directory path: outputDir/<project-name>/<language>/
	outputPath := filepath.Join(s.outputDir, parsedSpec.Name, input.Language)

	return nil, GetProjectStructureOutput{
		Files:     files,
		OutputDir: outputPath,
	}, nil
}

func (s *Server) handleImportSpecFromSource(ctx context.Context, req *mcp.CallToolRequest, input ImportSpecFromSourceInput) (*mcp.CallToolResult, ImportSpecFromSourceOutput, error) {
	// Expand ~ to home directory if present
	inputPath := input.InputPath
	if strings.HasPrefix(inputPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			inputPath = filepath.Join(homeDir, inputPath[2:])
		}
	}

	// Convert to absolute path if relative
	if !filepath.IsAbs(inputPath) {
		absPath, err := filepath.Abs(inputPath)
		if err == nil {
			inputPath = absPath
		}
	}

	// Validate input path exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Input directory does not exist: %s", inputPath)},
			},
		}, ImportSpecFromSourceOutput{}, nil
	}

	// Collect all project files using the new AI-powered approach
	files, err := importer.CollectProjectFiles(inputPath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to collect project files: %v", err)},
			},
		}, ImportSpecFromSourceOutput{}, nil
	}

	// Override name if provided
	if input.Name != "" {
		files.Name = input.Name
	}

	// Check if any files were found
	if files.GetTotalFileCount() == 0 {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("No source files found in: %s", inputPath)},
			},
		}, ImportSpecFromSourceOutput{}, nil
	}

	// Build the analysis prompt for AI
	analysisPrompt := importer.BuildAnalysisPrompt(files)

	// Determine the spec output path
	specsDir := "specs"
	specPath := filepath.Join(specsDir, files.Name+".spec.md")

	// Create specs directory if needed (for later use by AI)
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to create specs directory: %v", err)},
			},
		}, ImportSpecFromSourceOutput{}, nil
	}

	return nil, ImportSpecFromSourceOutput{
		ProjectName:      files.Name,
		DetectedLanguage: files.Language,
		SourceFileCount:  len(files.SourceFiles),
		TestFileCount:    len(files.TestFiles),
		APISpecCount:     len(files.APISpecs),
		ConfigFileCount:  len(files.ConfigFiles),
		DocFileCount:     len(files.DocFiles),
		TotalContentSize: files.GetTotalContentSize(),
		AnalysisPrompt:   analysisPrompt,
		SpecOutputPath:   specPath,
	}, nil
}

func (s *Server) handleEnsureParity(ctx context.Context, req *mcp.CallToolRequest, input EnsureParityInput) (*mcp.CallToolResult, EnsureParityOutput, error) {
	if len(input.Projects) < 2 {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: "At least 2 projects are required for parity comparison"},
			},
		}, EnsureParityOutput{}, nil
	}

	// Read the spec for context
	specContent, err := os.ReadFile(input.SpecPath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to read spec file: %v", err)},
			},
		}, EnsureParityOutput{}, nil
	}

	// First project is the reference implementation
	reference := input.Projects[0]

	// Extract features from all projects
	allFeatures := make(map[string]map[string]Implementation)
	projectFiles := make(map[string]map[string]string) // language -> filename -> content

	for _, project := range input.Projects {
		features, files, err := s.extractProjectFeatures(project.Path, project.Language)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to analyze %s project: %v", project.Language, err)},
				},
			}, EnsureParityOutput{}, nil
		}

		for featureID, impl := range features {
			if allFeatures[featureID] == nil {
				allFeatures[featureID] = make(map[string]Implementation)
			}
			allFeatures[featureID][project.Language] = impl
		}
		projectFiles[project.Language] = files
	}

	// Build feature matrix and identify gaps
	featureMatrix := []FeatureStatus{}
	gaps := []ParityGap{}

	for featureID, implementations := range allFeatures {
		status := FeatureStatus{
			ID:              featureID,
			Name:            featureID,
			Category:        categorizeFeature(featureID),
			Implementations: implementations,
		}
		featureMatrix = append(featureMatrix, status)

		// Check if reference has this feature
		refImpl, refHas := implementations[reference.Language]
		if !refHas || !refImpl.Present {
			continue
		}

		// Check other projects for this feature
		for _, project := range input.Projects[1:] {
			impl, has := implementations[project.Language]
			if !has || !impl.Present {
				gap := ParityGap{
					FeatureID:     featureID,
					FeatureName:   featureID,
					Category:      categorizeFeature(featureID),
					MissingIn:     project.Language,
					ReferenceFile: refImpl.FilePath,
					ReferenceCode: refImpl.CodeSnippet,
					TargetFile:    suggestTargetFile(featureID, project.Language, projectFiles[project.Language]),
				}
				gaps = append(gaps, gap)
			}
		}
	}

	// Calculate parity score
	totalChecks := len(featureMatrix) * (len(input.Projects) - 1)
	gapCount := len(gaps)
	parityScore := 1.0
	if totalChecks > 0 {
		parityScore = float64(totalChecks-gapCount) / float64(totalChecks)
	}

	// Generate fix instructions for Claude
	fixInstructions := s.generateFixInstructions(input.SpecPath, string(specContent), reference.Language, gaps, projectFiles)

	return nil, EnsureParityOutput{
		ParityScore:       parityScore,
		ReferenceLanguage: reference.Language,
		FeatureMatrix:     featureMatrix,
		Gaps:              gaps,
		FixInstructions:   fixInstructions,
	}, nil
}

// extractProjectFeatures analyzes a project directory and extracts features
func (s *Server) extractProjectFeatures(projectPath string, language string) (map[string]Implementation, map[string]string, error) {
	features := make(map[string]Implementation)
	files := make(map[string]string)

	// Get file extensions for this language
	extensions := getLanguageExtensions(language)

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		// Check if this is a source file
		ext := filepath.Ext(path)
		isSource := false
		for _, e := range extensions {
			if ext == e {
				isSource = true
				break
			}
		}
		if !isSource {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(projectPath, path)
		files[relPath] = string(content)

		// Extract features from this file
		fileFeatures := extractFeaturesFromCode(string(content), language, relPath)
		for id, impl := range fileFeatures {
			features[id] = impl
		}

		return nil
	})

	return features, files, err
}

// extractFeaturesFromCode extracts feature markers from source code
func extractFeaturesFromCode(content string, language string, filePath string) map[string]Implementation {
	features := make(map[string]Implementation)

	// Feature detection patterns based on common patterns
	patterns := []struct {
		id      string
		markers []string
	}{
		{"url_validation", []string{"Uri.TryCreate", "new URI(", "URL(", "url.Parse"}},
		{"tag_normalization", []string{"ToLowerInvariant", "toLowerCase", "lower()"}},
		{"empty_tag_filter", []string{"filter(t -> !t.isBlank())", "Where(t => !string.IsNullOrWhiteSpace", "filter(lambda t: t.strip())"}},
		{"async_metadata_fetch", []string{"Task.Run", "@Async", "async def", "go func"}},
		{"og_title_extraction", []string{"og:title", "property='og:title'", "property=\"og:title\""}},
		{"og_description_extraction", []string{"og:description", "property='og:description'"}},
		{"favicon_extraction", []string{"favicon", "rel='icon'", "rel=\"icon\"", "shortcut icon"}},
		{"html_entity_decode", []string{"HtmlDecode", "StringEscapeUtils", "html.unescape", "htmlspecialchars_decode"}},
		{"port_config", []string{"PORT", "server.port", "port ="}},
		{"database_path_config", []string{"DATABASE_PATH", "datasource.url", "database_url"}},
		{"fetch_timeout_config", []string{"FETCH_TIMEOUT", "fetch-timeout", "timeout"}},
		{"unique_url_constraint", []string{".IsUnique()", "unique = true", "UNIQUE", "@Column(unique"}},
		{"cancellation_support", []string{"CancellationToken", "Context ctx", "context.Context"}},
		{"error_logging", []string{"LogWarning", "LogError", "log.warn", "log.error", "logger.warning", "logger.error"}},
		{"search_title", []string{"Title?.Contains", "getTitle().toLowerCase().contains", "title.lower() in"}},
		{"search_description", []string{"Description?.Contains", "getDescription()", "description.lower()"}},
		{"search_url", []string{"Url.Contains", "getUrl().toLowerCase()", "url.lower()"}},
		{"search_tags", []string{"Tags.Any", "getTags().stream()", "any(t in tags"}},
	}

	for _, p := range patterns {
		for _, marker := range p.markers {
			if containsPattern(content, marker) {
				lineNum := findLineNumber(content, marker)
				snippet := extractSnippetAround(content, marker, 5)
				features[p.id] = Implementation{
					Present:     true,
					FilePath:    filePath,
					LineNumber:  lineNum,
					CodeSnippet: snippet,
				}
				break
			}
		}
	}

	return features
}

// containsPattern checks if content contains a pattern (case-insensitive for some)
func containsPattern(content, pattern string) bool {
	return strings.Contains(content, pattern)
}

// findLineNumber finds the line number where a pattern appears
func findLineNumber(content, pattern string) int {
	idx := strings.Index(content, pattern)
	if idx == -1 {
		return 0
	}
	return strings.Count(content[:idx], "\n") + 1
}

// extractSnippetAround extracts lines around a pattern
func extractSnippetAround(content, pattern string, contextLines int) string {
	lines := strings.Split(content, "\n")
	lineNum := findLineNumber(content, pattern) - 1

	start := lineNum - contextLines
	if start < 0 {
		start = 0
	}
	end := lineNum + contextLines + 1
	if end > len(lines) {
		end = len(lines)
	}

	return strings.Join(lines[start:end], "\n")
}

// getLanguageExtensions returns file extensions for a language
func getLanguageExtensions(language string) []string {
	switch language {
	case "go":
		return []string{".go"}
	case "java":
		return []string{".java"}
	case "csharp":
		return []string{".cs"}
	case "python":
		return []string{".py"}
	case "typescript":
		return []string{".ts", ".tsx"}
	case "rust":
		return []string{".rs"}
	default:
		return []string{}
	}
}

// categorizeFeature returns the category for a feature
func categorizeFeature(featureID string) string {
	switch {
	case strings.Contains(featureID, "validation") || strings.Contains(featureID, "filter"):
		return "validation"
	case strings.Contains(featureID, "config"):
		return "configuration"
	case strings.Contains(featureID, "search"):
		return "search"
	case strings.Contains(featureID, "extraction") || strings.Contains(featureID, "fetch"):
		return "metadata"
	case strings.Contains(featureID, "async") || strings.Contains(featureID, "cancellation"):
		return "async"
	default:
		return "core"
	}
}

// suggestTargetFile suggests where a missing feature should be implemented
func suggestTargetFile(featureID string, language string, files map[string]string) string {
	// Find the most likely file based on feature category
	category := categorizeFeature(featureID)

	for filePath := range files {
		lowerPath := strings.ToLower(filePath)
		switch category {
		case "validation", "core":
			if strings.Contains(lowerPath, "service") {
				return filePath
			}
		case "configuration":
			if strings.Contains(lowerPath, "program") || strings.Contains(lowerPath, "application") || strings.Contains(lowerPath, "main") {
				return filePath
			}
		case "metadata":
			if strings.Contains(lowerPath, "metadata") {
				return filePath
			}
		}
	}

	// Return first service file as default
	for filePath := range files {
		if strings.Contains(strings.ToLower(filePath), "service") {
			return filePath
		}
	}

	return ""
}

// generateFixInstructions creates detailed instructions for Claude to fix gaps
func (s *Server) generateFixInstructions(specPath string, specContent string, refLang string, gaps []ParityGap, projectFiles map[string]map[string]string) string {
	if len(gaps) == 0 {
		return "All projects have 100% feature parity. No fixes needed."
	}

	var sb strings.Builder
	sb.WriteString("## Parity Fix Instructions\n\n")
	sb.WriteString(fmt.Sprintf("Reference implementation: **%s**\n\n", refLang))
	sb.WriteString("The following features are missing and need to be implemented:\n\n")

	// Group gaps by target language
	gapsByLang := make(map[string][]ParityGap)
	for _, gap := range gaps {
		gapsByLang[gap.MissingIn] = append(gapsByLang[gap.MissingIn], gap)
	}

	for lang, langGaps := range gapsByLang {
		sb.WriteString(fmt.Sprintf("### %s Implementation Gaps\n\n", strings.ToUpper(lang)))

		for _, gap := range langGaps {
			sb.WriteString(fmt.Sprintf("#### Missing: %s (%s)\n\n", gap.FeatureName, gap.Category))
			sb.WriteString(fmt.Sprintf("**Target file:** `%s`\n\n", gap.TargetFile))
			sb.WriteString(fmt.Sprintf("**Reference implementation** (from %s `%s`):\n", refLang, gap.ReferenceFile))
			sb.WriteString("```\n")
			sb.WriteString(gap.ReferenceCode)
			sb.WriteString("\n```\n\n")
			sb.WriteString("**Action:** Implement equivalent functionality following ")
			sb.WriteString(lang)
			sb.WriteString(" idioms and conventions.\n\n")
		}
	}

	sb.WriteString("---\n\n")
	sb.WriteString("After fixing, call `ensure_parity` again to verify 100% parity.\n")

	return sb.String()
}

func (s *Server) handleExampleResource(exampleName string) func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		content := getExampleSpec(exampleName)
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{Text: content},
			},
		}, nil
	}
}

func (s *Server) handleLanguageConventions(langID string) func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		adapter, err := s.registry.Get(langID)
		if err != nil {
			return nil, err
		}

		lang := adapter.GetLanguage()
		content, _ := json.MarshalIndent(map[string]interface{}{
			"id":               lang.ID,
			"name":             lang.Name,
			"version":          lang.Version,
			"conventions":      lang.Conventions,
			"idioms":           lang.Idioms,
			"projectStructure": lang.ProjectStructure,
			"promptContext":    adapter.GetPromptContext(),
		}, "", "  ")

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{Text: string(content)},
			},
		}, nil
	}
}

// getExampleSpec returns the content of an example spec
func getExampleSpec(name string) string {
	switch name {
	case "simple-function":
		return simpleFunctionExample
	case "module":
		return moduleExample
	case "full-project":
		return fullProjectExample
	default:
		return "# Example not found"
	}
}

// Example specs (embedded)
const simpleFunctionExample = `# slugify

A utility function to convert text into URL-friendly slugs.

## Target Languages

- go
- python
- typescript

## Functions

### slugify

Converts a string into a URL-friendly slug.

**accepts:**
- text: Text
- separator: Text (defaults to "-")

**returns:** Text

**logic:**
` + "```" + `
convert text to lowercase
replace all whitespace with the separator
remove all characters that are not letters, numbers, or the separator
collapse multiple consecutive separators into one
trim separators from start and end
return the result
` + "```" + `

## Tests

### slugify

#### test: converts simple text
given: "Hello World"
expect: "hello-world"

#### test: handles special characters
given: "Hello, World! How are you?"
expect: "hello-world-how-are-you"

#### test: uses custom separator
given:
- text: "Hello World"
- separator: "_"
expect: "hello_world"
`

const moduleExample = `# validation-utils

A module for validating common data formats.

## Meta

- version: 2.0.0
- license: MIT

## Target Languages

- typescript
- python
- go

## Types

### ValidationResult
is one of:
- Valid
- Invalid with errors: List of Text

### EmailParts
contains:
- local: Text
- domain: Text

## Functions

### validate_email

Validates an email address format.

**accepts:**
- email: Text

**returns:** ValidationResult

**logic:**
` + "```" + `
set errors to empty list

if email is empty:
    add "Email cannot be empty" to errors
    return Invalid with errors

if email does not contain exactly one "@":
    add "Email must contain exactly one @ symbol" to errors
    return Invalid with errors

split email by "@" into local_part and domain

if local_part is empty:
    add "Local part cannot be empty" to errors

if domain is empty:
    add "Domain cannot be empty" to errors

if domain does not contain ".":
    add "Domain must contain at least one dot" to errors

if errors is not empty:
    return Invalid with errors

return Valid
` + "```" + `

### parse_email

Parses a valid email into its components.

**accepts:**
- email: Text

**returns:** Optional EmailParts

**logic:**
` + "```" + `
set validation to validate_email(email)

if validation is Invalid:
    return Nothing

split email by "@" into local_part and domain

return EmailParts with:
    local: local_part
    domain: domain
` + "```" + `

## Tests

### validate_email

#### test: accepts valid email
given: "user@example.com"
expect: Valid

#### test: rejects empty email
given: ""
expect: Invalid with errors containing "cannot be empty"

#### test: rejects email without @
given: "userexample.com"
expect: Invalid with errors containing "@ symbol"
`

const fullProjectExample = `# task-api

A RESTful API for managing tasks with user authentication.

## Meta

- version: 1.0.0
- author: Development Team
- license: Apache-2.0
- description: Task management API with JWT authentication

## Target Languages

- typescript
- python
- go

## Dependencies

### Required
- HTTP server framework
- JSON Web Token (JWT) library
- Database ORM/query builder
- Password hashing (bcrypt or argon2)
- UUID generator

## Types

### User
contains:
- id: UUID
- email: Text (unique)
- password_hash: Text
- name: Text
- created_at: Timestamp
- updated_at: Timestamp

### Task
contains:
- id: UUID
- title: Text
- description: Optional Text
- status: TaskStatus
- priority: Priority
- owner_id: UUID (references User)
- due_date: Optional Timestamp
- created_at: Timestamp
- updated_at: Timestamp

### TaskStatus
is one of:
- Pending
- InProgress
- Completed
- Cancelled

### Priority
is one of:
- Low
- Medium
- High
- Urgent

### CreateTaskRequest
contains:
- title: Text
- description: Optional Text
- priority: Priority (defaults to Medium)
- due_date: Optional Timestamp

## Functions

### authenticate_user [async]

Authenticates a user and returns a JWT token.

**accepts:**
- email: Text
- password: Text

**returns:** Result of AuthToken

**logic:**
` + "```" + `
await user from find user by email in database

if user does not exist:
    return Failure with AuthError "Invalid credentials"

set password_valid to verify password against user.password_hash

if password_valid is false:
    return Failure with AuthError "Invalid credentials"

set token to generate JWT with:
    subject: user.id
    expiration: 24 hours from now

return Success with AuthToken:
    access_token: token
    token_type: "Bearer"
    expires_in: 86400
` + "```" + `

**errors:**
- AuthError: invalid credentials or expired token

### create_task [async]

Creates a new task for the authenticated user.

**accepts:**
- request: CreateTaskRequest
- user_id: UUID (from auth context)

**returns:** Task

**logic:**
` + "```" + `
set new_task to Task with:
    id: generate new UUID
    title: request.title
    description: request.description
    status: Pending
    priority: request.priority or Medium
    owner_id: user_id
    due_date: request.due_date
    created_at: current timestamp
    updated_at: current timestamp

await save new_task to database

return new_task
` + "```" + `

## Tests

### authenticate_user

#### test: returns token for valid credentials
given:
- email: "test@example.com"
- password: "correct_password"
setup:
- user exists with email and matching password hash
expect: Success with AuthToken containing access_token

#### test: fails for wrong password
given:
- email: "test@example.com"
- password: "wrong_password"
setup:
- user exists with email
expect: Failure with AuthError

## Project Structure

` + "```" + `
/
├── src/
│   ├── main.[ext]              # Application entry point
│   ├── config.[ext]            # Configuration loading
│   ├── routes/
│   │   ├── auth.[ext]          # Authentication endpoints
│   │   └── tasks.[ext]         # Task CRUD endpoints
│   ├── middleware/
│   │   ├── auth.[ext]          # JWT validation middleware
│   │   └── error_handler.[ext] # Global error handling
│   ├── services/
│   │   ├── auth_service.[ext]  # Authentication logic
│   │   └── task_service.[ext]  # Task business logic
│   ├── models/
│   │   └── index.[ext]         # Database models
│   └── types/
│       └── index.[ext]         # Type definitions
├── tests/
│   ├── auth_test.[ext]
│   └── tasks_test.[ext]
└── [package manifest]
` + "```" + `
`
