package languages

import "fmt"

// PythonAdapter implements LanguageAdapter for Python.
type PythonAdapter struct {
	language Language
}

// NewPythonAdapter creates a new Python language adapter.
func NewPythonAdapter() *PythonAdapter {
	return &PythonAdapter{
		language: Language{
			ID:            "python",
			Name:          "Python",
			Version:       "3.11+",
			FileExtension: ".py",
			Conventions: Conventions{
				Naming: NamingConventions{
					Functions:  "snake_case",
					Variables:  "snake_case",
					Constants:  "SCREAMING_SNAKE_CASE",
					Types:      "PascalCase (classes)",
					Packages:   "lowercase_with_underscores",
					Private:    "_single_underscore prefix",
				},
				ErrorHandling: "Raise exceptions, use try/except/finally",
				FileNaming:    "snake_case.py",
				Imports:       "Group by: stdlib, third-party, local (separated by blank lines)",
				DocStyle:      "Docstrings (triple quotes) with Google or NumPy style",
			},
			Idioms: []string{
				"Use type hints everywhere (PEP 484, 585, 604)",
				"Use dataclasses or Pydantic for data structures",
				"Use context managers (with) for resource management",
				"Use pathlib for file paths",
				"Use logging module (not print)",
				"Use f-strings for formatting",
				"Use list/dict comprehensions appropriately",
				"Use generators for large sequences",
				"Follow PEP 8 style guide",
				"Use black, ruff, and mypy",
			},
			ProjectStructure: ProjectStructure{
				SourceDir:   "src/",
				TestDir:     "tests/",
				TestSuffix:  "_test.py",
				PackageFile: "pyproject.toml",
				EntryPoint:  "__main__.py",
				CommonDirs:  []string{"src/", "tests/", "docs/"},
			},
			ErrorPatterns: ErrorPatterns{
				Style:       "exceptions",
				CustomError: "class MyError(Exception):\n    def __init__(self, message: str) -> None:\n        super().__init__(message)",
				WrapError:   "raise MyError(\"context\") from original_error",
			},
			Dependencies: DependencyInfo{
				Manager:    "pip/poetry/uv",
				InstallCmd: "pip install -e . or poetry install",
				AddCmd:     "pip install or poetry add",
				LockFile:   "poetry.lock or requirements.txt",
			},
		},
	}
}

// GetLanguage returns the Python language configuration.
func (a *PythonAdapter) GetLanguage() Language {
	return a.language
}

// GetPromptContext returns Python-specific prompt instructions.
func (a *PythonAdapter) GetPromptContext() string {
	return `Generate idiomatic Python code following these conventions:

## Language: Python 3.11+

## Naming Conventions
- Functions/variables: snake_case
- Classes: PascalCase
- Constants: SCREAMING_SNAKE_CASE
- Private: _single_underscore prefix
- Modules/packages: lowercase_with_underscores

## Type Hints
- Use type hints for all functions (PEP 484)
- Use modern syntax: list[str] instead of List[str]
- Use | for unions: str | None instead of Optional[str]
- Use TypeVar for generics
- Use dataclasses or Pydantic for structured data

## Error Handling
- Define custom exceptions inheriting from Exception
- Use try/except/finally appropriately
- Chain exceptions with 'from' keyword
- Use logging, not print, for errors

## Code Style
- Follow PEP 8
- Use dataclasses with @dataclass(frozen=True) for immutability
- Use context managers (with) for resources
- Use pathlib for file paths
- Use f-strings for formatting
- Use comprehensions appropriately
- Use generators for large sequences

## Documentation
- Use docstrings for all public functions/classes
- Use Google or NumPy docstring style
- Include type information in docstrings

## Output Requirements
- Generate complete, runnable Python code
- Include all imports at the top
- Add comprehensive type hints
- Add docstrings for public APIs
- Use modern Python 3.11+ features`
}

// GetProjectStructure returns the recommended Python project structure.
func (a *PythonAdapter) GetProjectStructure(specName string, hasTests bool) []ProjectFile {
	files := []ProjectFile{
		{
			Path:        "pyproject.toml",
			Purpose:     "config",
			Description: "Project configuration and dependencies",
		},
		{
			Path:        fmt.Sprintf("src/%s/__init__.py", specName),
			Purpose:     "source",
			Description: "Package initialization",
		},
		{
			Path:        fmt.Sprintf("src/%s/%s.py", specName, specName),
			Purpose:     "source",
			Description: "Main implementation",
		},
	}

	if hasTests {
		files = append(files, ProjectFile{
			Path:        fmt.Sprintf("tests/test_%s.py", specName),
			Purpose:     "test",
			Description: "pytest tests",
		})
	}

	return files
}

// MapType maps a pseudo-code type to Python's equivalent.
func (a *PythonAdapter) MapType(pseudoType string) string {
	typeMap := map[string]string{
		"Text":      "str",
		"String":    "str",
		"Integer":   "int",
		"Int":       "int",
		"Number":    "int",
		"Float":     "float",
		"Decimal":   "Decimal",
		"Boolean":   "bool",
		"Bool":      "bool",
		"Timestamp": "datetime",
		"Duration":  "timedelta",
		"UUID":      "UUID",
		"Any":       "Any",
		"Nothing":   "None",
		"Void":      "None",
		"None":      "None",
	}

	if pyType, ok := typeMap[pseudoType]; ok {
		return pyType
	}
	return pseudoType
}
