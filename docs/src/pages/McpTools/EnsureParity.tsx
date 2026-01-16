import Tabs from '../../components/ui/Tabs';

export default function EnsureParity() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              I generated the URL shortener in Go, Rust, and TypeScript. Can you check if they all have the same features?
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll compare all three implementations against the spec to check for feature parity.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using ensure_parity with specPath and 3 projects (Go, Rust, TypeScript)
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Here's the parity report:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-2">
              <li><strong>Go (reference)</strong> - Complete implementation</li>
              <li><strong>Rust</strong> - Missing rate limiting. Add <code>rate_limit</code> middleware in src/middleware.rs</li>
              <li><strong>TypeScript</strong> - Different error handling. Consider using custom Error classes instead of string errors</li>
            </ul>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Overall parity: <strong>85%</strong>. The Rust and TypeScript implementations need minor adjustments to match the reference.
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
      ),
    },
  ];

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
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

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
