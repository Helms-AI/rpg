export default function Installation() {
  return (
    <div className="prose-docs">
      <h1>Installation</h1>
      <p className="lead">Get RPG installed and ready to use in your development environment.</p>

      <h2>Prerequisites</h2>
      <ul>
        <li>Go 1.21 or later</li>
        <li>Git</li>
        <li>An MCP-compatible AI assistant (Claude Desktop, Cursor, etc.)</li>
      </ul>

      <h2>Quick Install</h2>
      <pre><code>{`git clone https://github.com/kon1790/rpg.git
cd rpg
./run.sh --build`}</code></pre>

      <h2>Setup AI Clients</h2>
      <pre><code>{`./scripts/setup-clients.sh`}</code></pre>
      <p>This will automatically configure RPG for all detected AI coding assistants.</p>
    </div>
  );
}
