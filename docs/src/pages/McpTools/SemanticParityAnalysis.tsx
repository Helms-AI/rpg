import Tabs from '../../components/ui/Tabs';

export default function SemanticParityAnalysis() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              Compare my Go source code with the TypeScript and Rust ports. I need a detailed analysis of what's different.
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll perform a semantic parity analysis comparing your implementations.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using semantic_parity_analysis with sourcePath: "./go-project", generatedProjects: [typescript, rust]
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Analysis complete. Here are the parity scores:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li>TypeScript: <strong>94%</strong> parity (missing 2 edge cases)</li>
              <li>Rust: <strong>87%</strong> parity (missing error variants, async handling)</li>
            </ul>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              The main gaps in Rust are:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li><code>ValidationError</code> missing <code>FieldTooLong</code> variant</li>
              <li><code>process_batch</code> not using async streams</li>
              <li>Test for concurrent access not implemented</li>
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
{`// Compare source against generated implementations
const analysis = await mcp.callTool('rpg', 'semantic_parity_analysis', {
  sourcePath: './original-go-project',
  sourceLanguage: 'go',
  generatedProjects: [
    { language: 'typescript', path: './generated/typescript' },
    { language: 'rust', path: './generated/rust' },
    { language: 'python', path: './generated/python' }
  ]
});

// Custom weights for comparison priorities
const weightedAnalysis = await mcp.callTool('rpg', 'semantic_parity_analysis', {
  sourcePath: './source',
  generatedProjects: [
    { language: 'typescript', path: './ts-impl' }
  ],
  comparisonWeights: {
    type: 0.3,        // Type definitions match
    structural: 0.2,  // Function signatures
    behavioral: 0.3,  // Logic and edge cases
    test: 0.15,       // Test coverage
    idiomatic: 0.05   // Language idioms
  }
});`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>semantic_parity_analysis</h1>
      <p className="lead">Deep semantic comparison between source and generated implementations.</p>

      <h2>Description</h2>
      <p>
        Performs AST-based semantic analysis to compare source code against one or more
        generated implementations. Goes beyond textual comparison to understand type
        equivalence, behavioral parity, and idiomatic differences across languages.
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
            <td><code>sourcePath</code></td>
            <td>string</td>
            <td>Yes</td>
            <td>Path to the original source code</td>
          </tr>
          <tr>
            <td><code>sourceLanguage</code></td>
            <td>string</td>
            <td>No</td>
            <td>Language of source code (auto-detected if not specified)</td>
          </tr>
          <tr>
            <td><code>generatedProjects</code></td>
            <td>array</td>
            <td>Yes</td>
            <td>List of generated projects to compare</td>
          </tr>
          <tr>
            <td><code>comparisonWeights</code></td>
            <td>object</td>
            <td>No</td>
            <td>Custom weights for comparison categories</td>
          </tr>
        </tbody>
      </table>

      <h3>generatedProjects Schema</h3>
      <pre className="code-block">
{`[
  {
    "language": "typescript",  // Target language
    "path": "./generated/ts"   // Path to generated code
  }
]`}
      </pre>

      <h3>comparisonWeights Schema</h3>
      <pre className="code-block">
{`{
  "type": 0.25,        // Type definitions (structs, classes)
  "structural": 0.20,  // Function signatures
  "behavioral": 0.30,  // Logic, edge cases, error handling
  "test": 0.15,        // Test coverage parity
  "idiomatic": 0.10    // Language-specific patterns
}`}
      </pre>

      <h2>Response</h2>
      <p>Returns detailed parity analysis:</p>

      <pre className="code-block">
{`{
  "sourceLanguage": "go",
  "sourceSummary": {
    "types": 12,
    "functions": 34,
    "testCases": 45
  },
  "comparisons": [
    {
      "language": "typescript",
      "path": "./generated/ts",
      "overallParity": 0.94,
      "scores": {
        "type": 0.98,
        "structural": 0.95,
        "behavioral": 0.91,
        "test": 0.92,
        "idiomatic": 0.96
      },
      "gaps": [
        {
          "category": "behavioral",
          "description": "Missing edge case for empty input",
          "sourceLocation": "validator.go:45",
          "severity": "medium",
          "fix": "Add check for empty string in Validate()"
        }
      ],
      "suggestions": [
        "Consider using discriminated unions for error types"
      ]
    }
  ],
  "fixInstructions": "## TypeScript Fixes\\n\\n### Missing Edge Cases..."
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Comparison Categories</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Category</th>
            <th>What It Checks</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><strong>Type Parity</strong></td>
            <td>Struct/class fields match, enums have same variants, generics preserved</td>
          </tr>
          <tr>
            <td><strong>Structural</strong></td>
            <td>Function signatures, parameter types, return types, visibility</td>
          </tr>
          <tr>
            <td><strong>Behavioral</strong></td>
            <td>Error handling, edge cases, validation logic, business rules</td>
          </tr>
          <tr>
            <td><strong>Test Coverage</strong></td>
            <td>Test cases exist, same scenarios covered, assertions equivalent</td>
          </tr>
          <tr>
            <td><strong>Idiomatic</strong></td>
            <td>Uses language-appropriate patterns, naming conventions, error handling</td>
          </tr>
        </tbody>
      </table>

      <h2>Use Cases</h2>
      <ul>
        <li>
          <strong>Quality assurance</strong> - Verify generated code matches
          source behavior before deployment
        </li>
        <li>
          <strong>Incremental porting</strong> - Track progress when porting
          large codebases
        </li>
        <li>
          <strong>Regression detection</strong> - Ensure spec changes are
          reflected in all implementations
        </li>
        <li>
          <strong>Documentation</strong> - Generate reports on implementation
          differences across languages
        </li>
      </ul>
    </div>
  );
}
