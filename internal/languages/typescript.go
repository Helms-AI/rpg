package languages

import "fmt"

// TypeScriptAdapter implements LanguageAdapter for TypeScript.
type TypeScriptAdapter struct {
	language Language
}

// NewTypeScriptAdapter creates a new TypeScript language adapter.
func NewTypeScriptAdapter() *TypeScriptAdapter {
	return &TypeScriptAdapter{
		language: Language{
			ID:            "typescript",
			Name:          "TypeScript",
			Version:       "5.0+",
			FileExtension: ".ts",
			Conventions: Conventions{
				Naming: NamingConventions{
					Functions:  "camelCase",
					Variables:  "camelCase",
					Constants:  "SCREAMING_SNAKE_CASE or PascalCase",
					Types:      "PascalCase (interfaces, types, classes)",
					Packages:   "kebab-case (npm packages)",
					Private:    "camelCase with # prefix (private class fields) or _prefix",
				},
				ErrorHandling: "throw Error, try/catch, or Result pattern with discriminated unions",
				FileNaming:    "kebab-case.ts or camelCase.ts",
				Imports:       "ESM imports, group by: node built-ins, external, internal",
				DocStyle:      "TSDoc (/** ... */)",
			},
			Idioms: []string{
				"Enable strict mode in tsconfig",
				"Use interfaces for object shapes, types for unions/aliases",
				"Use discriminated unions for type-safe variants",
				"Prefer readonly for immutability",
				"Use const assertions for literal types",
				"Use unknown over any, narrow types properly",
				"Use optional chaining (?.) and nullish coalescing (??)",
				"Use generics for reusable, type-safe code",
				"Avoid type assertions (as), prefer type guards",
				"Use ESM (import/export), not CommonJS",
			},
			ProjectStructure: ProjectStructure{
				SourceDir:   "src/",
				TestDir:     "tests/ or __tests__/",
				TestSuffix:  ".test.ts or .spec.ts",
				PackageFile: "package.json",
				EntryPoint:  "src/index.ts",
				CommonDirs:  []string{"src/", "tests/", "dist/"},
			},
			ErrorPatterns: ErrorPatterns{
				Style:       "exceptions",
				CustomError: "class MyError extends Error {\n  constructor(message: string) {\n    super(message);\n    this.name = 'MyError';\n  }\n}",
				WrapError:   "throw new MyError(`context: ${original.message}`);",
			},
			Dependencies: DependencyInfo{
				Manager:    "npm/yarn/pnpm/bun",
				InstallCmd: "npm install",
				AddCmd:     "npm install <package>",
				LockFile:   "package-lock.json or yarn.lock",
			},
		},
	}
}

// GetLanguage returns the TypeScript language configuration.
func (a *TypeScriptAdapter) GetLanguage() Language {
	return a.language
}

// GetPromptContext returns TypeScript-specific prompt instructions.
func (a *TypeScriptAdapter) GetPromptContext() string {
	return `Generate idiomatic TypeScript code following these conventions:

## Language: TypeScript 5.0+

## Naming Conventions
- Functions/variables: camelCase
- Interfaces/types/classes: PascalCase
- Constants: SCREAMING_SNAKE_CASE
- Private class fields: # prefix or _prefix

## Type System
- Enable strict mode
- Use interfaces for object shapes
- Use type aliases for unions, intersections, mapped types
- Use discriminated unions for type-safe variants
- Prefer unknown over any, narrow with type guards
- Use generics for reusable code
- Use const assertions: as const

## Code Style
- Use ESM imports/exports
- Use optional chaining (?.) and nullish coalescing (??)
- Prefer readonly for immutability
- Use satisfies operator for type checking
- Avoid type assertions (as), prefer type guards
- Use async/await for promises

## Error Handling
- Use Error subclasses for custom errors
- Consider Result pattern with discriminated unions
- Handle promise rejections properly

## Documentation
- Use TSDoc comments (/** ... */)
- Include @param, @returns, @throws tags
- Add @example for usage examples

## Output Requirements
- Generate complete, compilable TypeScript code
- Use strict TypeScript features
- Include all imports
- Add TSDoc comments for exports
- Use modern ESM syntax`
}

// GetProjectStructure returns the recommended TypeScript project structure.
func (a *TypeScriptAdapter) GetProjectStructure(specName string, hasTests bool) []ProjectFile {
	files := []ProjectFile{
		{
			Path:        "package.json",
			Purpose:     "config",
			Description: "npm package configuration",
		},
		{
			Path:        "tsconfig.json",
			Purpose:     "config",
			Description: "TypeScript compiler configuration",
		},
		{
			Path:        fmt.Sprintf("src/%s.ts", specName),
			Purpose:     "source",
			Description: "Main implementation",
		},
		{
			Path:        "src/index.ts",
			Purpose:     "source",
			Description: "Public exports",
		},
	}

	if hasTests {
		files = append(files, ProjectFile{
			Path:        fmt.Sprintf("tests/%s.test.ts", specName),
			Purpose:     "test",
			Description: "Jest/Vitest tests",
		})
	}

	return files
}

// MapType maps a pseudo-code type to TypeScript's equivalent.
func (a *TypeScriptAdapter) MapType(pseudoType string) string {
	typeMap := map[string]string{
		"Text":      "string",
		"String":    "string",
		"Integer":   "number",
		"Int":       "number",
		"Number":    "number",
		"Float":     "number",
		"Decimal":   "number",
		"Boolean":   "boolean",
		"Bool":      "boolean",
		"Timestamp": "Date",
		"Duration":  "number",
		"UUID":      "string",
		"Any":       "unknown",
		"Nothing":   "void",
		"Void":      "void",
		"None":      "undefined",
	}

	if tsType, ok := typeMap[pseudoType]; ok {
		return tsType
	}
	return pseudoType
}
