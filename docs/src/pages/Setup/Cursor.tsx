export default function Cursor() {
  return (
    <div className="prose-docs">
      <h1>Cursor Setup</h1>
      <p className="lead">Configure RPG as an MCP server in Cursor IDE.</p>

      <h2>Prerequisites</h2>
      <ul>
        <li>Cursor IDE installed</li>
        <li>RPG binary built: <code>make build</code></li>
      </ul>

      <h2>Configuration File Location</h2>
      <p>
        Cursor supports both project-level and global MCP configuration.
      </p>

      <h3>Project Level (Recommended)</h3>
      <p>Create <code>.cursor/mcp.json</code> in your project root:</p>

      <pre className="code-block">
{`{
  "mcpServers": {
    "rpg": {
      "command": "./bin/rpg",
      "args": []
    }
  }
}`}
      </pre>

      <h3>Global Configuration</h3>
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
            <td><code>~/.cursor/mcp.json</code></td>
          </tr>
          <tr>
            <td>Windows</td>
            <td><code>%USERPROFILE%\.cursor\mcp.json</code></td>
          </tr>
          <tr>
            <td>Linux</td>
            <td><code>~/.cursor/mcp.json</code></td>
          </tr>
        </tbody>
      </table>

      <pre className="code-block">
{`{
  "mcpServers": {
    "rpg": {
      "command": "/absolute/path/to/rpg/bin/rpg",
      "args": []
    }
  }
}`}
      </pre>

      <h2>Automatic Setup</h2>
      <pre className="code-block">
{`./scripts/setup-clients.sh cursor`}
      </pre>

      <h2>Usage in Cursor</h2>
      <ol>
        <li>Open Cursor's AI chat (Cmd/Ctrl + K)</li>
        <li>RPG tools will be automatically available</li>
        <li>Reference spec files and ask for code generation</li>
      </ol>

      <h3>Example Prompts</h3>
      <pre className="code-block">
{`// In Cursor chat
"@rpg list languages"

"Generate a Go implementation from @specs/my-api.spec.md"

"Use RPG to ensure parity between my implementations"`}
      </pre>

      <h2>Troubleshooting</h2>
      <h3>MCP not loading</h3>
      <ul>
        <li>Reload Cursor window (Cmd/Ctrl + Shift + P → "Reload Window")</li>
        <li>Check that the config file is valid JSON</li>
        <li>Verify the binary path is correct</li>
      </ul>

      <h3>Tools not appearing</h3>
      <ul>
        <li>MCP servers are loaded when Cursor starts</li>
        <li>Try closing and reopening the project</li>
        <li>Check Cursor's developer console for errors</li>
      </ul>
    </div>
  );
}
