import Tabs from '../../components/ui/Tabs';

export default function GenerateSourceFromSpec() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              Generate a complete Rust implementation from my API spec. Handle everything automatically.
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll use autonomous code generation to create the Rust implementation.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using generate_source_from_spec with specPath: "./my-api.spec.md", language: "rust"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Generation complete! Here's what was created:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li><code>Cargo.toml</code> - Project manifest with dependencies</li>
              <li><code>src/lib.rs</code> - Main library with public API</li>
              <li><code>src/types.rs</code> - Type definitions</li>
              <li><code>src/handlers.rs</code> - Request handlers</li>
              <li><code>tests/integration_test.rs</code> - Test suite</li>
            </ul>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Parity validation passed: <strong>98%</strong> match with spec requirements.
              The code is ready at <code>./output/my-api/rust/</code>.
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
{`// Basic autonomous generation
const result = await mcp.callTool('rpg', 'generate_source_from_spec', {
  specPath: './my-api.spec.md',
  language: 'rust'
});

// With custom output directory
const custom = await mcp.callTool('rpg', 'generate_source_from_spec', {
  specPath: './specs/user-service.spec.md',
  language: 'typescript',
  outputDir: './services/user'
});

// The tool handles the complete workflow:
// 1. Parses the spec
// 2. Gets language conventions
// 3. Generates source code
// 4. Validates parity with spec
// 5. Iterates if needed to fix gaps
console.log(result.outputPath);   // Where code was generated
console.log(result.parityScore);  // How well it matches spec
console.log(result.files);        // List of generated files`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>generate_source_from_spec</h1>
      <p className="lead">Autonomous code generation from spec with automatic parity validation.</p>

      <h2>Description</h2>
      <p>
        Fully autonomous code generation that parses a spec, generates complete source code
        for the target language, validates the output against the spec, and iterates to fix
        any gaps. No manual intervention required - the tool handles the complete workflow.
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
            <td>Path to the specification file (.spec.md)</td>
          </tr>
          <tr>
            <td><code>language</code></td>
            <td>string</td>
            <td>Yes</td>
            <td>Target language (go, rust, typescript, python, java, csharp)</td>
          </tr>
          <tr>
            <td><code>outputDir</code></td>
            <td>string</td>
            <td>No</td>
            <td>Output directory (default: ./output/&lt;spec-name&gt;/&lt;language&gt;/)</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns generation results:</p>

      <pre className="code-block">
{`{
  "success": true,
  "language": "rust",
  "outputPath": "./output/my-api/rust",
  "parityScore": 0.98,
  "files": [
    { "path": "Cargo.toml", "type": "manifest" },
    { "path": "src/lib.rs", "type": "source", "lines": 245 },
    { "path": "src/types.rs", "type": "source", "lines": 89 },
    { "path": "src/handlers.rs", "type": "source", "lines": 312 },
    { "path": "tests/integration_test.rs", "type": "test", "lines": 156 }
  ],
  "iterations": 2,
  "specSummary": {
    "name": "my-api",
    "types": 5,
    "functions": 12,
    "endpoints": 8
  }
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>What Gets Generated</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Component</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><strong>Project Manifest</strong></td>
            <td>Cargo.toml, package.json, go.mod, etc. with dependencies</td>
          </tr>
          <tr>
            <td><strong>Type Definitions</strong></td>
            <td>Structs, classes, interfaces matching spec types</td>
          </tr>
          <tr>
            <td><strong>Implementation</strong></td>
            <td>Functions, methods, handlers with full logic</td>
          </tr>
          <tr>
            <td><strong>Tests</strong></td>
            <td>Unit and integration tests covering spec requirements</td>
          </tr>
          <tr>
            <td><strong>Configuration</strong></td>
            <td>Config files, env templates if specified in spec</td>
          </tr>
        </tbody>
      </table>

      <h2>Use Cases</h2>
      <ul>
        <li>
          <strong>Rapid prototyping</strong> - Generate working code from a
          spec in seconds
        </li>
        <li>
          <strong>Multi-language libraries</strong> - Create consistent
          implementations across all supported languages
        </li>
        <li>
          <strong>API clients</strong> - Generate type-safe clients from
          API specifications
        </li>
        <li>
          <strong>Microservices</strong> - Bootstrap new services from
          architectural specs
        </li>
      </ul>

      <h2>Generation Workflow</h2>
      <ol className="list-decimal list-inside space-y-2">
        <li>Parse the specification file to understand requirements</li>
        <li>Load language conventions and idioms</li>
        <li>Generate initial source code following conventions</li>
        <li>Validate generated code against spec requirements</li>
        <li>If parity is below threshold, iterate with fixes</li>
        <li>Return completed project with parity report</li>
      </ol>
    </div>
  );
}
