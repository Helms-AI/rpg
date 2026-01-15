package languages

import "fmt"

// JavaAdapter implements LanguageAdapter for Java.
type JavaAdapter struct {
	language Language
}

// NewJavaAdapter creates a new Java language adapter.
func NewJavaAdapter() *JavaAdapter {
	return &JavaAdapter{
		language: Language{
			ID:            "java",
			Name:          "Java",
			Version:       "17+",
			FileExtension: ".java",
			Conventions: Conventions{
				Naming: NamingConventions{
					Functions:  "camelCase",
					Variables:  "camelCase",
					Constants:  "SCREAMING_SNAKE_CASE",
					Types:      "PascalCase (classes, interfaces, enums)",
					Packages:   "all lowercase, reverse domain (com.company.project)",
					Private:    "camelCase with no prefix",
				},
				ErrorHandling: "Checked exceptions for recoverable errors, runtime exceptions for bugs",
				FileNaming:    "PascalCase.java (must match public class name)",
				Imports:       "Group by: java.*, javax.*, third-party, project packages",
				DocStyle:      "Javadoc (/** ... */)",
			},
			Idioms: []string{
				"Use Optional<T> for nullable returns",
				"Prefer immutable objects (final fields, no setters)",
				"Use records (Java 17+) for data classes",
				"Use streams for collection processing",
				"Use try-with-resources for AutoCloseable",
				"Prefer interface types in declarations",
				"Use dependency injection",
				"Use builders for complex object construction",
				"Use sealed interfaces for restricted hierarchies",
				"Document with Javadoc",
			},
			ProjectStructure: ProjectStructure{
				SourceDir:   "src/main/java/",
				TestDir:     "src/test/java/",
				TestSuffix:  "Test.java",
				PackageFile: "pom.xml or build.gradle",
				EntryPoint:  "Main.java",
				CommonDirs:  []string{"src/main/java/", "src/test/java/", "src/main/resources/"},
			},
			ErrorPatterns: ErrorPatterns{
				Style:       "exceptions",
				CustomError: "public class MyException extends Exception { ... }",
				WrapError:   "throw new MyException(\"context\", cause);",
			},
			Dependencies: DependencyInfo{
				Manager:    "Maven or Gradle",
				InstallCmd: "mvn install or gradle build",
				AddCmd:     "Add to pom.xml or build.gradle",
				LockFile:   "",
			},
		},
	}
}

// GetLanguage returns the Java language configuration.
func (a *JavaAdapter) GetLanguage() Language {
	return a.language
}

// GetPromptContext returns Java-specific prompt instructions.
func (a *JavaAdapter) GetPromptContext() string {
	return `Generate idiomatic Java code following these conventions:

## Language: Java 17+

## Naming Conventions
- Classes/interfaces/enums: PascalCase
- Methods/variables: camelCase
- Constants: SCREAMING_SNAKE_CASE
- Packages: all lowercase (com.company.project)

## Modern Java Features (17+)
- Use records for data classes: public record User(String id, String name) {}
- Use sealed interfaces for type hierarchies
- Use pattern matching in instanceof
- Use switch expressions with pattern matching
- Use text blocks for multi-line strings

## Error Handling
- Checked exceptions for recoverable errors client should handle
- Runtime exceptions for programming errors
- Use try-with-resources for AutoCloseable
- Include cause in exception chaining

## Code Style
- Use Optional<T> for nullable returns
- Prefer immutability (final fields)
- Use streams for collection processing
- Use builders for complex object construction
- Prefer interface types in declarations (List over ArrayList)
- Use dependency injection

## Documentation
- Use Javadoc (/** ... */) for all public APIs
- Include @param, @return, @throws tags

## Output Requirements
- Generate complete, compilable Java code
- Include package declaration
- Include all imports
- Add Javadoc for public classes and methods
- Use modern Java 17+ features where appropriate`
}

// GetProjectStructure returns the recommended Java project structure.
func (a *JavaAdapter) GetProjectStructure(specName string, hasTests bool) []ProjectFile {
	files := []ProjectFile{
		{
			Path:        "pom.xml",
			Purpose:     "config",
			Description: "Maven project configuration",
		},
		{
			Path:        fmt.Sprintf("src/main/java/com/example/%s.java", specName),
			Purpose:     "source",
			Description: "Main implementation",
		},
	}

	if hasTests {
		files = append(files, ProjectFile{
			Path:        fmt.Sprintf("src/test/java/com/example/%sTest.java", specName),
			Purpose:     "test",
			Description: "JUnit tests",
		})
	}

	return files
}

// MapType maps a pseudo-code type to Java's equivalent.
func (a *JavaAdapter) MapType(pseudoType string) string {
	typeMap := map[string]string{
		"Text":      "String",
		"String":    "String",
		"Integer":   "int",
		"Int":       "int",
		"Number":    "int",
		"Float":     "double",
		"Decimal":   "BigDecimal",
		"Boolean":   "boolean",
		"Bool":      "boolean",
		"Timestamp": "Instant",
		"Duration":  "Duration",
		"UUID":      "UUID",
		"Any":       "Object",
		"Nothing":   "void",
		"Void":      "void",
		"None":      "void",
	}

	if javaType, ok := typeMap[pseudoType]; ok {
		return javaType
	}
	return pseudoType
}
