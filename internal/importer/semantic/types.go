// Package semantic provides deep semantic analysis for source code
// beyond what tree-sitter's structural parsing can provide.
package semantic

import (
	"github.com/kon1790/rpg/internal/importer/treesitter"
)

// Analyzer provides semantic analysis capabilities for a specific language
type Analyzer interface {
	// Language returns the language this analyzer handles
	Language() treesitter.Language

	// Analyze performs semantic analysis on a directory
	Analyze(dir string) (*Analysis, error)

	// AnalyzeFile performs semantic analysis on a single file
	AnalyzeFile(path string, content []byte) (*FileAnalysis, error)

	// IsAvailable checks if the analyzer's dependencies are available
	IsAvailable() bool
}

// Analysis represents the complete semantic analysis of a project
type Analysis struct {
	// Language is the primary language analyzed
	Language treesitter.Language `json:"language"`

	// Name is the project name
	Name string `json:"name"`

	// Files contains per-file analysis
	Files []*FileAnalysis `json:"files"`

	// Types contains all type definitions with resolved references
	Types []ResolvedType `json:"types"`

	// Functions contains all functions with resolved types
	Functions []ResolvedFunction `json:"functions"`

	// CallGraph maps function names to functions they call
	CallGraph map[string][]string `json:"call_graph"`

	// TypeGraph maps types to interfaces they implement
	TypeGraph map[string][]string `json:"type_graph"`

	// Dependencies contains external package dependencies
	Dependencies []Dependency `json:"dependencies"`

	// Errors contains any analysis errors
	Errors []AnalysisError `json:"errors,omitempty"`
}

// FileAnalysis represents semantic analysis of a single file
type FileAnalysis struct {
	// Path is the file path
	Path string `json:"path"`

	// Package is the package/module name
	Package string `json:"package"`

	// Types in this file
	Types []ResolvedType `json:"types"`

	// Functions in this file
	Functions []ResolvedFunction `json:"functions"`

	// Imports with resolved packages
	Imports []ResolvedImport `json:"imports"`

	// Errors during analysis
	Errors []AnalysisError `json:"errors,omitempty"`
}

// ResolvedType represents a type with full semantic information
type ResolvedType struct {
	treesitter.TypeDef

	// ResolvedFields contains fields with resolved types
	ResolvedFields []ResolvedField `json:"resolved_fields,omitempty"`

	// MethodSignatures contains full method signatures
	MethodSignatures []MethodSignature `json:"method_signatures,omitempty"`

	// ImplementsInterfaces contains interfaces this type implements
	ImplementsInterfaces []string `json:"implements_interfaces,omitempty"`

	// UsedBy contains types that reference this type
	UsedBy []string `json:"used_by,omitempty"`

	// PackagePath is the full package path
	PackagePath string `json:"package_path,omitempty"`
}

// ResolvedField represents a field with resolved type information
type ResolvedField struct {
	treesitter.Field

	// ResolvedType is the fully qualified type
	ResolvedType string `json:"resolved_type"`

	// IsPointer indicates if the field is a pointer
	IsPointer bool `json:"is_pointer"`

	// IsSlice indicates if the field is a slice/array
	IsSlice bool `json:"is_slice"`

	// IsMap indicates if the field is a map
	IsMap bool `json:"is_map"`

	// ElementType is the element type for slices/maps
	ElementType string `json:"element_type,omitempty"`

	// KeyType is the key type for maps
	KeyType string `json:"key_type,omitempty"`

	// Tags contains struct tags (Go-specific but useful for all)
	Tags string `json:"tags,omitempty"`
}

// MethodSignature represents a method's full signature
type MethodSignature struct {
	// Name is the method name
	Name string `json:"name"`

	// Receiver is the receiver type (for methods)
	Receiver string `json:"receiver,omitempty"`

	// Parameters with resolved types
	Parameters []ResolvedParameter `json:"parameters"`

	// ReturnTypes with resolved types
	ReturnTypes []string `json:"return_types"`

	// IsExported indicates if the method is exported
	IsExported bool `json:"is_exported"`
}

// ResolvedParameter represents a parameter with resolved type
type ResolvedParameter struct {
	treesitter.Parameter

	// ResolvedType is the fully qualified type
	ResolvedType string `json:"resolved_type"`
}

// ResolvedFunction represents a function with full semantic information
type ResolvedFunction struct {
	treesitter.FunctionDef

	// ResolvedParameters contains parameters with resolved types
	ResolvedParameters []ResolvedParameter `json:"resolved_parameters,omitempty"`

	// ResolvedReturnTypes contains resolved return types
	ResolvedReturnTypes []string `json:"resolved_return_types,omitempty"`

	// ResolvedCalls contains resolved function references
	ResolvedCalls []CallReference `json:"resolved_calls,omitempty"`

	// LocalVariables contains local variable declarations
	LocalVariables []Variable `json:"local_variables,omitempty"`

	// PackagePath is the full package path
	PackagePath string `json:"package_path,omitempty"`
}

// CallReference represents a resolved function call
type CallReference struct {
	// Name is the function name as called
	Name string `json:"name"`

	// ResolvedName is the fully qualified name
	ResolvedName string `json:"resolved_name"`

	// Package is the package containing the function
	Package string `json:"package,omitempty"`

	// IsMethod indicates if this is a method call
	IsMethod bool `json:"is_method"`

	// ReceiverType is the receiver type for method calls
	ReceiverType string `json:"receiver_type,omitempty"`
}

// Variable represents a local variable
type Variable struct {
	// Name is the variable name
	Name string `json:"name"`

	// Type is the resolved type
	Type string `json:"type"`

	// Line is the declaration line
	Line int `json:"line"`
}

// ResolvedImport represents an import with resolved information
type ResolvedImport struct {
	treesitter.Import

	// PackageName is the resolved package name
	PackageName string `json:"package_name"`

	// ExportedTypes lists types exported by this package
	ExportedTypes []string `json:"exported_types,omitempty"`

	// ExportedFunctions lists functions exported by this package
	ExportedFunctions []string `json:"exported_functions,omitempty"`
}

// Dependency represents an external dependency
type Dependency struct {
	// Path is the import path
	Path string `json:"path"`

	// Version is the version if known
	Version string `json:"version,omitempty"`

	// IsStdLib indicates if this is a standard library package
	IsStdLib bool `json:"is_stdlib"`

	// IsLocal indicates if this is a local package
	IsLocal bool `json:"is_local"`
}

// AnalysisError represents an error during analysis
type AnalysisError struct {
	// File is the file where the error occurred
	File string `json:"file,omitempty"`

	// Line is the line number
	Line int `json:"line,omitempty"`

	// Message is the error message
	Message string `json:"message"`

	// Severity is the error severity
	Severity Severity `json:"severity"`
}

// Severity represents error severity levels
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// AnalyzerRegistry manages available semantic analyzers
type AnalyzerRegistry struct {
	analyzers map[treesitter.Language]Analyzer
}

// NewAnalyzerRegistry creates a new analyzer registry
func NewAnalyzerRegistry() *AnalyzerRegistry {
	return &AnalyzerRegistry{
		analyzers: make(map[treesitter.Language]Analyzer),
	}
}

// Register registers an analyzer for a language
func (r *AnalyzerRegistry) Register(a Analyzer) {
	r.analyzers[a.Language()] = a
}

// Get returns the analyzer for a language
func (r *AnalyzerRegistry) Get(lang treesitter.Language) (Analyzer, bool) {
	a, ok := r.analyzers[lang]
	return a, ok
}

// GetAvailable returns all analyzers that are available
func (r *AnalyzerRegistry) GetAvailable() []Analyzer {
	var available []Analyzer
	for _, a := range r.analyzers {
		if a.IsAvailable() {
			available = append(available, a)
		}
	}
	return available
}

// Languages returns all registered languages
func (r *AnalyzerRegistry) Languages() []treesitter.Language {
	var langs []treesitter.Language
	for lang := range r.analyzers {
		langs = append(langs, lang)
	}
	return langs
}
