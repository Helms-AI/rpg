// Package treesitter provides AST parsing capabilities using tree-sitter
// for cross-language source code analysis.
package treesitter

// Language represents supported programming languages
type Language string

const (
	LanguageGo         Language = "go"
	LanguageTypeScript Language = "typescript"
	LanguagePython     Language = "python"
	LanguageJava       Language = "java"
	LanguageRust       Language = "rust"
	LanguageCSharp     Language = "csharp"
)

// SourceLocation represents a position in source code
type SourceLocation struct {
	File      string `json:"file"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	StartCol  int    `json:"start_col"`
	EndCol    int    `json:"end_col"`
}

// FunctionDef represents an extracted function definition
type FunctionDef struct {
	Name        string         `json:"name"`
	Signature   string         `json:"signature"`
	Parameters  []Parameter    `json:"parameters"`
	ReturnType  string         `json:"return_type"`
	IsAsync     bool           `json:"is_async"`
	IsPublic    bool           `json:"is_public"`
	IsStatic    bool           `json:"is_static"`
	DocComment  string         `json:"doc_comment,omitempty"`
	Body        string         `json:"body,omitempty"`
	Location    SourceLocation `json:"location"`
	ASTHash     string         `json:"ast_hash,omitempty"`
	Calls       []string       `json:"calls,omitempty"`       // Functions called within this function
	Complexity  int            `json:"complexity,omitempty"`  // Cyclomatic complexity
}

// Parameter represents a function parameter
type Parameter struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"default_value,omitempty"`
	IsOptional   bool   `json:"is_optional"`
	IsVariadic   bool   `json:"is_variadic"`
}

// TypeDef represents an extracted type definition
type TypeDef struct {
	Name       string         `json:"name"`
	Kind       TypeKind       `json:"kind"`
	Fields     []Field        `json:"fields,omitempty"`
	Methods    []string       `json:"methods,omitempty"`
	Implements []string       `json:"implements,omitempty"`
	Extends    string         `json:"extends,omitempty"`
	Variants   []string       `json:"variants,omitempty"` // For enums
	AliasOf    string         `json:"alias_of,omitempty"` // For type aliases
	Generic    []string       `json:"generic,omitempty"`  // Generic type parameters
	DocComment string         `json:"doc_comment,omitempty"`
	IsPublic   bool           `json:"is_public"`
	Location   SourceLocation `json:"location"`
	ASTHash    string         `json:"ast_hash,omitempty"`
}

// TypeKind represents the kind of type definition
type TypeKind string

const (
	TypeKindStruct    TypeKind = "struct"
	TypeKindClass     TypeKind = "class"
	TypeKindInterface TypeKind = "interface"
	TypeKindEnum      TypeKind = "enum"
	TypeKindAlias     TypeKind = "alias"
	TypeKindUnion     TypeKind = "union"
)

// Field represents a field in a struct/class
type Field struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Tags       string `json:"tags,omitempty"` // e.g., JSON tags in Go
	IsOptional bool   `json:"is_optional"`
	IsReadonly bool   `json:"is_readonly"`
	Default    string `json:"default,omitempty"`
	DocComment string `json:"doc_comment,omitempty"`
}

// Import represents an import/dependency
type Import struct {
	Path    string `json:"path"`
	Alias   string `json:"alias,omitempty"`
	IsLocal bool   `json:"is_local"` // Local vs external dependency
	Items   []string `json:"items,omitempty"` // For named imports
}

// ParseResult contains the complete AST analysis result
type ParseResult struct {
	Language    Language      `json:"language"`
	FileName    string        `json:"file_name"`
	Functions   []FunctionDef `json:"functions"`
	Types       []TypeDef     `json:"types"`
	Imports     []Import      `json:"imports"`
	Constants   []Constant    `json:"constants"`
	Errors      []ParseError  `json:"errors,omitempty"`
	RawAST      interface{}   `json:"-"` // Internal: the raw tree-sitter tree
}

// Constant represents a constant definition
type Constant struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Value      string         `json:"value"`
	DocComment string         `json:"doc_comment,omitempty"`
	IsPublic   bool           `json:"is_public"`
	Location   SourceLocation `json:"location"`
}

// ParseError represents a parsing error
type ParseError struct {
	Message  string         `json:"message"`
	Location SourceLocation `json:"location"`
}

// NormalizedSignature is a language-agnostic function signature
// used for cross-language comparison
type NormalizedSignature struct {
	Name           string   `json:"name"`
	ParameterTypes []string `json:"parameter_types"`
	ReturnTypes    []string `json:"return_types"`
	IsAsync        bool     `json:"is_async"`
	HasError       bool     `json:"has_error"` // Returns error (Go) or throws (Java/TS)
}

// NormalizedType is a language-agnostic type representation
type NormalizedType struct {
	Name       string            `json:"name"`
	Kind       TypeKind          `json:"kind"`
	Fields     map[string]string `json:"fields"`     // field name -> normalized type
	Methods    []string          `json:"methods"`    // method names
	Implements []string          `json:"implements"` // interface names
}

// ProjectAnalysis represents the complete analysis of a project
type ProjectAnalysis struct {
	Name         string           `json:"name"`
	Language     Language         `json:"language"`
	Files        []ParseResult    `json:"files"`
	AllFunctions []FunctionDef    `json:"all_functions"`
	AllTypes     []TypeDef        `json:"all_types"`
	AllImports   []Import         `json:"all_imports"`
	CallGraph    map[string][]string `json:"call_graph,omitempty"`    // func -> called funcs
	TypeGraph    map[string][]string `json:"type_graph,omitempty"`    // type -> implemented interfaces
	Dependencies []string         `json:"dependencies"`              // External dependencies
}
