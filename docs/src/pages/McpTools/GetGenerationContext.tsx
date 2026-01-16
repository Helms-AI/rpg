import Tabs from '../../components/ui/Tabs';

export default function GetGenerationContext() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              I have a spec file at ./specs/url-shortener.spec.md. Generate Go code for it.
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll get the generation context for your URL shortener spec targeting Go.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using get_generation_context with specPath: "./specs/url-shortener.spec.md", language: "go"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              I've retrieved the spec content and Go conventions. Based on the spec, I'll now generate
              idiomatic Go code following the project structure conventions...
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
{`// Get context for Go code generation
const context = await mcp.callTool('rpg', 'get_generation_context', {
  specPath: './specs/url-shortener.spec.md',
  language: 'go'
});

// The AI can now use this context to generate code
// The specContent is raw markdown for AI interpretation
// The conventions guide idiomatic code generation`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>get_generation_context</h1>
      <p className="lead">Get full context for code generation including spec and language conventions.</p>

      <h2>Description</h2>
      <p>
        This is the primary tool for code generation. It returns the raw spec content along with
        language-specific conventions and a prompt template. The AI interprets the spec content
        directly to generate idiomatic code in the target language.
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
            <td><code>language</code></td>
            <td>string</td>
            <td>Yes</td>
            <td>Target language ID (go, rust, java, csharp, python, typescript)</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns a comprehensive generation context:</p>

      <pre className="code-block">
{`{
  "specContent": "# URL Shortener\\n\\nA service that...",
  "language": {
    "id": "go",
    "name": "Go",
    "conventions": {
      "namingConventions": {...},
      "errorHandling": "...",
      "testFramework": "..."
    }
  },
  "promptTemplate": "Generate idiomatic Go code..."
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Why Raw Spec Content?</h2>
      <p>
        RPG intentionally returns the raw markdown spec rather than structured data.
        This allows the AI to interpret natural language descriptions, understand context,
        and make intelligent decisions about implementation details - just like a human
        developer reading a design document.
      </p>
    </div>
  );
}
