package semantic

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kon1790/rpg/internal/importer/treesitter"
)

// CSharpAnalyzer provides semantic analysis for C# code
type CSharpAnalyzer struct {
	*SubprocessAnalyzer
}

// NewCSharpAnalyzer creates a new C# semantic analyzer
func NewCSharpAnalyzer() *CSharpAnalyzer {
	return &CSharpAnalyzer{
		SubprocessAnalyzer: NewSubprocessAnalyzer(SubprocessConfig{
			Language: treesitter.LanguageCSharp,
			Command:  "dotnet",
			Args:     []string{"--version"},
		}),
	}
}

// IsAvailable checks if .NET CLI is available
func (a *CSharpAnalyzer) IsAvailable() bool {
	return a.SubprocessAnalyzer.IsAvailable()
}

// Analyze performs semantic analysis on a C# project directory
func (a *CSharpAnalyzer) Analyze(dir string) (*Analysis, error) {
	analysis := &Analysis{
		Language:  treesitter.LanguageCSharp,
		CallGraph: make(map[string][]string),
		TypeGraph: make(map[string][]string),
	}

	// Find project name
	analysis.Name = a.findProjectName(dir)

	// Find all C# source files
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "bin" || name == "obj" || name == "packages" {
				return filepath.SkipDir
			}
		}
		if !info.IsDir() && strings.HasSuffix(path, ".cs") {
			// Skip test files and generated files
			if !strings.Contains(path, ".Tests") &&
				!strings.HasSuffix(path, ".g.cs") &&
				!strings.HasSuffix(path, ".Designer.cs") {
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

	// Try semantic enrichment via dotnet build
	if a.IsAvailable() {
		a.enrichWithDotnet(dir, analysis)
	}

	// Build graphs
	a.buildCallGraph(analysis)
	a.buildTypeGraph(analysis)
	a.extractDependencies(dir, analysis)

	return analysis, nil
}

// AnalyzeFile performs semantic analysis on a single C# file
func (a *CSharpAnalyzer) AnalyzeFile(path string, content []byte) (*FileAnalysis, error) {
	// Use tree-sitter for structural parsing
	return a.TreeSitterAnalysis(content, path)
}

// enrichWithDotnet enriches analysis with dotnet build diagnostic information
func (a *CSharpAnalyzer) enrichWithDotnet(dir string, analysis *Analysis) {
	ctx := context.Background()

	// Find project file
	projectFile := a.findProjectFile(dir)
	if projectFile == "" {
		return
	}

	// Run dotnet build with diagnostics
	a.SubprocessAnalyzer.command = "dotnet"
	a.SubprocessAnalyzer.args = nil

	output, err := a.RunCommand(ctx, "build", projectFile, "--no-restore", "-v", "quiet", "-nologo")
	if err != nil {
		// Build can fail, we still want to parse output
	}

	// Parse build output for warnings/errors
	if output != nil {
		a.parseDotnetOutput(string(output), analysis)
	}
}

// findProjectFile finds a .csproj or .sln file
func (a *CSharpAnalyzer) findProjectFile(dir string) string {
	// Look for .sln first
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".sln") {
			return filepath.Join(dir, entry.Name())
		}
	}

	// Look for .csproj
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".csproj") {
			return filepath.Join(dir, entry.Name())
		}
	}

	return ""
}

// parseDotnetOutput parses dotnet build output
func (a *CSharpAnalyzer) parseDotnetOutput(output string, analysis *Analysis) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// MSBuild format: File(line,col): warning/error CODE: message
		if strings.Contains(line, ": warning ") {
			analysis.Errors = append(analysis.Errors, AnalysisError{
				Message:  line,
				Severity: SeverityWarning,
			})
		} else if strings.Contains(line, ": error ") {
			analysis.Errors = append(analysis.Errors, AnalysisError{
				Message:  line,
				Severity: SeverityError,
			})
		}
	}
}

// findProjectName finds the project name from .csproj
func (a *CSharpAnalyzer) findProjectName(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return filepath.Base(dir)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".csproj") {
			projPath := filepath.Join(dir, entry.Name())
			if data, err := os.ReadFile(projPath); err == nil {
				var proj struct {
					XMLName       xml.Name `xml:"Project"`
					PropertyGroup []struct {
						AssemblyName string `xml:"AssemblyName"`
						PackageId    string `xml:"PackageId"`
					} `xml:"PropertyGroup"`
				}
				if xml.Unmarshal(data, &proj) == nil {
					for _, pg := range proj.PropertyGroup {
						if pg.AssemblyName != "" {
							return pg.AssemblyName
						}
						if pg.PackageId != "" {
							return pg.PackageId
						}
					}
				}
			}
			// Use csproj filename without extension
			return strings.TrimSuffix(entry.Name(), ".csproj")
		}
	}

	// Fall back to directory name
	return filepath.Base(dir)
}

// buildCallGraph builds the call graph from analysis
func (a *CSharpAnalyzer) buildCallGraph(analysis *Analysis) {
	for _, fn := range analysis.Functions {
		analysis.CallGraph[fn.Name] = fn.Calls
	}
}

// buildTypeGraph builds the type graph from analysis
func (a *CSharpAnalyzer) buildTypeGraph(analysis *Analysis) {
	for _, typ := range analysis.Types {
		if len(typ.ImplementsInterfaces) > 0 {
			analysis.TypeGraph[typ.Name] = typ.ImplementsInterfaces
		}
	}
}

// extractDependencies extracts dependencies from .csproj
func (a *CSharpAnalyzer) extractDependencies(dir string, analysis *Analysis) {
	// Parse imports from analysis
	seen := make(map[string]bool)

	for _, file := range analysis.Files {
		for _, imp := range file.Imports {
			if seen[imp.Path] {
				continue
			}
			seen[imp.Path] = true

			// Check if it's a System namespace (standard library)
			isStdLib := strings.HasPrefix(imp.Path, "System") ||
				strings.HasPrefix(imp.Path, "Microsoft")

			dep := Dependency{
				Path:     imp.Path,
				IsStdLib: isStdLib,
				IsLocal:  false,
			}
			analysis.Dependencies = append(analysis.Dependencies, dep)
		}
	}

	// Also parse .csproj for PackageReference
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".csproj") {
			continue
		}

		projPath := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(projPath)
		if err != nil {
			continue
		}

		var proj struct {
			XMLName   xml.Name `xml:"Project"`
			ItemGroup []struct {
				PackageReference []struct {
					Include string `xml:"Include,attr"`
					Version string `xml:"Version,attr"`
				} `xml:"PackageReference"`
			} `xml:"ItemGroup"`
		}

		if err := xml.Unmarshal(data, &proj); err != nil {
			continue
		}

		for _, ig := range proj.ItemGroup {
			for _, pkg := range ig.PackageReference {
				if pkg.Include != "" && !seen[pkg.Include] {
					seen[pkg.Include] = true
					analysis.Dependencies = append(analysis.Dependencies, Dependency{
						Path:     pkg.Include,
						Version:  pkg.Version,
						IsStdLib: false,
						IsLocal:  false,
					})
				}
			}
		}
	}
}
