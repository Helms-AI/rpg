export default function ImportSpecFromSource() {
  return (
    <div className="prose-docs">
      <h1>import_spec_from_source</h1>
      <p className="lead">Analyze source code for AI-powered spec generation.</p>

      <h2>Description</h2>
      <p>
        Collects and analyzes source code from any directory to help generate a specification.
        Returns an analysis prompt containing all source files, tests, API specs, and configurations.
        The AI uses this context to generate a comprehensive .spec.md file.
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
            <td><code>inputPath</code></td>
            <td>string</td>
            <td>Yes</td>
            <td>Path to the source code directory to analyze</td>
          </tr>
          <tr>
            <td><code>name</code></td>
            <td>string</td>
            <td>No</td>
            <td>Optional name for the generated spec</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns an analysis prompt for spec generation:</p>

      <pre className="code-block">
{`{
  "analysisPrompt": "Based on the following source code analysis...\\n\\n## Source Files\\n...",
  "sourceFiles": [
    { "path": "src/main.go", "language": "go", "lines": 150 },
    { "path": "src/handler.go", "language": "go", "lines": 200 }
  ],
  "detectedLanguage": "go",
  "suggestedName": "url-shortener"
}`}
      </pre>

      <h2>Example Usage</h2>
      <pre className="code-block">
{`// Reverse-engineer a spec from existing code
const analysis = await mcp.callTool('rpg', 'import_spec_from_source', {
  inputPath: './existing-project',
  name: 'my-service'
});

// The AI can use analysisPrompt to generate a spec
// The spec can then be used to generate code in other languages`}
      </pre>

      <h2>Use Cases</h2>
      <ul>
        <li>
          <strong>Port to new language</strong> - Analyze existing Go code, generate spec,
          then generate TypeScript implementation
        </li>
        <li>
          <strong>Documentation</strong> - Generate specification documentation from
          undocumented legacy code
        </li>
        <li>
          <strong>Refactoring</strong> - Create a spec to understand current behavior
          before making changes
        </li>
        <li>
          <strong>Multi-language libraries</strong> - Start with one implementation,
          generate spec, create versions in all supported languages
        </li>
      </ul>

      <h2>Analyzed Content</h2>
      <p>The importer examines:</p>
      <ul>
        <li>Source code files (functions, types, interfaces)</li>
        <li>Test files (test cases, expected behaviors)</li>
        <li>Configuration files (environment variables, settings)</li>
        <li>API specifications (OpenAPI, GraphQL schemas)</li>
        <li>Package manifests (dependencies, versions)</li>
      </ul>
    </div>
  );
}
