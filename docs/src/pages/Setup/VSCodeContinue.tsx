export default function VSCodeContinue() {
  return (
    <div className="prose-docs">
      <h1>VS Code Continue Setup</h1>
      <p className="lead">Configure RPG with the Continue AI assistant extension for VS Code.</p>

      <h2>Prerequisites</h2>
      <ul>
        <li>VS Code with Continue extension installed</li>
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
            <td><code>~/.continue/config.json</code></td>
          </tr>
          <tr>
            <td>Windows</td>
            <td><code>%USERPROFILE%\.continue\config.json</code></td>
          </tr>
          <tr>
            <td>Linux</td>
            <td><code>~/.continue/config.json</code></td>
          </tr>
        </tbody>
      </table>

      <h2>Configuration</h2>
      <p>Add RPG to the <code>mcpServers</code> section of your Continue config:</p>

      <pre className="code-block">
{`{
  "mcpServers": [
    {
      "name": "rpg",
      "command": "/absolute/path/to/rpg/bin/rpg",
      "args": []
    }
  ]
}`}
      </pre>

      <h2>Automatic Setup</h2>
      <pre className="code-block">
{`./scripts/setup-clients.sh vscode-continue`}
      </pre>

      <h2>Usage in VS Code</h2>
      <ol>
        <li>Open the Continue sidebar (Cmd/Ctrl + L)</li>
        <li>RPG tools will be available in the tools menu</li>
        <li>Ask Continue to generate code using your specs</li>
      </ol>

      <h3>Example Prompts</h3>
      <pre className="code-block">
{`// In Continue chat
"Use RPG to list supported languages"

"Generate TypeScript code from specs/url-shortener.spec.md"

"Check if my Go and Rust implementations have feature parity"`}
      </pre>

      <h2>Troubleshooting</h2>
      <h3>Server not connecting</h3>
      <ul>
        <li>Restart VS Code after config changes</li>
        <li>Check Continue's output panel for errors</li>
        <li>Verify the binary path is correct and absolute</li>
      </ul>

      <h3>Config file issues</h3>
      <ul>
        <li>Continue may create the config directory on first run</li>
        <li>Ensure valid JSON syntax</li>
        <li>The <code>mcpServers</code> key is an array, not an object</li>
      </ul>
    </div>
  );
}
