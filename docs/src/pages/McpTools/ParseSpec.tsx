export default function ParseSpec() {
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
      <pre className="code-block">
{`// Parse a specification file
const result = await mcp.callTool('rpg', 'parse_spec', {
  specPath: './specs/my-project.spec.md'
});

console.log(result.name);        // "my-project"
console.log(result.functions);   // Array of function definitions`}
      </pre>

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
