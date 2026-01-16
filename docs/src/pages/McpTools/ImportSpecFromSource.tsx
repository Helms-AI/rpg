import Tabs from '../../components/ui/Tabs';

export default function ImportSpecFromSource() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              I have an existing Go project at ./legacy-api. Can you help me create a spec so I can port it to TypeScript?
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll analyze your existing code to generate a specification.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using import_spec_from_source with inputPath: "./legacy-api", name: "legacy-api"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              I've analyzed your project. Here's what I found:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li>Detected language: <strong>Go</strong></li>
              <li>Source files: 8 files (1,250 lines)</li>
              <li>Functions: 12 exported functions</li>
              <li>Types: 5 structs</li>
            </ul>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Based on this analysis, I can now generate a specification file. Would you like me to create
              <code>./specs/legacy-api.spec.md</code>?
            </p>
          </div>
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              Yes, create the spec and then generate TypeScript code.
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200">
              I'll create the spec and then generate the TypeScript implementation...
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
{`// Reverse-engineer a spec from existing code
const analysis = await mcp.callTool('rpg', 'import_spec_from_source', {
  inputPath: './existing-project',
  name: 'my-service'
});

// The AI can use analysisPrompt to generate a spec
// The spec can then be used to generate code in other languages`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>import_spec_from_source</h1>
      <p className="lead">Analyze source code for AI-powered spec generation.</p>

      <h2>Description</h2>
      <p>
        Collects and analyzes source code from any directory to help generate a specification.
        Returns an analysis prompt containing all source files, tests, API specs, and configurations.
        The AI uses this context to generate a comprehensive .spec.md file.
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
            <td><code>inputPath</code></td>
            <td>string</td>
            <td>Yes</td>
            <td>Path to the source code directory to analyze</td>
          </tr>
          <tr>
            <td><code>name</code></td>
            <td>string</td>
            <td>No</td>
            <td>Optional name for the generated spec</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns an analysis prompt for spec generation:</p>

      <pre className="code-block">
{`{
  "analysisPrompt": "Based on the following source code analysis...\\n\\n## Source Files\\n...",
  "sourceFiles": [
    { "path": "src/main.go", "language": "go", "lines": 150 },
    { "path": "src/handler.go", "language": "go", "lines": 200 }
  ],
  "detectedLanguage": "go",
  "suggestedName": "url-shortener"
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Use Cases</h2>
      <ul>
        <li>
          <strong>Port to new language</strong> - Analyze existing Go code, generate spec,
          then generate TypeScript implementation
        </li>
        <li>
          <strong>Documentation</strong> - Generate specification documentation from
          undocumented legacy code
        </li>
        <li>
          <strong>Refactoring</strong> - Create a spec to understand current behavior
          before making changes
        </li>
        <li>
          <strong>Multi-language libraries</strong> - Start with one implementation,
          generate spec, create versions in all supported languages
        </li>
      </ul>

      <h2>Analyzed Content</h2>
      <p>The importer examines:</p>
      <ul>
        <li>Source code files (functions, types, interfaces)</li>
        <li>Test files (test cases, expected behaviors)</li>
        <li>Configuration files (environment variables, settings)</li>
        <li>API specifications (OpenAPI, GraphQL schemas)</li>
        <li>Package manifests (dependencies, versions)</li>
      </ul>
    </div>
  );
}
