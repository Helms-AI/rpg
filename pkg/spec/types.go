// Package spec defines the data structures for parsed specifications.
package spec

// Spec represents a fully parsed specification file.
type Spec struct {
	Name            string     `json:"name"`
	Description     string     `json:"description,omitempty"`
	Version         string     `json:"version,omitempty"`
	Author          string     `json:"author,omitempty"`
	License         string     `json:"license,omitempty"`
	TargetLanguages []string   `json:"targetLanguages"`
	Dependencies    []string   `json:"dependencies,omitempty"`
	Types           []TypeDef  `json:"types,omitempty"`
	Functions       []Function `json:"functions"`
	Tests           []TestCase `json:"tests,omitempty"`
}

// TypeDef represents a type definition in the spec.
type TypeDef struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Kind        string  `json:"kind"` // "struct", "enum", "alias"
	Fields      []Field `json:"fields,omitempty"`
	Variants    []Field `json:"variants,omitempty"` // For enums
	AliasOf     string  `json:"aliasOf,omitempty"`  // For type aliases
}

// Field represents a field in a struct or enum variant.
type Field struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Default     string `json:"default,omitempty"`
	Optional    bool   `json:"optional,omitempty"`
}

// Function represents a function definition in the spec.
type Function struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Accepts     []Param  `json:"accepts,omitempty"`
	Returns     *Return  `json:"returns,omitempty"`
	Logic       string   `json:"logic"`
	Errors      []string `json:"errors,omitempty"`
	Async       bool     `json:"async,omitempty"`
	Pure        bool     `json:"pure,omitempty"`
}

// Param represents a function parameter.
type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Default     string `json:"default,omitempty"`
}

// Return represents a function's return value.
type Return struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// TestCase represents a test case definition.
type TestCase struct {
	Function string      `json:"function"`
	Name     string      `json:"name"`
	Given    interface{} `json:"given,omitempty"` // Can be simple value or structured
	When     string      `json:"when,omitempty"`
	Expect   interface{} `json:"expect"`
}

// ValidationError represents an error found during spec validation.
type ValidationError struct {
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
	Severity string `json:"severity"` // "error", "warning", "info"
	Code     string `json:"code"`
	Message  string `json:"message"`
}

// ValidationResult contains the results of spec validation.
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// Normalize ensures all slice fields are non-nil (empty slices instead of nil).
// This prevents JSON serialization from producing null instead of [] which
// can cause MCP schema validation errors.
func (s *Spec) Normalize() {
	if s == nil {
		return
	}
	if s.TargetLanguages == nil {
		s.TargetLanguages = []string{}
	}
	if s.Dependencies == nil {
		s.Dependencies = []string{}
	}
	if s.Types == nil {
		s.Types = []TypeDef{}
	}
	if s.Functions == nil {
		s.Functions = []Function{}
	}
	if s.Tests == nil {
		s.Tests = []TestCase{}
	}
	// Normalize nested types
	for i := range s.Types {
		if s.Types[i].Fields == nil {
			s.Types[i].Fields = []Field{}
		}
		if s.Types[i].Variants == nil {
			s.Types[i].Variants = []Field{}
		}
	}
	// Normalize nested functions
	for i := range s.Functions {
		if s.Functions[i].Accepts == nil {
			s.Functions[i].Accepts = []Param{}
		}
		if s.Functions[i].Errors == nil {
			s.Functions[i].Errors = []string{}
		}
	}
}
