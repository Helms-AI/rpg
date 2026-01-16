import Tabs from '../../components/ui/Tabs';

export default function ValidateSpec() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              Is my spec file valid? It's at ./specs/my-api.spec.md
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll validate your specification file.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using validate_spec with specPath: "./specs/my-api.spec.md"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Your spec file is valid and contains content. You can proceed with code generation
              or parsing the spec for more details.
            </p>
          </div>
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              What about ./specs/missing.spec.md?
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using validate_spec with specPath: "./specs/missing.spec.md"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              The spec file was not found at that path. Please check the file location or create
              a new spec file.
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
{`// Validate before parsing
const validation = await mcp.callTool('rpg', 'validate_spec', {
  specPath: './specs/my-project.spec.md'
});

if (validation.valid) {
  // Safe to parse
  const spec = await mcp.callTool('rpg', 'parse_spec', {
    specPath: './specs/my-project.spec.md'
  });
} else {
  console.error('Invalid spec:', validation.message);
}`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>validate_spec</h1>
      <p className="lead">Check if a spec file exists and contains valid content.</p>

      <h2>Description</h2>
      <p>
        Validates that a specification file exists at the given path and contains meaningful content.
        Use this tool before attempting to parse or generate code from a spec to ensure
        the file is accessible and non-empty.
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
            <td>Path to the specification file to validate</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns a validation result object:</p>

      <pre className="code-block">
{`// Success
{
  "valid": true,
  "message": "Spec file is valid and contains content"
}

// Failure
{
  "valid": false,
  "message": "File not found: ./specs/missing.spec.md"
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Validation Checks</h2>
      <ul>
        <li>File exists at the specified path</li>
        <li>File is readable</li>
        <li>File contains non-empty content</li>
      </ul>
    </div>
  );
}
