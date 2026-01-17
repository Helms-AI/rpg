package semantic

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kon1790/rpg/internal/importer/treesitter"
)

// TypeScriptAnalyzer provides semantic analysis for TypeScript code
type TypeScriptAnalyzer struct {
	*SubprocessAnalyzer
}

// NewTypeScriptAnalyzer creates a new TypeScript semantic analyzer
func NewTypeScriptAnalyzer() *TypeScriptAnalyzer {
	return &TypeScriptAnalyzer{
		SubprocessAnalyzer: NewSubprocessAnalyzer(SubprocessConfig{
			Language: treesitter.LanguageTypeScript,
			Command:  "npx",
			Args:     []string{"tsc"},
		}),
	}
}

// IsAvailable checks if TypeScript compiler is available
func (a *TypeScriptAnalyzer) IsAvailable() bool {
	// Check for npx (comes with npm)
	if a.SubprocessAnalyzer.IsAvailable() {
		return true
	}
	// Also check for direct tsc
	a.SubprocessAnalyzer.command = "tsc"
	a.SubprocessAnalyzer.args = nil
	return a.SubprocessAnalyzer.IsAvailable()
}

// Analyze performs semantic analysis on a TypeScript project directory
func (a *TypeScriptAnalyzer) Analyze(dir string) (*Analysis, error) {
	analysis := &Analysis{
		Language:  treesitter.LanguageTypeScript,
		CallGraph: make(map[string][]string),
		TypeGraph: make(map[string][]string),
	}

	// Find project name from package.json or directory name
	analysis.Name = a.findProjectName(dir)

	// Find all TypeScript files
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && (info.Name() == "node_modules" || info.Name() == ".git") {
			return filepath.SkipDir
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".tsx")) {
			if !strings.HasSuffix(path, ".d.ts") && !strings.Contains(path, ".test.") && !strings.Contains(path, ".spec.") {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	// Analyze each file
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			analysis.Errors = append(analysis.Errors, AnalysisError{
				File:     file,
				Message:  err.Error(),
				Severity: SeverityWarning,
			})
			continue
		}

		fileAnalysis, err := a.AnalyzeFile(file, content)
		if err != nil {
			analysis.Errors = append(analysis.Errors, AnalysisError{
				File:     file,
				Message:  err.Error(),
				Severity: SeverityWarning,
			})
			continue
		}

		analysis.Files = append(analysis.Files, fileAnalysis)
		analysis.Types = append(analysis.Types, fileAnalysis.Types...)
		analysis.Functions = append(analysis.Functions, fileAnalysis.Functions...)
	}

	// Try semantic enrichment via tsc
	if a.IsAvailable() {
		a.enrichWithTSC(dir, analysis)
	}

	// Build graphs
	a.buildCallGraph(analysis)
	a.buildTypeGraph(analysis)
	a.extractDependencies(analysis)

	return analysis, nil
}

// AnalyzeFile performs semantic analysis on a single TypeScript file
func (a *TypeScriptAnalyzer) AnalyzeFile(path string, content []byte) (*FileAnalysis, error) {
	// Use tree-sitter for structural parsing
	return a.TreeSitterAnalysis(content, path)
}

// enrichWithTSC enriches analysis with TypeScript compiler information
func (a *TypeScriptAnalyzer) enrichWithTSC(dir string, analysis *Analysis) {
	ctx := context.Background()

	// Try to get diagnostics from tsc
	output, err := a.RunCommand(ctx, "--noEmit", "--pretty", "false", "--project", dir)
	if err != nil {
		// Not a fatal error - we have tree-sitter analysis
		analysis.Errors = append(analysis.Errors, AnalysisError{
			Message:  fmt.Sprintf("tsc diagnostics unavailable: %v", err),
			Severity: SeverityInfo,
		})
		return
	}

	// Parse tsc output for type information
	a.parseTSCOutput(string(output), analysis)
}

// parseTSCOutput parses TypeScript compiler output
func (a *TypeScriptAnalyzer) parseTSCOutput(output string, analysis *Analysis) {
	// TSC outputs errors in format: file(line,col): error TSxxxx: message
	// We could parse these to find type issues
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, ": error TS") {
			analysis.Errors = append(analysis.Errors, AnalysisError{
				Message:  line,
				Severity: SeverityWarning,
			})
		}
	}
}

// findProjectName finds the project name from package.json or directory
func (a *TypeScriptAnalyzer) findProjectName(dir string) string {
	// Try to read package.json
	pkgPath := filepath.Join(dir, "package.json")
	if data, err := os.ReadFile(pkgPath); err == nil {
		var pkg struct {
			Name string `json:"name"`
		}
		if json.Unmarshal(data, &pkg) == nil && pkg.Name != "" {
			return pkg.Name
		}
	}

	// Fall back to directory name
	return filepath.Base(dir)
}

// buildCallGraph builds the call graph from analysis
func (a *TypeScriptAnalyzer) buildCallGraph(analysis *Analysis) {
	for _, fn := range analysis.Functions {
		analysis.CallGraph[fn.Name] = fn.Calls
	}
}

// buildTypeGraph builds the type graph from analysis
func (a *TypeScriptAnalyzer) buildTypeGraph(analysis *Analysis) {
	for _, typ := range analysis.Types {
		if len(typ.ImplementsInterfaces) > 0 {
			analysis.TypeGraph[typ.Name] = typ.ImplementsInterfaces
		}
	}
}

// extractDependencies extracts external dependencies
func (a *TypeScriptAnalyzer) extractDependencies(analysis *Analysis) {
	seen := make(map[string]bool)

	for _, file := range analysis.Files {
		for _, imp := range file.Imports {
			if seen[imp.Path] {
				continue
			}
			seen[imp.Path] = true

			dep := Dependency{
				Path:     imp.Path,
				IsStdLib: false, // TypeScript has no stdlib in the same sense
				IsLocal:  strings.HasPrefix(imp.Path, ".") || strings.HasPrefix(imp.Path, "/"),
			}
			analysis.Dependencies = append(analysis.Dependencies, dep)
		}
	}
}
