import Tabs from '../../components/ui/Tabs';

export default function ParseSpec() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              Can you read my spec file and tell me what functions it defines?
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll parse your specification file to extract its structure.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using parse_spec with specPath: "./specs/my-project.spec.md"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Your spec "my-project" defines 3 functions:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2">
              <li><code>createUser</code> - Creates a new user account</li>
              <li><code>getUser</code> - Retrieves a user by ID</li>
              <li><code>deleteUser</code> - Removes a user from the system</li>
            </ul>
          </div>
        </div>
      ),
    },
    {
      id: 'code',
      label: 'Code',
      content: (
        <pre className="code-block">
{`// Parse a specification file
const result = await mcp.callTool('rpg', 'parse_spec', {
  specPath: './specs/my-project.spec.md'
});

console.log(result.name);        // "my-project"
console.log(result.functions);   // Array of function definitions`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>parse_spec</h1>
      <p className="lead">Parse a markdown specification file and return its structured content.</p>

      <h2>Description</h2>
      <p>
        Reads a markdown specification file and extracts its content into a structured format.
        The spec can be in any narrative format - architecture docs, API designs, or feature descriptions.
        The AI interprets the spec content directly to generate idiomatic code.
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
            <td>Path to the markdown specification file</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns the parsed specification with extracted sections:</p>

      <pre className="code-block">
{`{
  "name": "url-shortener",
  "meta": {
    "version": "1.0.0",
    "author": "Your Name"
  },
  "targetLanguages": ["go", "rust", "typescript"],
  "types": [...],
  "functions": [...],
  "tests": [...],
  "dependencies": [...],
  "configuration": [...]
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Spec Sections</h2>
      <p>The parser extracts these markdown sections (case-insensitive):</p>
      <ul>
        <li><code># Title</code> (H1) - Spec name</li>
        <li><code>## Meta</code> - Version, author, license</li>
        <li><code>## Target Languages</code> - Language IDs</li>
        <li><code>## Types</code> - Data structure definitions</li>
        <li><code>## Functions</code> - Function signatures with accepts/returns/logic</li>
        <li><code>## Tests</code> - Test cases with given/expect format</li>
        <li><code>## Dependencies</code> - External dependencies</li>
        <li><code>## Configuration</code> - Environment variables</li>
      </ul>
    </div>
  );
}
