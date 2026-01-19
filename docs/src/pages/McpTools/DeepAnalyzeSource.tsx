import Tabs from '../../components/ui/Tabs';

export default function DeepAnalyzeSource() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              I need to understand the type system in my Go project before porting it. Can you do a deep analysis?
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll perform a deep semantic analysis of your Go code using AST parsing.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using deep_analyze_source with sourcePath: "./my-project", language: "go"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              I've completed the semantic analysis. Here's what I found:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li><strong>12 structs</strong> with 45 fields total</li>
              <li><strong>8 interfaces</strong> defining 23 methods</li>
              <li><strong>34 functions</strong> with full signatures</li>
              <li><strong>Call graph</strong> showing dependencies between packages</li>
            </ul>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              The main complexity is in the <code>UserService</code> struct which implements 3 interfaces
              and has dependencies on the database and cache layers.
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
{`// Deep analysis with Go's native AST parser
const analysis = await mcp.callTool('rpg', 'deep_analyze_source', {
  sourcePath: './my-go-project',
  language: 'go',
  analysisDepth: 'full',
  includeTests: true
});

// Analysis returns structured semantic data
console.log(analysis.types);     // Type definitions
console.log(analysis.functions); // Function signatures
console.log(analysis.callGraph); // Dependency relationships

// TypeScript analysis
const tsAnalysis = await mcp.callTool('rpg', 'deep_analyze_source', {
  sourcePath: './my-ts-project',
  language: 'typescript'
});`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>deep_analyze_source</h1>
      <p className="lead">Perform deep semantic analysis using AST parsing and type resolution.</p>

      <h2>Description</h2>
      <p>
        Performs AST-based semantic analysis on source code to extract structured information
        about types, functions, interfaces, and their relationships. Go uses native go/ast
        for precise parsing; other languages use pattern-based analysis with varying depth.
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
            <td>Path to the source code directory</td>
          </tr>
          <tr>
            <td><code>language</code></td>
            <td>string</td>
            <td>No</td>
            <td>Target language (auto-detected if not specified)</td>
          </tr>
          <tr>
            <td><code>analysisDepth</code></td>
            <td>string</td>
            <td>No</td>
            <td>Analysis depth: "basic", "standard", or "full" (default: "standard")</td>
          </tr>
          <tr>
            <td><code>includeTests</code></td>
            <td>boolean</td>
            <td>No</td>
            <td>Include test files in analysis (default: false)</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns structured semantic analysis:</p>

      <pre className="code-block">
{`{
  "language": "go",
  "analysisDepth": "full",
  "types": [
    {
      "name": "UserService",
      "kind": "struct",
      "fields": [
        { "name": "db", "type": "*sql.DB", "tags": "json:\\"db\\"" },
        { "name": "cache", "type": "Cache", "tags": "" }
      ],
      "methods": ["Create", "Get", "Update", "Delete"],
      "implements": ["Service", "Cacheable"]
    }
  ],
  "functions": [
    {
      "name": "NewUserService",
      "receiver": "",
      "params": [{ "name": "db", "type": "*sql.DB" }],
      "returns": ["*UserService", "error"],
      "exported": true
    }
  ],
  "interfaces": [
    {
      "name": "Service",
      "methods": [
        { "name": "Create", "params": ["ctx", "data"], "returns": ["error"] }
      ]
    }
  ],
  "callGraph": {
    "UserService.Create": ["db.Exec", "cache.Set", "log.Info"]
  },
  "dependencies": ["database/sql", "encoding/json"]
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Supported Languages</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Language</th>
            <th>Parser</th>
            <th>Semantic Depth</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Go</td>
            <td>Native go/ast</td>
            <td>Full (types, interfaces, call graphs)</td>
          </tr>
          <tr>
            <td>TypeScript</td>
            <td>Pattern-based</td>
            <td>High (types, interfaces, generics)</td>
          </tr>
          <tr>
            <td>Python</td>
            <td>Pattern-based</td>
            <td>Medium (classes, functions, type hints)</td>
          </tr>
          <tr>
            <td>Java</td>
            <td>Pattern-based</td>
            <td>Medium (classes, interfaces, generics)</td>
          </tr>
          <tr>
            <td>Rust</td>
            <td>Pattern-based</td>
            <td>Medium (structs, traits, impls)</td>
          </tr>
          <tr>
            <td>C#</td>
            <td>Pattern-based</td>
            <td>Medium (classes, interfaces, records)</td>
          </tr>
        </tbody>
      </table>

      <h2>Use Cases</h2>
      <ul>
        <li>
          <strong>Pre-port analysis</strong> - Understand type relationships before
          porting to another language
        </li>
        <li>
          <strong>Documentation generation</strong> - Extract API surface for
          generating documentation
        </li>
        <li>
          <strong>Refactoring planning</strong> - Identify dependencies and coupling
          before making changes
        </li>
        <li>
          <strong>Code review</strong> - Quickly understand the structure of
          unfamiliar code
        </li>
      </ul>
    </div>
  );
}
