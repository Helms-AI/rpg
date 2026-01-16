export default function GetGenerationContext() {
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
