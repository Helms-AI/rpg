package languages

import "fmt"

// RustAdapter implements LanguageAdapter for Rust.
type RustAdapter struct {
	language Language
}

// NewRustAdapter creates a new Rust language adapter.
func NewRustAdapter() *RustAdapter {
	return &RustAdapter{
		language: Language{
			ID:            "rust",
			Name:          "Rust",
			Version:       "2021 edition",
			FileExtension: ".rs",
			Conventions: Conventions{
				Naming: NamingConventions{
					Functions:  "snake_case",
					Variables:  "snake_case",
					Constants:  "SCREAMING_SNAKE_CASE",
					Types:      "PascalCase (structs, enums, traits)",
					Packages:   "snake_case (crate/module names)",
					Private:    "No prefix, private by default",
				},
				ErrorHandling: "Use Result<T, E> for recoverable errors, panic!() only for bugs",
				FileNaming:    "snake_case.rs",
				Imports:       "use statements at top, group by: std, external crates, crate modules",
				DocStyle:      "/// for items, //! for modules (markdown supported)",
			},
			Idioms: []string{
				"Use Result<T, E> for operations that can fail",
				"Use Option<T> for nullable values",
				"Use the ? operator for error propagation",
				"Prefer borrowing over ownership transfer",
				"Use derive macros: Debug, Clone, PartialEq, Serialize, Deserialize",
				"Implement From/Into for type conversions",
				"Use pattern matching exhaustively",
				"Prefer iterators over manual loops",
				"Use clippy and rustfmt",
				"Document with /// for public items",
			},
			ProjectStructure: ProjectStructure{
				SourceDir:   "src/",
				TestDir:     "tests/",
				TestSuffix:  "_test.rs",
				PackageFile: "Cargo.toml",
				EntryPoint:  "src/lib.rs or src/main.rs",
				CommonDirs:  []string{"src/", "tests/", "benches/", "examples/"},
			},
			ErrorPatterns: ErrorPatterns{
				Style:       "result",
				CustomError: "#[derive(Error, Debug)]\npub enum MyError { ... }",
				WrapError:   "map_err or thiserror's #[from]",
			},
			Dependencies: DependencyInfo{
				Manager:    "cargo",
				InstallCmd: "cargo build",
				AddCmd:     "cargo add",
				LockFile:   "Cargo.lock",
			},
		},
	}
}

// GetLanguage returns the Rust language configuration.
func (a *RustAdapter) GetLanguage() Language {
	return a.language
}

// GetPromptContext returns Rust-specific prompt instructions.
func (a *RustAdapter) GetPromptContext() string {
	return `Generate idiomatic Rust code following these conventions:

## Language: Rust 2021 Edition

## Naming Conventions
- Functions/variables: snake_case
- Types (structs, enums, traits): PascalCase
- Constants: SCREAMING_SNAKE_CASE
- Modules/crates: snake_case
- Lifetimes: short lowercase ('a, 'b)

## Error Handling
- Use Result<T, E> for recoverable errors
- Use Option<T> for nullable values
- Use ? operator for error propagation
- Use thiserror for custom error types
- Use anyhow for application errors

## Ownership & Borrowing
- Prefer borrowing (&T, &mut T) over ownership transfer
- Use owned types when storing or returning
- Clone explicitly when needed
- Use Cow<'_, T> for flexibility

## Code Style
- Use derive macros: #[derive(Debug, Clone, PartialEq)]
- Add serde derives for serialization: #[derive(Serialize, Deserialize)]
- Use pattern matching exhaustively
- Prefer iterators over manual loops
- Document with /// (markdown supported)

## Output Requirements
- Generate complete, compilable Rust code
- Include all use statements
- Add documentation comments for public items
- Use proper module structure
- Handle all Result/Option cases`
}

// GetProjectStructure returns the recommended Rust project structure.
func (a *RustAdapter) GetProjectStructure(specName string, hasTests bool) []ProjectFile {
	files := []ProjectFile{
		{
			Path:        "Cargo.toml",
			Purpose:     "config",
			Description: "Cargo package manifest",
		},
		{
			Path:        "src/lib.rs",
			Purpose:     "source",
			Description: "Library root with re-exports",
		},
	}

	if hasTests {
		files = append(files, ProjectFile{
			Path:        fmt.Sprintf("tests/%s_test.rs", specName),
			Purpose:     "test",
			Description: "Integration tests",
		})
	}

	return files
}

// MapType maps a pseudo-code type to Rust's equivalent.
func (a *RustAdapter) MapType(pseudoType string) string {
	typeMap := map[string]string{
		"Text":      "String",
		"String":    "String",
		"Integer":   "i64",
		"Int":       "i64",
		"Number":    "i64",
		"Float":     "f64",
		"Decimal":   "f64",
		"Boolean":   "bool",
		"Bool":      "bool",
		"Timestamp": "chrono::DateTime<chrono::Utc>",
		"Duration":  "std::time::Duration",
		"UUID":      "uuid::Uuid",
		"Any":       "serde_json::Value",
		"Nothing":   "()",
		"Void":      "()",
		"None":      "()",
	}

	if rustType, ok := typeMap[pseudoType]; ok {
		return rustType
	}
	return pseudoType
}
