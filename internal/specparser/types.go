// Package specparser provides types and utilities for parsing specification files.
package specparser

// SpecAnalysis contains the fully parsed specification ready for code generation.
type SpecAnalysis struct {
	// Name is the project/spec name
	Name string `json:"name"`

	// Overview provides a description of what this spec defines
	Overview string `json:"overview"`

	// Types contains all type definitions extracted from the spec
	Types []SpecType `json:"types"`

	// Functions contains all function definitions extracted from the spec
	Functions []SpecFunction `json:"functions"`

	// Tests contains test case definitions
	Tests []SpecTest `json:"tests"`

	// Dependencies contains external dependency requirements
	Dependencies []SpecDependency `json:"dependencies"`

	// Configuration contains environment variables and settings
	Configuration []SpecConfig `json:"configuration"`

	// TotalItems is the sum of types, functions, tests, and dependencies
	TotalItems int `json:"totalItems"`
}

// SpecType represents a type definition from the spec.
type SpecType struct {
	// Name of the type
	Name string `json:"name"`

	// Kind is the type kind: struct, interface, enum, alias
	Kind string `json:"kind"`

	// Description explains the purpose of this type
	Description string `json:"description"`

	// Fields for struct types
	Fields []SpecField `json:"fields,omitempty"`

	// Methods for interfaces
	Methods []string `json:"methods,omitempty"`

	// Values for enum types
	Values []SpecEnumValue `json:"values,omitempty"`

	// Implements lists interfaces this type implements
	Implements []string `json:"implements,omitempty"`

	// Generic type parameters
	Generic []string `json:"generic,omitempty"`

	// IsPublic indicates if the type is exported/public
	IsPublic bool `json:"isPublic"`
}

// SpecField represents a field in a type definition.
type SpecField struct {
	// Name of the field
	Name string `json:"name"`

	// Type of the field (using pseudo-types)
	Type string `json:"type"`

	// Description explains the field's purpose
	Description string `json:"description"`

	// Required indicates if the field is required
	Required bool `json:"required"`

	// Tags contains serialization tags (json, xml, etc.)
	Tags map[string]string `json:"tags,omitempty"`

	// Default value if any
	Default string `json:"default,omitempty"`
}

// SpecEnumValue represents a value in an enum type.
type SpecEnumValue struct {
	// Name of the enum value
	Name string `json:"name"`

	// Value is the underlying value
	Value string `json:"value,omitempty"`

	// Description explains this enum value
	Description string `json:"description,omitempty"`
}

// SpecFunction represents a function definition from the spec.
type SpecFunction struct {
	// Name of the function
	Name string `json:"name"`

	// Receiver type if this is a method
	Receiver string `json:"receiver,omitempty"`

	// Description explains what the function does
	Description string `json:"description"`

	// Parameters lists the function parameters
	Parameters []SpecParameter `json:"parameters,omitempty"`

	// Returns lists the return types
	Returns []SpecReturn `json:"returns,omitempty"`

	// Logic describes the implementation logic
	Logic string `json:"logic,omitempty"`

	// Errors lists possible error conditions
	Errors []SpecError `json:"errors,omitempty"`

	// IsAsync indicates if the function is asynchronous
	IsAsync bool `json:"isAsync"`

	// IsPublic indicates if the function is exported/public
	IsPublic bool `json:"isPublic"`

	// Complexity is an optional complexity indicator
	Complexity int `json:"complexity,omitempty"`
}

// SpecParameter represents a function parameter.
type SpecParameter struct {
	// Name of the parameter
	Name string `json:"name"`

	// Type of the parameter (using pseudo-types)
	Type string `json:"type"`

	// Description explains the parameter
	Description string `json:"description"`

	// Required indicates if the parameter is required
	Required bool `json:"required"`

	// Default value if optional
	Default string `json:"default,omitempty"`
}

// SpecReturn represents a function return value.
type SpecReturn struct {
	// Type of the return value
	Type string `json:"type"`

	// Description explains what is returned
	Description string `json:"description"`
}

// SpecError represents an error condition.
type SpecError struct {
	// Condition describes when this error occurs
	Condition string `json:"condition"`

	// Type is the error type or code
	Type string `json:"type"`

	// Message is the error message
	Message string `json:"message"`
}

// SpecTest represents a test case definition.
type SpecTest struct {
	// Name of the test case
	Name string `json:"name"`

	// Description explains what is being tested
	Description string `json:"description"`

	// Target is the function/type being tested
	Target string `json:"target"`

	// Given describes the initial conditions
	Given []SpecCondition `json:"given,omitempty"`

	// When describes the action being tested
	When string `json:"when"`

	// Then describes the expected outcome
	Then []SpecAssertion `json:"then,omitempty"`
}

// SpecCondition represents a test precondition.
type SpecCondition struct {
	// Description of the condition
	Description string `json:"description"`

	// Value is the setup value
	Value string `json:"value,omitempty"`
}

// SpecAssertion represents a test assertion.
type SpecAssertion struct {
	// Description of what is being asserted
	Description string `json:"description"`

	// Expected value
	Expected string `json:"expected,omitempty"`

	// Comparison operator
	Operator string `json:"operator,omitempty"`
}

// SpecDependency represents an external dependency.
type SpecDependency struct {
	// Name of the dependency
	Name string `json:"name"`

	// Version constraint
	Version string `json:"version,omitempty"`

	// Description explains why this dependency is needed
	Description string `json:"description"`

	// Language-specific package names
	Packages map[string]string `json:"packages,omitempty"`
}

// SpecConfig represents a configuration item.
type SpecConfig struct {
	// Name of the configuration (e.g., environment variable name)
	Name string `json:"name"`

	// Type of the value
	Type string `json:"type"`

	// Description explains the configuration
	Description string `json:"description"`

	// Default value
	Default string `json:"default,omitempty"`

	// Required indicates if this configuration must be provided
	Required bool `json:"required"`
}

// CalculateTotals updates the TotalItems field.
func (s *SpecAnalysis) CalculateTotals() {
	s.TotalItems = len(s.Types) + len(s.Functions) + len(s.Tests) + len(s.Dependencies)
}
