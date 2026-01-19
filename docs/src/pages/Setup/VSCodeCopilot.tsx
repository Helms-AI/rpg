export default function VSCodeCopilot() {
  return (
    <div className="prose-docs">
      <h1>VS Code + GitHub Copilot Setup</h1>
      <p className="lead">Configure RPG with GitHub Copilot's native MCP support in VS Code.</p>

      <div className="not-prose my-6 p-4 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800">
        <p className="text-emerald-800 dark:text-emerald-300 text-sm">
          <strong>Native Integration:</strong> VS Code 1.102+ includes built-in MCP support for GitHub Copilot.
          No additional extensions required.
        </p>
      </div>

      <h2>Prerequisites</h2>
      <ul>
        <li>VS Code 1.102 or later</li>
        <li>GitHub Copilot extension with an active subscription</li>
        <li>RPG binary built: <code>make build</code></li>
      </ul>

      <h2>Configuration Options</h2>
      <p>You can configure MCP servers at the workspace or user level.</p>

      <h3>Option 1: Workspace Configuration (Recommended)</h3>
      <p>Create <code>.vscode/mcp.json</code> in your project root:</p>

      <pre className="code-block">
{`{
  "servers": {
    "rpg": {
      "type": "stdio",
      "command": "/absolute/path/to/rpg/bin/rpg",
      "args": ["--output", "./generated"]
    }
  }
}`}
      </pre>

      <p>For relative paths within a workspace, use the <code>$&#123;workspaceFolder&#125;</code> variable:</p>

      <pre className="code-block">
{`{
  "servers": {
    "rpg": {
      "type": "stdio",
      "command": "\${workspaceFolder}/bin/rpg",
      "args": ["--output", "./generated"]
    }
  }
}`}
      </pre>

      <h3>Option 2: User-Level Configuration</h3>
      <p>Add to your VS Code <code>settings.json</code> (Cmd/Ctrl+Shift+P → "Preferences: Open User Settings (JSON)"):</p>

      <pre className="code-block">
{`{
  "github.copilot.chat.mcp.servers": {
    "rpg": {
      "type": "stdio",
      "command": "/absolute/path/to/rpg/bin/rpg",
      "args": ["--output", "./generated"]
    }
  }
}`}
      </pre>

      <h2>Environment Variables</h2>
      <p>You can pass environment variables to the MCP server:</p>

      <pre className="code-block">
{`{
  "servers": {
    "rpg": {
      "type": "stdio",
      "command": "/path/to/rpg/bin/rpg",
      "args": ["--output", "./generated"],
      "env": {
        "RPG_LOG_LEVEL": "debug"
      }
    }
  }
}`}
      </pre>

      <h2>Using RPG with Copilot</h2>
      <ol>
        <li>Open Copilot Chat (Cmd/Ctrl+Shift+I or click the chat icon)</li>
        <li>Switch to <strong>Agent mode</strong> using the dropdown at the top</li>
        <li>RPG tools will be available automatically</li>
      </ol>

      <h3>Example Prompts</h3>
      <pre className="code-block">
{`// List available languages
"Use RPG to list supported languages"

// Generate code from a spec
"Generate TypeScript code from specs/url-shortener.spec.md using RPG"

// Import from GitHub
"Use RPG to import a spec from github.com/example/repo"

// Check parity between implementations
"Use RPG to verify parity between my TypeScript and Go implementations"`}
      </pre>

      <h2>Verifying the Connection</h2>
      <ol>
        <li>Open Command Palette (Cmd/Ctrl+Shift+P)</li>
        <li>Run <strong>"MCP: List Servers"</strong></li>
        <li>RPG should appear in the list with a connected status</li>
      </ol>

      <h2>Available Tools</h2>
      <p>Once connected, these RPG tools are available to Copilot:</p>

      <table className="w-full">
        <thead>
          <tr>
            <th>Tool</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>list_languages</code></td>
            <td>List supported target languages</td>
          </tr>
          <tr>
            <td><code>parse_spec</code></td>
            <td>Parse a markdown specification file</td>
          </tr>
          <tr>
            <td><code>get_generation_context</code></td>
            <td>Get context for code generation</td>
          </tr>
          <tr>
            <td><code>get_project_structure</code></td>
            <td>Get recommended project structure</td>
          </tr>
          <tr>
            <td><code>import_spec_from_source</code></td>
            <td>Generate spec from source code</td>
          </tr>
          <tr>
            <td><code>import_spec_from_github</code></td>
            <td>Import and analyze a GitHub repository</td>
          </tr>
          <tr>
            <td><code>ensure_parity</code></td>
            <td>Check feature parity across implementations</td>
          </tr>
          <tr>
            <td><code>deep_analyze_source</code></td>
            <td>Deep semantic analysis of source code</td>
          </tr>
        </tbody>
      </table>

      <h2>Troubleshooting</h2>

      <h3>Server not appearing</h3>
      <ul>
        <li>Reload VS Code window after adding the configuration</li>
        <li>Ensure the binary path is absolute (or uses <code>$&#123;workspaceFolder&#125;</code>)</li>
        <li>Verify the RPG binary is built: <code>make build</code></li>
        <li>Check that Copilot is in Agent mode, not Chat mode</li>
      </ul>

      <h3>Connection errors</h3>
      <ul>
        <li>Run the binary manually to check for errors: <code>./bin/rpg</code></li>
        <li>Check the VS Code Output panel (View → Output → select "GitHub Copilot")</li>
        <li>Verify JSON syntax in your configuration file</li>
      </ul>

      <h3>Tools not available</h3>
      <ul>
        <li>Ensure you're using <strong>Agent mode</strong> in Copilot Chat</li>
        <li>Some organizations may have MCP disabled—check with your admin</li>
        <li>Try restarting the MCP server: Command Palette → "MCP: Restart Server"</li>
      </ul>

      <h2>Enterprise Considerations</h2>
      <p>
        Organizations can enable or disable MCP support via the "MCP servers in Copilot" policy.
        If you can't connect to MCP servers, contact your GitHub organization administrator.
      </p>
    </div>
  );
}
