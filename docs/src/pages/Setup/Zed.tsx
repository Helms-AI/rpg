export default function Zed() {
  return (
    <div className="prose-docs">
      <h1>Zed Setup</h1>
      <p className="lead">Configure RPG as an MCP server in Zed editor.</p>

      <div className="not-prose my-6 p-4 rounded-lg bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800">
        <p className="text-amber-800 dark:text-amber-300 text-sm">
          <strong>Note:</strong> Zed's MCP support is currently in development.
          Check Zed's documentation for the latest configuration format.
        </p>
      </div>

      <h2>Prerequisites</h2>
      <ul>
        <li>Zed editor installed (latest version)</li>
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
            <td><code>~/.config/zed/settings.json</code></td>
          </tr>
          <tr>
            <td>Linux</td>
            <td><code>~/.config/zed/settings.json</code></td>
          </tr>
        </tbody>
      </table>

      <h2>Configuration</h2>
      <p>Add RPG to Zed's settings:</p>

      <pre className="code-block">
{`{
  "mcp": {
    "servers": {
      "rpg": {
        "command": "/absolute/path/to/rpg/bin/rpg",
        "args": []
      }
    }
  }
}`}
      </pre>

      <h2>Automatic Setup</h2>
      <pre className="code-block">
{`./scripts/setup-clients.sh zed`}
      </pre>

      <h2>Usage in Zed</h2>
      <ol>
        <li>Open Zed's AI assistant panel</li>
        <li>RPG tools should be available for code generation</li>
        <li>Reference spec files when requesting implementations</li>
      </ol>

      <h3>Example Prompts</h3>
      <pre className="code-block">
{`// In Zed assistant
"List RPG's supported programming languages"

"Generate Python code from my spec file"

"Help me create a spec for my existing Go code"`}
      </pre>

      <h2>Troubleshooting</h2>
      <h3>MCP not available</h3>
      <ul>
        <li>Ensure you're using a Zed version that supports MCP</li>
        <li>Check Zed's release notes for MCP feature availability</li>
        <li>Restart Zed after configuration changes</li>
      </ul>

      <h3>Configuration issues</h3>
      <ul>
        <li>Zed's settings.json must be valid JSON</li>
        <li>The MCP section format may change - check Zed docs</li>
        <li>Use absolute paths for the binary location</li>
      </ul>
    </div>
  );
}
