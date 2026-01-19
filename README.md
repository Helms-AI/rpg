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

RPG exposes 13 MCP tools organized into three categories:

### Core Generation

| Tool | Description |
|------|-------------|
| `list_languages` | List all supported languages with their conventions and idioms |
| `parse_spec` | Read markdown specification file content |
| `get_generation_context` | Get spec + language conventions + prompt template for code generation |
| `get_project_structure` | Get recommended file structure for a project in the target language |
| `generate_source_from_spec` | Autonomous code generation with automatic parity validation |

### Import & Analysis

| Tool | Description |
|------|-------------|
| `import_spec_from_source` | Analyze local source code for AI-powered spec generation |
| `import_spec_from_github` | Clone and analyze a GitHub repository for spec generation |
| `deep_analyze_source` | AST-based semantic analysis (types, functions, call graphs) |
| `list_project_languages` | Detect all programming languages in a project |
| `get_files_for_language` | Get raw file contents for AI-driven analysis |

### Parity & Refinement

| Tool | Description |
|------|-------------|
| `ensure_parity` | Compare implementations across languages with fix instructions |
| `semantic_parity_analysis` | Deep semantic comparison using AST-based analysis |
| `iterative_refinement_loop` | Automated refinement until parity threshold is reached |

### Tool Usage Examples

**Import from GitHub and port to new languages:**
```
Use rpg import_spec_from_github on owner/repo to analyze the codebase
Generate a spec from the analysis, then port to TypeScript and Go
```

**Autonomous code generation:**
```
Use rpg generate_source_from_spec to create a Rust implementation from my-api.spec.md
The tool will automatically validate parity and iterate until complete
```

**Deep analysis for complex codebases:**
```
Use rpg list_project_languages on ./my-project to see what languages are used
Use rpg deep_analyze_source for semantic analysis of the Go code
Use rpg get_files_for_language for SQL and Protobuf files
```

**Iterative refinement for production quality:**
```
Use rpg iterative_refinement_loop to port my Go project to TypeScript
Keep refining until 95% parity is achieved
```

**Basic generation workflow:**
```
Use rpg get_generation_context for examples/my-api.spec.md in typescript
Generate the slugify spec in go, rust, and python
Use rpg ensure_parity to compare implementations
```

## Importing Specs from Existing Code

Have existing code? Let RPG reverse-engineer it into a spec for multi-language generation.

### The Import Workflow

```
Your Code           RPG Analysis          AI Generation        Multi-Language
    │                    │                     │                     │
    ├── src/*       ────►│                     │                     │
    ├── tests/*     ────►│ import_spec_from   │ Creates natural    │ Go
    ├── config      ────►│ _source/_github ──►│ language spec  ───►│ Rust
    └── README      ────►│                     │                     │ TypeScript
                         └── Analysis Prompt ──┘                     └─ Python...
```

### Quick Start

**From GitHub repository:**
```bash
# Import directly from GitHub (supports owner/repo shorthand)
Use rpg import_spec_from_github on owner/repo

# With specific branch
Use rpg import_spec_from_github on owner/repo@main

# Full URL also works
Use rpg import_spec_from_github on https://github.com/owner/repo
```

**From local directory:**
```bash
# Import from local path
Use rpg import_spec_from_source on ./legacy-api
```

**After importing:**
```bash
# AI generates spec (review and refine as needed)
# Then generate in new languages
Use rpg get_generation_context for ./specs/my-project.spec.md in go
Use rpg get_generation_context for ./specs/my-project.spec.md in rust

# Verify parity
Use rpg ensure_parity to compare implementations
```

### What Gets Analyzed

| File Type | What It Reveals |
|-----------|-----------------|
| Source files | Functions, types, classes, business logic |
| Test files | Expected behavior, edge cases, usage examples |
| Config files | Dependencies, environment variables |
| Documentation | README, comments provide context |
| API specs | OpenAPI/GraphQL schemas if present |

### import_spec_from_github Tool

| Parameter | Required | Description |
|-----------|----------|-------------|
| `repository` | Yes | GitHub URL or shorthand (`owner/repo`, `owner/repo@branch`) |
| `ref` | No | Branch, tag, or commit SHA (overrides ref in URL) |
| `token` | No | GitHub PAT for private repos (or use `GITHUB_TOKEN` env var) |
| `name` | No | Optional name for the spec |
| `shallow` | No | Use shallow clone (default: true) |

### import_spec_from_source Tool

| Parameter | Required | Description |
|-----------|----------|-------------|
| `inputPath` | Yes | Path to source code directory |
| `name` | No | Optional name for the spec |

**Returns**: Analysis prompt for AI-powered spec generation

> **Note**: This is AI-assisted, not automatic. Review generated specs for accuracy and enhance with context the code doesn't capture.

## Writing Specs from Scratch

RPG specs are **natural language markdown**. Describe what you want—the AI interprets your intent.

### No Format Rules

The AI adapts to your style. All of these work:

**Minimal:**
```markdown
# url-shortener
A service that shortens URLs and redirects them.
Store mappings in memory.
```

**Structured:**
```markdown
# url-shortener

## Functions
### shorten(url: string): string
Takes a long URL and returns a 6-character code.

### resolve(code: string): string
Returns the original URL for a code, or throws if not found.
```

**Conversational:**
```markdown
# url-shortener
I need a URL shortening service like bit.ly. When someone gives us a long URL,
we generate a short code. Later they can use that code to get redirected.
Keep it simple - just store everything in memory for now.
```

### Quick Start

```bash
# 1. Write your spec (any style)
cat > my-api.spec.md << 'EOF'
# my-api
A REST API for managing tasks. Users can create, list, and complete tasks.
Use SQLite for storage.
EOF

# 2. Generate in any language
Use rpg get_generation_context for my-api.spec.md in typescript
```

### Common Patterns (All Optional)

| Section | Purpose |
|---------|---------|
| `## Types` | Define data structures |
| `## Functions` | Describe behavior, inputs, outputs |
| `## API Endpoints` | REST routes or GraphQL operations |
| `## Configuration` | Environment variables and defaults |
| `## Tests` | Given/expect scenarios |

### The Key: Describe Intent, Not Implementation

| ❌ Too Specific | ✅ Intent-Focused |
|-----------------|-------------------|
| "Use a HashMap<String, User>" | "Store users by ID for quick lookup" |
| "Iterate with a for loop and filter" | "Filter users matching criteria" |
| "Return an ArrayList" | "Return matching users as a list" |

The AI chooses idiomatic implementations for each target language.

### Portable Types

Use these generic types—the AI maps them to language-specific types:

| Type | Maps To |
|------|---------|
| `string`, `text` | String types |
| `int`, `integer` | Integer types |
| `float`, `decimal` | Floating point |
| `bool`, `boolean` | Boolean |
| `timestamp`, `datetime` | Date/time types |
| `uuid` | UUID/GUID |
| `list of X` | Array/slice of X |
| `optional X` | Nullable X |

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

### Missing Language

Check supported languages:
```
Use rpg list_languages
```

## License

MIT
