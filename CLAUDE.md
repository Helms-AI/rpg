# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

RPG (Rosetta Project Generator) is a markdown-driven multi-language code generation tool exposed as an MCP server. It translates natural language specifications into idiomatic code across 6 languages: Go, Rust, Java, Python, TypeScript, and C#.

## Build & Development Commands

```bash
make build          # Build binary to bin/rpg
make test           # Run all tests with verbose output
make test-coverage  # Generate coverage report (HTML)
make run            # Build and run the MCP server
make vet            # Run go vet static analysis
make tidy           # Tidy go.mod dependencies
make clean          # Remove build artifacts
```

Run a single test:
```bash
go test -v -run TestFunctionName ./internal/parser/
```

## Architecture

### Core Data Flow
1. Client provides spec path + target language
2. Parser converts markdown spec to structured `Spec` object
3. Language adapter provides conventions + prompt template
4. Returns complete context for AI code generation

### Package Structure

- **cmd/rpg/** - CLI entry point with flag parsing and graceful shutdown
- **internal/server/** - MCP server core: tool registration (`server.go`) and handlers (`handlers.go`)
- **internal/parser/** - Goldmark-based markdown parser that extracts spec sections by heading level
- **internal/languages/** - Plugin architecture for language adapters implementing `LanguageAdapter` interface
- **internal/importer/** - Source code analysis for reverse-engineering specs from existing code
- **pkg/spec/** - Public spec types (Spec, TypeDef, Function, TestCase)

### MCP Tools Exposed

| Tool | Purpose |
|------|---------|
| `list_languages` | List supported languages with conventions |
| `parse_spec` | Parse markdown spec into structured data |
| `validate_spec` | Validate spec for errors/warnings |
| `get_generation_context` | Get parsed spec + language conventions + prompt template |
| `get_project_structure` | Get recommended file structure for target language |
| `ensure_parity` | Compare implementations across languages |
| `import_spec_from_source` | Analyze source code for AI-powered spec generation |

### Language Adapter Pattern

Each adapter in `internal/languages/` implements:
```go
type LanguageAdapter interface {
    GetLanguage() Language              // Metadata: naming conventions, error patterns
    GetPromptContext() string           // Language-specific generation prompt
    GetProjectStructure(name, hasTests bool) []ProjectFile
    MapType(pseudoType string) string   // Map spec types to language types
}
```

To add a new language: create adapter file, implement interface, register in `registry.go`.

### Spec Parser Sections

The parser extracts these markdown sections (case-insensitive):
- `# Title` (H1) - Spec name
- `## Meta` - Version, author, license
- `## Target Languages` - Language IDs
- `## Types` - Data structure definitions
- `## Functions` - Function signatures with accepts/returns/logic
- `## Tests` - Test cases with given/expect format
- `## Dependencies` - External dependencies
- `## Configuration` - Environment variables

### Type Normalization

The `Spec.Normalize()` method converts nil slices to empty slices for JSON schema compliance. Always call after parsing.

## Key Files

- `internal/server/handlers.go` - All MCP tool implementations
- `internal/parser/parser.go` - Markdown parsing with goldmark AST
- `internal/languages/registry.go` - Language adapter registration
- `pkg/spec/types.go` - Core data structures
- `examples/*.spec.md` - Example specification files

## Configuration

**.mcp.json** configures the MCP server:
```json
{
  "mcpServers": {
    "rpg": {
      "command": "./bin/rpg",
      "args": ["--output", "./generated"]
    }
  }
}
```

## Error Handling Patterns

- Wrap errors with context: `fmt.Errorf("parsing spec: %w", err)`
- MCP tool errors use `CallToolResult` with `IsError: true`
- Validation errors return structured `ValidationResult` with error details
