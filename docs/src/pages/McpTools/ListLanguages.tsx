import Tabs from '../../components/ui/Tabs';

export default function ListLanguages() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              What programming languages can RPG generate code for?
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              Let me check the supported languages.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using list_languages
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              RPG supports 6 languages:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li><strong>Go</strong> - Error handling via multiple returns</li>
              <li><strong>Rust</strong> - Result/Option types, ownership</li>
              <li><strong>Java</strong> - Records, streams, modern Java 17+</li>
              <li><strong>C#</strong> - .NET 6+, nullable reference types</li>
              <li><strong>Python</strong> - Type hints, Python 3.10+</li>
              <li><strong>TypeScript</strong> - Strict mode, full type safety</li>
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
{`// In your AI assistant
const result = await mcp.callTool('rpg', 'list_languages', {});
console.log(result.languages);`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>list_languages</h1>
      <p className="lead">List all supported target languages with their conventions and idioms.</p>

      <h2>Description</h2>
      <p>
        Returns a comprehensive list of all programming languages that RPG can generate code for,
        including each language's naming conventions, error handling patterns, and idiomatic practices.
      </p>

      <h2>Parameters</h2>
      <p>This tool takes no parameters.</p>

      <h2>Response</h2>
      <p>Returns an array of language objects with the following structure:</p>

      <pre className="code-block">
{`{
  "languages": [
    {
      "id": "go",
      "name": "Go",
      "fileExtension": ".go",
      "namingConventions": {
        "functions": "camelCase",
        "types": "PascalCase",
        "constants": "PascalCase or ALL_CAPS",
        "packages": "lowercase"
      },
      "errorHandling": "Multiple return values with error type",
      "testFramework": "testing package (built-in)"
    },
    // ... more languages
  ]
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Supported Languages</h2>
      <ul>
        <li><strong>Go</strong> - Idiomatic Go with error handling via multiple returns</li>
        <li><strong>Rust</strong> - Safe Rust with Result/Option types</li>
        <li><strong>Java</strong> - Modern Java with records and streams</li>
        <li><strong>C#</strong> - .NET 6+ with nullable reference types</li>
        <li><strong>Python</strong> - Type-hinted Python 3.10+</li>
        <li><strong>TypeScript</strong> - Strict TypeScript with full type safety</li>
      </ul>
    </div>
  );
}
