package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kon1790/rpg/internal/importer/extractors"
)

// BuildAnalysisPrompt creates a comprehensive prompt for AI analysis of a project
func BuildAnalysisPrompt(files *ProjectFiles) string {
	var sb strings.Builder

	sb.WriteString("# Source Code Analysis Task\n\n")
	sb.WriteString("Analyze the following project and generate a comprehensive .spec.md file.\n\n")
	sb.WriteString(fmt.Sprintf("## Project: %s\n", files.Name))
	sb.WriteString(fmt.Sprintf("## Detected Language: %s\n", files.Language))
	sb.WriteString(fmt.Sprintf("## Total Files: %d\n\n", files.GetTotalFileCount()))

	// Source Files
	if len(files.SourceFiles) > 0 {
		sb.WriteString("---\n\n## Source Files\n\n")
		for _, f := range files.SourceFiles {
			sb.WriteString(fmt.Sprintf("### FILE: %s\n", f.Path))
			sb.WriteString("```\n")
			sb.WriteString(f.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	// Test Files
	if len(files.TestFiles) > 0 {
		sb.WriteString("---\n\n## Test Files\n\n")
		for _, f := range files.TestFiles {
			sb.WriteString(fmt.Sprintf("### FILE: %s\n", f.Path))
			sb.WriteString("```\n")
			sb.WriteString(f.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	// API Specifications
	if len(files.APISpecs) > 0 {
		sb.WriteString("---\n\n## API Specifications\n\n")
		for _, f := range files.APISpecs {
			sb.WriteString(fmt.Sprintf("### FILE: %s\n", f.Path))
			sb.WriteString("```\n")
			sb.WriteString(f.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	// Configuration Files
	if len(files.ConfigFiles) > 0 {
		sb.WriteString("---\n\n## Configuration Files\n\n")
		for _, f := range files.ConfigFiles {
			sb.WriteString(fmt.Sprintf("### FILE: %s\n", f.Path))
			sb.WriteString("```\n")
			sb.WriteString(f.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	// Documentation Files
	if len(files.DocFiles) > 0 {
		sb.WriteString("---\n\n## Documentation Files\n\n")
		for _, f := range files.DocFiles {
			sb.WriteString(fmt.Sprintf("### FILE: %s\n", f.Path))
			sb.WriteString("```\n")
			sb.WriteString(f.Content)
			sb.WriteString("\n```\n\n")
		}
	}

	// Instructions for AI
	sb.WriteString("---\n\n## Instructions\n\n")
	sb.WriteString("Generate a complete .spec.md file that:\n\n")
	sb.WriteString("1. **Captures ALL types** - structs, classes, interfaces, enums with:\n")
	sb.WriteString("   - Field names and types\n")
	sb.WriteString("   - Optional/required status\n")
	sb.WriteString("   - Inheritance relationships\n")
	sb.WriteString("   - Generic type parameters\n\n")
	sb.WriteString("2. **Documents ALL functions** with:\n")
	sb.WriteString("   - Clear description of what it does\n")
	sb.WriteString("   - All parameters with types and defaults\n")
	sb.WriteString("   - Return type\n")
	sb.WriteString("   - Detailed logic as pseudo-code (step-by-step what the function does)\n")
	sb.WriteString("   - Error conditions\n")
	sb.WriteString("   - Whether it's async\n\n")
	sb.WriteString("3. **Extracts ALL tests** with:\n")
	sb.WriteString("   - Test name\n")
	sb.WriteString("   - Given inputs (actual values from test)\n")
	sb.WriteString("   - Expected outputs (actual values from test)\n")
	sb.WriteString("   - Any setup/mock requirements\n\n")
	sb.WriteString("4. **Includes API endpoints** (from OpenAPI/code):\n")
	sb.WriteString("   - HTTP method and path\n")
	sb.WriteString("   - Request parameters and body\n")
	sb.WriteString("   - Response types\n")
	sb.WriteString("   - Error responses\n\n")
	sb.WriteString("5. **Documents configuration**:\n")
	sb.WriteString("   - Environment variables with defaults\n")
	sb.WriteString("   - Config file settings\n")
	sb.WriteString("   - Required vs optional settings\n\n")
	sb.WriteString("6. **Lists dependencies**:\n")
	sb.WriteString("   - External packages\n")
	sb.WriteString("   - Framework detection\n\n")
	sb.WriteString("## Output Format\n\n")
	sb.WriteString("Return ONLY the .spec.md content in the standard format:\n")
	sb.WriteString("- # Project Name\n")
	sb.WriteString("- Description\n")
	sb.WriteString("- ## Target Languages\n")
	sb.WriteString("- ## Configuration (if applicable)\n")
	sb.WriteString("- ## Types\n")
	sb.WriteString("- ## Functions (with full logic blocks)\n")
	sb.WriteString("- ## Tests (with given/expect values)\n")
	sb.WriteString("- ## Dependencies\n")

	return sb.String()
}

// Re-export types from extractors package for backward compatibility
type ExtractedProject = extractors.ExtractedProject
type ExtractedType = extractors.ExtractedType
type ExtractedField = extractors.ExtractedField
type ExtractedFunction = extractors.ExtractedFunction
type ExtractedParam = extractors.ExtractedParam
type ExtractedTest = extractors.ExtractedTest

// Importer coordinates source code analysis and spec generation.
type Importer struct {
	registry *extractors.Registry
}

// New creates a new Importer.
func New() *Importer {
	return &Importer{
		registry: extractors.NewRegistry(),
	}
}

// DetectLanguage detects the primary language in a directory by counting file extensions.
func (i *Importer) DetectLanguage(dir string) (string, error) {
	counts := make(map[string]int)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		ext := filepath.Ext(path)
		if extractor, ok := i.registry.GetByExtension(ext); ok {
			counts[extractor.LanguageID()]++
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(counts) == 0 {
		return "", fmt.Errorf("no supported source files found in %s", dir)
	}

	// Find the language with the most files
	maxCount := 0
	detectedLang := ""
	for lang, count := range counts {
		if count > maxCount {
			maxCount = count
			detectedLang = lang
		}
	}

	return detectedLang, nil
}

// Extract analyzes source code in a directory and returns extracted information.
func (i *Importer) Extract(dir string, language string) (*ExtractedProject, error) {
	extractor, ok := i.registry.Get(language)
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	project := &ExtractedProject{
		Name:             filepath.Base(dir),
		DetectedLanguage: language,
	}

	// Walk the directory and extract from each file
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		ext := filepath.Ext(path)
		validExt := false
		for _, e := range extractor.Extensions() {
			if ext == e {
				validExt = true
				break
			}
		}
		if !validExt {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			project.Warnings = append(project.Warnings, fmt.Sprintf("Failed to read %s: %v", path, err))
			return nil
		}

		relPath, _ := filepath.Rel(dir, path)
		contentStr := string(content)

		// Extract package description from first file
		if project.Description == "" {
			project.Description = extractor.ExtractPackageDescription(contentStr)
		}

		if extractor.IsTestFile(relPath) {
			// Extract tests
			tests := extractor.ExtractTests(contentStr, relPath)
			project.Tests = append(project.Tests, tests...)
		} else {
			// Extract types and functions
			types := extractor.ExtractTypes(contentStr, relPath)
			project.Types = append(project.Types, types...)

			functions := extractor.ExtractFunctions(contentStr, relPath)
			project.Functions = append(project.Functions, functions...)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to extract from directory: %w", err)
	}

	// Deduplicate types and functions by name
	project.Types = deduplicateTypes(project.Types)
	project.Functions = deduplicateFunctions(project.Functions)

	return project, nil
}

// deduplicateTypes removes duplicate types by name.
func deduplicateTypes(types []ExtractedType) []ExtractedType {
	seen := make(map[string]bool)
	var result []ExtractedType

	for _, t := range types {
		if !seen[t.Name] {
			seen[t.Name] = true
			result = append(result, t)
		}
	}

	return result
}

// deduplicateFunctions removes duplicate functions by name.
func deduplicateFunctions(functions []ExtractedFunction) []ExtractedFunction {
	seen := make(map[string]bool)
	var result []ExtractedFunction

	for _, f := range functions {
		if !seen[f.Name] {
			seen[f.Name] = true
			result = append(result, f)
		}
	}

	return result
}

// GenerateSpec generates a markdown spec from extracted project information.
func (i *Importer) GenerateSpec(project *ExtractedProject) string {
	return GenerateMarkdown(project)
}

// ImportFromDirectory is a convenience method that detects, extracts, and generates a spec.
func (i *Importer) ImportFromDirectory(inputDir string, outputDir string, specName string) (string, *ExtractedProject, error) {
	// Detect language
	language, err := i.DetectLanguage(inputDir)
	if err != nil {
		return "", nil, err
	}

	// Extract project information
	project, err := i.Extract(inputDir, language)
	if err != nil {
		return "", nil, err
	}

	// Override name if provided
	if specName != "" {
		project.Name = specName
	}

	// Generate spec markdown
	specContent := i.GenerateSpec(project)

	// Determine output path
	specPath := filepath.Join(outputDir, project.Name+".spec.md")
	if _, err := os.Stat(specPath); err == nil {
		// File exists, add -imported suffix
		specPath = filepath.Join(outputDir, project.Name+"-imported.spec.md")
	}

	// Create output directory if needed
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write spec file
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		return "", nil, fmt.Errorf("failed to write spec file: %w", err)
	}

	return specPath, project, nil
}

// GetSupportedLanguages returns a list of supported language IDs.
func (i *Importer) GetSupportedLanguages() []string {
	var languages []string
	for _, e := range i.registry.List() {
		languages = append(languages, e.LanguageID())
	}
	return languages
}

// slugify converts a name to a URL-friendly slug.
func slugify(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return name
}
