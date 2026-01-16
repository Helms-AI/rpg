export default function Windsurf() {
  return (
    <div className="prose-docs">
      <h1>Windsurf Setup</h1>
      <p className="lead">Configure RPG as an MCP server in Codeium's Windsurf IDE.</p>

      <h2>Prerequisites</h2>
      <ul>
        <li>Windsurf IDE installed</li>
        <li>RPG binary built: <code>make build</code></li>
      </ul>

      <h2>Configuration File Location</h2>

      <h3>Project Level (Recommended)</h3>
      <p>Create <code>.windsurf/mcp.json</code> in your project root:</p>

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
            <td><code>~/.windsurf/mcp.json</code></td>
          </tr>
          <tr>
            <td>Windows</td>
            <td><code>%USERPROFILE%\.windsurf\mcp.json</code></td>
          </tr>
          <tr>
            <td>Linux</td>
            <td><code>~/.windsurf/mcp.json</code></td>
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
{`./scripts/setup-clients.sh windsurf`}
      </pre>

      <h2>Usage in Windsurf</h2>
      <ol>
        <li>Open Windsurf's Cascade AI assistant</li>
        <li>RPG tools will be available for code generation tasks</li>
        <li>Reference your spec files when asking for implementations</li>
      </ol>

      <h3>Example Prompts</h3>
      <pre className="code-block">
{`// In Windsurf Cascade
"What languages does RPG support?"

"Generate Rust code from the url-shortener spec"

"Analyze my existing code and create a spec file"`}
      </pre>

      <h2>Troubleshooting</h2>
      <h3>Server not detected</h3>
      <ul>
        <li>Restart Windsurf after adding configuration</li>
        <li>Ensure the mcp.json file is in the correct location</li>
        <li>Verify JSON syntax is valid</li>
      </ul>

      <h3>Connection issues</h3>
      <ul>
        <li>Check that the RPG binary is built and executable</li>
        <li>Test the binary manually: <code>./bin/rpg</code></li>
        <li>Review Windsurf logs for MCP connection errors</li>
      </ul>
    </div>
  );
}
