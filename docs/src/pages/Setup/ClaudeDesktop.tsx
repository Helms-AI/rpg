export default function ClaudeDesktop() {
  return (
    <div className="prose-docs">
      <h1>Claude Desktop Setup</h1>
      <p className="lead">Configure RPG as an MCP server in Anthropic's Claude Desktop application.</p>

      <h2>Prerequisites</h2>
      <ul>
        <li>Claude Desktop installed (macOS, Windows, or Linux)</li>
        <li>RPG binary built: <code>make build</code></li>
      </ul>

      <h2>Configuration File Location</h2>
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
            <td><code>~/Library/Application Support/Claude/claude_desktop_config.json</code></td>
          </tr>
          <tr>
            <td>Windows</td>
            <td><code>%APPDATA%\Claude\claude_desktop_config.json</code></td>
          </tr>
          <tr>
            <td>Linux</td>
            <td><code>~/.config/Claude/claude_desktop_config.json</code></td>
          </tr>
        </tbody>
      </table>

      <h2>Configuration</h2>
      <p>Add RPG to your <code>claude_desktop_config.json</code>:</p>

      <pre className="code-block">
{`{
  "mcpServers": {
    "rpg": {
      "command": "/absolute/path/to/rpg/bin/rpg",
      "args": [],
      "env": {}
    }
  }
}`}
      </pre>

      <h2>Automatic Setup</h2>
      <p>Use the setup script to automatically configure Claude Desktop:</p>

      <pre className="code-block">
{`# From the RPG project root
./scripts/setup-clients.sh claude-desktop`}
      </pre>

      <h2>Verify Installation</h2>
      <ol>
        <li>Restart Claude Desktop</li>
        <li>Look for the MCP tools icon in the chat interface</li>
        <li>Click it to see "rpg" listed as an available server</li>
        <li>Try asking Claude to "list available languages for RPG"</li>
      </ol>

      <h2>Troubleshooting</h2>
      <h3>Server not appearing</h3>
      <ul>
        <li>Ensure the path is absolute (starts with <code>/</code> on macOS/Linux or drive letter on Windows)</li>
        <li>Verify the binary exists and is executable</li>
        <li>Check Claude Desktop logs for connection errors</li>
      </ul>

      <h3>Tools not working</h3>
      <ul>
        <li>Restart Claude Desktop after config changes</li>
        <li>Verify JSON syntax is valid</li>
        <li>Test the binary manually: <code>./bin/rpg</code></li>
      </ul>
    </div>
  );
}
