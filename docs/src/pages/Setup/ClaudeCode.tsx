export default function ClaudeCode() {
  return (
    <div className="prose-docs">
      <h1>Claude Code Setup</h1>
      <p className="lead">Configure RPG as an MCP server in Claude Code (CLI).</p>

      <h2>Prerequisites</h2>
      <ul>
        <li>Claude Code CLI installed</li>
        <li>RPG binary built: <code>make build</code></li>
      </ul>

      <h2>Configuration</h2>
      <p>
        Claude Code uses a <code>.mcp.json</code> file in your project root for workspace-level
        MCP configuration. RPG includes this file by default.
      </p>

      <h3>Project Configuration (Recommended)</h3>
      <p>The <code>.mcp.json</code> in the RPG project root:</p>

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
      <p>To use RPG from any project, add to <code>~/.claude/settings.json</code>:</p>

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
{`# Configure for current project
./scripts/setup-clients.sh claude-code`}
      </pre>

      <h2>Usage</h2>
      <p>Once configured, Claude Code can access RPG tools directly:</p>

      <pre className="code-block">
{`# In your terminal with Claude Code
claude

# Then ask:
> List the languages supported by RPG
> Generate Go code from specs/my-spec.md
> Check parity between my Go and Rust implementations`}
      </pre>

      <h2>Verify Installation</h2>
      <ol>
        <li>Open a terminal in your RPG project directory</li>
        <li>Run <code>claude</code></li>
        <li>Ask "What MCP tools are available?"</li>
        <li>You should see RPG's tools listed</li>
      </ol>

      <h2>Troubleshooting</h2>
      <h3>Tools not available</h3>
      <ul>
        <li>Make sure you're in a directory with <code>.mcp.json</code> or have global config</li>
        <li>Run <code>make build</code> to ensure the binary exists</li>
        <li>Check that <code>./bin/rpg</code> is executable</li>
      </ul>

      <h3>Path issues</h3>
      <ul>
        <li>Project config can use relative paths (<code>./bin/rpg</code>)</li>
        <li>Global config must use absolute paths</li>
      </ul>
    </div>
  );
}
