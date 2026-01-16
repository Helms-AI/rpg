import { Link } from 'react-router-dom';
import { FileJson, Wrench } from 'lucide-react';

export default function ApiIndex() {
  return (
    <div className="prose-docs">
      <h1>API Reference</h1>
      <p className="lead">
        Complete technical reference for RPG's MCP server interface and data schemas.
      </p>

      <div className="not-prose mt-8 grid gap-4 md:grid-cols-2">
        <Link
          to="/api/schemas"
          className="group flex flex-col gap-3 rounded-lg border border-gray-200 dark:border-gray-700 p-6 hover:border-primary-500 dark:hover:border-primary-400 transition-colors"
        >
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400">
              <FileJson className="w-5 h-5" />
            </div>
            <h3 className="font-semibold text-gray-900 dark:text-white group-hover:text-primary-600 dark:group-hover:text-primary-400">
              Data Schemas
            </h3>
          </div>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            JSON schemas for specs, types, functions, and tool responses
          </p>
        </Link>

        <Link
          to="/mcp-tools"
          className="group flex flex-col gap-3 rounded-lg border border-gray-200 dark:border-gray-700 p-6 hover:border-primary-500 dark:hover:border-primary-400 transition-colors"
        >
          <div className="flex items-center gap-3">
            <div className="p-2 rounded-lg bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400">
              <Wrench className="w-5 h-5" />
            </div>
            <h3 className="font-semibold text-gray-900 dark:text-white group-hover:text-primary-600 dark:group-hover:text-primary-400">
              MCP Tools
            </h3>
          </div>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            Reference for all 7 MCP tools with parameters and responses
          </p>
        </Link>
      </div>

      <h2>MCP Protocol</h2>
      <p>
        RPG implements the{' '}
        <a href="https://modelcontextprotocol.io" target="_blank" rel="noopener noreferrer">
          Model Context Protocol (MCP)
        </a>
        , a standard for connecting AI assistants to external tools and data sources.
      </p>

      <h3>Server Configuration</h3>
      <pre className="code-block">
{`{
  "mcpServers": {
    "rpg": {
      "command": "./bin/rpg",
      "args": [],
      "env": {}
    }
  }
}`}
      </pre>

      <h3>Tool Invocation</h3>
      <p>AI assistants call RPG tools using the MCP protocol:</p>
      <pre className="code-block">
{`// Request
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_generation_context",
    "arguments": {
      "specPath": "./specs/my-project.spec.md",
      "language": "go"
    }
  },
  "id": 1
}

// Response
{
  "jsonrpc": "2.0",
  "result": {
    "content": [{
      "type": "text",
      "text": "..."
    }]
  },
  "id": 1
}`}
      </pre>

      <h2>Error Handling</h2>
      <p>RPG returns structured errors in MCP format:</p>

      <pre className="code-block">
{`{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": {
      "field": "specPath",
      "error": "file not found: ./specs/missing.spec.md"
    }
  },
  "id": 1
}`}
      </pre>

      <h3>Error Codes</h3>
      <table className="w-full">
        <thead>
          <tr>
            <th>Code</th>
            <th>Meaning</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>-32602</code></td>
            <td>Invalid parameters (missing required field, invalid path)</td>
          </tr>
          <tr>
            <td><code>-32603</code></td>
            <td>Internal error (parsing failure, unexpected error)</td>
          </tr>
          <tr>
            <td><code>-32601</code></td>
            <td>Method not found (unknown tool name)</td>
          </tr>
        </tbody>
      </table>

      <h2>Rate Limiting</h2>
      <p>
        RPG has no built-in rate limiting. It processes requests as fast as the
        underlying system allows. For production deployments, consider adding
        rate limiting at the network layer.
      </p>
    </div>
  );
}
