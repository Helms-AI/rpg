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

// PythonAnalyzer provides semantic analysis for Python code
type PythonAnalyzer struct {
	*SubprocessAnalyzer
}

// NewPythonAnalyzer creates a new Python semantic analyzer
func NewPythonAnalyzer() *PythonAnalyzer {
	return &PythonAnalyzer{
		SubprocessAnalyzer: NewSubprocessAnalyzer(SubprocessConfig{
			Language: treesitter.LanguagePython,
			Command:  "python3",
			Args:     []string{"-c"},
		}),
	}
}

// IsAvailable checks if Python is available
func (a *PythonAnalyzer) IsAvailable() bool {
	if a.SubprocessAnalyzer.IsAvailable() {
		return true
	}
	// Try python instead of python3
	a.SubprocessAnalyzer.command = "python"
	return a.SubprocessAnalyzer.IsAvailable()
}

// Analyze performs semantic analysis on a Python project directory
func (a *PythonAnalyzer) Analyze(dir string) (*Analysis, error) {
	analysis := &Analysis{
		Language:  treesitter.LanguagePython,
		CallGraph: make(map[string][]string),
		TypeGraph: make(map[string][]string),
	}

	// Find project name
	analysis.Name = a.findProjectName(dir)

	// Find all Python files
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			name := info.Name()
			if name == "__pycache__" || name == ".git" || name == "venv" || name == ".venv" || name == "env" {
				return filepath.SkipDir
			}
		}
		if !info.IsDir() && strings.HasSuffix(path, ".py") {
			if !strings.HasSuffix(path, "_test.py") && !strings.Contains(path, "test_") {
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

	// Try semantic enrichment via Python ast + mypy
	if a.IsAvailable() {
		a.enrichWithPythonAST(dir, analysis)
	}

	// Build graphs
	a.buildCallGraph(analysis)
	a.buildTypeGraph(analysis)
	a.extractDependencies(analysis)

	return analysis, nil
}

// AnalyzeFile performs semantic analysis on a single Python file
func (a *PythonAnalyzer) AnalyzeFile(path string, content []byte) (*FileAnalysis, error) {
	// Use tree-sitter for structural parsing
	return a.TreeSitterAnalysis(content, path)
}

// Python script for deep AST analysis
const pythonASTScript = `
import ast
import json
import sys

def analyze_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()

    try:
        tree = ast.parse(content, filename=filepath)
    except SyntaxError as e:
        return {"error": str(e)}

    result = {
        "functions": [],
        "classes": [],
        "imports": [],
        "calls": []
    }

    for node in ast.walk(tree):
        if isinstance(node, ast.FunctionDef) or isinstance(node, ast.AsyncFunctionDef):
            func = {
                "name": node.name,
                "line": node.lineno,
                "is_async": isinstance(node, ast.AsyncFunctionDef),
                "decorators": [ast.dump(d) for d in node.decorator_list],
                "args": [],
                "returns": None
            }
            for arg in node.args.args:
                func["args"].append({
                    "name": arg.arg,
                    "annotation": ast.dump(arg.annotation) if arg.annotation else None
                })
            if node.returns:
                func["returns"] = ast.dump(node.returns)
            result["functions"].append(func)

        elif isinstance(node, ast.ClassDef):
            cls = {
                "name": node.name,
                "line": node.lineno,
                "bases": [ast.dump(b) for b in node.bases],
                "methods": []
            }
            for item in node.body:
                if isinstance(item, (ast.FunctionDef, ast.AsyncFunctionDef)):
                    cls["methods"].append(item.name)
            result["classes"].append(cls)

        elif isinstance(node, ast.Import):
            for alias in node.names:
                result["imports"].append({
                    "module": alias.name,
                    "alias": alias.asname
                })

        elif isinstance(node, ast.ImportFrom):
            for alias in node.names:
                result["imports"].append({
                    "module": node.module or "",
                    "name": alias.name,
                    "alias": alias.asname,
                    "level": node.level
                })

        elif isinstance(node, ast.Call):
            if isinstance(node.func, ast.Name):
                result["calls"].append(node.func.id)
            elif isinstance(node.func, ast.Attribute):
                result["calls"].append(node.func.attr)

    return result

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"error": "no file provided"}))
    else:
        print(json.dumps(analyze_file(sys.argv[1])))
`

// enrichWithPythonAST enriches analysis with Python AST information
func (a *PythonAnalyzer) enrichWithPythonAST(dir string, analysis *Analysis) {
	ctx := context.Background()

	for _, file := range analysis.Files {
		// Run Python AST script on each file
		output, err := a.RunCommand(ctx, pythonASTScript, file.Path)
		if err != nil {
			analysis.Errors = append(analysis.Errors, AnalysisError{
				File:     file.Path,
				Message:  fmt.Sprintf("Python AST analysis failed: %v", err),
				Severity: SeverityInfo,
			})
			continue
		}

		var astResult struct {
			Functions []struct {
				Name       string `json:"name"`
				Line       int    `json:"line"`
				IsAsync    bool   `json:"is_async"`
				Decorators []string `json:"decorators"`
				Args       []struct {
					Name       string `json:"name"`
					Annotation string `json:"annotation"`
				} `json:"args"`
				Returns string `json:"returns"`
			} `json:"functions"`
			Classes []struct {
				Name    string   `json:"name"`
				Line    int      `json:"line"`
				Bases   []string `json:"bases"`
				Methods []string `json:"methods"`
			} `json:"classes"`
			Imports []struct {
				Module string `json:"module"`
				Name   string `json:"name"`
				Alias  string `json:"alias"`
				Level  int    `json:"level"`
			} `json:"imports"`
			Calls []string `json:"calls"`
			Error string   `json:"error"`
		}

		if err := json.Unmarshal(output, &astResult); err != nil {
			continue
		}

		if astResult.Error != "" {
			analysis.Errors = append(analysis.Errors, AnalysisError{
				File:     file.Path,
				Message:  astResult.Error,
				Severity: SeverityWarning,
			})
			continue
		}

		// Enrich functions with async info and type annotations
		for i, fn := range file.Functions {
			for _, astFn := range astResult.Functions {
				if fn.Name == astFn.Name {
					file.Functions[i].IsAsync = astFn.IsAsync
					// Enrich with type annotations from decorators if available
					break
				}
			}
		}
	}

	// Try mypy for additional type checking
	a.enrichWithMypy(dir, analysis)
}

// enrichWithMypy enriches analysis with mypy type information
func (a *PythonAnalyzer) enrichWithMypy(dir string, analysis *Analysis) {
	ctx := context.Background()

	// Check if mypy is available
	_, err := a.RunCommand(ctx, "import mypy; print('ok')")
	if err != nil {
		return // mypy not available
	}

	// Run mypy
	output, err := a.RunCommand(ctx, "-m", "mypy", "--show-error-codes", "--no-error-summary", dir)
	if err != nil {
		// mypy returns non-zero on type errors, which is expected
		// Parse output for type information
	}

	// Parse mypy output for type errors/hints
	if output != nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if strings.Contains(line, ": error:") {
				analysis.Errors = append(analysis.Errors, AnalysisError{
					Message:  line,
					Severity: SeverityWarning,
				})
			}
		}
	}
}

// findProjectName finds the project name
func (a *PythonAnalyzer) findProjectName(dir string) string {
	// Try setup.py
	setupPath := filepath.Join(dir, "setup.py")
	if _, err := os.Stat(setupPath); err == nil {
		// Could parse setup.py for name but that's complex
	}

	// Try pyproject.toml
	pyprojectPath := filepath.Join(dir, "pyproject.toml")
	if data, err := os.ReadFile(pyprojectPath); err == nil {
		// Simple extraction - look for name = "..."
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "name") {
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
func (a *PythonAnalyzer) buildCallGraph(analysis *Analysis) {
	for _, fn := range analysis.Functions {
		analysis.CallGraph[fn.Name] = fn.Calls
	}
}

// buildTypeGraph builds the type graph from analysis
func (a *PythonAnalyzer) buildTypeGraph(analysis *Analysis) {
	for _, typ := range analysis.Types {
		if len(typ.ImplementsInterfaces) > 0 {
			analysis.TypeGraph[typ.Name] = typ.ImplementsInterfaces
		}
	}
}

// extractDependencies extracts external dependencies
func (a *PythonAnalyzer) extractDependencies(analysis *Analysis) {
	seen := make(map[string]bool)

	for _, file := range analysis.Files {
		for _, imp := range file.Imports {
			if seen[imp.Path] {
				continue
			}
			seen[imp.Path] = true

			// Check if it's a stdlib module
			isStdLib := isPythonStdLib(imp.Path)

			dep := Dependency{
				Path:     imp.Path,
				IsStdLib: isStdLib,
				IsLocal:  strings.HasPrefix(imp.Path, "."),
			}
			analysis.Dependencies = append(analysis.Dependencies, dep)
		}
	}
}

// isPythonStdLib checks if a module is part of Python standard library
func isPythonStdLib(module string) bool {
	stdlib := map[string]bool{
		"os": true, "sys": true, "re": true, "json": true, "typing": true,
		"collections": true, "itertools": true, "functools": true, "pathlib": true,
		"dataclasses": true, "abc": true, "asyncio": true, "logging": true,
		"unittest": true, "datetime": true, "time": true, "math": true,
		"random": true, "string": true, "io": true, "subprocess": true,
		"threading": true, "multiprocessing": true, "contextlib": true,
		"copy": true, "enum": true, "hashlib": true, "http": true,
		"urllib": true, "socket": true, "ssl": true, "email": true,
	}

	// Get the top-level module
	parts := strings.Split(module, ".")
	return stdlib[parts[0]]
}
