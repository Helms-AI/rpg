package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kon1790/rpg/internal/github"
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
	SpecPath string `json:"specPath" jsonschema:"required" jsonschema_description:"Path to the markdown spec file"`
}

// ParseSpecOutput contains the raw spec content
type ParseSpecOutput struct {
	Content string `json:"content"` // Raw markdown content
}

// ValidateSpecInput contains the path to a spec file
type ValidateSpecInput struct {
	SpecPath string `json:"specPath" jsonschema:"required" jsonschema_description:"Path to the markdown spec file"`
}

// ValidateSpecOutput contains validation results
type ValidateSpecOutput struct {
	Valid  bool                   `json:"valid"`
	Errors []spec.ValidationError `json:"errors,omitempty"`
}

// GetGenerationContextInput contains spec path and target language
type GetGenerationContextInput struct {
	SpecPath string `json:"specPath" jsonschema:"required" jsonschema_description:"Path to the markdown spec file"`
	Language string `json:"language" jsonschema:"required" jsonschema_description:"Target language ID (go, rust, java, python, typescript, csharp)"`
}

// GetGenerationContextOutput contains full generation context
type GetGenerationContextOutput struct {
	SpecContent    string             `json:"specContent"`    // Raw markdown content
	Language       languages.Language `json:"language"`       // Target language conventions
	PromptTemplate string             `json:"promptTemplate"` // Generation prompt
	OutputDir      string             `json:"outputDir"`      // Where to write output
}

// GetProjectStructureInput contains project name and target language
type GetProjectStructureInput struct {
	ProjectName string `json:"projectName" jsonschema:"required" jsonschema_description:"Name for the project (used for output directory and file naming)"`
	Language    string `json:"language" jsonschema:"required" jsonschema_description:"Target language ID (go, rust, java, python, typescript, csharp)"`
}

// GetProjectStructureOutput contains recommended file structure
type GetProjectStructureOutput struct {
	Files     []languages.ProjectFile `json:"files"`
	OutputDir string                  `json:"outputDir"`
}

// EnsureParityInput contains projects to check for parity
type EnsureParityInput struct {
	SpecPath string        `json:"specPath" jsonschema:"required" jsonschema_description:"Path to the original spec file"`
	Projects []ProjectInfo `json:"projects" jsonschema:"required" jsonschema_description:"List of generated projects to compare"`
}

// ProjectInfo identifies a generated project
type ProjectInfo struct {
	Language string `json:"language" jsonschema:"required" jsonschema_description:"Language ID (go, java, csharp, etc.)"`
	Path     string `json:"path" jsonschema:"required" jsonschema_description:"Path to the generated project directory"`
}

// EnsureParityOutput contains the parity analysis and fix instructions
type EnsureParityOutput struct {
	ParityScore       float64         `json:"parityScore"`
	ReferenceLanguage string          `json:"referenceLanguage"`
	FeatureMatrix     []FeatureStatus `json:"featureMatrix"`
	Gaps              []ParityGap     `json:"gaps"`
	FixInstructions   string          `json:"fixInstructions"`
}

// FeatureStatus tracks a feature across all implementations
type FeatureStatus struct {
	ID              string                     `json:"id"`
	Name            string                     `json:"name"`
	Category        string                     `json:"category"`
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
	FeatureID     string `json:"featureId"`
	FeatureName   string `json:"featureName"`
	Category      string `json:"category"`
	MissingIn     string `json:"missingIn"`
	ReferenceFile string `json:"referenceFile"`
	ReferenceCode string `json:"referenceCode"`
	SuggestedFix  string `json:"suggestedFix"`
	TargetFile    string `json:"targetFile"`
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

// ImportSpecFromGitHubInput contains parameters for importing a spec from a GitHub repository
type ImportSpecFromGitHubInput struct {
	Repository string `json:"repository" jsonschema:"required" jsonschema_description:"GitHub repository URL or shorthand (e.g., 'owner/repo', 'https://github.com/owner/repo', 'owner/repo@branch')"`
	Ref        string `json:"ref,omitempty" jsonschema_description:"Optional branch, tag, or commit SHA to checkout (overrides ref in URL)"`
	Token      string `json:"token,omitempty" jsonschema_description:"Optional GitHub personal access token for private repos (uses GITHUB_TOKEN env var if not provided)"`
	Name       string `json:"name,omitempty" jsonschema_description:"Optional name for the generated spec (defaults to repository name)"`
	Shallow    *bool  `json:"shallow,omitempty" jsonschema_description:"Use shallow clone for faster operation (default: true)"`
}

// ImportSpecFromGitHubOutput contains the analysis prompt and repository information for AI-powered spec generation
type ImportSpecFromGitHubOutput struct {
	ProjectName      string `json:"projectName"`
	Repository       string `json:"repository"`
	Branch           string `json:"branch"`
	CommitSHA        string `json:"commitSha"`
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

	return nil, ParseSpecOutput{Content: string(content)}, nil
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

	// Build output directory path
	outputPath := filepath.Join(s.outputDir, input.Language)

	// Build the response
	return nil, GetGenerationContextOutput{
		SpecContent:    string(content),
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

	files := adapter.GetProjectStructure(input.ProjectName, false)

	// Build output directory path: outputDir/<project-name>/<language>/
	outputPath := filepath.Join(s.outputDir, input.ProjectName, input.Language)

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

func (s *Server) handleImportSpecFromGitHub(ctx context.Context, req *mcp.CallToolRequest, input ImportSpecFromGitHubInput) (*mcp.CallToolResult, ImportSpecFromGitHubOutput, error) {
	// Parse the repository URL/shorthand
	repoInfo, err := github.ParseRepository(input.Repository)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Invalid repository: %v", err)},
			},
		}, ImportSpecFromGitHubOutput{}, nil
	}

	// Override ref if provided as separate parameter
	if input.Ref != "" {
		repoInfo.Ref = input.Ref
	}

	// Create cloner with settings
	cloner := github.NewCloner()
	cloner.Token = input.Token

	// Handle shallow clone setting (default true)
	if input.Shallow != nil {
		cloner.Shallow = *input.Shallow
	}

	// Clone the repository
	cloneResult, err := cloner.Clone(repoInfo)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to clone repository: %v", err)},
			},
		}, ImportSpecFromGitHubOutput{}, nil
	}

	// Ensure cleanup of temp directory
	defer github.Cleanup(cloneResult.LocalPath)

	// Collect all project files using existing import logic
	files, err := importer.CollectProjectFiles(cloneResult.LocalPath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to collect project files: %v", err)},
			},
		}, ImportSpecFromGitHubOutput{}, nil
	}

	// Override name if provided, otherwise use repo name
	if input.Name != "" {
		files.Name = input.Name
	} else {
		files.Name = repoInfo.Name
	}

	// Check if any files were found
	if files.GetTotalFileCount() == 0 {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("No source files found in repository: %s/%s", repoInfo.Owner, repoInfo.Name)},
			},
		}, ImportSpecFromGitHubOutput{}, nil
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
		}, ImportSpecFromGitHubOutput{}, nil
	}

	return nil, ImportSpecFromGitHubOutput{
		ProjectName:      files.Name,
		Repository:       fmt.Sprintf("%s/%s", repoInfo.Owner, repoInfo.Name),
		Branch:           cloneResult.Branch,
		CommitSHA:        cloneResult.CommitSHA,
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

// Example specs (embedded) - narrative style
const simpleFunctionExample = `# slugify

A utility function to convert text into URL-friendly slugs.

## Overview

This function takes any text input and transforms it into a clean, URL-safe string by:
- Converting to lowercase
- Replacing spaces and special characters with a separator (default: hyphen)
- Removing consecutive separators
- Trimming separators from the start and end

## Parameters

- **text** (required): The input string to slugify
- **separator** (optional, default: "-"): Character to use between words

## Expected Behavior

| Input | Output |
|-------|--------|
| "Hello World" | "hello-world" |
| "Hello, World! How are you?" | "hello-world-how-are-you" |
| "Café Au Lait" | "cafe-au-lait" |
| "Hello World" with separator "_" | "hello_world" |
`

const moduleExample = `# validation-utils

A module for validating common data formats like email addresses, URLs, and phone numbers.

## Overview

Provides reusable validation functions that return clear success/failure results with error messages. Designed for form validation and data sanitization.

## Email Validation

Should validate that an email address:
- Is not empty
- Contains exactly one @ symbol
- Has a non-empty local part (before @)
- Has a non-empty domain (after @)
- Domain contains at least one dot

Returns either "valid" or "invalid" with a list of specific error messages.

## Email Parsing

If an email is valid, extract its components:
- Local part (before @)
- Domain (after @)

If invalid, return nothing/null.
`

const fullProjectExample = `# task-api

A RESTful API for managing tasks with user authentication.

## Overview

A backend service that provides:
- User registration and JWT-based authentication
- CRUD operations for tasks
- Task filtering by status and priority
- Due date tracking

## Authentication

Users authenticate with email/password and receive a JWT token valid for 24 hours. All task endpoints require a valid token in the Authorization header.

## Data Model

### User
- Unique identifier
- Email (unique)
- Hashed password
- Display name
- Timestamps (created, updated)

### Task
- Unique identifier
- Title (required)
- Description (optional)
- Status: pending, in_progress, completed, cancelled
- Priority: low, medium, high, urgent
- Owner (references user)
- Due date (optional)
- Timestamps (created, updated)

## API Endpoints

### Authentication
- POST /auth/register - Create new user
- POST /auth/login - Get JWT token

### Tasks (requires authentication)
- GET /tasks - List user's tasks (supports filtering)
- POST /tasks - Create new task
- GET /tasks/:id - Get task details
- PUT /tasks/:id - Update task
- DELETE /tasks/:id - Delete task

## Configuration

Environment variables:
- PORT (default: 3000)
- JWT_SECRET (required)
- DATABASE_URL (required)
`
