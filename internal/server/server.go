// Package server provides the MCP server implementation.
package server

import (
	"context"
	"log/slog"
	"os"

	"github.com/kon1790/rpg/internal/languages"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP server with our dependencies.
type Server struct {
	mcpServer *mcp.Server
	registry  *languages.Registry
	outputDir string
}

// New creates a new MCP server with all tools and resources registered.
// outputDir specifies the base directory for generated projects (e.g., "./output").
// Generated projects will be placed in outputDir/<project-name>/<language>/.
func New(outputDir string) *Server {
	// Create language registry
	registry := languages.NewRegistry()

	// Create MCP server
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "rpg",
			Version: "1.0.0",
		},
		&mcp.ServerOptions{
			Instructions: "A markdown-driven multi-language code generation tool. " +
				"Parse spec files and get generation context with language-specific conventions and prompts. " +
				"After generating code in multiple languages, use ensure_parity to verify feature consistency " +
				"across implementations. The first project passed becomes the reference implementation.",
			Logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})),
		},
	)

	s := &Server{
		mcpServer: mcpServer,
		registry:  registry,
		outputDir: outputDir,
	}

	// Register all tools
	s.registerTools()

	// Register all resources
	s.registerResources()

	return s
}

// Run starts the MCP server with stdio transport.
func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

// registerTools registers all MCP tools.
func (s *Server) registerTools() {
	// Tool: list_languages
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "list_languages",
		Description: "List all supported target languages with their conventions and idioms",
	}, s.handleListLanguages)

	// Tool: parse_spec
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "parse_spec",
		Description: "Parse and validate a markdown specification file",
	}, s.handleParseSpec)

	// Tool: validate_spec
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "validate_spec",
		Description: "Validate a spec file and return any errors or warnings",
	}, s.handleValidateSpec)

	// Tool: get_generation_context
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_generation_context",
		Description: "Get full context for code generation including parsed spec, language conventions, and prompt template",
	}, s.handleGetGenerationContext)

	// Tool: get_project_structure
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_project_structure",
		Description: "Get recommended file structure for a spec in a target language",
	}, s.handleGetProjectStructure)

	// Tool: ensure_parity
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "ensure_parity",
		Description: "Check feature parity across generated projects and provide fix instructions. Compares implementations against a reference (first project) and identifies missing features with suggested fixes.",
	}, s.handleEnsureParity)

	// Tool: import_spec_from_source
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "import_spec_from_source",
		Description: "Collect and analyze source code from any directory for AI-powered spec generation. Returns an analysis prompt containing all source files, tests, API specs, and configurations. The AI should use this prompt to generate a comprehensive .spec.md file at the specified output path.",
	}, s.handleImportSpecFromSource)
}

// registerResources registers all MCP resources.
func (s *Server) registerResources() {
	// Example spec resources
	s.mcpServer.AddResource(&mcp.Resource{
		Name:        "simple-function-example",
		URI:         "spec://examples/simple-function",
		Description: "Example spec for a simple function (slugify)",
		MIMEType:    "text/markdown",
	}, s.handleExampleResource("simple-function"))

	s.mcpServer.AddResource(&mcp.Resource{
		Name:        "module-example",
		URI:         "spec://examples/module",
		Description: "Example spec for a module with multiple functions and types",
		MIMEType:    "text/markdown",
	}, s.handleExampleResource("module"))

	s.mcpServer.AddResource(&mcp.Resource{
		Name:        "full-project-example",
		URI:         "spec://examples/full-project",
		Description: "Example spec for a full project scaffold",
		MIMEType:    "text/markdown",
	}, s.handleExampleResource("full-project"))

	// Language convention resources
	for _, lang := range s.registry.List() {
		s.mcpServer.AddResource(&mcp.Resource{
			Name:        lang.ID + "-conventions",
			URI:         "lang://" + lang.ID + "/conventions",
			Description: lang.Name + " language conventions and idioms",
			MIMEType:    "text/plain",
		}, s.handleLanguageConventions(lang.ID))
	}
}
