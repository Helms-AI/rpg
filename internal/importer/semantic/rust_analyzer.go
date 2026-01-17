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

// RustAnalyzer provides semantic analysis for Rust code
type RustAnalyzer struct {
	*SubprocessAnalyzer
}

// NewRustAnalyzer creates a new Rust semantic analyzer
func NewRustAnalyzer() *RustAnalyzer {
	return &RustAnalyzer{
		SubprocessAnalyzer: NewSubprocessAnalyzer(SubprocessConfig{
			Language: treesitter.LanguageRust,
			Command:  "cargo",
			Args:     []string{"--version"},
		}),
	}
}

// IsAvailable checks if Rust/Cargo is available
func (a *RustAnalyzer) IsAvailable() bool {
	return a.SubprocessAnalyzer.IsAvailable()
}

// Analyze performs semantic analysis on a Rust project directory
func (a *RustAnalyzer) Analyze(dir string) (*Analysis, error) {
	analysis := &Analysis{
		Language:  treesitter.LanguageRust,
		CallGraph: make(map[string][]string),
		TypeGraph: make(map[string][]string),
	}

	// Find project name
	analysis.Name = a.findProjectName(dir)

	// Find all Rust source files
	var files []string
	srcDir := filepath.Join(dir, "src")
	if _, err := os.Stat(srcDir); err == nil {
		err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".rs") {
				// Skip test files
				if !strings.HasSuffix(path, "_test.rs") && !strings.Contains(path, "/tests/") {
					files = append(files, path)
				}
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walking src directory: %w", err)
		}
	} else {
		// No src directory, look for .rs files in root
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && (info.Name() == ".git" || info.Name() == "target") {
				return filepath.SkipDir
			}
			if !info.IsDir() && strings.HasSuffix(path, ".rs") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walking directory: %w", err)
		}
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

	// Try semantic enrichment via cargo/rustc
	if a.IsAvailable() {
		a.enrichWithCargo(dir, analysis)
	}

	// Build graphs
	a.buildCallGraph(analysis)
	a.buildTypeGraph(analysis)
	a.extractDependencies(dir, analysis)

	return analysis, nil
}

// AnalyzeFile performs semantic analysis on a single Rust file
func (a *RustAnalyzer) AnalyzeFile(path string, content []byte) (*FileAnalysis, error) {
	// Use tree-sitter for structural parsing
	return a.TreeSitterAnalysis(content, path)
}

// enrichWithCargo enriches analysis with cargo/rustc diagnostic information
func (a *RustAnalyzer) enrichWithCargo(dir string, analysis *Analysis) {
	ctx := context.Background()

	// Check if this is a Cargo project
	cargoToml := filepath.Join(dir, "Cargo.toml")
	if _, err := os.Stat(cargoToml); err != nil {
		return // Not a Cargo project
	}

	// Run cargo check with JSON output
	a.SubprocessAnalyzer.command = "cargo"
	a.SubprocessAnalyzer.args = nil

	output, err := a.RunCommand(ctx, "check", "--message-format=json", "--manifest-path", cargoToml)
	if err != nil {
		// Cargo check returns non-zero on errors
		// Parse output for diagnostic information
	}

	// Parse cargo JSON output
	if output != nil {
		a.parseCargoOutput(output, analysis)
	}
}

// CargoMessage represents a message from cargo's JSON output
type CargoMessage struct {
	Reason  string `json:"reason"`
	Message *struct {
		Code    *struct {
			Code string `json:"code"`
		} `json:"code"`
		Level   string `json:"level"`
		Message string `json:"message"`
		Spans   []struct {
			FileName string `json:"file_name"`
			LineStart int `json:"line_start"`
			LineEnd   int `json:"line_end"`
		} `json:"spans"`
	} `json:"message"`
}

// parseCargoOutput parses cargo's JSON output
func (a *RustAnalyzer) parseCargoOutput(output []byte, analysis *Analysis) {
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}

		var msg CargoMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}

		if msg.Reason != "compiler-message" || msg.Message == nil {
			continue
		}

		severity := SeverityInfo
		switch msg.Message.Level {
		case "error":
			severity = SeverityError
		case "warning":
			severity = SeverityWarning
		}

		errMsg := AnalysisError{
			Message:  msg.Message.Message,
			Severity: severity,
		}

		if len(msg.Message.Spans) > 0 {
			span := msg.Message.Spans[0]
			errMsg.File = span.FileName
			errMsg.Line = span.LineStart
		}

		analysis.Errors = append(analysis.Errors, errMsg)
	}
}

// findProjectName finds the project name from Cargo.toml
func (a *RustAnalyzer) findProjectName(dir string) string {
	cargoToml := filepath.Join(dir, "Cargo.toml")
	if data, err := os.ReadFile(cargoToml); err == nil {
		// Simple TOML parsing for name
		lines := strings.Split(string(data), "\n")
		inPackage := false
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "[package]" {
				inPackage = true
				continue
			}
			if strings.HasPrefix(line, "[") {
				inPackage = false
				continue
			}
			if inPackage && strings.HasPrefix(line, "name") {
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
func (a *RustAnalyzer) buildCallGraph(analysis *Analysis) {
	for _, fn := range analysis.Functions {
		analysis.CallGraph[fn.Name] = fn.Calls
	}
}

// buildTypeGraph builds the type graph from analysis
func (a *RustAnalyzer) buildTypeGraph(analysis *Analysis) {
	for _, typ := range analysis.Types {
		if len(typ.ImplementsInterfaces) > 0 {
			analysis.TypeGraph[typ.Name] = typ.ImplementsInterfaces
		}
	}
}

// extractDependencies extracts dependencies from Cargo.toml
func (a *RustAnalyzer) extractDependencies(dir string, analysis *Analysis) {
	// Parse imports from analysis
	seen := make(map[string]bool)

	for _, file := range analysis.Files {
		for _, imp := range file.Imports {
			if seen[imp.Path] {
				continue
			}
			seen[imp.Path] = true

			// Check if it's a std crate
			isStdLib := isRustStdLib(imp.Path)
			isLocal := strings.HasPrefix(imp.Path, "crate::") ||
				strings.HasPrefix(imp.Path, "self::") ||
				strings.HasPrefix(imp.Path, "super::")

			dep := Dependency{
				Path:     imp.Path,
				IsStdLib: isStdLib,
				IsLocal:  isLocal,
			}
			analysis.Dependencies = append(analysis.Dependencies, dep)
		}
	}

	// Also parse Cargo.toml for declared dependencies
	cargoToml := filepath.Join(dir, "Cargo.toml")
	if data, err := os.ReadFile(cargoToml); err == nil {
		lines := strings.Split(string(data), "\n")
		inDeps := false
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "[dependencies]" || line == "[dev-dependencies]" {
				inDeps = true
				continue
			}
			if strings.HasPrefix(line, "[") {
				inDeps = false
				continue
			}
			if inDeps && strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
				parts := strings.Split(line, "=")
				if len(parts) >= 1 {
					name := strings.TrimSpace(parts[0])
					if name != "" && !seen[name] {
						seen[name] = true
						analysis.Dependencies = append(analysis.Dependencies, Dependency{
							Path:     name,
							IsStdLib: false,
							IsLocal:  false,
						})
					}
				}
			}
		}
	}
}

// isRustStdLib checks if a crate is part of Rust standard library
func isRustStdLib(path string) bool {
	stdCrates := map[string]bool{
		"std": true, "core": true, "alloc": true, "collections": true,
		"proc_macro": true, "test": true,
	}

	// Get the top-level crate
	parts := strings.Split(path, "::")
	return stdCrates[parts[0]]
}
