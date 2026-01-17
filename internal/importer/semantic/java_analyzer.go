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

// JavaAnalyzer provides semantic analysis for Java code
type JavaAnalyzer struct {
	*SubprocessAnalyzer
}

// NewJavaAnalyzer creates a new Java semantic analyzer
func NewJavaAnalyzer() *JavaAnalyzer {
	return &JavaAnalyzer{
		SubprocessAnalyzer: NewSubprocessAnalyzer(SubprocessConfig{
			Language: treesitter.LanguageJava,
			Command:  "javac",
			Args:     []string{"-version"},
		}),
	}
}

// IsAvailable checks if Java compiler is available
func (a *JavaAnalyzer) IsAvailable() bool {
	return a.SubprocessAnalyzer.IsAvailable()
}

// Analyze performs semantic analysis on a Java project directory
func (a *JavaAnalyzer) Analyze(dir string) (*Analysis, error) {
	analysis := &Analysis{
		Language:  treesitter.LanguageJava,
		CallGraph: make(map[string][]string),
		TypeGraph: make(map[string][]string),
	}

	// Find project name
	analysis.Name = a.findProjectName(dir)

	// Find all Java files
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "target" || name == "build" || name == ".gradle" {
				return filepath.SkipDir
			}
		}
		if !info.IsDir() && strings.HasSuffix(path, ".java") {
			if !strings.Contains(path, "Test.java") && !strings.Contains(path, "test/") {
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

	// Try semantic enrichment via javac
	if a.IsAvailable() && len(files) > 0 {
		a.enrichWithJavac(dir, files, analysis)
	}

	// Build graphs
	a.buildCallGraph(analysis)
	a.buildTypeGraph(analysis)
	a.extractDependencies(analysis)

	return analysis, nil
}

// AnalyzeFile performs semantic analysis on a single Java file
func (a *JavaAnalyzer) AnalyzeFile(path string, content []byte) (*FileAnalysis, error) {
	// Use tree-sitter for structural parsing
	return a.TreeSitterAnalysis(content, path)
}

// enrichWithJavac enriches analysis with javac diagnostic information
func (a *JavaAnalyzer) enrichWithJavac(dir string, files []string, analysis *Analysis) {
	ctx := context.Background()

	// Build classpath from common locations
	classpath := a.buildClasspath(dir)

	// Run javac with diagnostics
	args := []string{
		"-d", "/tmp/rpg-java-out",
		"-Xlint:all",
		"-proc:none", // Skip annotation processing
	}

	if classpath != "" {
		args = append(args, "-cp", classpath)
	}

	// Create output directory
	os.MkdirAll("/tmp/rpg-java-out", 0755)
	defer os.RemoveAll("/tmp/rpg-java-out")

	// Add source files
	args = append(args, files...)

	a.SubprocessAnalyzer.args = nil
	output, err := a.RunCommand(ctx, args...)
	if err != nil {
		// Javac returns non-zero on errors, which is expected
		// Parse output for diagnostic information
	}

	// Parse javac output for warnings/errors
	if output != nil {
		a.parseJavacOutput(string(output), analysis)
	}
}

// buildClasspath builds a classpath from common locations
func (a *JavaAnalyzer) buildClasspath(dir string) string {
	var paths []string

	// Check for Maven dependencies
	m2Repo := filepath.Join(os.Getenv("HOME"), ".m2", "repository")
	if _, err := os.Stat(m2Repo); err == nil {
		// Don't add entire m2 repo - too large
	}

	// Check for target/classes (Maven)
	targetClasses := filepath.Join(dir, "target", "classes")
	if _, err := os.Stat(targetClasses); err == nil {
		paths = append(paths, targetClasses)
	}

	// Check for build/classes (Gradle)
	buildClasses := filepath.Join(dir, "build", "classes", "java", "main")
	if _, err := os.Stat(buildClasses); err == nil {
		paths = append(paths, buildClasses)
	}

	// Check for lib directory
	libDir := filepath.Join(dir, "lib")
	if info, err := os.Stat(libDir); err == nil && info.IsDir() {
		// Add all JARs in lib
		filepath.Walk(libDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && strings.HasSuffix(path, ".jar") {
				paths = append(paths, path)
			}
			return nil
		})
	}

	return strings.Join(paths, string(os.PathListSeparator))
}

// parseJavacOutput parses javac compiler output
func (a *JavaAnalyzer) parseJavacOutput(output string, analysis *Analysis) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Javac format: File.java:line: warning/error: message
		if strings.Contains(line, ": warning:") {
			analysis.Errors = append(analysis.Errors, AnalysisError{
				Message:  line,
				Severity: SeverityWarning,
			})
		} else if strings.Contains(line, ": error:") {
			analysis.Errors = append(analysis.Errors, AnalysisError{
				Message:  line,
				Severity: SeverityError,
			})
		}
	}
}

// findProjectName finds the project name from build files
func (a *JavaAnalyzer) findProjectName(dir string) string {
	// Try pom.xml (Maven)
	pomPath := filepath.Join(dir, "pom.xml")
	if data, err := os.ReadFile(pomPath); err == nil {
		var pom struct {
			XMLName    xml.Name `xml:"project"`
			ArtifactId string   `xml:"artifactId"`
			Name       string   `xml:"name"`
		}
		if xml.Unmarshal(data, &pom) == nil {
			if pom.Name != "" {
				return pom.Name
			}
			if pom.ArtifactId != "" {
				return pom.ArtifactId
			}
		}
	}

	// Try build.gradle (Gradle)
	gradlePath := filepath.Join(dir, "build.gradle")
	if data, err := os.ReadFile(gradlePath); err == nil {
		// Simple extraction - look for rootProject.name or archivesBaseName
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "archivesBaseName") || strings.Contains(line, "rootProject.name") {
				parts := strings.Split(line, "=")
				if len(parts) == 2 {
					name := strings.TrimSpace(parts[1])
					name = strings.Trim(name, `"'`)
					if name != "" {
						return name
					}
				}
			}
		}
	}

	// Fall back to directory name
	return filepath.Base(dir)
}

// buildCallGraph builds the call graph from analysis
func (a *JavaAnalyzer) buildCallGraph(analysis *Analysis) {
	for _, fn := range analysis.Functions {
		analysis.CallGraph[fn.Name] = fn.Calls
	}
}

// buildTypeGraph builds the type graph from analysis
func (a *JavaAnalyzer) buildTypeGraph(analysis *Analysis) {
	for _, typ := range analysis.Types {
		if len(typ.ImplementsInterfaces) > 0 {
			analysis.TypeGraph[typ.Name] = typ.ImplementsInterfaces
		}
	}
}

// extractDependencies extracts external dependencies
func (a *JavaAnalyzer) extractDependencies(analysis *Analysis) {
	seen := make(map[string]bool)

	for _, file := range analysis.Files {
		for _, imp := range file.Imports {
			if seen[imp.Path] {
				continue
			}
			seen[imp.Path] = true

			// Check if it's a java.* or javax.* (standard library)
			isStdLib := strings.HasPrefix(imp.Path, "java.") || strings.HasPrefix(imp.Path, "javax.")

			dep := Dependency{
				Path:     imp.Path,
				IsStdLib: isStdLib,
				IsLocal:  false, // Java doesn't have relative imports in the same sense
			}
			analysis.Dependencies = append(analysis.Dependencies, dep)
		}
	}
}
