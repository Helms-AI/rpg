// Package extractors provides language-specific source code extraction.
package extractors

import (
	"regexp"
	"strings"
)

// ExtractedProject holds all information extracted from source code.
type ExtractedProject struct {
	Name             string
	Description      string
	DetectedLanguage string
	Types            []ExtractedType
	Functions        []ExtractedFunction
	Tests            []ExtractedTest
	Dependencies     []string
	Warnings         []string
}

// ExtractedType represents a discovered type definition.
type ExtractedType struct {
	Name        string
	Kind        string // "struct", "enum", "interface", "alias"
	Description string
	Fields      []ExtractedField
	Variants    []string // For enums
	AliasOf     string   // For type aliases
	SourceFile  string
	LineNumber  int
}

// ExtractedField represents a field in a struct/class.
type ExtractedField struct {
	Name        string
	Type        string
	Description string
	Optional    bool
	Default     string
}

// ExtractedFunction represents a discovered function/method.
type ExtractedFunction struct {
	Name        string
	Description string
	Parameters  []ExtractedParam
	Returns     string
	IsAsync     bool
	IsPure      bool
	Logic       string
	Errors      []string
	SourceFile  string
	LineNumber  int
	Body        string // Raw function body for logic inference
}

// ExtractedParam represents a function parameter.
type ExtractedParam struct {
	Name        string
	Type        string
	Description string
	Default     string
}

// ExtractedTest represents a discovered test case.
type ExtractedTest struct {
	Function   string
	Name       string
	Given      interface{}
	When       string
	Expect     interface{}
	SourceFile string
}

// LanguageExtractor defines language-specific extraction behavior.
type LanguageExtractor interface {
	// LanguageID returns the language identifier (go, typescript, python, etc.)
	LanguageID() string

	// Extensions returns file extensions for this language (e.g., [".go"])
	Extensions() []string

	// IsTestFile returns true if the file is a test file
	IsTestFile(filename string) bool

	// ExtractTypes extracts type definitions from source code
	ExtractTypes(content string, filename string) []ExtractedType

	// ExtractFunctions extracts function signatures and bodies
	ExtractFunctions(content string, filename string) []ExtractedFunction

	// ExtractTests extracts test cases from test files
	ExtractTests(content string, filename string) []ExtractedTest

	// MapTypeToSpec converts language type to spec pseudo-type
	MapTypeToSpec(langType string) string

	// ExtractPackageDescription extracts package/module description
	ExtractPackageDescription(content string) string
}

// Registry holds all registered language extractors.
type Registry struct {
	extractors map[string]LanguageExtractor
}

// NewRegistry creates a new extractor registry with all languages.
func NewRegistry() *Registry {
	r := &Registry{
		extractors: make(map[string]LanguageExtractor),
	}

	// Register all extractors
	r.Register(NewGoExtractor())
	r.Register(NewTypeScriptExtractor())
	r.Register(NewPythonExtractor())
	r.Register(NewJavaExtractor())
	r.Register(NewRustExtractor())
	r.Register(NewCSharpExtractor())

	return r
}

// Register adds an extractor to the registry.
func (r *Registry) Register(e LanguageExtractor) {
	r.extractors[e.LanguageID()] = e
}

// Get returns an extractor by language ID.
func (r *Registry) Get(languageID string) (LanguageExtractor, bool) {
	e, ok := r.extractors[strings.ToLower(languageID)]
	return e, ok
}

// GetByExtension returns an extractor that handles the given file extension.
func (r *Registry) GetByExtension(ext string) (LanguageExtractor, bool) {
	for _, e := range r.extractors {
		for _, langExt := range e.Extensions() {
			if langExt == ext {
				return e, true
			}
		}
	}
	return nil, false
}

// List returns all registered extractors.
func (r *Registry) List() []LanguageExtractor {
	var list []LanguageExtractor
	for _, e := range r.extractors {
		list = append(list, e)
	}
	return list
}

// Common utility functions for all extractors

// ExtractCommentAbove extracts the comment block immediately above a line number.
func ExtractCommentAbove(lines []string, lineNum int, singleLinePrefix, blockStart, blockEnd string) string {
	if lineNum <= 0 || lineNum > len(lines) {
		return ""
	}

	var comments []string
	idx := lineNum - 2 // Convert to 0-indexed, then go to previous line

	// Check for block comment
	for idx >= 0 {
		line := strings.TrimSpace(lines[idx])
		if strings.HasSuffix(line, blockEnd) && blockEnd != "" {
			// Found end of block comment, search for start
			for idx >= 0 {
				line = strings.TrimSpace(lines[idx])
				comments = append([]string{line}, comments...)
				if strings.HasPrefix(line, blockStart) {
					break
				}
				idx--
			}
			break
		} else if strings.HasPrefix(line, singleLinePrefix) {
			// Single line comment
			comment := strings.TrimPrefix(line, singleLinePrefix)
			comment = strings.TrimSpace(comment)
			comments = append([]string{comment}, comments...)
			idx--
		} else if line == "" {
			// Skip empty lines between comment and definition
			idx--
		} else {
			// Not a comment, stop
			break
		}
	}

	return strings.Join(comments, " ")
}

// InferLogicFromBody infers pseudo-code logic description from function body.
func InferLogicFromBody(body string) string {
	var steps []string

	// Define patterns and their descriptions
	patterns := []struct {
		pattern     *regexp.Regexp
		description string
	}{
		{regexp.MustCompile(`(?i)ToLower|toLowerCase|\.lower\(\)`), "convert to lowercase"},
		{regexp.MustCompile(`(?i)ToUpper|toUpperCase|\.upper\(\)`), "convert to uppercase"},
		{regexp.MustCompile(`(?i)\.Trim|\.trim\(\)|\.strip\(\)`), "trim whitespace"},
		{regexp.MustCompile(`(?i)for\s+|\.forEach|\.map\(|\.each\(`), "iterate over elements"},
		{regexp.MustCompile(`(?i)if\s+.*==\s*nil|if\s+.*==\s*null|if\s+.*is\s+None|if\s+!\w+`), "check for null/empty"},
		{regexp.MustCompile(`(?i)regexp|Regex|re\.|Pattern\.compile`), "apply regex pattern"},
		{regexp.MustCompile(`(?i)append|\.push\(|\.add\(|\.append\(`), "add to collection"},
		{regexp.MustCompile(`(?i)\.filter\(|\.Where\(|filter\(`), "filter elements"},
		{regexp.MustCompile(`(?i)\.sort\(|\.Sort\(|sorted\(`), "sort elements"},
		{regexp.MustCompile(`(?i)\.split\(|\.Split\(|split\(`), "split string"},
		{regexp.MustCompile(`(?i)\.join\(|\.Join\(|join\(`), "join elements"},
		{regexp.MustCompile(`(?i)\.replace\(|\.Replace\(|replace\(`), "replace text"},
		{regexp.MustCompile(`(?i)\.contains\(|\.Contains\(|\.includes\(|\sin\s`), "check if contains"},
		{regexp.MustCompile(`(?i)len\(|\.length|\.Length|\.size\(\)|\.count\(`), "get length/count"},
		{regexp.MustCompile(`(?i)return\s+`), "return the result"},
	}

	seen := make(map[string]bool)
	for _, p := range patterns {
		if p.pattern.MatchString(body) && !seen[p.description] {
			steps = append(steps, p.description)
			seen[p.description] = true
		}
	}

	if len(steps) == 0 {
		return "// TODO: Add logic description"
	}

	return strings.Join(steps, "\n")
}

// FindLineNumber finds the line number where a pattern first appears.
func FindLineNumber(content string, pattern string) int {
	idx := strings.Index(content, pattern)
	if idx == -1 {
		return 0
	}
	return strings.Count(content[:idx], "\n") + 1
}

// ExtractFunctionBody extracts the body of a function given its start position.
// It handles brace matching to find the complete body.
func ExtractFunctionBody(content string, startIdx int) string {
	if startIdx < 0 || startIdx >= len(content) {
		return ""
	}

	// Find the opening brace
	braceIdx := strings.Index(content[startIdx:], "{")
	if braceIdx == -1 {
		return ""
	}
	braceIdx += startIdx

	// Match braces to find the end
	depth := 1
	endIdx := braceIdx + 1
	inString := false
	stringChar := rune(0)

	for endIdx < len(content) && depth > 0 {
		ch := rune(content[endIdx])

		// Handle string literals
		if !inString && (ch == '"' || ch == '\'' || ch == '`') {
			inString = true
			stringChar = ch
		} else if inString && ch == stringChar && (endIdx == 0 || content[endIdx-1] != '\\') {
			inString = false
		}

		if !inString {
			if ch == '{' {
				depth++
			} else if ch == '}' {
				depth--
			}
		}
		endIdx++
	}

	if depth != 0 {
		return ""
	}

	return content[braceIdx+1 : endIdx-1]
}

// CleanType removes common type annotations and normalizes type strings.
func CleanType(typeStr string) string {
	typeStr = strings.TrimSpace(typeStr)
	// Remove pointer/reference markers
	typeStr = strings.TrimPrefix(typeStr, "*")
	typeStr = strings.TrimPrefix(typeStr, "&")
	// Remove array brackets
	typeStr = strings.TrimSuffix(typeStr, "[]")
	typeStr = strings.TrimPrefix(typeStr, "[]")
	return typeStr
}
