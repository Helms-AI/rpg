package languages

import "fmt"

// CSharpAdapter implements LanguageAdapter for C#.
type CSharpAdapter struct {
	language Language
}

// NewCSharpAdapter creates a new C# language adapter.
func NewCSharpAdapter() *CSharpAdapter {
	return &CSharpAdapter{
		language: Language{
			ID:            "csharp",
			Name:          "C#",
			Version:       "12+",
			FileExtension: ".cs",
			Conventions: Conventions{
				Naming: NamingConventions{
					Functions:  "PascalCase (methods)",
					Variables:  "camelCase (locals), _camelCase (private fields)",
					Constants:  "PascalCase",
					Types:      "PascalCase (classes, structs, records)",
					Packages:   "PascalCase (namespaces: Company.Project.Module)",
					Private:    "_camelCase prefix for fields",
				},
				ErrorHandling: "Exceptions with try/catch/finally, ArgumentException for validation",
				FileNaming:    "PascalCase.cs",
				Imports:       "using statements at top, group by: System, third-party, project",
				DocStyle:      "XML documentation comments (/// <summary>)",
			},
			Idioms: []string{
				"Use primary constructors (C# 12)",
				"Use records for DTOs and value objects",
				"Use init-only setters for immutability",
				"Use pattern matching in switch expressions",
				"Use async/await for all I/O operations",
				"Use nullable reference types (enable in project)",
				"Use LINQ for collection operations",
				"Use dependency injection",
				"Use IOptions<T> for configuration",
				"Use structured logging (Serilog, Microsoft.Extensions.Logging)",
			},
			ProjectStructure: ProjectStructure{
				SourceDir:   "",
				TestDir:     "Tests/",
				TestSuffix:  "Tests.cs",
				PackageFile: ".csproj",
				EntryPoint:  "Program.cs",
				CommonDirs:  []string{"src/", "tests/"},
			},
			ErrorPatterns: ErrorPatterns{
				Style:       "exceptions",
				CustomError: "public class MyException : Exception\n{\n    public MyException(string message) : base(message) { }\n}",
				WrapError:   "throw new MyException(\"context\", innerException);",
			},
			Dependencies: DependencyInfo{
				Manager:    "NuGet",
				InstallCmd: "dotnet restore",
				AddCmd:     "dotnet add package",
				LockFile:   "packages.lock.json",
			},
		},
	}
}

// GetLanguage returns the C# language configuration.
func (a *CSharpAdapter) GetLanguage() Language {
	return a.language
}

// GetPromptContext returns C#-specific prompt instructions.
func (a *CSharpAdapter) GetPromptContext() string {
	return `Generate idiomatic C# code following these conventions:

## Language: C# 12+

## Naming Conventions
- Classes/records/structs: PascalCase
- Methods/properties: PascalCase
- Private fields: _camelCase
- Local variables: camelCase
- Interfaces: IPascalCase (I prefix)
- Constants: PascalCase

## Modern C# Features (12+)
- Use primary constructors: public class Service(IRepository repo)
- Use records for data types: public record User(string Id, string Name);
- Use collection expressions: int[] arr = [1, 2, 3];
- Use required members for validation
- Use raw string literals for multi-line text

## Async/Await
- Use async/await for all I/O operations
- Return Task or Task<T>, not void
- Use CancellationToken for cancellation
- Use ValueTask for hot paths

## Nullable Reference Types
- Enable nullable reference types
- Use ? for nullable, no suffix for non-nullable
- Use ArgumentNullException.ThrowIfNull for validation

## Code Style
- Use LINQ for collection operations
- Use pattern matching in switch expressions
- Use dependency injection
- Use init-only setters for immutability
- Use IOptions<T> for configuration

## Documentation
- Use XML documentation (/// <summary>)
- Include <param>, <returns>, <exception> tags

## Output Requirements
- Generate complete, compilable C# code
- Include using statements
- Use nullable reference types
- Add XML documentation for public APIs
- Use modern C# 12+ features`
}

// GetProjectStructure returns the recommended C# project structure.
func (a *CSharpAdapter) GetProjectStructure(specName string, hasTests bool) []ProjectFile {
	files := []ProjectFile{
		{
			Path:        fmt.Sprintf("%s.csproj", specName),
			Purpose:     "config",
			Description: ".NET project file",
		},
		{
			Path:        fmt.Sprintf("%s.cs", specName),
			Purpose:     "source",
			Description: "Main implementation",
		},
	}

	if hasTests {
		files = append(files, ProjectFile{
			Path:        fmt.Sprintf("%s.Tests/%s.Tests.csproj", specName, specName),
			Purpose:     "config",
			Description: "Test project file",
		})
		files = append(files, ProjectFile{
			Path:        fmt.Sprintf("%s.Tests/%sTests.cs", specName, specName),
			Purpose:     "test",
			Description: "xUnit tests",
		})
	}

	return files
}

// MapType maps a pseudo-code type to C#'s equivalent.
func (a *CSharpAdapter) MapType(pseudoType string) string {
	typeMap := map[string]string{
		"Text":      "string",
		"String":    "string",
		"Integer":   "int",
		"Int":       "int",
		"Number":    "int",
		"Float":     "double",
		"Decimal":   "decimal",
		"Boolean":   "bool",
		"Bool":      "bool",
		"Timestamp": "DateTimeOffset",
		"Duration":  "TimeSpan",
		"UUID":      "Guid",
		"Any":       "object",
		"Nothing":   "void",
		"Void":      "void",
		"None":      "void",
	}

	if csType, ok := typeMap[pseudoType]; ok {
		return csType
	}
	return pseudoType
}
