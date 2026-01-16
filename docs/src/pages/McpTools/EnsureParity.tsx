export default function EnsureParity() {
  return (
    <div className="prose-docs">
      <h1>ensure_parity</h1>
      <p className="lead">Check feature parity across generated projects and provide fix instructions.</p>

      <h2>Description</h2>
      <p>
        Compares implementations across multiple languages against a reference (first project)
        and identifies missing features with suggested fixes. This is essential for maintaining
        consistent behavior when generating the same spec in multiple languages.
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
            <td><code>specPath</code></td>
            <td>string</td>
            <td>Yes</td>
            <td>Path to the specification file</td>
          </tr>
          <tr>
            <td><code>projects</code></td>
            <td>array</td>
            <td>Yes</td>
            <td>Array of project objects with language and path</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns a parity report with findings and suggestions:</p>

      <pre className="code-block">
{`{
  "reference": {
    "language": "go",
    "path": "./generated/go"
  },
  "comparisons": [
    {
      "language": "rust",
      "path": "./generated/rust",
      "status": "partial",
      "missing": [
        {
          "feature": "rate limiting",
          "suggestion": "Add rate_limit middleware in src/middleware.rs"
        }
      ],
      "differences": [
        {
          "feature": "error handling",
          "reference": "Returns structured Error type",
          "current": "Returns string errors",
          "suggestion": "Create custom Error enum with From implementations"
        }
      ]
    }
  ],
  "overallParity": 85
}`}
      </pre>

      <h2>Example Usage</h2>
      <pre className="code-block">
{`// Check parity across three implementations
const report = await mcp.callTool('rpg', 'ensure_parity', {
  specPath: './specs/url-shortener.spec.md',
  projects: [
    { language: 'go', path: './generated/go' },
    { language: 'rust', path: './generated/rust' },
    { language: 'typescript', path: './generated/typescript' }
  ]
});

// Review findings
for (const comparison of report.comparisons) {
  if (comparison.missing.length > 0) {
    console.log(\`\${comparison.language} is missing features:\`);
    comparison.missing.forEach(m => console.log(\`  - \${m.feature}\`));
  }
}`}
      </pre>

      <h2>Parity Checks</h2>
      <ul>
        <li><strong>Function coverage</strong> - All specified functions implemented</li>
        <li><strong>Type definitions</strong> - All types defined with correct fields</li>
        <li><strong>Error handling</strong> - Consistent error patterns</li>
        <li><strong>Test coverage</strong> - All test cases implemented</li>
        <li><strong>Configuration</strong> - Environment variables handled</li>
      </ul>
    </div>
  );
}
