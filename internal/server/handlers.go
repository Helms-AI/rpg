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
	"github.com/kon1790/rpg/internal/importer/semantic"
	"github.com/kon1790/rpg/internal/importer/treesitter"
	"github.com/kon1790/rpg/internal/languages"
	"github.com/kon1790/rpg/internal/parity"
	"github.com/kon1790/rpg/internal/refinement"
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
	DetectedLanguage string `json:"detectedLanguage"` // Primary language (most files)
	SourceFileCount  int    `json:"sourceFileCount"`
	TestFileCount    int    `json:"testFileCount"`
	APISpecCount     int    `json:"apiSpecCount"`
	ConfigFileCount  int    `json:"configFileCount"`
	DocFileCount     int    `json:"docFileCount"`
	TotalContentSize int    `json:"totalContentSize"`
	AnalysisPrompt   string `json:"analysisPrompt"`
	SpecOutputPath   string `json:"specOutputPath"`
	SpecGenerated    bool   `json:"specGenerated"`    // Whether spec was successfully written
	GeneratedSpec    string `json:"generatedSpec"`    // The generated spec content

	// Multi-language detection and semantic analysis (AI orchestration integration)
	DetectedLanguages  []LanguageInfo              `json:"detectedLanguages"`           // All detected languages with metadata
	SemanticAnalyses   map[string]*SemanticSummary `json:"semanticAnalyses,omitempty"`  // Deep analysis for languages with parsers
	RawFilesByLanguage map[string][]FileContent    `json:"rawFilesByLanguage,omitempty"` // Raw files for languages without parsers
}

// SemanticSummary is a condensed version of DeepAnalyzeSourceOutput for embedding
type SemanticSummary struct {
	Language     string             `json:"language"`
	TypeCount    int                `json:"typeCount"`
	FunctionCount int               `json:"functionCount"`
	Types        []SemanticType     `json:"types"`
	Functions    []SemanticFunction `json:"functions"`
	Dependencies []DependencyInfo   `json:"dependencies"`
	FileCount    int                `json:"fileCount"`
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

// DeepAnalyzeSourceInput contains parameters for deep semantic analysis
type DeepAnalyzeSourceInput struct {
	SourcePath    string `json:"sourcePath" jsonschema:"required" jsonschema_description:"Path to the source code directory to analyze"`
	Language      string `json:"language,omitempty" jsonschema_description:"Optional language override (go, typescript, python, java, rust, csharp). Auto-detected if not provided."`
	AnalysisDepth string `json:"analysisDepth,omitempty" jsonschema_description:"Analysis depth: 'quick' (structure only), 'standard' (default, types+functions), 'deep' (full semantic with call graphs)"`
	IncludeTests  bool   `json:"includeTests,omitempty" jsonschema_description:"Include test files in analysis (default: false)"`
}

// DeepAnalyzeSourceOutput contains the full semantic analysis results
type DeepAnalyzeSourceOutput struct {
	ProjectName    string                   `json:"projectName"`
	Language       string                   `json:"language"`
	AnalysisDepth  string                   `json:"analysisDepth"`
	Types          []SemanticType           `json:"types"`
	Functions      []SemanticFunction       `json:"functions"`
	CallGraph      map[string][]string      `json:"callGraph,omitempty"`
	TypeGraph      map[string][]string      `json:"typeGraph,omitempty"`
	Dependencies   []DependencyInfo         `json:"dependencies"`
	FileCount      int                      `json:"fileCount"`
	Errors         []string                 `json:"errors,omitempty"`
	AnalysisPrompt string                   `json:"analysisPrompt"`
}

// SemanticType represents a type definition with semantic information
type SemanticType struct {
	Name                 string            `json:"name"`
	Kind                 string            `json:"kind"`
	IsPublic             bool              `json:"isPublic"`
	Fields               []SemanticField   `json:"fields,omitempty"`
	Methods              []string          `json:"methods,omitempty"`
	ImplementsInterfaces []string          `json:"implementsInterfaces,omitempty"`
	Generic              []string          `json:"generic,omitempty"`
	DocComment           string            `json:"docComment,omitempty"`
	Location             LocationInfo      `json:"location"`
}

// SemanticField represents a field with resolved type information
type SemanticField struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	ResolvedType string `json:"resolvedType"`
	Tags         string `json:"tags,omitempty"`
	IsPointer    bool   `json:"isPointer,omitempty"`
	IsSlice      bool   `json:"isSlice,omitempty"`
	IsMap        bool   `json:"isMap,omitempty"`
}

// SemanticFunction represents a function with semantic information
type SemanticFunction struct {
	Name        string              `json:"name"`
	Signature   string              `json:"signature"`
	IsPublic    bool                `json:"isPublic"`
	IsAsync     bool                `json:"isAsync,omitempty"`
	Parameters  []SemanticParameter `json:"parameters,omitempty"`
	ReturnTypes []string            `json:"returnTypes,omitempty"`
	Calls       []string            `json:"calls,omitempty"`
	Complexity  int                 `json:"complexity,omitempty"`
	DocComment  string              `json:"docComment,omitempty"`
	Location    LocationInfo        `json:"location"`
}

// SemanticParameter represents a function parameter
type SemanticParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// DependencyInfo represents a project dependency
type DependencyInfo struct {
	Path     string `json:"path"`
	Version  string `json:"version,omitempty"`
	IsStdLib bool   `json:"isStdLib"`
	IsLocal  bool   `json:"isLocal"`
}

// LocationInfo represents a source location
type LocationInfo struct {
	File      string `json:"file"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine,omitempty"`
}

// SemanticParityAnalysisInput contains parameters for semantic parity analysis
type SemanticParityAnalysisInput struct {
	SourcePath        string               `json:"sourcePath" jsonschema:"required" jsonschema_description:"Path to the source code directory (reference implementation)"`
	SourceLanguage    string               `json:"sourceLanguage,omitempty" jsonschema_description:"Language of source code (auto-detected if not provided)"`
	GeneratedProjects []GeneratedProject   `json:"generatedProjects" jsonschema:"required" jsonschema_description:"List of generated projects to compare against source"`
	ComparisonWeights *ComparisonWeights   `json:"comparisonWeights,omitempty" jsonschema_description:"Optional weights for parity dimensions (must sum to 1.0)"`
}

// GeneratedProject represents a generated project for parity analysis
type GeneratedProject struct {
	Language string `json:"language" jsonschema:"required" jsonschema_description:"Language of the generated project"`
	Path     string `json:"path" jsonschema:"required" jsonschema_description:"Path to the generated project directory"`
}

// ComparisonWeights configures the weights for parity dimensions
type ComparisonWeights struct {
	Structural float64 `json:"structural,omitempty" jsonschema_description:"Weight for structural parity (default: 0.20)"`
	Type       float64 `json:"type,omitempty" jsonschema_description:"Weight for type parity (default: 0.25)"`
	Behavioral float64 `json:"behavioral,omitempty" jsonschema_description:"Weight for behavioral parity (default: 0.35)"`
	Test       float64 `json:"test,omitempty" jsonschema_description:"Weight for test parity (default: 0.15)"`
	Idiomatic  float64 `json:"idiomatic,omitempty" jsonschema_description:"Weight for idiomatic code (default: 0.05)"`
}

// SemanticParityAnalysisOutput contains the full parity analysis results
type SemanticParityAnalysisOutput struct {
	OverallScore     float64                        `json:"overallScore"`
	Converged        bool                           `json:"converged"`
	ByDimension      DimensionScoresOutput          `json:"byDimension"`
	ByLanguage       map[string]LanguageParityResult `json:"byLanguage"`
	Gaps             []SemanticParityGap            `json:"gaps"`
	FixInstructions  string                         `json:"fixInstructions"`
}

// DimensionScoresOutput contains scores for each parity dimension
type DimensionScoresOutput struct {
	Structural float64 `json:"structural"`
	Type       float64 `json:"type"`
	Behavioral float64 `json:"behavioral"`
	Test       float64 `json:"test"`
	Idiomatic  float64 `json:"idiomatic"`
}

// LanguageParityResult contains parity analysis for a specific language
type LanguageParityResult struct {
	Language     string                `json:"language"`
	OverallScore float64               `json:"overallScore"`
	ByDimension  DimensionScoresOutput `json:"byDimension"`
	MissingTypes []string              `json:"missingTypes,omitempty"`
	MissingFuncs []string              `json:"missingFuncs,omitempty"`
}

// SemanticParityGap represents a specific parity issue
type SemanticParityGap struct {
	Dimension     string `json:"dimension"`
	Severity      string `json:"severity"`
	SourceItem    string `json:"sourceItem"`
	GeneratedItem string `json:"generatedItem,omitempty"`
	Discrepancy   string `json:"discrepancy"`
	SuggestedFix  string `json:"suggestedFix"`
	Language      string `json:"language"`
}

// IterativeRefinementLoopInput contains parameters for the refinement loop
type IterativeRefinementLoopInput struct {
	SourcePath           string   `json:"sourcePath" jsonschema:"required" jsonschema_description:"Path to the source code directory"`
	SourceLanguage       string   `json:"sourceLanguage,omitempty" jsonschema_description:"Language of source code (auto-detected if not provided)"`
	TargetLanguages      []string `json:"targetLanguages" jsonschema:"required" jsonschema_description:"Target languages to generate and compare"`
	OutputDir            string   `json:"outputDir" jsonschema:"required" jsonschema_description:"Directory for generated projects"`
	SpecPath             string   `json:"specPath,omitempty" jsonschema_description:"Path to spec file (generated if not exists)"`
	ConvergenceThreshold float64  `json:"convergenceThreshold,omitempty" jsonschema_description:"Minimum parity score to consider converged (default: 0.95)"`
	MaxIterations        int      `json:"maxIterations,omitempty" jsonschema_description:"Maximum refinement iterations (default: 5)"`
	RefinementStrategy   string   `json:"refinementStrategy,omitempty" jsonschema_description:"Strategy: 'spec-first', 'code-first', 'balanced', or 'adaptive' (default: 'balanced')"`
}

// IterativeRefinementLoopOutput contains the refinement loop results
type IterativeRefinementLoopOutput struct {
	Converged          bool                      `json:"converged"`
	FinalScore         float64                   `json:"finalScore"`
	IterationsUsed     int                       `json:"iterationsUsed"`
	IterationHistory   []IterationSummary        `json:"iterationHistory"`
	FinalSpec          string                    `json:"finalSpec,omitempty"`
	GeneratedProjects  map[string]string         `json:"generatedProjects,omitempty"`
	UnresolvedGaps     []SemanticParityGap       `json:"unresolvedGaps,omitempty"`
	RefinementSummary  string                    `json:"refinementSummary"`
	RefinementPrompt   string                    `json:"refinementPrompt,omitempty"`
}

// IterationSummary contains a summary of a single refinement iteration
type IterationSummary struct {
	Number             int     `json:"number"`
	ParityScore        float64 `json:"parityScore"`
	Phase              string  `json:"phase"`
	RefinementsApplied int     `json:"refinementsApplied"`
	ScoreImprovement   float64 `json:"scoreImprovement"`
}

// =============================================================================
// AUTONOMOUS CODE GENERATION - Spec-to-source with automatic parity loop
// =============================================================================

// GenerateSourceFromSpecInput contains parameters for AI-driven code generation
type GenerateSourceFromSpecInput struct {
	SpecPath  string `json:"specPath" jsonschema:"required" jsonschema_description:"Path to the markdown specification file"`
	Language  string `json:"language" jsonschema:"required" jsonschema_description:"Target language ID (go, rust, java, python, typescript, csharp)"`
	OutputDir string `json:"outputDir,omitempty" jsonschema_description:"Output directory for generated code"`
}

// GenerateSourceFromSpecOutput contains the spec content and generation prompt for AI
type GenerateSourceFromSpecOutput struct {
	// Spec content
	SpecContent string `json:"specContent"` // Full markdown content for AI to interpret
	SpecPath    string `json:"specPath"`    // Original spec path

	// Language conventions
	Language       languages.Language `json:"language"`       // Target language with conventions
	PromptTemplate string             `json:"promptTemplate"` // Language-specific generation hints

	// Output configuration
	OutputDir       string                  `json:"outputDir"`       // Where to write generated files
	ProjectStructure []languages.ProjectFile `json:"projectStructure"` // Recommended file structure

	// Generation instructions
	GenerationPrompt string `json:"generationPrompt"` // Complete instructions for AI
}

// =============================================================================
// AI ORCHESTRATION TOOLS - For multi-language analysis orchestrated by Claude
// =============================================================================

// ListProjectLanguagesInput contains the path to analyze
type ListProjectLanguagesInput struct {
	SourcePath string `json:"sourcePath" jsonschema:"required" jsonschema_description:"Path to the source code directory to scan for languages"`
}

// ListProjectLanguagesOutput contains all detected languages with metadata
type ListProjectLanguagesOutput struct {
	ProjectName string             `json:"projectName"`
	Languages   []LanguageInfo     `json:"languages"`
	TotalFiles  int                `json:"totalFiles"`
	Recommendation string          `json:"recommendation"`
}

// LanguageInfo contains information about a detected language
type LanguageInfo struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	FileCount     int      `json:"fileCount"`
	Extensions    []string `json:"extensions"`
	HasParser     bool     `json:"hasParser"`
	SampleFiles   []string `json:"sampleFiles"`
}

// GetFilesForLanguageInput contains parameters to get files for a specific language
type GetFilesForLanguageInput struct {
	SourcePath   string `json:"sourcePath" jsonschema:"required" jsonschema_description:"Path to the source code directory"`
	Language     string `json:"language" jsonschema:"required" jsonschema_description:"Language ID to get files for (e.g., 'go', 'typescript', 'python', 'sql', 'protobuf')"`
	IncludeTests bool   `json:"includeTests,omitempty" jsonschema_description:"Include test files (default: false)"`
	MaxFiles     int    `json:"maxFiles,omitempty" jsonschema_description:"Maximum number of files to return (default: 50, max: 200)"`
}

// GetFilesForLanguageOutput contains raw file contents for AI interpretation
type GetFilesForLanguageOutput struct {
	Language    string        `json:"language"`
	FileCount   int           `json:"fileCount"`
	TotalSize   int           `json:"totalSize"`
	Files       []FileContent `json:"files"`
	Truncated   bool          `json:"truncated"`
	AIPrompt    string        `json:"aiPrompt"`
}

// FileContent contains a file's path and content
type FileContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int    `json:"size"`
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

	// Collect all project files using the basic collector (for file counts)
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

	// =========================================================================
	// AI-DRIVEN SPEC GENERATION - Full source code analysis without parsers
	// =========================================================================

	// Step 1: Detect all languages in the project (for statistics only)
	detectedLanguages := s.detectAllLanguages(inputPath)

	// Step 2: Collect ALL source files for AI analysis (no semantic parsing)
	// We'll send all source code to the AI for comprehensive analysis
	allSourceFiles := make(map[string][]FileContent)
	for _, lang := range detectedLanguages {
		if lang.FileCount > 0 {
			// Collect ALL files for this language (increased limit for AI analysis)
			sourceFiles := s.collectRawFilesForLanguage(inputPath, lang.ID, 200) // Increased limit
			if len(sourceFiles) > 0 {
				allSourceFiles[lang.ID] = sourceFiles
			}
		}
	}

	// Step 3: Build AI-driven analysis prompt with full source code
	analysisPrompt := s.buildAIAnalysisPrompt(files, detectedLanguages, allSourceFiles)

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

	// Determine primary language (most files)
	primaryLang := files.Language
	if len(detectedLanguages) > 0 {
		primaryLang = detectedLanguages[0].ID // Already sorted by file count
	}

	// Step 4: Generate the spec content from AI analysis
	generatedSpec := s.generateAISpec(files.Name, primaryLang, detectedLanguages, allSourceFiles, files, analysisPrompt)

	// Step 6: Write the spec to the output file
	specGenerated := false
	if err := os.WriteFile(specPath, []byte(generatedSpec), 0644); err == nil {
		specGenerated = true
	}

	return nil, ImportSpecFromSourceOutput{
		ProjectName:        files.Name,
		DetectedLanguage:   primaryLang,
		SourceFileCount:    len(files.SourceFiles),
		TestFileCount:      len(files.TestFiles),
		APISpecCount:       len(files.APISpecs),
		ConfigFileCount:    len(files.ConfigFiles),
		DocFileCount:       len(files.DocFiles),
		TotalContentSize:   files.GetTotalContentSize(),
		AnalysisPrompt:     analysisPrompt,
		SpecOutputPath:     specPath,
		SpecGenerated:      specGenerated,
		GeneratedSpec:      generatedSpec,
		DetectedLanguages:  detectedLanguages,
		SemanticAnalyses:   nil, // No longer using semantic analysis
		RawFilesByLanguage: allSourceFiles,
	}, nil
}

// detectAllLanguages scans a directory and returns all detected languages with metadata
func (s *Server) detectAllLanguages(sourcePath string) []LanguageInfo {
	// Extended language detection
	langExtensions := map[string][]string{
		"go":         {".go"},
		"typescript": {".ts", ".tsx"},
		"javascript": {".js", ".jsx", ".mjs", ".cjs"},
		"python":     {".py", ".pyi"},
		"java":       {".java"},
		"rust":       {".rs"},
		"csharp":     {".cs"},
		"kotlin":     {".kt", ".kts"},
		"swift":      {".swift"},
		"ruby":       {".rb"},
		"php":        {".php"},
		"scala":      {".scala"},
		"c":          {".c", ".h"},
		"cpp":        {".cpp", ".cc", ".cxx", ".hpp", ".hxx"},
		"sql":        {".sql"},
		"protobuf":   {".proto"},
		"graphql":    {".graphql", ".gql"},
		"yaml":       {".yaml", ".yml"},
		"json":       {".json"},
		"markdown":   {".md"},
		"html":       {".html", ".htm"},
		"css":        {".css", ".scss", ".sass", ".less"},
		"shell":      {".sh", ".bash", ".zsh"},
	}

	parsersAvailable := map[string]bool{
		"go": true, "typescript": true, "python": true,
		"java": true, "rust": true, "csharp": true,
	}

	langNames := map[string]string{
		"go": "Go", "typescript": "TypeScript", "javascript": "JavaScript",
		"python": "Python", "java": "Java", "rust": "Rust", "csharp": "C#",
		"kotlin": "Kotlin", "swift": "Swift", "ruby": "Ruby", "php": "PHP",
		"scala": "Scala", "c": "C", "cpp": "C++", "sql": "SQL",
		"protobuf": "Protocol Buffers", "graphql": "GraphQL", "yaml": "YAML",
		"json": "JSON", "markdown": "Markdown", "html": "HTML", "css": "CSS",
		"shell": "Shell",
	}

	skipDirs := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, "target": true,
		"build": true, "dist": true, "out": true, ".idea": true, ".vscode": true,
		"__pycache__": true, ".pytest_cache": true, ".mypy_cache": true,
		".next": true, ".nuxt": true, "coverage": true, ".gradle": true,
		".mvn": true, "bin": true, "obj": true, "testdata": true,
	}

	langFiles := make(map[string][]string)
	langExts := make(map[string]map[string]bool)

	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || skipDirs[name] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		relPath, _ := filepath.Rel(sourcePath, path)

		for langID, exts := range langExtensions {
			for _, e := range exts {
				if ext == e {
					langFiles[langID] = append(langFiles[langID], relPath)
					if langExts[langID] == nil {
						langExts[langID] = make(map[string]bool)
					}
					langExts[langID][ext] = true
					break
				}
			}
		}
		return nil
	})

	var languages []LanguageInfo
	for langID, files := range langFiles {
		if len(files) == 0 {
			continue
		}

		var exts []string
		for ext := range langExts[langID] {
			exts = append(exts, ext)
		}

		sampleFiles := files
		if len(sampleFiles) > 5 {
			sampleFiles = sampleFiles[:5]
		}

		languages = append(languages, LanguageInfo{
			ID:          langID,
			Name:        langNames[langID],
			FileCount:   len(files),
			Extensions:  exts,
			HasParser:   parsersAvailable[langID],
			SampleFiles: sampleFiles,
		})
	}

	// Sort by file count (descending)
	for i := 0; i < len(languages); i++ {
		for j := i + 1; j < len(languages); j++ {
			if languages[j].FileCount > languages[i].FileCount {
				languages[i], languages[j] = languages[j], languages[i]
			}
		}
	}

	return languages
}

// performSemanticAnalysis performs deep semantic analysis for a specific language
func (s *Server) performSemanticAnalysis(sourcePath, language string) *SemanticSummary {
	registry := semantic.DefaultRegistry()
	lang := treesitter.Language(language)

	analyzer, ok := registry.Get(lang)
	if !ok || !analyzer.IsAvailable() {
		return nil
	}

	analysis, err := analyzer.Analyze(sourcePath)
	if err != nil {
		return nil
	}

	// Convert to SemanticSummary
	summary := &SemanticSummary{
		Language:      language,
		TypeCount:     len(analysis.Types),
		FunctionCount: len(analysis.Functions),
		Types:         []SemanticType{},
		Functions:     []SemanticFunction{},
		Dependencies:  []DependencyInfo{},
		FileCount:     len(analysis.Files),
	}

	// Convert types
	for _, t := range analysis.Types {
		st := SemanticType{
			Name:                 t.Name,
			Kind:                 string(t.Kind),
			IsPublic:             t.IsPublic,
			Methods:              t.Methods,
			ImplementsInterfaces: t.ImplementsInterfaces,
			Generic:              t.Generic,
			DocComment:           t.DocComment,
			Location: LocationInfo{
				File:      t.Location.File,
				StartLine: t.Location.StartLine,
				EndLine:   t.Location.EndLine,
			},
		}

		for _, f := range t.ResolvedFields {
			st.Fields = append(st.Fields, SemanticField{
				Name:         f.Name,
				Type:         f.Type,
				ResolvedType: f.ResolvedType,
				Tags:         f.Tags,
				IsPointer:    f.IsPointer,
				IsSlice:      f.IsSlice,
				IsMap:        f.IsMap,
			})
		}

		summary.Types = append(summary.Types, st)
	}

	// Convert functions
	for _, f := range analysis.Functions {
		sf := SemanticFunction{
			Name:        f.Name,
			Signature:   f.Signature,
			IsPublic:    f.IsPublic,
			IsAsync:     f.IsAsync,
			ReturnTypes: f.ResolvedReturnTypes,
			Calls:       f.Calls,
			Complexity:  f.Complexity,
			DocComment:  f.DocComment,
			Location: LocationInfo{
				File:      f.Location.File,
				StartLine: f.Location.StartLine,
				EndLine:   f.Location.EndLine,
			},
		}

		for _, p := range f.ResolvedParameters {
			sf.Parameters = append(sf.Parameters, SemanticParameter{
				Name: p.Name,
				Type: p.ResolvedType,
			})
		}

		summary.Functions = append(summary.Functions, sf)
	}

	// Convert dependencies
	for _, d := range analysis.Dependencies {
		summary.Dependencies = append(summary.Dependencies, DependencyInfo{
			Path:     d.Path,
			Version:  d.Version,
			IsStdLib: d.IsStdLib,
			IsLocal:  d.IsLocal,
		})
	}

	return summary
}

// collectRawFilesForLanguage collects raw file contents for a specific language
func (s *Server) collectRawFilesForLanguage(sourcePath, language string, maxFiles int) []FileContent {
	langExtensions := map[string][]string{
		"sql":        {".sql"},
		"protobuf":   {".proto"},
		"graphql":    {".graphql", ".gql"},
		"yaml":       {".yaml", ".yml"},
		"json":       {".json"},
		"markdown":   {".md"},
		"html":       {".html", ".htm"},
		"css":        {".css", ".scss", ".sass", ".less"},
		"shell":      {".sh", ".bash", ".zsh"},
		"kotlin":     {".kt", ".kts"},
		"swift":      {".swift"},
		"ruby":       {".rb"},
		"php":        {".php"},
		"scala":      {".scala"},
		"c":          {".c", ".h"},
		"cpp":        {".cpp", ".cc", ".cxx", ".hpp", ".hxx"},
		"javascript": {".js", ".jsx", ".mjs", ".cjs"},
	}

	exts, ok := langExtensions[language]
	if !ok {
		return nil
	}

	skipDirs := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, "target": true,
		"build": true, "dist": true, "out": true, ".idea": true, ".vscode": true,
		"__pycache__": true, ".pytest_cache": true, ".mypy_cache": true,
		".next": true, ".nuxt": true, "coverage": true, ".gradle": true,
		".mvn": true, "bin": true, "obj": true, "testdata": true,
	}

	var files []FileContent

	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || skipDirs[name] {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip large files
		if info.Size() > 512*1024 { // 512KB limit for raw files
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		relPath, _ := filepath.Rel(sourcePath, path)

		matches := false
		for _, e := range exts {
			if ext == e {
				matches = true
				break
			}
		}
		if !matches {
			return nil
		}

		// Check limit
		if len(files) >= maxFiles {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		files = append(files, FileContent{
			Path:    relPath,
			Content: string(content),
			Size:    len(content),
		})

		return nil
	})

	return files
}

// buildAIAnalysisPrompt builds a comprehensive prompt for AI-driven spec generation
func (s *Server) buildAIAnalysisPrompt(
	files *importer.ProjectFiles,
	detectedLanguages []LanguageInfo,
	allSourceFiles map[string][]FileContent,
) string {
	var sb strings.Builder

	sb.WriteString("# AI-Driven Spec Generation\n\n")
	sb.WriteString("Analyze the following source code and generate a comprehensive specification that captures:\n")
	sb.WriteString("1. **ALL functions/methods** (public AND private) with their complete signatures, parameters, return types\n")
	sb.WriteString("2. **Behavioral logic** - What each function actually does, not just its signature\n")
	sb.WriteString("3. **Data flow** - How data moves through the system\n")
	sb.WriteString("4. **HTTP endpoints** - Routes, methods, request/response formats\n")
	sb.WriteString("5. **Types/Classes** - All data structures with their fields and relationships\n")
	sb.WriteString("6. **Dependencies** - External libraries and their usage\n")
	sb.WriteString("7. **Configuration** - Settings, environment variables, default values\n")
	sb.WriteString("8. **Entry points** - Main functions, server startup, initialization logic\n\n")

	sb.WriteString("## Project Information\n\n")
	sb.WriteString(fmt.Sprintf("**Name:** %s\n", files.Name))
	sb.WriteString(fmt.Sprintf("**Total Files:** %d source, %d test\n\n", len(files.SourceFiles), len(files.TestFiles)))

	// Language summary
	sb.WriteString("## Languages Detected\n\n")
	for _, lang := range detectedLanguages {
		sb.WriteString(fmt.Sprintf("- **%s**: %d files\n", lang.Name, lang.FileCount))
	}
	sb.WriteString("\n")

	// ALL source files for AI analysis
	sb.WriteString("## Complete Source Code\n\n")
	sb.WriteString("Analyze ALL the following source files to extract specifications:\n\n")

	for langID, langFiles := range allSourceFiles {
		langName := langID
		for _, lang := range detectedLanguages {
			if lang.ID == langID {
				langName = lang.Name
				break
			}
		}

		sb.WriteString(fmt.Sprintf("### %s Source Files\n\n", langName))

		for _, f := range langFiles {
			sb.WriteString(fmt.Sprintf("#### File: `%s`\n", f.Path))
			sb.WriteString(fmt.Sprintf("```%s\n", langID))
			sb.WriteString(f.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	// Configuration files
	if len(files.ConfigFiles) > 0 {
		sb.WriteString("## Configuration Files\n\n")
		for _, f := range files.ConfigFiles {
			sb.WriteString(fmt.Sprintf("### `%s`\n", f.Path))
			sb.WriteString("```\n")
			sb.WriteString(f.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	// API specs if present
	if len(files.APISpecs) > 0 {
		sb.WriteString("## API Specifications\n\n")
		for _, f := range files.APISpecs {
			sb.WriteString(fmt.Sprintf("### `%s`\n", f.Path))
			sb.WriteString("```\n")
			sb.WriteString(f.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	// Documentation files
	if len(files.DocFiles) > 0 {
		sb.WriteString("## Documentation\n\n")
		for _, f := range files.DocFiles {
			sb.WriteString(fmt.Sprintf("### `%s`\n", f.Path))
			sb.WriteString("```markdown\n")
			sb.WriteString(f.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	sb.WriteString("---\n\n")
	sb.WriteString("## Specification Requirements\n\n")
	sb.WriteString("Generate a `.spec.md` file that includes:\n\n")
	sb.WriteString("### Functions Section\n")
	sb.WriteString("For EACH function found (including main, private functions, handlers):\n")
	sb.WriteString("- **Name**: The function name\n")
	sb.WriteString("- **Signature**: Complete signature with types\n")
	sb.WriteString("- **Parameters**: List of parameters with types and descriptions\n")
	sb.WriteString("- **Returns**: Return type and what it represents\n")
	sb.WriteString("- **Logic**: Step-by-step description of what the function does\n")
	sb.WriteString("- **HTTP Mapping**: If it's an HTTP handler, specify the route and method\n\n")

	sb.WriteString("### Types Section\n")
	sb.WriteString("For EACH type/class/struct:\n")
	sb.WriteString("- **Name**: The type name\n")
	sb.WriteString("- **Fields**: All fields with their types\n")
	sb.WriteString("- **Methods**: Associated methods\n")
	sb.WriteString("- **Usage**: How this type is used in the system\n\n")

	sb.WriteString("### Configuration Section\n")
	sb.WriteString("- Default values (ports, addresses, etc.)\n")
	sb.WriteString("- Environment variables\n")
	sb.WriteString("- Runtime settings\n\n")

	sb.WriteString("The spec must be detailed enough to regenerate functionally equivalent code in any language.\n")

	return sb.String()
}

// generateAISpec creates a spec from AI analysis (placeholder for AI integration)
func (s *Server) generateAISpec(
	projectName string,
	primaryLang string,
	detectedLanguages []LanguageInfo,
	allSourceFiles map[string][]FileContent,
	files *importer.ProjectFiles,
	analysisPrompt string,
) string {
	// This is a placeholder that creates a basic spec
	// In production, this would call an AI model with the analysisPrompt
	// For now, we'll create a structured spec that the AI can fill in

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", projectName))
	sb.WriteString("## Overview\n\n")
	sb.WriteString("*[AI: Describe the project's purpose and architecture based on the source code analysis]*\n\n")

	sb.WriteString("## Configuration\n\n")
	sb.WriteString("*[AI: Extract configuration settings, environment variables, and defaults from the code]*\n\n")

	sb.WriteString("## Types\n\n")
	sb.WriteString("*[AI: List all data structures, classes, and types with their fields and relationships]*\n\n")

	sb.WriteString("## Functions\n\n")
	sb.WriteString("*[AI: List ALL functions (public and private) with:]*\n")
	sb.WriteString("- *Name and signature*\n")
	sb.WriteString("- *Parameters with types*\n")
	sb.WriteString("- *Return values*\n")
	sb.WriteString("- *Behavioral logic (what it does)*\n")
	sb.WriteString("- *HTTP endpoints if applicable*\n\n")

	sb.WriteString("### Entry Points\n\n")
	sb.WriteString("*[AI: Identify main functions and server initialization]*\n\n")

	sb.WriteString("### HTTP Endpoints\n\n")
	sb.WriteString("*[AI: List all HTTP routes with methods and handlers]*\n\n")

	sb.WriteString("## Dependencies\n\n")
	sb.WriteString("*[AI: List external libraries and their usage]*\n\n")

	sb.WriteString("## Tests\n\n")
	sb.WriteString("*[AI: Describe test cases if present]*\n\n")

	sb.WriteString("---\n")
	sb.WriteString("*Note: This spec should be processed by an AI model using the analysis prompt to fill in the actual details.*\n")

	return sb.String()
}

// buildComprehensiveAnalysisPrompt builds an analysis prompt that includes all semantic and raw file analysis
func (s *Server) buildComprehensiveAnalysisPrompt(
	files *importer.ProjectFiles,
	detectedLanguages []LanguageInfo,
	semanticAnalyses map[string]*SemanticSummary,
	rawFilesByLanguage map[string][]FileContent,
) string {
	var sb strings.Builder

	sb.WriteString("# Project Analysis for Spec Generation\n\n")
	sb.WriteString(fmt.Sprintf("**Project:** %s\n\n", files.Name))

	// Language summary
	sb.WriteString("## Detected Languages\n\n")
	sb.WriteString("| Language | Files | Analysis Method |\n")
	sb.WriteString("|----------|-------|----------------|\n")
	for _, lang := range detectedLanguages {
		method := "Raw files"
		if lang.HasParser {
			method = "Semantic AST analysis"
		}
		sb.WriteString(fmt.Sprintf("| %s | %d | %s |\n", lang.Name, lang.FileCount, method))
	}
	sb.WriteString("\n")

	// Semantic analysis results (structured data from AST parsing)
	if len(semanticAnalyses) > 0 {
		sb.WriteString("## Semantic Analysis Results\n\n")
		sb.WriteString("The following languages were analyzed using AST parsing for precise type and function extraction:\n\n")

		for langID, analysis := range semanticAnalyses {
			sb.WriteString(fmt.Sprintf("### %s Analysis\n\n", strings.ToUpper(langID)))
			sb.WriteString(fmt.Sprintf("- **Files:** %d\n", analysis.FileCount))
			sb.WriteString(fmt.Sprintf("- **Types:** %d\n", analysis.TypeCount))
			sb.WriteString(fmt.Sprintf("- **Functions:** %d\n\n", analysis.FunctionCount))

			// Types
			if len(analysis.Types) > 0 {
				sb.WriteString("#### Types\n\n")
				for _, t := range analysis.Types {
					visibility := "private"
					if t.IsPublic {
						visibility = "public"
					}
					sb.WriteString(fmt.Sprintf("- **%s** `%s` (%s) - `%s:%d`\n", t.Kind, t.Name, visibility, t.Location.File, t.Location.StartLine))
					if len(t.Fields) > 0 {
						for _, f := range t.Fields {
							sb.WriteString(fmt.Sprintf("  - `%s`: %s\n", f.Name, f.ResolvedType))
						}
					}
				}
				sb.WriteString("\n")
			}

			// Functions (summarized)
			if len(analysis.Functions) > 0 {
				sb.WriteString("#### Functions\n\n")
				publicFuncs := 0
				for _, f := range analysis.Functions {
					if f.IsPublic {
						publicFuncs++
						sb.WriteString(fmt.Sprintf("- `%s` - `%s:%d`\n", f.Signature, f.Location.File, f.Location.StartLine))
					}
				}
				if publicFuncs < len(analysis.Functions) {
					sb.WriteString(fmt.Sprintf("\n*Plus %d private functions*\n", len(analysis.Functions)-publicFuncs))
				}
				sb.WriteString("\n")
			}

			// External dependencies
			if len(analysis.Dependencies) > 0 {
				sb.WriteString("#### External Dependencies\n\n")
				for _, d := range analysis.Dependencies {
					if !d.IsStdLib && !d.IsLocal {
						sb.WriteString(fmt.Sprintf("- %s", d.Path))
						if d.Version != "" {
							sb.WriteString(fmt.Sprintf(" (%s)", d.Version))
						}
						sb.WriteString("\n")
					}
				}
				sb.WriteString("\n")
			}
		}
	}

	// Raw file contents for languages without parsers
	if len(rawFilesByLanguage) > 0 {
		sb.WriteString("## Raw File Contents\n\n")
		sb.WriteString("The following files are provided for AI interpretation (no AST parser available):\n\n")

		for langID, langFiles := range rawFilesByLanguage {
			sb.WriteString(fmt.Sprintf("### %s Files\n\n", strings.ToUpper(langID)))
			for _, f := range langFiles {
				sb.WriteString(fmt.Sprintf("#### `%s`\n", f.Path))
				sb.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", langID, f.Content))
			}
		}
	}

	// Original collected files (config, docs, API specs)
	if len(files.ConfigFiles) > 0 {
		sb.WriteString("## Configuration Files\n\n")
		for _, f := range files.ConfigFiles {
			sb.WriteString(fmt.Sprintf("### `%s`\n", f.Path))
			sb.WriteString(fmt.Sprintf("```\n%s\n```\n\n", f.Content))
		}
	}

	if len(files.APISpecs) > 0 {
		sb.WriteString("## API Specifications\n\n")
		for _, f := range files.APISpecs {
			sb.WriteString(fmt.Sprintf("### `%s`\n", f.Path))
			sb.WriteString(fmt.Sprintf("```\n%s\n```\n\n", f.Content))
		}
	}

	if len(files.DocFiles) > 0 {
		sb.WriteString("## Documentation\n\n")
		for _, f := range files.DocFiles {
			sb.WriteString(fmt.Sprintf("### `%s`\n", f.Path))
			sb.WriteString(fmt.Sprintf("```markdown\n%s\n```\n\n", f.Content))
		}
	}

	// Instructions for AI
	sb.WriteString("---\n\n")
	sb.WriteString("## Spec Generation Instructions\n\n")
	sb.WriteString("Based on the above analysis, generate a comprehensive `.spec.md` file that includes:\n\n")
	sb.WriteString("1. **Project Overview**: Name, purpose, and architecture\n")
	sb.WriteString("2. **Data Model**: All types, entities, and their relationships\n")
	sb.WriteString("3. **API Surface**: All public functions, endpoints, and interfaces\n")
	sb.WriteString("4. **Dependencies**: External libraries and their purpose\n")
	sb.WriteString("5. **Configuration**: Environment variables and settings\n")
	sb.WriteString("6. **Behavioral Specifications**: Key workflows and business logic\n\n")
	sb.WriteString("The spec should be detailed enough to regenerate equivalent code in any target language.\n")

	return sb.String()
}

// generateSpecFromAnalysis creates a .spec.md file from the semantic analysis results
func (s *Server) generateSpecFromAnalysis(
	projectName string,
	primaryLang string,
	detectedLanguages []LanguageInfo,
	semanticAnalyses map[string]*SemanticSummary,
	rawFilesByLanguage map[string][]FileContent,
	files *importer.ProjectFiles,
) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# %s\n\n", projectName))

	// Overview section - try to extract from README if available
	sb.WriteString("## Overview\n\n")
	if readme := findReadmeContent(rawFilesByLanguage); readme != "" {
		// Extract first paragraph from README
		lines := strings.Split(readme, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				sb.WriteString(line + "\n\n")
				break
			}
		}
	} else {
		sb.WriteString(fmt.Sprintf("A %s project with %d source files.\n\n", strings.ToUpper(primaryLang), len(files.SourceFiles)))
	}

	// Languages section
	sb.WriteString("## Languages\n\n")
	for _, lang := range detectedLanguages {
		sb.WriteString(fmt.Sprintf("- **%s**: %d files\n", lang.Name, lang.FileCount))
	}
	sb.WriteString("\n")

	// Types section from semantic analysis
	for langID, analysis := range semanticAnalyses {
		if len(analysis.Types) > 0 {
			sb.WriteString(fmt.Sprintf("## Types (%s)\n\n", strings.ToUpper(langID)))

			// Group by kind
			structs := []SemanticType{}
			interfaces := []SemanticType{}
			enums := []SemanticType{}
			others := []SemanticType{}

			for _, t := range analysis.Types {
				if !t.IsPublic {
					continue // Only include public types
				}
				switch t.Kind {
				case "struct":
					structs = append(structs, t)
				case "interface":
					interfaces = append(interfaces, t)
				case "enum":
					enums = append(enums, t)
				default:
					others = append(others, t)
				}
			}

			// Structs/Classes
			if len(structs) > 0 {
				sb.WriteString("### Data Structures\n\n")
				for _, t := range structs {
					sb.WriteString(fmt.Sprintf("#### %s\n\n", t.Name))
					if t.DocComment != "" {
						sb.WriteString(fmt.Sprintf("%s\n\n", strings.TrimSpace(t.DocComment)))
					}
					if len(t.Fields) > 0 {
						sb.WriteString("| Field | Type | Description |\n")
						sb.WriteString("|-------|------|-------------|\n")
						for _, f := range t.Fields {
							fieldType := f.ResolvedType
							if f.IsPointer {
								fieldType = "*" + fieldType
							}
							if f.IsSlice {
								fieldType = "[]" + fieldType
							}
							if f.IsMap {
								fieldType = "map[...]" + fieldType
							}
							sb.WriteString(fmt.Sprintf("| %s | %s | |\n", f.Name, fieldType))
						}
						sb.WriteString("\n")
					}
				}
			}

			// Interfaces
			if len(interfaces) > 0 {
				sb.WriteString("### Interfaces\n\n")
				for _, t := range interfaces {
					sb.WriteString(fmt.Sprintf("#### %s\n\n", t.Name))
					if t.DocComment != "" {
						sb.WriteString(fmt.Sprintf("%s\n\n", strings.TrimSpace(t.DocComment)))
					}
					if len(t.Methods) > 0 {
						sb.WriteString("**Methods:**\n")
						for _, m := range t.Methods {
							sb.WriteString(fmt.Sprintf("- `%s`\n", m))
						}
						sb.WriteString("\n")
					}
				}
			}

			// Enums
			if len(enums) > 0 {
				sb.WriteString("### Enums\n\n")
				for _, t := range enums {
					sb.WriteString(fmt.Sprintf("- **%s**\n", t.Name))
				}
				sb.WriteString("\n")
			}
		}

		// Functions section
		if len(analysis.Functions) > 0 {
			sb.WriteString(fmt.Sprintf("## Functions (%s)\n\n", strings.ToUpper(langID)))

			// Group public functions by file/module
			publicFuncs := []SemanticFunction{}
			for _, f := range analysis.Functions {
				if f.IsPublic {
					publicFuncs = append(publicFuncs, f)
				}
			}

			if len(publicFuncs) > 0 {
				sb.WriteString("### Public API\n\n")
				for _, f := range publicFuncs {
					sb.WriteString(fmt.Sprintf("#### `%s`\n\n", f.Name))
					sb.WriteString(fmt.Sprintf("**Signature:** `%s`\n\n", f.Signature))
					if f.DocComment != "" {
						sb.WriteString(fmt.Sprintf("%s\n\n", strings.TrimSpace(f.DocComment)))
					}
					if len(f.Parameters) > 0 {
						sb.WriteString("**Parameters:**\n")
						for _, p := range f.Parameters {
							sb.WriteString(fmt.Sprintf("- `%s`: %s\n", p.Name, p.Type))
						}
						sb.WriteString("\n")
					}
					if len(f.ReturnTypes) > 0 {
						sb.WriteString(fmt.Sprintf("**Returns:** %s\n\n", strings.Join(f.ReturnTypes, ", ")))
					}
				}
			}

			sb.WriteString(fmt.Sprintf("*Total: %d public functions, %d private functions*\n\n", len(publicFuncs), len(analysis.Functions)-len(publicFuncs)))
		}

		// Dependencies section
		if len(analysis.Dependencies) > 0 {
			sb.WriteString(fmt.Sprintf("## Dependencies (%s)\n\n", strings.ToUpper(langID)))

			// External dependencies
			sb.WriteString("### External\n\n")
			for _, d := range analysis.Dependencies {
				if !d.IsStdLib && !d.IsLocal {
					sb.WriteString(fmt.Sprintf("- `%s`", d.Path))
					if d.Version != "" {
						sb.WriteString(fmt.Sprintf(" (%s)", d.Version))
					}
					sb.WriteString("\n")
				}
			}
			sb.WriteString("\n")
		}
	}

	// Database Schema from SQL files
	if sqlFiles, ok := rawFilesByLanguage["sql"]; ok && len(sqlFiles) > 0 {
		sb.WriteString("## Database Schema\n\n")
		for _, f := range sqlFiles {
			if strings.Contains(f.Path, "migration") || strings.Contains(f.Path, "schema") {
				sb.WriteString(fmt.Sprintf("### %s\n\n", filepath.Base(f.Path)))
				sb.WriteString("```sql\n")
				sb.WriteString(f.Content)
				sb.WriteString("\n```\n\n")
			}
		}
	}

	// Configuration from YAML/config files
	if len(files.ConfigFiles) > 0 {
		sb.WriteString("## Configuration\n\n")
		for _, f := range files.ConfigFiles {
			if strings.Contains(f.Path, "config") || strings.HasSuffix(f.Path, ".env.example") {
				sb.WriteString(fmt.Sprintf("### %s\n\n", f.Path))
				sb.WriteString("```\n")
				sb.WriteString(f.Content)
				sb.WriteString("\n```\n\n")
			}
		}
	}

	// API Specifications
	if len(files.APISpecs) > 0 {
		sb.WriteString("## API Specifications\n\n")
		for _, f := range files.APISpecs {
			sb.WriteString(fmt.Sprintf("### %s\n\n", f.Path))
			sb.WriteString("```yaml\n")
			sb.WriteString(f.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	return sb.String()
}

// findReadmeContent extracts README content from raw files
func findReadmeContent(rawFilesByLanguage map[string][]FileContent) string {
	if mdFiles, ok := rawFilesByLanguage["markdown"]; ok {
		for _, f := range mdFiles {
			if strings.ToLower(filepath.Base(f.Path)) == "readme.md" {
				return f.Content
			}
		}
	}
	return ""
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

	// Delegate to handleImportSpecFromSource for the actual import logic
	sourceInput := ImportSpecFromSourceInput{
		InputPath: cloneResult.LocalPath,
		Name:      input.Name,
	}

	// Use repo name as default if no name provided
	if sourceInput.Name == "" {
		sourceInput.Name = repoInfo.Name
	}

	// Call the source import handler
	result, sourceOutput, err := s.handleImportSpecFromSource(ctx, req, sourceInput)
	if err != nil {
		return nil, ImportSpecFromGitHubOutput{}, err
	}

	// If source import returned an error result, propagate it
	if result != nil && result.IsError {
		return result, ImportSpecFromGitHubOutput{}, nil
	}

	// Map the source output to GitHub output, adding GitHub-specific fields
	return nil, ImportSpecFromGitHubOutput{
		ProjectName:      sourceOutput.ProjectName,
		Repository:       fmt.Sprintf("%s/%s", repoInfo.Owner, repoInfo.Name),
		Branch:           cloneResult.Branch,
		CommitSHA:        cloneResult.CommitSHA,
		DetectedLanguage: sourceOutput.DetectedLanguage,
		SourceFileCount:  sourceOutput.SourceFileCount,
		TestFileCount:    sourceOutput.TestFileCount,
		APISpecCount:     sourceOutput.APISpecCount,
		ConfigFileCount:  sourceOutput.ConfigFileCount,
		DocFileCount:     sourceOutput.DocFileCount,
		TotalContentSize: sourceOutput.TotalContentSize,
		AnalysisPrompt:   sourceOutput.AnalysisPrompt,
		SpecOutputPath:   sourceOutput.SpecOutputPath,
	}, nil
}

func (s *Server) handleDeepAnalyzeSource(ctx context.Context, req *mcp.CallToolRequest, input DeepAnalyzeSourceInput) (*mcp.CallToolResult, DeepAnalyzeSourceOutput, error) {
	// Expand ~ to home directory if present
	sourcePath := input.SourcePath
	if strings.HasPrefix(sourcePath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			sourcePath = filepath.Join(homeDir, sourcePath[2:])
		}
	}

	// Convert to absolute path if relative
	if !filepath.IsAbs(sourcePath) {
		absPath, err := filepath.Abs(sourcePath)
		if err == nil {
			sourcePath = absPath
		}
	}

	// Validate source path exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Source directory does not exist: %s", sourcePath)},
			},
		}, DeepAnalyzeSourceOutput{}, nil
	}

	// Determine analysis depth (default: standard)
	analysisDepth := input.AnalysisDepth
	if analysisDepth == "" {
		analysisDepth = "standard"
	}

	// Detect or use provided language
	language := input.Language
	if language == "" {
		language = s.detectPrimaryLanguage(sourcePath)
	}

	// Create semantic analyzer registry
	registry := semantic.DefaultRegistry()

	// Get the appropriate analyzer
	lang := treesitter.Language(language)
	analyzer, ok := registry.Get(lang)
	if !ok {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("No analyzer available for language: %s", language)},
			},
		}, DeepAnalyzeSourceOutput{}, nil
	}

	// Check if analyzer is available
	if !analyzer.IsAvailable() {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Analyzer for %s is not available (missing dependencies)", language)},
			},
		}, DeepAnalyzeSourceOutput{}, nil
	}

	// Perform semantic analysis
	analysis, err := analyzer.Analyze(sourcePath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Analysis failed: %v", err)},
			},
		}, DeepAnalyzeSourceOutput{}, nil
	}

	// Convert to output format
	output := s.convertAnalysisToOutput(analysis, analysisDepth)

	// Generate analysis prompt for AI
	output.AnalysisPrompt = s.buildDeepAnalysisPrompt(analysis, analysisDepth)

	return nil, output, nil
}

// detectPrimaryLanguage detects the primary language of a source directory
func (s *Server) detectPrimaryLanguage(dir string) string {
	extensionCounts := make(map[string]int)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".go":
			extensionCounts["go"]++
		case ".ts", ".tsx":
			extensionCounts["typescript"]++
		case ".py":
			extensionCounts["python"]++
		case ".java":
			extensionCounts["java"]++
		case ".rs":
			extensionCounts["rust"]++
		case ".cs":
			extensionCounts["csharp"]++
		}
		return nil
	})

	// Find language with most files
	maxCount := 0
	primaryLang := "go" // default
	for lang, count := range extensionCounts {
		if count > maxCount {
			maxCount = count
			primaryLang = lang
		}
	}

	return primaryLang
}

// convertAnalysisToOutput converts semantic analysis to MCP output format
func (s *Server) convertAnalysisToOutput(analysis *semantic.Analysis, depth string) DeepAnalyzeSourceOutput {
	output := DeepAnalyzeSourceOutput{
		ProjectName:   analysis.Name,
		Language:      string(analysis.Language),
		AnalysisDepth: depth,
		FileCount:     len(analysis.Files),
		Types:         []SemanticType{},      // Initialize to empty array
		Functions:     []SemanticFunction{},  // Initialize to empty array
		Dependencies:  []DependencyInfo{},    // Initialize to empty array
		Errors:        []string{},            // Initialize to empty array
	}

	// Convert types
	for _, t := range analysis.Types {
		st := SemanticType{
			Name:                 t.Name,
			Kind:                 string(t.Kind),
			IsPublic:             t.IsPublic,
			Methods:              t.Methods,
			ImplementsInterfaces: t.ImplementsInterfaces,
			Generic:              t.Generic,
			DocComment:           t.DocComment,
			Location: LocationInfo{
				File:      t.Location.File,
				StartLine: t.Location.StartLine,
				EndLine:   t.Location.EndLine,
			},
		}

		// Convert fields
		for _, f := range t.ResolvedFields {
			st.Fields = append(st.Fields, SemanticField{
				Name:         f.Name,
				Type:         f.Type,
				ResolvedType: f.ResolvedType,
				Tags:         f.Tags,
				IsPointer:    f.IsPointer,
				IsSlice:      f.IsSlice,
				IsMap:        f.IsMap,
			})
		}

		output.Types = append(output.Types, st)
	}

	// Convert functions
	for _, f := range analysis.Functions {
		sf := SemanticFunction{
			Name:        f.Name,
			Signature:   f.Signature,
			IsPublic:    f.IsPublic,
			IsAsync:     f.IsAsync,
			ReturnTypes: f.ResolvedReturnTypes,
			Calls:       f.Calls,
			Complexity:  f.Complexity,
			DocComment:  f.DocComment,
			Location: LocationInfo{
				File:      f.Location.File,
				StartLine: f.Location.StartLine,
				EndLine:   f.Location.EndLine,
			},
		}

		// Convert parameters
		for _, p := range f.ResolvedParameters {
			sf.Parameters = append(sf.Parameters, SemanticParameter{
				Name: p.Name,
				Type: p.ResolvedType,
			})
		}

		output.Functions = append(output.Functions, sf)
	}

	// Include graphs for deep analysis
	if depth == "deep" {
		// Filter out nil values to ensure valid JSON schema
		if analysis.CallGraph != nil {
			output.CallGraph = make(map[string][]string)
			for k, v := range analysis.CallGraph {
				if v != nil {
					output.CallGraph[k] = v
				}
			}
		}
		if analysis.TypeGraph != nil {
			output.TypeGraph = make(map[string][]string)
			for k, v := range analysis.TypeGraph {
				if v != nil {
					output.TypeGraph[k] = v
				}
			}
		}
	}

	// Convert dependencies
	for _, d := range analysis.Dependencies {
		output.Dependencies = append(output.Dependencies, DependencyInfo{
			Path:     d.Path,
			Version:  d.Version,
			IsStdLib: d.IsStdLib,
			IsLocal:  d.IsLocal,
		})
	}

	// Convert errors
	for _, e := range analysis.Errors {
		output.Errors = append(output.Errors, fmt.Sprintf("[%s] %s: %s", e.Severity, e.File, e.Message))
	}

	return output
}

// buildDeepAnalysisPrompt generates an AI prompt from the semantic analysis
func (s *Server) buildDeepAnalysisPrompt(analysis *semantic.Analysis, depth string) string {
	var sb strings.Builder

	sb.WriteString("# Deep Semantic Analysis Results\n\n")
	sb.WriteString(fmt.Sprintf("**Project:** %s\n", analysis.Name))
	sb.WriteString(fmt.Sprintf("**Language:** %s\n", analysis.Language))
	sb.WriteString(fmt.Sprintf("**Files Analyzed:** %d\n\n", len(analysis.Files)))

	// Types section
	if len(analysis.Types) > 0 {
		sb.WriteString("## Type Definitions\n\n")
		for _, t := range analysis.Types {
			visibility := "private"
			if t.IsPublic {
				visibility = "public"
			}
			sb.WriteString(fmt.Sprintf("### %s %s (%s)\n", t.Kind, t.Name, visibility))
			if t.DocComment != "" {
				sb.WriteString(fmt.Sprintf("*%s*\n", strings.TrimSpace(t.DocComment)))
			}
			sb.WriteString(fmt.Sprintf("Location: `%s:%d`\n\n", t.Location.File, t.Location.StartLine))

			if len(t.ResolvedFields) > 0 {
				sb.WriteString("**Fields:**\n")
				for _, f := range t.ResolvedFields {
					typeInfo := f.ResolvedType
					if f.IsPointer {
						typeInfo = "*" + typeInfo
					}
					if f.IsSlice {
						typeInfo = "[]" + typeInfo
					}
					sb.WriteString(fmt.Sprintf("- `%s`: %s", f.Name, typeInfo))
					if f.Tags != "" {
						sb.WriteString(fmt.Sprintf(" `%s`", f.Tags))
					}
					sb.WriteString("\n")
				}
				sb.WriteString("\n")
			}

			if len(t.Methods) > 0 {
				sb.WriteString(fmt.Sprintf("**Methods:** %s\n\n", strings.Join(t.Methods, ", ")))
			}

			if len(t.ImplementsInterfaces) > 0 {
				sb.WriteString(fmt.Sprintf("**Implements:** %s\n\n", strings.Join(t.ImplementsInterfaces, ", ")))
			}
		}
	}

	// Functions section
	if len(analysis.Functions) > 0 {
		sb.WriteString("## Functions\n\n")
		for _, f := range analysis.Functions {
			visibility := "private"
			if f.IsPublic {
				visibility = "public"
			}
			asyncStr := ""
			if f.IsAsync {
				asyncStr = "async "
			}
			sb.WriteString(fmt.Sprintf("### %s%s (%s)\n", asyncStr, f.Name, visibility))
			sb.WriteString(fmt.Sprintf("**Signature:** `%s`\n", f.Signature))
			if f.DocComment != "" {
				sb.WriteString(fmt.Sprintf("*%s*\n", strings.TrimSpace(f.DocComment)))
			}
			sb.WriteString(fmt.Sprintf("Location: `%s:%d`\n", f.Location.File, f.Location.StartLine))
			if f.Complexity > 1 {
				sb.WriteString(fmt.Sprintf("Complexity: %d\n", f.Complexity))
			}
			if len(f.Calls) > 0 && depth == "deep" {
				sb.WriteString(fmt.Sprintf("Calls: %s\n", strings.Join(f.Calls, ", ")))
			}
			sb.WriteString("\n")
		}
	}

	// Dependencies section
	if len(analysis.Dependencies) > 0 {
		sb.WriteString("## Dependencies\n\n")
		sb.WriteString("**External:**\n")
		for _, d := range analysis.Dependencies {
			if !d.IsStdLib && !d.IsLocal {
				sb.WriteString(fmt.Sprintf("- %s", d.Path))
				if d.Version != "" {
					sb.WriteString(fmt.Sprintf(" (%s)", d.Version))
				}
				sb.WriteString("\n")
			}
		}
		sb.WriteString("\n**Standard Library:**\n")
		for _, d := range analysis.Dependencies {
			if d.IsStdLib {
				sb.WriteString(fmt.Sprintf("- %s\n", d.Path))
			}
		}
		sb.WriteString("\n")
	}

	// Call graph for deep analysis
	if depth == "deep" && len(analysis.CallGraph) > 0 {
		sb.WriteString("## Call Graph\n\n")
		for fn, calls := range analysis.CallGraph {
			if len(calls) > 0 {
				sb.WriteString(fmt.Sprintf("- **%s** → %s\n", fn, strings.Join(calls, ", ")))
			}
		}
		sb.WriteString("\n")
	}

	// Type graph for deep analysis
	if depth == "deep" && len(analysis.TypeGraph) > 0 {
		sb.WriteString("## Type Relationships\n\n")
		for typ, interfaces := range analysis.TypeGraph {
			sb.WriteString(fmt.Sprintf("- **%s** implements %s\n", typ, strings.Join(interfaces, ", ")))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	sb.WriteString("Use this semantic analysis to generate an accurate spec file or to compare against other implementations.\n")

	return sb.String()
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

// handleSemanticParityAnalysis performs deep semantic parity analysis between source and generated code
func (s *Server) handleSemanticParityAnalysis(ctx context.Context, req *mcp.CallToolRequest, input SemanticParityAnalysisInput) (*mcp.CallToolResult, SemanticParityAnalysisOutput, error) {
	// Expand and validate source path
	sourcePath := expandPath(input.SourcePath)
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Source directory does not exist: %s", sourcePath)},
			},
		}, SemanticParityAnalysisOutput{}, nil
	}

	// Detect source language if not provided
	sourceLang := input.SourceLanguage
	if sourceLang == "" {
		sourceLang = s.detectPrimaryLanguage(sourcePath)
	}

	// Create semantic analyzer registry
	registry := semantic.DefaultRegistry()

	// Analyze source code
	sourceAnalyzer, ok := registry.Get(treesitter.Language(sourceLang))
	if !ok {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("No analyzer available for source language: %s", sourceLang)},
			},
		}, SemanticParityAnalysisOutput{}, nil
	}

	sourceAnalysis, err := sourceAnalyzer.Analyze(sourcePath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to analyze source: %v", err)},
			},
		}, SemanticParityAnalysisOutput{}, nil
	}

	// Analyze generated projects
	generatedAnalyses := make(map[string]*semantic.Analysis)
	for _, proj := range input.GeneratedProjects {
		projPath := expandPath(proj.Path)
		if _, err := os.Stat(projPath); os.IsNotExist(err) {
			continue // Skip non-existent projects
		}

		analyzer, ok := registry.Get(treesitter.Language(proj.Language))
		if !ok {
			continue
		}

		analysis, err := analyzer.Analyze(projPath)
		if err != nil {
			continue
		}
		generatedAnalyses[proj.Language] = analysis
	}

	// Configure parity comparison
	parityConfig := parity.DefaultConfig()
	if input.ComparisonWeights != nil {
		if input.ComparisonWeights.Structural > 0 {
			parityConfig.Weights.Structural = input.ComparisonWeights.Structural
		}
		if input.ComparisonWeights.Type > 0 {
			parityConfig.Weights.Type = input.ComparisonWeights.Type
		}
		if input.ComparisonWeights.Behavioral > 0 {
			parityConfig.Weights.Behavioral = input.ComparisonWeights.Behavioral
		}
		if input.ComparisonWeights.Test > 0 {
			parityConfig.Weights.Test = input.ComparisonWeights.Test
		}
		if input.ComparisonWeights.Idiomatic > 0 {
			parityConfig.Weights.Idiomatic = input.ComparisonWeights.Idiomatic
		}
	}

	// Perform parity comparison
	comparator := parity.NewComparator(parityConfig)
	result := comparator.Compare(sourceAnalysis, generatedAnalyses)

	// Convert to output format
	output := SemanticParityAnalysisOutput{
		OverallScore: result.OverallScore,
		Converged:    result.Converged,
		ByDimension: DimensionScoresOutput{
			Structural: result.ByDimension.Structural,
			Type:       result.ByDimension.Type,
			Behavioral: result.ByDimension.Behavioral,
			Test:       result.ByDimension.Test,
			Idiomatic:  result.ByDimension.Idiomatic,
		},
		ByLanguage: make(map[string]LanguageParityResult),
		Gaps:       []SemanticParityGap{},
	}

	// Convert per-language results
	for lang, lr := range result.ByLanguage {
		output.ByLanguage[lang] = LanguageParityResult{
			Language:     lang,
			OverallScore: lr.OverallScore,
			ByDimension: DimensionScoresOutput{
				Structural: lr.ByDimension.Structural,
				Type:       lr.ByDimension.Type,
				Behavioral: lr.ByDimension.Behavioral,
				Test:       lr.ByDimension.Test,
				Idiomatic:  lr.ByDimension.Idiomatic,
			},
			MissingTypes: lr.MissingTypes,
			MissingFuncs: lr.MissingFuncs,
		}
	}

	// Convert gaps
	for _, gap := range result.Gaps {
		lang := ""
		genItem := ""
		if gap.GeneratedItem != nil {
			lang = gap.GeneratedItem.Language
			genItem = gap.GeneratedItem.Name
		}
		output.Gaps = append(output.Gaps, SemanticParityGap{
			Dimension:     gap.Dimension,
			Severity:      gap.Severity,
			SourceItem:    gap.SourceItem.Name,
			GeneratedItem: genItem,
			Discrepancy:   gap.Discrepancy,
			SuggestedFix:  gap.SuggestedFix,
			Language:      lang,
		})
	}

	// Generate fix instructions
	output.FixInstructions = parity.GenerateFixInstructions(result, sourceLang)

	return nil, output, nil
}

// handleIterativeRefinementLoop orchestrates the full refinement loop
func (s *Server) handleIterativeRefinementLoop(ctx context.Context, req *mcp.CallToolRequest, input IterativeRefinementLoopInput) (*mcp.CallToolResult, IterativeRefinementLoopOutput, error) {
	// Expand and validate source path
	sourcePath := expandPath(input.SourcePath)
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Source directory does not exist: %s", sourcePath)},
			},
		}, IterativeRefinementLoopOutput{}, nil
	}

	// Expand output directory
	outputDir := expandPath(input.OutputDir)

	// Configure refinement loop
	config := refinement.DefaultLoopConfig()
	if input.ConvergenceThreshold > 0 {
		config.ConvergenceThreshold = input.ConvergenceThreshold
	}
	if input.MaxIterations > 0 {
		config.MaxIterations = input.MaxIterations
	}
	if input.RefinementStrategy != "" {
		config.RefinementStrategy = refinement.Strategy(input.RefinementStrategy)
	}

	// Create refinement engine
	registry := semantic.DefaultRegistry()
	engine := refinement.NewEngine(config, registry)

	// Create loop input
	loopInput := &refinement.LoopInput{
		SourcePath:      sourcePath,
		SourceLanguage:  input.SourceLanguage,
		TargetLanguages: input.TargetLanguages,
		SpecPath:        input.SpecPath,
		OutputDir:       outputDir,
	}

	// Run refinement loop
	result, err := engine.Run(ctx, loopInput)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Refinement loop failed: %v", err)},
			},
		}, IterativeRefinementLoopOutput{}, nil
	}

	// Convert to output format
	output := IterativeRefinementLoopOutput{
		Converged:         result.Converged,
		FinalScore:        result.FinalScore,
		IterationsUsed:    result.IterationsUsed,
		FinalSpec:         result.FinalSpec,
		GeneratedProjects: result.GeneratedProjects,
		RefinementSummary: result.RefinementSummary,
	}

	// Convert iteration history
	for _, iter := range result.IterationHistory {
		output.IterationHistory = append(output.IterationHistory, IterationSummary{
			Number:             iter.Number,
			ParityScore:        iter.ParityScore,
			Phase:              iter.Phase,
			RefinementsApplied: iter.RefinementsApplied,
			ScoreImprovement:   iter.ScoreImprovement,
		})
	}

	// Convert unresolved gaps
	for _, gap := range result.UnresolvedGaps {
		lang := ""
		genItem := ""
		if gap.GeneratedItem != nil {
			lang = gap.GeneratedItem.Language
			genItem = gap.GeneratedItem.Name
		}
		output.UnresolvedGaps = append(output.UnresolvedGaps, SemanticParityGap{
			Dimension:     gap.Dimension,
			Severity:      gap.Severity,
			SourceItem:    gap.SourceItem.Name,
			GeneratedItem: genItem,
			Discrepancy:   gap.Discrepancy,
			SuggestedFix:  gap.SuggestedFix,
			Language:      lang,
		})
	}

	// Generate refinement prompt if not converged
	if !result.Converged && len(result.IterationHistory) > 0 {
		lastIter := result.IterationHistory[len(result.IterationHistory)-1]
		if lastIter.ParityResult != nil {
			instructions := &refinement.RefinementInstructions{
				Summary:  fmt.Sprintf("Refinement needed: current score %.1f%%, target %.1f%%", result.FinalScore*100, config.ConvergenceThreshold*100),
				Priority: lastIter.Phase,
			}
			output.RefinementPrompt = engine.GenerateRefinementPrompt(instructions, input.SourceLanguage)
		}
	}

	return nil, output, nil
}

// =============================================================================
// AUTONOMOUS CODE GENERATION HANDLER
// =============================================================================

// handleGenerateSourceFromSpec reads the spec and returns it with generation instructions for AI
func (s *Server) handleGenerateSourceFromSpec(ctx context.Context, req *mcp.CallToolRequest, input GenerateSourceFromSpecInput) (*mcp.CallToolResult, GenerateSourceFromSpecOutput, error) {
	output := GenerateSourceFromSpecOutput{}

	// Read spec file
	specPath := expandPath(input.SpecPath)
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to read spec file: %v", err)},
			},
		}, output, nil
	}

	output.SpecContent = string(specContent)
	output.SpecPath = specPath

	// Get language adapter
	adapter, err := s.registry.Get(input.Language)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Unsupported language: %s. Use list_languages to see supported languages.", input.Language)},
			},
		}, output, nil
	}

	output.Language = adapter.GetLanguage()
	output.PromptTemplate = adapter.GetPromptContext()

	// Determine output directory
	specName := strings.TrimSuffix(filepath.Base(specPath), ".spec.md")
	specName = strings.TrimSuffix(specName, ".md")
	outputDir := input.OutputDir
	if outputDir == "" {
		outputDir = filepath.Join(s.outputDir, specName, input.Language)
	}
	outputDir = expandPath(outputDir)
	output.OutputDir = outputDir

	// Get recommended project structure (ensure non-nil for JSON)
	output.ProjectStructure = adapter.GetProjectStructure(specName, true)
	if output.ProjectStructure == nil {
		output.ProjectStructure = []languages.ProjectFile{}
	}

	// Build comprehensive generation prompt
	output.GenerationPrompt = buildAIGenerationPrompt(output.SpecContent, output.Language, output.OutputDir, output.ProjectStructure)

	return nil, output, nil
}

// buildAIGenerationPrompt creates a comprehensive prompt for AI to generate code from the spec
func buildAIGenerationPrompt(specContent string, lang languages.Language, outputDir string, structure []languages.ProjectFile) string {
	var sb strings.Builder

	sb.WriteString("# Code Generation Task with Parity Loop\n\n")
	sb.WriteString("Generate complete, idiomatic ")
	sb.WriteString(lang.Name)
	sb.WriteString(" code from the specification below. You MUST achieve 100% parity with the spec.\n\n")

	sb.WriteString("## Target Language: ")
	sb.WriteString(lang.Name)
	sb.WriteString("\n\n")

	sb.WriteString("### Language Conventions\n")
	sb.WriteString("- **Types**: ")
	sb.WriteString(lang.Conventions.Naming.Types)
	sb.WriteString("\n")
	sb.WriteString("- **Functions**: ")
	sb.WriteString(lang.Conventions.Naming.Functions)
	sb.WriteString("\n")
	sb.WriteString("- **Error Handling**: ")
	sb.WriteString(lang.Conventions.ErrorHandling)
	sb.WriteString("\n")
	sb.WriteString("- **File Extension**: ")
	sb.WriteString(lang.FileExtension)
	sb.WriteString("\n\n")

	sb.WriteString("## Output Directory\n")
	sb.WriteString("`")
	sb.WriteString(outputDir)
	sb.WriteString("`\n\n")

	sb.WriteString("## Recommended Project Structure\n")
	for _, f := range structure {
		sb.WriteString("- `")
		sb.WriteString(f.Path)
		sb.WriteString("` - ")
		sb.WriteString(f.Description)
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	sb.WriteString("## Generation Loop Process\n\n")
	sb.WriteString("You MUST follow this iterative process until 100% parity is achieved:\n\n")
	sb.WriteString("### Step 1: Generate Initial Code\n")
	sb.WriteString("Read the entire specification and generate ALL code:\n")
	sb.WriteString("- Project configuration (")
	switch lang.ID {
	case "go":
		sb.WriteString("go.mod")
	case "typescript":
		sb.WriteString("package.json, tsconfig.json")
	case "python":
		sb.WriteString("pyproject.toml")
	case "java":
		sb.WriteString("pom.xml or build.gradle")
	case "rust":
		sb.WriteString("Cargo.toml, src/lib.rs")
	case "csharp":
		sb.WriteString(".csproj file")
	}
	sb.WriteString(")\n")
	sb.WriteString("- Type definitions for EVERY type in the spec\n")
	sb.WriteString("- Function implementations for EVERY function in the spec\n")
	sb.WriteString("- Test files for EVERY test case in the spec\n\n")

	sb.WriteString("### Step 2: Verify Parity\n")
	sb.WriteString("After generating code, call `semantic_parity_analysis` with:\n")
	sb.WriteString("- sourcePath: The spec file path\n")
	sb.WriteString("- generatedProjects: [{language: \"")
	sb.WriteString(lang.ID)
	sb.WriteString("\", path: \"")
	sb.WriteString(outputDir)
	sb.WriteString("\"}]\n\n")

	sb.WriteString("### Step 3: Fix Gaps\n")
	sb.WriteString("If parity < 100%, examine the gaps returned and:\n")
	sb.WriteString("- Add missing types\n")
	sb.WriteString("- Add missing functions\n")
	sb.WriteString("- Add missing tests\n")
	sb.WriteString("- Fix any type mismatches\n\n")

	sb.WriteString("### Step 4: Loop\n")
	sb.WriteString("Repeat Steps 2-3 until parity reaches 100% or no more gaps can be fixed.\n\n")

	sb.WriteString("## Specification Content\n\n")
	sb.WriteString("Read this specification carefully. Extract ALL:\n")
	sb.WriteString("- Types (structs, classes, interfaces, enums)\n")
	sb.WriteString("- Functions (with parameters, return types, and logic)\n")
	sb.WriteString("- Tests (with given/when/then structure)\n")
	sb.WriteString("- Dependencies\n")
	sb.WriteString("- Configuration\n\n")
	sb.WriteString("```markdown\n")
	sb.WriteString(specContent)
	sb.WriteString("\n```\n\n")

	sb.WriteString("## BEGIN GENERATION\n\n")
	sb.WriteString("Now generate all the code files using the Write tool. Create each file in the output directory.\n")

	return sb.String()
}


// expandPath expands ~ to home directory and converts to absolute path
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(homeDir, path[2:])
		}
	}
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err == nil {
			path = absPath
		}
	}
	return path
}

// =============================================================================
// AI ORCHESTRATION HANDLERS
// =============================================================================

// handleListProjectLanguages scans a directory and returns all detected languages
func (s *Server) handleListProjectLanguages(ctx context.Context, req *mcp.CallToolRequest, input ListProjectLanguagesInput) (*mcp.CallToolResult, ListProjectLanguagesOutput, error) {
	sourcePath := expandPath(input.SourcePath)

	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Source directory does not exist: %s", sourcePath)},
			},
		}, ListProjectLanguagesOutput{}, nil
	}

	// Extended language detection with more file types
	langExtensions := map[string][]string{
		"go":         {".go"},
		"typescript": {".ts", ".tsx"},
		"javascript": {".js", ".jsx", ".mjs", ".cjs"},
		"python":     {".py", ".pyi"},
		"java":       {".java"},
		"rust":       {".rs"},
		"csharp":     {".cs"},
		"kotlin":     {".kt", ".kts"},
		"swift":      {".swift"},
		"ruby":       {".rb"},
		"php":        {".php"},
		"scala":      {".scala"},
		"c":          {".c", ".h"},
		"cpp":        {".cpp", ".cc", ".cxx", ".hpp", ".hxx"},
		"sql":        {".sql"},
		"protobuf":   {".proto"},
		"graphql":    {".graphql", ".gql"},
		"yaml":       {".yaml", ".yml"},
		"json":       {".json"},
		"markdown":   {".md"},
		"html":       {".html", ".htm"},
		"css":        {".css", ".scss", ".sass", ".less"},
		"shell":      {".sh", ".bash", ".zsh"},
	}

	// Languages with semantic parsers
	parsersAvailable := map[string]bool{
		"go":         true,
		"typescript": true,
		"python":     true,
		"java":       true,
		"rust":       true,
		"csharp":     true,
	}

	// Language display names
	langNames := map[string]string{
		"go": "Go", "typescript": "TypeScript", "javascript": "JavaScript",
		"python": "Python", "java": "Java", "rust": "Rust", "csharp": "C#",
		"kotlin": "Kotlin", "swift": "Swift", "ruby": "Ruby", "php": "PHP",
		"scala": "Scala", "c": "C", "cpp": "C++", "sql": "SQL",
		"protobuf": "Protocol Buffers", "graphql": "GraphQL", "yaml": "YAML",
		"json": "JSON", "markdown": "Markdown", "html": "HTML", "css": "CSS",
		"shell": "Shell",
	}

	// Skip directories
	skipDirs := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, "target": true,
		"build": true, "dist": true, "out": true, ".idea": true, ".vscode": true,
		"__pycache__": true, ".pytest_cache": true, ".mypy_cache": true,
		".next": true, ".nuxt": true, "coverage": true, ".gradle": true,
		".mvn": true, "bin": true, "obj": true, "testdata": true,
	}

	// Collect files by language
	langFiles := make(map[string][]string)
	langExts := make(map[string]map[string]bool)
	totalFiles := 0

	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || skipDirs[name] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		relPath, _ := filepath.Rel(sourcePath, path)

		for langID, exts := range langExtensions {
			for _, e := range exts {
				if ext == e {
					langFiles[langID] = append(langFiles[langID], relPath)
					if langExts[langID] == nil {
						langExts[langID] = make(map[string]bool)
					}
					langExts[langID][ext] = true
					totalFiles++
					break
				}
			}
		}
		return nil
	})

	// Build output
	var languages []LanguageInfo
	for langID, files := range langFiles {
		if len(files) == 0 {
			continue
		}

		// Get unique extensions
		var exts []string
		for ext := range langExts[langID] {
			exts = append(exts, ext)
		}

		// Sample files (max 5)
		sampleFiles := files
		if len(sampleFiles) > 5 {
			sampleFiles = sampleFiles[:5]
		}

		languages = append(languages, LanguageInfo{
			ID:          langID,
			Name:        langNames[langID],
			FileCount:   len(files),
			Extensions:  exts,
			HasParser:   parsersAvailable[langID],
			SampleFiles: sampleFiles,
		})
	}

	// Sort by file count (descending)
	for i := 0; i < len(languages); i++ {
		for j := i + 1; j < len(languages); j++ {
			if languages[j].FileCount > languages[i].FileCount {
				languages[i], languages[j] = languages[j], languages[i]
			}
		}
	}

	// Build recommendation
	var recommendation strings.Builder
	recommendation.WriteString("Recommended analysis approach:\n")

	var withParser, withoutParser []string
	for _, lang := range languages {
		if lang.HasParser {
			withParser = append(withParser, lang.ID)
		} else {
			withoutParser = append(withoutParser, lang.ID)
		}
	}

	if len(withParser) > 0 {
		recommendation.WriteString(fmt.Sprintf("1. Use deep_analyze_source for: %s\n", strings.Join(withParser, ", ")))
	}
	if len(withoutParser) > 0 {
		recommendation.WriteString(fmt.Sprintf("2. Use get_files_for_language + AI interpretation for: %s\n", strings.Join(withoutParser, ", ")))
	}
	recommendation.WriteString("3. Combine all analyses into a unified understanding of the project")

	return nil, ListProjectLanguagesOutput{
		ProjectName:    filepath.Base(sourcePath),
		Languages:      languages,
		TotalFiles:     totalFiles,
		Recommendation: recommendation.String(),
	}, nil
}

// handleGetFilesForLanguage returns raw file contents for AI interpretation
func (s *Server) handleGetFilesForLanguage(ctx context.Context, req *mcp.CallToolRequest, input GetFilesForLanguageInput) (*mcp.CallToolResult, GetFilesForLanguageOutput, error) {
	sourcePath := expandPath(input.SourcePath)

	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Source directory does not exist: %s", sourcePath)},
			},
		}, GetFilesForLanguageOutput{}, nil
	}

	// Set defaults
	maxFiles := input.MaxFiles
	if maxFiles <= 0 {
		maxFiles = 50
	}
	if maxFiles > 200 {
		maxFiles = 200
	}

	// Language extensions mapping
	langExtensions := map[string][]string{
		"go":         {".go"},
		"typescript": {".ts", ".tsx"},
		"javascript": {".js", ".jsx", ".mjs", ".cjs"},
		"python":     {".py", ".pyi"},
		"java":       {".java"},
		"rust":       {".rs"},
		"csharp":     {".cs"},
		"kotlin":     {".kt", ".kts"},
		"swift":      {".swift"},
		"ruby":       {".rb"},
		"php":        {".php"},
		"scala":      {".scala"},
		"c":          {".c", ".h"},
		"cpp":        {".cpp", ".cc", ".cxx", ".hpp", ".hxx"},
		"sql":        {".sql"},
		"protobuf":   {".proto"},
		"graphql":    {".graphql", ".gql"},
		"yaml":       {".yaml", ".yml"},
		"json":       {".json"},
		"markdown":   {".md"},
		"html":       {".html", ".htm"},
		"css":        {".css", ".scss", ".sass", ".less"},
		"shell":      {".sh", ".bash", ".zsh"},
	}

	// Test file patterns
	testPatterns := []string{
		"_test.go", ".test.ts", ".test.tsx", ".test.js", ".spec.ts", ".spec.js",
		"test_", "_test.py", "Test.java", "Tests.java", "_test.rs", "Tests.cs", "Test.cs",
	}

	isTestFile := func(path string) bool {
		lowerPath := strings.ToLower(path)
		for _, pat := range testPatterns {
			if strings.Contains(lowerPath, strings.ToLower(pat)) {
				return true
			}
		}
		return false
	}

	exts, ok := langExtensions[input.Language]
	if !ok {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Unknown language: %s. Supported: go, typescript, javascript, python, java, rust, csharp, kotlin, swift, ruby, php, scala, c, cpp, sql, protobuf, graphql, yaml, json, markdown, html, css, shell", input.Language)},
			},
		}, GetFilesForLanguageOutput{}, nil
	}

	// Skip directories
	skipDirs := map[string]bool{
		".git": true, "node_modules": true, "vendor": true, "target": true,
		"build": true, "dist": true, "out": true, ".idea": true, ".vscode": true,
		"__pycache__": true, ".pytest_cache": true, ".mypy_cache": true,
		".next": true, ".nuxt": true, "coverage": true, ".gradle": true,
		".mvn": true, "bin": true, "obj": true, "testdata": true,
	}

	var files []FileContent
	var totalSize int
	truncated := false
	totalFound := 0

	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || skipDirs[name] {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip large files
		if info.Size() > 1024*1024 {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		relPath, _ := filepath.Rel(sourcePath, path)

		// Check if this file matches the requested language
		matches := false
		for _, e := range exts {
			if ext == e {
				matches = true
				break
			}
		}
		if !matches {
			return nil
		}

		// Check test file filter
		if !input.IncludeTests && isTestFile(relPath) {
			return nil
		}

		totalFound++

		// Check if we've hit the limit
		if len(files) >= maxFiles {
			truncated = true
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		files = append(files, FileContent{
			Path:    relPath,
			Content: string(content),
			Size:    len(content),
		})
		totalSize += len(content)

		return nil
	})

	// Build AI prompt for interpretation
	var prompt strings.Builder
	prompt.WriteString(fmt.Sprintf("# %s Files Analysis\n\n", strings.Title(input.Language)))
	prompt.WriteString(fmt.Sprintf("Analyze the following %d %s files and extract:\n", len(files), input.Language))
	prompt.WriteString("1. **Types/Structures**: All data types, classes, structs, enums, interfaces\n")
	prompt.WriteString("2. **Functions/Methods**: All functions with their signatures and purpose\n")
	prompt.WriteString("3. **Dependencies**: External libraries and internal imports\n")
	prompt.WriteString("4. **Patterns**: Design patterns, architectural patterns observed\n")
	prompt.WriteString("5. **API Surface**: Public interfaces, endpoints, exported items\n\n")

	if truncated {
		prompt.WriteString(fmt.Sprintf("Note: Only showing %d of %d total files. Call again with different maxFiles if needed.\n\n", len(files), totalFound))
	}

	prompt.WriteString("---\n\n")
	for _, f := range files {
		prompt.WriteString(fmt.Sprintf("## File: %s\n```%s\n%s\n```\n\n", f.Path, input.Language, f.Content))
	}

	return nil, GetFilesForLanguageOutput{
		Language:  input.Language,
		FileCount: len(files),
		TotalSize: totalSize,
		Files:     files,
		Truncated: truncated,
		AIPrompt:  prompt.String(),
	}, nil
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
