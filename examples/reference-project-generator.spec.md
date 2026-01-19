# RPG - Rosetta Project Generator

A markdown-driven multi-language code generation tool exposed as an MCP (Model Context Protocol) server. RPG translates natural language specifications into idiomatic code across 6 supported languages. Specs can be written in any narrative format - architecture docs, API designs, or feature descriptions - and the AI interprets the content directly to generate code.

## Target Languages

- go
- rust
- java
- python
- typescript
- csharp

## Meta

- **Version**: 1.0.0
- **License**: MIT
- **Author**: kon1790

## Configuration

| Variable | Type | Default | Required | Description |
|----------|------|---------|----------|-------------|
| --output | string | ./output | No | Base output directory for generated projects |
| -o | string | ./output | No | Short flag for output directory |

## Types

### Spec

Represents a specification file containing raw markdown content for AI interpretation.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| rawContent | string | Yes | Full markdown content of the specification |

---

### ValidationError

Represents an error or warning found during spec validation.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| severity | string | Yes | Error level: "error", "warning", or "info" |
| code | string | Yes | Error code identifier (e.g., "EMPTY_CONTENT") |
| message | string | Yes | Human-readable error description |

---

### ValidationResult

Contains the results of spec validation.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| valid | bool | Yes | Whether the spec passed validation |
| errors | []ValidationError | No | List of validation errors if any |

---

### Language

Represents a supported target language with its conventions and metadata.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | Yes | Language identifier (e.g., "go", "rust", "csharp") |
| name | string | Yes | Human-readable name (e.g., "Go", "Rust", "C#") |
| version | string | Yes | Recommended language version (e.g., "1.21+", "12+") |
| fileExtension | string | Yes | Primary file extension (e.g., ".go", ".rs", ".cs") |
| conventions | Conventions | Yes | Naming and style conventions |
| idioms | []string | Yes | Language-specific best practices |
| projectStructure | ProjectStructure | Yes | Typical project layout |
| errorPatterns | ErrorPatterns | Yes | Error handling patterns |
| dependencies | DependencyInfo | Yes | Package manager information |

---

### Conventions

Defines naming and style conventions for a language.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| naming | NamingConventions | Yes | Naming patterns for identifiers |
| errorHandling | string | Yes | Error handling approach description |
| fileNaming | string | Yes | File naming convention (e.g., "snake_case.go") |
| imports | string | Yes | Import organization style |
| docStyle | string | Yes | Documentation comment style |

---

### NamingConventions

Defines naming patterns for different identifier types.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| functions | string | Yes | Function naming convention |
| variables | string | Yes | Variable naming convention |
| constants | string | Yes | Constant naming convention |
| types | string | Yes | Type naming convention |
| packages | string | Yes | Package/namespace naming convention |
| private | string | No | Private member naming convention |

---

### ProjectStructure

Defines the typical project layout for a language.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| sourceDir | string | Yes | Source code directory (empty for root) |
| testDir | string | No | Test directory |
| testSuffix | string | Yes | Test file suffix (e.g., "_test.go", "Tests.cs") |
| packageFile | string | Yes | Package manifest file (e.g., "go.mod", ".csproj") |
| entryPoint | string | No | Main entry point file |
| commonDirs | []string | No | Common directory structure |

---

### ErrorPatterns

Defines how errors are handled in the language.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| style | string | Yes | Error style: "exceptions", "result", "tuple", "optional" |
| customError | string | No | Custom error type definition example |
| wrapError | string | No | Error wrapping pattern example |

---

### DependencyInfo

Provides package manager information.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| manager | string | Yes | Package manager name (e.g., "go mod", "NuGet", "cargo") |
| installCmd | string | Yes | Command to install dependencies |
| addCmd | string | Yes | Command to add a new dependency |
| lockFile | string | No | Lock file name if applicable |

---

### ProjectFile

Represents a file in a recommended project structure.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| path | string | Yes | Relative file path |
| purpose | string | Yes | File purpose: "source", "test", "config" |
| description | string | No | Human-readable description |

---

### LanguageAdapter

Interface that language adapters must implement.

| Method | Returns | Description |
|--------|---------|-------------|
| GetLanguage() | Language | Returns the language configuration |
| GetPromptContext() | string | Returns language-specific prompt instructions |
| GetProjectStructure(specName, hasTests) | []ProjectFile | Returns recommended project structure |
| MapType(pseudoType) | string | Maps pseudo-code types to language equivalents |

---

### FileCategory

Enumeration of file categories for collection.

| Value | Description |
|-------|-------------|
| source | Source code files |
| test | Test files |
| api_spec | API specification files (OpenAPI, AsyncAPI) |
| config | Configuration files |
| doc | Documentation files |

---

### FileContent

Holds the path and content of a collected file.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| path | string | Yes | Relative file path |
| category | FileCategory | Yes | Category of the file |
| content | string | Yes | File content |

---

### ProjectFiles

Contains all collected files from a project for analysis.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Project name |
| language | string | Yes | Detected primary language |
| rootPath | string | Yes | Root directory path |
| sourceFiles | []FileContent | Yes | Collected source files |
| testFiles | []FileContent | Yes | Collected test files |
| apiSpecs | []FileContent | Yes | Collected API specification files |
| configFiles | []FileContent | Yes | Collected configuration files |
| docFiles | []FileContent | Yes | Collected documentation files |

---

### Server

MCP server wrapper with dependencies.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| mcpServer | *mcp.Server | Yes | Underlying MCP server |
| registry | *Registry | Yes | Language adapter registry |
| outputDir | string | Yes | Base output directory |

---

### Registry

Holds all registered language adapters.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| adapters | map[string]LanguageAdapter | Yes | Map of language ID to adapter |

---

### FeatureStatus

Tracks a feature across all implementations for parity checking.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | string | Yes | Feature identifier |
| name | string | Yes | Feature name |
| category | string | Yes | Feature category |
| implementations | map[string]Implementation | Yes | Implementation status by language |

---

### Implementation

Details for a feature in a specific language.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| present | bool | Yes | Whether feature is implemented |
| filePath | string | No | File where feature is found |
| lineNumber | int | No | Line number of implementation |
| codeSnippet | string | No | Code snippet showing implementation |

---

### ParityGap

Represents a missing or different feature between implementations.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| featureId | string | Yes | Feature identifier |
| featureName | string | Yes | Feature name |
| category | string | Yes | Feature category |
| missingIn | string | Yes | Language missing the feature |
| referenceFile | string | Yes | Reference implementation file |
| referenceCode | string | Yes | Reference code snippet |
| suggestedFix | string | Yes | Fix suggestion |
| targetFile | string | Yes | Target file for implementation |

## Functions

### main

Application entry point that parses flags and starts the MCP server.

**Accepts**: Command-line flags

**Returns**: None (exits on error)

**Logic**:
```
1. Parse --output/-o flag with default "./output"
2. Create context with cancellation support
3. Set up signal handler for SIGINT and SIGTERM
4. Create new MCP server with output directory
5. Run server with stdio transport
6. If server returns error, log fatal and exit
```

---

### Server.New

Creates a new MCP server with all tools and resources registered.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| outputDir | string | Yes | Base output directory for generated projects |

**Returns**: *Server

**Logic**:
```
1. Create new language registry with all built-in adapters
2. Create MCP server with name "rpg" and version "1.0.0"
3. Set server instructions describing the tool's purpose
4. Configure logger to stderr with info level
5. Create Server struct with mcpServer, registry, and outputDir
6. Register all MCP tools (list_languages, parse_spec, validate_spec, etc.)
7. Register all MCP resources (examples and language conventions)
8. Return server instance
```

---

### Server.Run

Starts the MCP server with stdio transport.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| ctx | context.Context | Yes | Context for cancellation |

**Returns**: error

**Logic**:
```
1. Call mcpServer.Run with the provided context
2. Use StdioTransport for communication
3. Return any error from the server
```

---

### handleListLanguages

MCP tool handler that returns all supported languages.

**Accepts**: ListLanguagesInput (empty)

**Returns**: ListLanguagesOutput with Languages array

**Logic**:
```
1. Call registry.List() to get all registered languages
2. Return ListLanguagesOutput containing the language list
```

---

### handleParseSpec

MCP tool handler that reads and returns spec file content.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| specPath | string | Yes | Path to the markdown spec file |

**Returns**: ParseSpecOutput with Content

**Logic**:
```
1. Read file at specPath using os.ReadFile
2. If error, return CallToolResult with IsError=true and error message
3. Return ParseSpecOutput with file content as string
```

---

### handleValidateSpec

MCP tool handler that validates a spec file.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| specPath | string | Yes | Path to the markdown spec file |

**Returns**: ValidateSpecOutput with Valid and Errors

**Logic**:
```
1. Read file at specPath
2. If read error, return invalid result with FILE_READ_ERROR
3. Call parser.Validate with content
4. Return validation result with valid status and any errors
```

---

### handleGetGenerationContext

MCP tool handler that returns full context for code generation.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| specPath | string | Yes | Path to the markdown spec file |
| language | string | Yes | Target language ID |

**Returns**: GetGenerationContextOutput

**Logic**:
```
1. Get language adapter from registry by language ID
2. If unsupported language, return error with suggestion to use list_languages
3. Read spec file content
4. If read error, return error result
5. Build output path as outputDir/language
6. Return context with specContent, language config, prompt template, and outputDir
```

---

### handleGetProjectStructure

MCP tool handler that returns recommended file structure.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| projectName | string | Yes | Name for the project |
| language | string | Yes | Target language ID |

**Returns**: GetProjectStructureOutput

**Logic**:
```
1. Get language adapter from registry
2. If unsupported language, return error
3. Get project structure from adapter using project name
4. Build output path as outputDir/projectName/language
5. Return structure with files and outputDir
```

---

### handleImportSpecFromSource

MCP tool handler that analyzes source code and returns an analysis prompt for AI-powered spec generation.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| inputPath | string | Yes | Path to source directory |
| name | string | No | Optional name for the spec |

**Returns**: ImportSpecFromSourceOutput

**Logic**:
```
1. Expand ~ to home directory if present
2. Convert relative path to absolute
3. Validate input path exists
4. Call CollectProjectFiles to gather all source, test, API, config, and doc files
5. Override project name if provided
6. Check if any files were found, return error if not
7. Build analysis prompt using BuildAnalysisPrompt
8. Determine spec output path as specs/<name>.spec.md
9. Create specs directory if needed
10. Return output with project stats and analysis prompt
```

---

### handleEnsureParity

MCP tool handler that checks feature parity across implementations.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| specPath | string | Yes | Path to original spec |
| projects | []ProjectInfo | Yes | List of projects to compare (min 2) |

**Returns**: EnsureParityOutput

**Logic**:
```
1. Validate at least 2 projects provided
2. Read spec file for context
3. First project is reference implementation
4. For each project, extract features and files
5. Build feature matrix tracking each feature across languages
6. Identify gaps where reference has feature but other doesn't
7. Calculate parity score as (total - gaps) / total
8. Generate fix instructions for Claude
9. Return parity score, reference language, feature matrix, gaps, and instructions
```

---

### CollectProjectFiles

Scans a directory and collects all relevant files with their contents.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| dir | string | Yes | Directory to scan |

**Returns**: (*ProjectFiles, error)

**Logic**:
```
1. Create ProjectFiles struct with name from directory base
2. Initialize language count map
3. Walk directory recursively
4. Skip hidden directories and known skip dirs (node_modules, vendor, etc.)
5. Skip hidden files (except .env)
6. Skip files larger than 1MB
7. For each file, determine category using CategorizeFile
8. Read file content and add to appropriate slice
9. Track language by file extension
10. After walk, determine primary language by highest count
11. Return populated ProjectFiles
```

---

### CategorizeFile

Determines the category of a file based on its path and name.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| path | string | Yes | File path to categorize |

**Returns**: FileCategory (or empty string if not relevant)

**Logic**:
```
1. Check for API spec files (openapi.yaml, swagger.json, asyncapi.yaml)
2. Check for documentation files (README.md, CHANGELOG.md, docs/)
3. Check for config files (go.mod, package.json, .csproj, etc.)
4. Check for test files using patterns (_test.go, .test.ts, Test.java, etc.)
5. Check for source files by extension (.go, .ts, .py, .java, .rs, .cs)
6. Return empty string if no match
```

---

### BuildAnalysisPrompt

Creates a comprehensive prompt for AI analysis of a project.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| files | *ProjectFiles | Yes | Collected project files |

**Returns**: string (analysis prompt)

**Logic**:
```
1. Write header with project name, detected language, total files
2. For each category (source, test, API, config, doc), write section
3. Include full file content in code blocks
4. Append instructions for generating spec.md
5. Specify output format requirements
6. Return complete prompt string
```

---

### Registry.NewRegistry

Creates a new language registry with all built-in adapters.

**Accepts**: None

**Returns**: *Registry

**Logic**:
```
1. Create Registry with empty adapters map
2. Register Go adapter
3. Register Rust adapter
4. Register Java adapter
5. Register Python adapter
6. Register TypeScript adapter
7. Register C# adapter
8. Return registry
```

---

### Registry.Get

Returns a language adapter by ID.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| languageID | string | Yes | Language identifier (case-insensitive) |

**Returns**: (LanguageAdapter, error)

**Logic**:
```
1. Convert languageID to lowercase
2. Look up in adapters map
3. If not found, return error "unsupported language: <id>"
4. Return adapter and nil error
```

---

### Registry.List

Returns all registered languages.

**Accepts**: None

**Returns**: []Language

**Logic**:
```
1. Create slice with capacity of adapters count
2. For each adapter, call GetLanguage() and append
3. Return languages slice
```

---

### parser.Parse

Reads a spec file and returns it as raw content.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| content | string | Yes | Markdown content |

**Returns**: (*Spec, error)

**Logic**:
```
1. Create Spec with RawContent set to content
2. Return spec and nil error
```

---

### parser.Validate

Checks if content is valid (non-empty).

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| content | string | Yes | Content to validate |

**Returns**: ValidationResult

**Logic**:
```
1. If trimmed content is empty, return invalid with EMPTY_CONTENT error
2. Otherwise return valid result
```

---

### extractFeaturesFromCode

Extracts feature markers from source code for parity checking.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| content | string | Yes | Source code content |
| language | string | Yes | Language identifier |
| filePath | string | Yes | File path for reference |

**Returns**: map[string]Implementation

**Logic**:
```
1. Define feature patterns (url_validation, tag_normalization, async_metadata_fetch, etc.)
2. For each pattern, check if any marker exists in content
3. If found, extract line number and code snippet
4. Add to features map with implementation details
5. Return features map
```

## MCP Tools

### list_languages

List all supported target languages with their conventions and idioms.

**Input**: None

**Output**: List of Language objects with full configuration

---

### parse_spec

Read a markdown specification file and return its content.

**Input**: specPath (required) - Path to markdown spec file

**Output**: Raw markdown content

---

### validate_spec

Check if a spec file exists and contains content.

**Input**: specPath (required) - Path to markdown spec file

**Output**: Validation result with valid boolean and any errors

---

### get_generation_context

Get full context for code generation including spec content, language conventions, and prompt template.

**Input**:
- specPath (required) - Path to markdown spec file
- language (required) - Target language ID

**Output**: Spec content, language config, prompt template, output directory

---

### get_project_structure

Get recommended file structure for a project in the target language.

**Input**:
- projectName (required) - Name for the project
- language (required) - Target language ID

**Output**: List of ProjectFile with paths and purposes, output directory

---

### ensure_parity

Check feature parity across generated projects and provide fix instructions.

**Input**:
- specPath (required) - Path to original spec
- projects (required) - List of projects with language and path

**Output**: Parity score, feature matrix, gaps, fix instructions

---

### import_spec_from_source

Collect and analyze source code for AI-powered spec generation.

**Input**:
- inputPath (required) - Path to source directory
- name (optional) - Name for generated spec

**Output**: Project stats, detected language, analysis prompt, spec output path

## MCP Resources

### spec://examples/simple-function

Example spec for a simple slugify function in narrative style.

### spec://examples/module

Example spec for a validation utilities module in narrative style.

### spec://examples/full-project

Example spec for a REST API project with authentication.

### lang://{language}/conventions

Language-specific conventions and idioms for each supported language (go, rust, java, python, typescript, csharp).

## Dependencies

### Go Standard Library

| Package | Purpose |
|---------|---------|
| context | Context for cancellation |
| encoding/json | JSON marshaling |
| fmt | String formatting and error wrapping |
| log | Logging |
| log/slog | Structured logging |
| os | File operations and environment |
| os/signal | Signal handling |
| path/filepath | Path manipulation |
| strings | String operations |
| syscall | System signals |

### External Dependencies

| Package | Purpose |
|---------|---------|
| github.com/modelcontextprotocol/go-sdk/mcp | MCP server SDK |

## Supported Languages

RPG includes built-in adapters for the following languages:

| Language | ID | File Extension | Package Manager | Test Suffix |
|----------|----|--------------| ----------------|-------------|
| Go | go | .go | go mod | _test.go |
| Rust | rust | .rs | cargo | _test.rs |
| Java | java | .java | Maven/Gradle | Test.java |
| Python | python | .py | pip | test_.py / _test.py |
| TypeScript | typescript | .ts | npm | .test.ts / .spec.ts |
| C# | csharp | .cs | NuGet | Tests.cs |
