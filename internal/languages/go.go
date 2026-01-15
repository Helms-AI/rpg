package languages

import "fmt"

// GoAdapter implements LanguageAdapter for Go.
type GoAdapter struct {
	language Language
}

// NewGoAdapter creates a new Go language adapter.
func NewGoAdapter() *GoAdapter {
	return &GoAdapter{
		language: Language{
			ID:            "go",
			Name:          "Go",
			Version:       "1.21+",
			FileExtension: ".go",
			Conventions: Conventions{
				Naming: NamingConventions{
					Functions:  "PascalCase for exported, camelCase for unexported",
					Variables:  "camelCase, short names in small scopes (i, v, k, err)",
					Constants:  "PascalCase for exported, camelCase for unexported",
					Types:      "PascalCase",
					Packages:   "lowercase, single word, no underscores",
					Private:    "camelCase (unexported)",
				},
				ErrorHandling: "Return (result, error), check err != nil, wrap errors with fmt.Errorf",
				FileNaming:    "snake_case.go",
				Imports:       "Group by: stdlib, external, internal (separated by blank lines)",
				DocStyle:      "godoc (// comment above declaration)",
			},
			Idioms: []string{
				"Use short variable names in small scopes (i, v, k, err)",
				"Return early on errors (guard clauses)",
				"Prefer composition over inheritance",
				"Use interfaces for abstraction, define where consumed",
				"Accept interfaces, return structs",
				"Use defer for cleanup (close files, unlock mutexes)",
				"Context as first parameter for cancellable operations",
				"Errors are values, handle them explicitly",
				"Make the zero value useful",
				"Use table-driven tests",
			},
			ProjectStructure: ProjectStructure{
				SourceDir:   "",
				TestDir:     "",
				TestSuffix:  "_test.go",
				PackageFile: "go.mod",
				EntryPoint:  "main.go",
				CommonDirs:  []string{"cmd/", "internal/", "pkg/"},
			},
			ErrorPatterns: ErrorPatterns{
				Style:       "tuple",
				CustomError: "type MyError struct { ... }\nfunc (e *MyError) Error() string { ... }",
				WrapError:   "fmt.Errorf(\"context: %w\", err)",
			},
			Dependencies: DependencyInfo{
				Manager:    "go mod",
				InstallCmd: "go mod download",
				AddCmd:     "go get",
				LockFile:   "go.sum",
			},
		},
	}
}

// GetLanguage returns the Go language configuration.
func (a *GoAdapter) GetLanguage() Language {
	return a.language
}

// GetPromptContext returns Go-specific prompt instructions.
func (a *GoAdapter) GetPromptContext() string {
	return `Generate idiomatic Go code following these conventions:

## Language: Go 1.21+

## Naming Conventions
- Exported functions/types: PascalCase (e.g., ProcessData, UserService)
- Unexported (private): camelCase (e.g., processData, userService)
- Variables: camelCase, short names in small scopes (i, v, k, err, ctx)
- Packages: lowercase, single word, no underscores

## Error Handling
- Return (result, error) for operations that can fail
- Check errors immediately: if err != nil { return ..., err }
- Wrap errors with context: fmt.Errorf("operation failed: %w", err)
- Return early on errors (guard clauses)

## Code Style
- Group imports: stdlib, then external, then internal (blank line between groups)
- Use defer for cleanup operations
- Accept interfaces, return concrete types
- Make zero values useful
- Context as first parameter for cancellable operations

## Documentation
- Use godoc style: // FunctionName does something.
- Document all exported functions, types, and packages
- Include usage examples in package documentation

## Output Requirements
- Generate complete, compilable Go code
- Include package declaration
- Include all necessary imports
- Add godoc comments for exported items
- Handle all error cases appropriately`
}

// GetProjectStructure returns the recommended Go project structure.
func (a *GoAdapter) GetProjectStructure(specName string, hasTests bool) []ProjectFile {
	files := []ProjectFile{
		{
			Path:        "go.mod",
			Purpose:     "config",
			Description: "Go module definition",
		},
		{
			Path:        fmt.Sprintf("%s.go", specName),
			Purpose:     "source",
			Description: "Main implementation",
		},
	}

	if hasTests {
		files = append(files, ProjectFile{
			Path:        fmt.Sprintf("%s_test.go", specName),
			Purpose:     "test",
			Description: "Unit tests",
		})
	}

	return files
}

// MapType maps a pseudo-code type to Go's equivalent.
func (a *GoAdapter) MapType(pseudoType string) string {
	typeMap := map[string]string{
		"Text":      "string",
		"String":    "string",
		"Integer":   "int",
		"Int":       "int",
		"Number":    "int",
		"Float":     "float64",
		"Decimal":   "float64",
		"Boolean":   "bool",
		"Bool":      "bool",
		"Timestamp": "time.Time",
		"Duration":  "time.Duration",
		"UUID":      "string",
		"Any":       "interface{}",
		"Nothing":   "",
		"Void":      "",
		"None":      "",
	}

	if goType, ok := typeMap[pseudoType]; ok {
		return goType
	}
	return pseudoType
}
