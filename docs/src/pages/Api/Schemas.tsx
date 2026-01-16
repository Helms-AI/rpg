export default function Schemas() {
  return (
    <div className="prose-docs">
      <h1>Data Schemas</h1>
      <p className="lead">
        JSON schemas for RPG's core data structures and tool responses.
      </p>

      <h2>Spec Schema</h2>
      <p>The main specification structure returned by <code>parse_spec</code>:</p>

      <pre className="code-block">
{`{
  "name": "string",              // Spec name from H1 heading
  "meta": {
    "version": "string",         // Semantic version
    "author": "string",          // Author name
    "license": "string"          // License identifier (optional)
  },
  "targetLanguages": ["string"], // Array of language IDs
  "types": [TypeDef],            // Type definitions
  "functions": [Function],       // Function definitions
  "tests": [TestCase],           // Test cases
  "dependencies": ["string"],    // External dependencies
  "configuration": [ConfigVar]   // Environment variables
}`}
      </pre>

      <h2>TypeDef Schema</h2>
      <p>Type definitions from the <code>## Types</code> section:</p>

      <pre className="code-block">
{`{
  "name": "string",              // Type name (PascalCase)
  "description": "string",       // What the type represents
  "fields": [
    {
      "name": "string",          // Field name (camelCase)
      "type": "string",          // Pseudo-type (string, number, etc.)
      "description": "string",   // Field description (optional)
      "optional": boolean        // Whether field is optional
    }
  ]
}`}
      </pre>

      <h2>Function Schema</h2>
      <p>Function definitions from the <code>## Functions</code> section:</p>

      <pre className="code-block">
{`{
  "name": "string",              // Function name
  "description": "string",       // What the function does
  "accepts": [
    {
      "name": "string",          // Parameter name
      "type": "string",          // Parameter type
      "description": "string"    // Parameter description
    }
  ],
  "returns": [
    {
      "type": "string",          // Return type
      "description": "string"    // Return description
    }
  ],
  "logic": ["string"],           // Step-by-step logic
  "errors": ["string"]           // Possible error conditions
}`}
      </pre>

      <h2>TestCase Schema</h2>
      <p>Test cases from the <code>## Tests</code> section:</p>

      <pre className="code-block">
{`{
  "name": "string",              // Test name
  "given": "string",             // Input/setup description
  "expect": "string"             // Expected output/behavior
}`}
      </pre>

      <h2>ConfigVar Schema</h2>
      <p>Configuration variables from <code>## Configuration</code>:</p>

      <pre className="code-block">
{`{
  "name": "string",              // Variable name (e.g., PORT)
  "description": "string",       // What it configures
  "default": "string",           // Default value (optional)
  "required": boolean            // Whether it must be set
}`}
      </pre>

      <h2>Language Schema</h2>
      <p>Language information from <code>list_languages</code>:</p>

      <pre className="code-block">
{`{
  "id": "string",                // Language identifier (go, rust, etc.)
  "name": "string",              // Display name (Go, Rust, etc.)
  "fileExtension": "string",     // Primary file extension (.go, .rs)
  "namingConventions": {
    "functions": "string",       // Function naming style
    "types": "string",           // Type naming style
    "constants": "string",       // Constant naming style
    "packages": "string"         // Package naming style
  },
  "errorHandling": "string",     // Error handling pattern
  "testFramework": "string"      // Testing framework used
}`}
      </pre>

      <h2>ProjectFile Schema</h2>
      <p>File structure from <code>get_project_structure</code>:</p>

      <pre className="code-block">
{`{
  "path": "string",              // Relative file path
  "purpose": "string"            // What the file is for
}`}
      </pre>

      <h2>ParityReport Schema</h2>
      <p>Parity check results from <code>ensure_parity</code>:</p>

      <pre className="code-block">
{`{
  "reference": {
    "language": "string",        // Reference language ID
    "path": "string"             // Reference project path
  },
  "comparisons": [
    {
      "language": "string",      // Compared language ID
      "path": "string",          // Compared project path
      "status": "string",        // "complete" | "partial" | "missing"
      "missing": [
        {
          "feature": "string",   // Missing feature name
          "suggestion": "string" // How to add it
        }
      ],
      "differences": [
        {
          "feature": "string",   // Different feature
          "reference": "string", // How reference implements it
          "current": "string",   // How this impl does it
          "suggestion": "string" // How to align
        }
      ]
    }
  ],
  "overallParity": number        // Percentage (0-100)
}`}
      </pre>

      <h2>Pseudo-Types</h2>
      <p>
        RPG uses pseudo-types in specs that get mapped to language-specific types:
      </p>

      <table className="w-full">
        <thead>
          <tr>
            <th>Pseudo-Type</th>
            <th>Go</th>
            <th>Rust</th>
            <th>TypeScript</th>
            <th>Python</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>string</code></td>
            <td>string</td>
            <td>String</td>
            <td>string</td>
            <td>str</td>
          </tr>
          <tr>
            <td><code>number</code></td>
            <td>float64</td>
            <td>f64</td>
            <td>number</td>
            <td>float</td>
          </tr>
          <tr>
            <td><code>integer</code></td>
            <td>int</td>
            <td>i64</td>
            <td>number</td>
            <td>int</td>
          </tr>
          <tr>
            <td><code>boolean</code></td>
            <td>bool</td>
            <td>bool</td>
            <td>boolean</td>
            <td>bool</td>
          </tr>
          <tr>
            <td><code>datetime</code></td>
            <td>time.Time</td>
            <td>DateTime</td>
            <td>Date</td>
            <td>datetime</td>
          </tr>
          <tr>
            <td><code>error</code></td>
            <td>error</td>
            <td>Error</td>
            <td>Error</td>
            <td>Exception</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
