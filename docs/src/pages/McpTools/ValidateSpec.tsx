export default function ValidateSpec() {
  return (
    <div className="prose-docs">
      <h1>validate_spec</h1>
      <p className="lead">Check if a spec file exists and contains valid content.</p>

      <h2>Description</h2>
      <p>
        Validates that a specification file exists at the given path and contains meaningful content.
        Use this tool before attempting to parse or generate code from a spec to ensure
        the file is accessible and non-empty.
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
            <td>Path to the specification file to validate</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns a validation result object:</p>

      <pre className="code-block">
{`// Success
{
  "valid": true,
  "message": "Spec file is valid and contains content"
}

// Failure
{
  "valid": false,
  "message": "File not found: ./specs/missing.spec.md"
}`}
      </pre>

      <h2>Example Usage</h2>
      <pre className="code-block">
{`// Validate before parsing
const validation = await mcp.callTool('rpg', 'validate_spec', {
  specPath: './specs/my-project.spec.md'
});

if (validation.valid) {
  // Safe to parse
  const spec = await mcp.callTool('rpg', 'parse_spec', {
    specPath: './specs/my-project.spec.md'
  });
} else {
  console.error('Invalid spec:', validation.message);
}`}
      </pre>

      <h2>Validation Checks</h2>
      <ul>
        <li>File exists at the specified path</li>
        <li>File is readable</li>
        <li>File contains non-empty content</li>
      </ul>
    </div>
  );
}
