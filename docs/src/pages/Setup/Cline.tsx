export default function Cline() {
  return (
    <div className="prose-docs">
      <h1>Cline Setup</h1>
      <p className="lead">Configure RPG with Cline, the autonomous coding agent for VS Code.</p>

      <h2>Prerequisites</h2>
      <ul>
        <li>VS Code with Cline extension installed</li>
        <li>RPG binary built: <code>make build</code></li>
      </ul>

      <h2>Configuration File Location</h2>

      <h3>Project Level (Recommended)</h3>
      <p>Create <code>.vscode/cline_mcp_settings.json</code> in your project:</p>

      <pre className="code-block">
{`{
  "mcpServers": {
    "rpg": {
      "command": "./bin/rpg",
      "args": [],
      "disabled": false
    }
  }
}`}
      </pre>

      <h3>Global Configuration</h3>
      <p>Cline stores global settings in VS Code's settings directory:</p>
      <table className="w-full">
        <thead>
          <tr>
            <th>Platform</th>
            <th>Path</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>macOS</td>
            <td><code>~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json</code></td>
          </tr>
          <tr>
            <td>Windows</td>
            <td><code>%APPDATA%\Code\User\globalStorage\saoudrizwan.claude-dev\settings\cline_mcp_settings.json</code></td>
          </tr>
          <tr>
            <td>Linux</td>
            <td><code>~/.config/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json</code></td>
          </tr>
        </tbody>
      </table>

      <h2>Automatic Setup</h2>
      <pre className="code-block">
{`./scripts/setup-clients.sh cline`}
      </pre>

      <h2>Usage with Cline</h2>
      <ol>
        <li>Open the Cline sidebar in VS Code</li>
        <li>Click the MCP servers icon to verify RPG is connected</li>
        <li>Ask Cline to use RPG for code generation</li>
      </ol>

      <h3>Example Tasks</h3>
      <pre className="code-block">
{`// In Cline chat
"Use RPG to generate a Go REST API from specs/api.spec.md"

"Create implementations in all supported languages"

"Check feature parity across my Go, Rust, and Python versions"

"Analyze the existing codebase and generate a spec file"`}
      </pre>

      <h2>Cline + RPG Workflow</h2>
      <p>
        Cline's autonomous capabilities work well with RPG for complex tasks:
      </p>
      <ol>
        <li>Cline can create spec files based on your requirements</li>
        <li>Use RPG tools to generate code in multiple languages</li>
        <li>Cline can then create the files and set up the project structure</li>
        <li>Finally, use <code>ensure_parity</code> to verify implementations match</li>
      </ol>

      <h2>Troubleshooting</h2>
      <h3>MCP server not appearing</h3>
      <ul>
        <li>Reload VS Code window after config changes</li>
        <li>Check Cline's MCP settings panel for errors</li>
        <li>Ensure <code>"disabled": false</code> is set</li>
      </ul>

      <h3>Permission issues</h3>
      <ul>
        <li>Cline may prompt for permission to use MCP tools</li>
        <li>Approve the RPG server when prompted</li>
        <li>Check Cline's settings for auto-approval options</li>
      </ul>
    </div>
  );
}
