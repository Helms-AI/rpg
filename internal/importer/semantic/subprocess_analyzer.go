package semantic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/kon1790/rpg/internal/importer/treesitter"
)

// SubprocessAnalyzer provides semantic analysis via external tools
type SubprocessAnalyzer struct {
	lang       treesitter.Language
	command    string
	args       []string
	timeout    time.Duration
	tsParser   *treesitter.Parser
	available  *bool // Cached availability check
}

// SubprocessConfig configures a subprocess analyzer
type SubprocessConfig struct {
	Language treesitter.Language
	Command  string
	Args     []string
	Timeout  time.Duration
}

// NewSubprocessAnalyzer creates a new subprocess-based analyzer
func NewSubprocessAnalyzer(config SubprocessConfig) *SubprocessAnalyzer {
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &SubprocessAnalyzer{
		lang:     config.Language,
		command:  config.Command,
		args:     config.Args,
		timeout:  timeout,
		tsParser: treesitter.NewParser(),
	}
}

// Language returns the language this analyzer handles
func (a *SubprocessAnalyzer) Language() treesitter.Language {
	return a.lang
}

// IsAvailable checks if the analyzer's command is available
func (a *SubprocessAnalyzer) IsAvailable() bool {
	if a.available != nil {
		return *a.available
	}

	_, err := exec.LookPath(a.command)
	available := err == nil
	a.available = &available
	return available
}

// RunCommand runs the analyzer command with the given args
func (a *SubprocessAnalyzer) RunCommand(ctx context.Context, extraArgs ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	args := append(a.args, extraArgs...)
	cmd := exec.CommandContext(ctx, a.command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("command timed out after %v", a.timeout)
		}
		return nil, fmt.Errorf("command failed: %w, stderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

// RunCommandWithInput runs the analyzer command with stdin input
func (a *SubprocessAnalyzer) RunCommandWithInput(ctx context.Context, input []byte, extraArgs ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	args := append(a.args, extraArgs...)
	cmd := exec.CommandContext(ctx, a.command, args...)
	cmd.Stdin = bytes.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("command timed out after %v", a.timeout)
		}
		return nil, fmt.Errorf("command failed: %w, stderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

// ParseJSON parses JSON output from a command into the given structure
func ParseJSON[T any](data []byte) (T, error) {
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("parsing JSON output: %w", err)
	}
	return result, nil
}

// TreeSitterAnalysis performs tree-sitter based analysis as a fallback
func (a *SubprocessAnalyzer) TreeSitterAnalysis(code []byte, filename string) (*FileAnalysis, error) {
	parser, ok := a.tsParser.GetParser(a.lang)
	if !ok {
		return nil, fmt.Errorf("no tree-sitter parser for language: %s", a.lang)
	}

	result, err := parser.Parse(code, filename)
	if err != nil {
		return nil, fmt.Errorf("tree-sitter parsing: %w", err)
	}

	return convertParseResult(result), nil
}

// convertParseResult converts a tree-sitter ParseResult to FileAnalysis
func convertParseResult(pr *treesitter.ParseResult) *FileAnalysis {
	analysis := &FileAnalysis{
		Path:    pr.FileName,
		Package: "", // Package is not available from tree-sitter parse
	}

	// Convert functions
	for _, fn := range pr.Functions {
		rf := ResolvedFunction{FunctionDef: fn}
		for _, p := range fn.Parameters {
			rf.ResolvedParameters = append(rf.ResolvedParameters, ResolvedParameter{
				Parameter:    p,
				ResolvedType: p.Type,
			})
		}
		if fn.ReturnType != "" {
			rf.ResolvedReturnTypes = []string{fn.ReturnType}
		}
		for _, call := range fn.Calls {
			rf.ResolvedCalls = append(rf.ResolvedCalls, CallReference{
				Name:         call,
				ResolvedName: call,
			})
		}
		analysis.Functions = append(analysis.Functions, rf)
	}

	// Convert types
	for _, typ := range pr.Types {
		rt := ResolvedType{TypeDef: typ}
		for _, f := range typ.Fields {
			rt.ResolvedFields = append(rt.ResolvedFields, ResolvedField{
				Field:        f,
				ResolvedType: f.Type,
			})
		}
		analysis.Types = append(analysis.Types, rt)
	}

	// Convert imports
	for _, imp := range pr.Imports {
		analysis.Imports = append(analysis.Imports, ResolvedImport{
			Import:      imp,
			PackageName: imp.Alias,
		})
	}

	return analysis
}
