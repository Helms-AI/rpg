import Tabs from '../../components/ui/Tabs';

export default function GetProjectStructure() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              What's the recommended project structure for a Rust project called "url-shortener"?
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll get the idiomatic Rust project structure for you.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using get_project_structure with projectName: "url-shortener", language: "rust"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Here's the recommended structure for your Rust project:
            </p>
            <pre className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300 mt-2">
{`url-shortener/
├── src/
│   ├── lib.rs          # Library entry point
│   └── main.rs         # Binary entry point
├── tests/
│   └── integration.rs  # Integration tests
├── Cargo.toml          # Package manifest
└── README.md           # Documentation`}
            </pre>
          </div>
        </div>
      ),
    },
    {
      id: 'code',
      label: 'Code',
      content: (
        <pre className="code-block">
{`// Get structure for a Rust project
const structure = await mcp.callTool('rpg', 'get_project_structure', {
  projectName: 'url-shortener',
  language: 'rust'
});

// Create the files
for (const file of structure.files) {
  console.log(\`Create: \${file.path} - \${file.purpose}\`);
}`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>get_project_structure</h1>
      <p className="lead">Get recommended file structure for a project in the target language.</p>

      <h2>Description</h2>
      <p>
        Returns the idiomatic project structure for the specified language.
        This includes recommended directories, file names, and organization patterns
        that follow community conventions for each language ecosystem.
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
            <td><code>projectName</code></td>
            <td>string</td>
            <td>Yes</td>
            <td>Name of the project for directory naming</td>
          </tr>
          <tr>
            <td><code>language</code></td>
            <td>string</td>
            <td>Yes</td>
            <td>Target language ID</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns an array of files with their paths and purposes:</p>

      <pre className="code-block">
{`// Go project structure
{
  "files": [
    { "path": "url-shortener/cmd/main.go", "purpose": "Application entry point" },
    { "path": "url-shortener/internal/shortener/shortener.go", "purpose": "Core logic" },
    { "path": "url-shortener/internal/shortener/shortener_test.go", "purpose": "Unit tests" },
    { "path": "url-shortener/go.mod", "purpose": "Module definition" },
    { "path": "url-shortener/README.md", "purpose": "Documentation" }
  ]
}

// TypeScript project structure
{
  "files": [
    { "path": "url-shortener/src/index.ts", "purpose": "Main export" },
    { "path": "url-shortener/src/shortener.ts", "purpose": "Core logic" },
    { "path": "url-shortener/src/shortener.test.ts", "purpose": "Unit tests" },
    { "path": "url-shortener/package.json", "purpose": "Package definition" },
    { "path": "url-shortener/tsconfig.json", "purpose": "TypeScript config" }
  ]
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Language-Specific Patterns</h2>
      <ul>
        <li><strong>Go</strong>: cmd/, internal/, pkg/ layout with go.mod</li>
        <li><strong>Rust</strong>: src/lib.rs, src/main.rs with Cargo.toml</li>
        <li><strong>Java</strong>: Maven/Gradle structure with src/main/java</li>
        <li><strong>C#</strong>: .NET project with .csproj and namespaces</li>
        <li><strong>Python</strong>: Package with __init__.py and pyproject.toml</li>
        <li><strong>TypeScript</strong>: src/ with package.json and tsconfig.json</li>
      </ul>
    </div>
  );
}
