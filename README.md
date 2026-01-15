# rpg - Reference Project Generator

A markdown-driven multi-language code generation tool exposed as an MCP (Model Context Protocol) server. Write specifications in natural language markdown, and AI agents translate them into idiomatic code for any supported language.

## Features

- **Multi-Language Support**: Generate code in Go, Rust, Java, C#, Python, and TypeScript
- **Language-Specific Conventions**: Each language adapter enforces idiomatic patterns, naming conventions, and project structure
- **Flexible Spec Format**: Write specs as detailed pseudo-code or high-level explanations
- **MCP Integration**: Works with Claude Code, Cursor, and other MCP-compatible tools
- **Parity Checking**: Automatically compare implementations across languages and generate fix instructions
- **Configurable Output**: Control where generated projects are saved

## Installation

### Prerequisites

- Go 1.21 or later
- An MCP-compatible AI tool (Claude Code, Cursor, etc.)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/kon1790/rpg.git
cd rpg

# Build
make build
# or
go build -o bin/rpg ./cmd/rpg

# Verify
./bin/rpg --help
```

## Configuration

### Claude Code Setup

Add rpg to your project's `.mcp.json`:

```json
{
  "mcpServers": {
    "rpg": {
      "command": "/absolute/path/to/rpg"
    }
  }
}
```

Or with custom output directory:

```json
{
  "mcpServers": {
    "rpg": {
      "command": "/absolute/path/to/rpg",
      "args": ["--output", "./generated"]
    }
  }
}
```

### CLI Options

```
rpg [options]

Options:
  -o, --output <dir>    Base output directory for generated projects (default: ./output)
                        Projects are saved to: <output>/<project-name>/<language>/
```

## MCP Tools

| Tool | Description |
|------|-------------|
| `list_languages` | List all supported languages with their conventions and idioms |
| `parse_spec` | Parse a markdown spec file and return structured data |
| `validate_spec` | Validate a spec file and return any errors or warnings |
| `get_generation_context` | Get full context for code generation including spec, language conventions, and prompt template |
| `get_project_structure` | Get recommended file structure for a spec in a target language |
| `ensure_parity` | Compare implementations across languages and generate fix instructions for gaps |

### Tool Usage Examples

**Get generation context for a spec:**
```
Use rpg to get the generation context for examples/fullstack-app.spec.md in csharp
```

**Generate in multiple languages:**
```
Generate the slugify spec in all available languages
```

**Check parity across implementations:**
```
Use rpg ensure_parity to compare the csharp and java implementations in output/bookmark-manager/
```

## MCP Resources

| Resource URI | Description |
|--------------|-------------|
| `spec://examples/simple-function` | Example: simple utility function |
| `spec://examples/module` | Example: module with types and validation |
| `spec://examples/full-project` | Example: complete REST API project |
| `lang://go/conventions` | Go language conventions |
| `lang://rust/conventions` | Rust language conventions |
| `lang://java/conventions` | Java language conventions |
| `lang://csharp/conventions` | C# language conventions |
| `lang://python/conventions` | Python language conventions |
| `lang://typescript/conventions` | TypeScript language conventions |

## Writing Spec Files

rpg supports two spec styles: **explicit** (structured pseudo-code) and **flexible** (natural language descriptions).

### Explicit Style

Best for utility functions, libraries, and well-defined APIs:

```markdown
# slugify

A utility function to convert text into URL-friendly slugs.

## Target Languages

- go
- python
- typescript

## Functions

### slugify

Converts a string into a URL-friendly slug.

**accepts:**
- text: Text
- separator: Text (defaults to "-")

**returns:** Text

**logic:**
```
convert text to lowercase
replace all whitespace with the separator
remove all characters that are not letters, numbers, or the separator
collapse multiple consecutive separators into one
trim separators from start and end
return the result
```

## Tests

### slugify

#### test: converts simple text
given: "Hello World"
expect: "hello-world"

#### test: handles special characters
given: "Hello, World! How are you?"
expect: "hello-world-how-are-you"
```

### Flexible Style

Best for fullstack applications and complex systems:

```markdown
# bookmark-manager

A personal bookmark manager with tagging, full-text search, and metadata extraction.

## Overview

This is a self-hosted web application for saving bookmarks. When a user saves a URL,
the system automatically fetches the page title, description, and favicon.

## Target Languages

- csharp
- java
- typescript

## Architecture

Traditional fullstack web application with REST API backend and SPA frontend.

### Backend
- REST API for bookmark CRUD
- Background metadata fetching
- SQLite with FTS5 for search

### Frontend
- Simple SPA, no build step
- Searchable bookmark list
- Tag filtering sidebar

## Data Model

### Bookmark
- `id` - unique identifier
- `url` - the saved URL (required, valid URL)
- `title` - page title (auto-fetched)
- `tags` - list of associated tags
- `created_at` - timestamp

## API Design

```
GET    /api/bookmarks          - list bookmarks (?q=search&tag=filter)
POST   /api/bookmarks          - create bookmark
DELETE /api/bookmarks/:id      - delete bookmark
```

## Key Behaviors

### Adding a Bookmark
1. Immediately create bookmark with URL
2. Return success (optimistic UI)
3. Background fetch metadata
4. Update bookmark with title, description

## Configuration

Environment variables:
- `PORT` - HTTP port (default: 3000)
- `DATABASE_PATH` - SQLite location (default: ./data/bookmarks.db)
```

### Spec Sections Reference

| Section | Required | Description |
|---------|----------|-------------|
| `# Name` | Yes | Project/function name (H1 heading) |
| `## Target Languages` | Yes | List of target language IDs |
| `## Overview` | No | High-level description |
| `## Architecture` | No | System design and components |
| `## Types` | No | Data structures and enums |
| `## Functions` | No | Function signatures and logic |
| `## Data Model` | No | Database entities |
| `## API Design` | No | REST/GraphQL endpoints |
| `## Tests` | No | Test cases with given/expect |
| `## Configuration` | No | Environment variables and settings |

### Pseudo-Code Types

| Type | Maps To |
|------|---------|
| `Text`, `String` | string types |
| `Integer`, `Int`, `Number` | integer types |
| `Float`, `Decimal` | floating point |
| `Boolean`, `Bool` | boolean |
| `Timestamp` | datetime types |
| `UUID` | string (UUID format) |
| `List of X` | array/slice of X |
| `Optional X` | nullable X |
| `Result of X` | result type with error |

## Supported Languages

| Language | Version | Key Conventions |
|----------|---------|-----------------|
| **Go** | 1.21+ | `PascalCase` exports, `(result, error)` returns, `go.mod` |
| **Rust** | 2021 | `snake_case`, `Result<T, E>`, ownership, `Cargo.toml` |
| **Java** | 17+ | `PascalCase` classes, records, `Optional<T>`, Maven/Gradle |
| **C#** | 12+ | `PascalCase`, primary constructors, `async/await`, NuGet |
| **Python** | 3.11+ | `snake_case`, type hints, dataclasses, `pyproject.toml` |
| **TypeScript** | 5.0+ | `camelCase`, strict mode, discriminated unions, npm |

## Complete Workflow

### 1. Create a Spec File

```bash
# Create spec
cat > examples/my-api.spec.md << 'EOF'
# my-api

A simple REST API for managing items.

## Target Languages
- go
- typescript

## Data Model
### Item
- id: UUID
- name: Text
- created_at: Timestamp

## API Design
```
GET  /items     - list all items
POST /items     - create item
```
EOF
```

### 2. Generate Code

In Claude Code:
```
Use rpg to get the generation context for examples/my-api.spec.md in go,
then generate the complete implementation.
```

### 3. Generate in Additional Languages

```
Now generate the same spec in typescript.
```

### 4. Check Parity

```
Use rpg ensure_parity to compare the go and typescript implementations
and fix any missing features.
```

## Parity Checking

The `ensure_parity` tool compares implementations across languages:

**Input:**
```json
{
  "specPath": "examples/my-app.spec.md",
  "projects": [
    {"language": "csharp", "path": "output/my-app/csharp"},
    {"language": "java", "path": "output/my-app/java"}
  ]
}
```

**Output includes:**
- `parityScore` - 0.0 to 1.0 (1.0 = perfect parity)
- `referenceLanguage` - First project is the reference
- `featureMatrix` - All detected features with per-language implementation status
- `gaps` - Missing features with reference code snippets
- `fixInstructions` - Detailed markdown instructions for fixing gaps

**Detected Features:**
- URL validation
- Tag normalization
- Async operations
- Metadata extraction (og:title, og:description, favicon)
- Configuration (port, database, timeouts)
- Error handling and logging
- Search functionality

## Output Directory Structure

Generated projects follow this structure:

```
output/
├── slugify/
│   ├── go/
│   │   ├── go.mod
│   │   ├── slugify.go
│   │   └── slugify_test.go
│   ├── python/
│   │   ├── pyproject.toml
│   │   └── slugify.py
│   └── typescript/
│       ├── package.json
│       └── src/slugify.ts
└── bookmark-manager/
    ├── csharp/
    │   ├── BookmarkManager.csproj
    │   ├── Program.cs
    │   └── ...
    └── java/
        ├── pom.xml
        ├── src/main/java/...
        └── ...
```

## Project Structure

```
rpg/
├── cmd/rpg/              # CLI entry point
├── internal/
│   ├── server/           # MCP server implementation
│   │   ├── server.go     # Server setup and registration
│   │   └── handlers.go   # Tool handlers
│   ├── parser/           # Spec file parser
│   └── languages/        # Language adapters
│       ├── go.go
│       ├── rust.go
│       ├── java.go
│       ├── csharp.go
│       ├── python.go
│       └── typescript.go
├── pkg/spec/             # Public types
├── examples/             # Example spec files
├── output/               # Generated projects (default)
├── Makefile
└── .mcp.json             # MCP configuration
```

## Troubleshooting

### MCP Server Not Responding

1. Rebuild the binary: `make build`
2. Restart Claude Code to reload MCP servers
3. Check the path in `.mcp.json` is absolute

### Validation Errors

Run validation to see specific issues:
```
Use rpg validate_spec on examples/my-spec.md
```

### Missing Language

Check supported languages:
```
Use rpg list_languages
```

## License

MIT
