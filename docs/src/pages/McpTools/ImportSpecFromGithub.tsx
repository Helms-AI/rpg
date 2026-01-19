import Tabs from '../../components/ui/Tabs';

export default function ImportSpecFromGithub() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              I want to port the hashicorp/consul project to Rust. Can you analyze it first?
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll clone and analyze the Consul repository from GitHub.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using import_spec_from_github with repository: "hashicorp/consul"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              I've analyzed the repository. Here's what I found:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li>Primary language: <strong>Go</strong></li>
              <li>Source files: 1,250+ files</li>
              <li>Key packages: agent, api, command, consul</li>
              <li>Test coverage: Extensive test suites</li>
            </ul>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Based on this analysis, I can generate a specification. Would you like me to focus on
              a specific component like the API client?
            </p>
          </div>
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              Yes, focus on the API client. Then generate Rust code.
            </p>
          </div>
        </div>
      ),
    },
    {
      id: 'code',
      label: 'Code',
      content: (
        <pre className="code-block">
{`// Clone and analyze a public GitHub repository
const analysis = await mcp.callTool('rpg', 'import_spec_from_github', {
  repository: 'owner/repo'
});

// With specific branch
const branchAnalysis = await mcp.callTool('rpg', 'import_spec_from_github', {
  repository: 'owner/repo',
  ref: 'develop'
});

// Private repository with token
const privateAnalysis = await mcp.callTool('rpg', 'import_spec_from_github', {
  repository: 'private-org/private-repo',
  token: process.env.GITHUB_TOKEN
});

// Full URL format also works
const urlAnalysis = await mcp.callTool('rpg', 'import_spec_from_github', {
  repository: 'https://github.com/owner/repo@v2.0.0'
});`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>import_spec_from_github</h1>
      <p className="lead">Clone and analyze a GitHub repository for AI-powered spec generation.</p>

      <h2>Description</h2>
      <p>
        Clones a GitHub repository and analyzes its source code to generate a comprehensive
        specification. Supports public repositories directly, and private repositories via
        a GitHub token. Accepts flexible input formats including shorthand notation.
      </p>

      <h2>Parameters</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Parameter</th>
            <th>Type</th>
            <th>Required</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>repository</code></td>
            <td>string</td>
            <td>Yes</td>
            <td>GitHub URL or shorthand (e.g., "owner/repo", "owner/repo@branch", full URL)</td>
          </tr>
          <tr>
            <td><code>ref</code></td>
            <td>string</td>
            <td>No</td>
            <td>Branch, tag, or commit SHA (overrides ref in URL)</td>
          </tr>
          <tr>
            <td><code>token</code></td>
            <td>string</td>
            <td>No</td>
            <td>GitHub PAT for private repos (or use GITHUB_TOKEN env var)</td>
          </tr>
          <tr>
            <td><code>name</code></td>
            <td>string</td>
            <td>No</td>
            <td>Optional name for the generated spec</td>
          </tr>
          <tr>
            <td><code>shallow</code></td>
            <td>boolean</td>
            <td>No</td>
            <td>Use shallow clone for faster download (default: true)</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns an analysis prompt for spec generation:</p>

      <pre className="code-block">
{`{
  "analysisPrompt": "Based on the following source code analysis...\\n\\n## Repository\\n...",
  "repository": "owner/repo",
  "ref": "main",
  "sourceFiles": [
    { "path": "cmd/main.go", "language": "go", "lines": 150 },
    { "path": "pkg/api/client.go", "language": "go", "lines": 450 }
  ],
  "detectedLanguage": "go",
  "suggestedName": "consul-api"
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Repository Formats</h2>
      <ul>
        <li><code>owner/repo</code> - Simple shorthand (uses default branch)</li>
        <li><code>owner/repo@branch</code> - Shorthand with branch/tag</li>
        <li><code>owner/repo@v1.2.3</code> - Shorthand with version tag</li>
        <li><code>https://github.com/owner/repo</code> - Full URL</li>
        <li><code>https://github.com/owner/repo@main</code> - Full URL with ref</li>
      </ul>

      <h2>Use Cases</h2>
      <ul>
        <li>
          <strong>Port open-source projects</strong> - Analyze popular Go libraries
          and port them to Rust or TypeScript
        </li>
        <li>
          <strong>Learn from existing code</strong> - Generate specs from well-designed
          projects to understand their architecture
        </li>
        <li>
          <strong>Fork and extend</strong> - Analyze a project before creating your own
          implementation with modifications
        </li>
        <li>
          <strong>Multi-language SDKs</strong> - Generate consistent SDKs in multiple
          languages from a single reference implementation
        </li>
      </ul>

      <h2>Authentication</h2>
      <p>For private repositories, provide authentication via:</p>
      <ul>
        <li><code>GITHUB_TOKEN</code> environment variable (recommended)</li>
        <li><code>token</code> parameter in the tool call</li>
      </ul>
      <p>
        The token needs <code>repo</code> scope for private repository access.
      </p>
    </div>
  );
}
