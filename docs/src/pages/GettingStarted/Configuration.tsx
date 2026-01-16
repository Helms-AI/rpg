export default function Configuration() {
  return (
    <div className="prose-docs">
      <h1>Configuration</h1>
      <p className="lead">Configure RPG for your development workflow.</p>

      <h2>MCP Configuration</h2>
      <p>RPG uses the Model Context Protocol (MCP) to communicate with AI assistants.</p>

      <h3>Claude Code</h3>
      <pre><code>{`// .mcp.json
{
  "mcpServers": {
    "rpg": {
      "command": "./bin/rpg",
      "args": []
    }
  }
}`}</code></pre>

      <h3>Output Directory</h3>
      <p>By default, generated code is placed in <code>./generated</code>. Configure with:</p>
      <pre><code>{`./bin/rpg --output /path/to/output`}</code></pre>

      <h2>Environment Variables</h2>
      <ul>
        <li><code>RPG_OUTPUT_DIR</code> - Default output directory</li>
        <li><code>RPG_BINARY_PATH</code> - Path to RPG binary</li>
      </ul>
    </div>
  );
}
