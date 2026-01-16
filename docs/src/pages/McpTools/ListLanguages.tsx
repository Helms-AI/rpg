export default function ListLanguages() {
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
      <pre className="code-block">
{`// In your AI assistant
const result = await mcp.callTool('rpg', 'list_languages', {});
console.log(result.languages);`}
      </pre>

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
