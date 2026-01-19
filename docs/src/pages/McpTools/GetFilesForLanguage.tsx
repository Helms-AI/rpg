import Tabs from '../../components/ui/Tabs';

export default function GetFilesForLanguage() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              I need to understand the SQL migrations and Protobuf schemas in this project. Can you analyze them?
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll get the raw files for AI interpretation since these don't have semantic parsers.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using get_files_for_language with language: "sql"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Found 15 SQL migration files. Analyzing the schema:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li><strong>users</strong> table - id, email, name, created_at</li>
              <li><strong>orders</strong> table - id, user_id (FK), total, status</li>
              <li><strong>products</strong> table - id, name, price, inventory</li>
              <li>Indexes on user_id, email, created_at</li>
            </ul>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Now let me get the Protobuf files...
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300 mt-2">
              Using get_files_for_language with language: "protobuf"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Found 8 .proto files defining your API:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li><code>user.proto</code> - User, UserRequest, UserResponse</li>
              <li><code>order.proto</code> - Order, OrderItem, OrderStatus enum</li>
              <li><code>service.proto</code> - gRPC service definitions</li>
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
{`// Get SQL files for AI analysis
const sqlFiles = await mcp.callTool('rpg', 'get_files_for_language', {
  sourcePath: './my-project',
  language: 'sql'
});

// Get Protobuf schemas
const protoFiles = await mcp.callTool('rpg', 'get_files_for_language', {
  sourcePath: './my-project',
  language: 'protobuf'
});

// Limit files for large projects
const limitedFiles = await mcp.callTool('rpg', 'get_files_for_language', {
  sourcePath: './large-project',
  language: 'graphql',
  maxFiles: 20,
  includeTests: false
});

// The response includes an AI prompt template
// Use this to guide your analysis of the raw files
console.log(sqlFiles.aiPrompt);
console.log(sqlFiles.files);`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>get_files_for_language</h1>
      <p className="lead">Get raw file contents for AI-driven analysis of languages without semantic parsers.</p>

      <h2>Description</h2>
      <p>
        Returns raw file contents for languages that don't have semantic parsers, along with
        an AI prompt template for extracting types, functions, and patterns. Use this for
        SQL, Protobuf, GraphQL, YAML configs, or when you need actual source code for any
        language.
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
            <td>Yes</td>
            <td>Language to get files for (sql, protobuf, graphql, etc.)</td>
          </tr>
          <tr>
            <td><code>maxFiles</code></td>
            <td>number</td>
            <td>No</td>
            <td>Maximum files to return (default: 50)</td>
          </tr>
          <tr>
            <td><code>includeTests</code></td>
            <td>boolean</td>
            <td>No</td>
            <td>Include test files (default: false)</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns files with AI analysis prompt:</p>

      <pre className="code-block">
{`{
  "language": "sql",
  "sourcePath": "./my-project",
  "files": [
    {
      "path": "migrations/001_create_users.sql",
      "content": "CREATE TABLE users (\\n  id UUID PRIMARY KEY,\\n  email VARCHAR(255)...",
      "lines": 25
    },
    {
      "path": "migrations/002_create_orders.sql",
      "content": "CREATE TABLE orders (\\n  id UUID PRIMARY KEY,\\n  user_id UUID REFERENCES...",
      "lines": 35
    }
  ],
  "totalFiles": 15,
  "totalLines": 450,
  "aiPrompt": "Analyze the following SQL files to extract:\\n\\n1. Table definitions (name, columns, types)\\n2. Relationships (foreign keys, references)\\n3. Indexes and constraints\\n4. Common patterns and conventions\\n\\n## Files\\n\\n..."
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Supported Languages</h2>
      <p>This tool is particularly useful for:</p>
      <table className="w-full">
        <thead>
          <tr>
            <th>Language</th>
            <th>Extensions</th>
            <th>What to Extract</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>SQL</td>
            <td>.sql</td>
            <td>Tables, columns, relationships, indexes</td>
          </tr>
          <tr>
            <td>Protobuf</td>
            <td>.proto</td>
            <td>Messages, services, enums, options</td>
          </tr>
          <tr>
            <td>GraphQL</td>
            <td>.graphql, .gql</td>
            <td>Types, queries, mutations, subscriptions</td>
          </tr>
          <tr>
            <td>OpenAPI</td>
            <td>.yaml, .json</td>
            <td>Endpoints, schemas, parameters</td>
          </tr>
          <tr>
            <td>Terraform</td>
            <td>.tf</td>
            <td>Resources, variables, outputs</td>
          </tr>
          <tr>
            <td>Docker</td>
            <td>Dockerfile</td>
            <td>Build steps, dependencies, configuration</td>
          </tr>
        </tbody>
      </table>

      <h2>Use Cases</h2>
      <ul>
        <li>
          <strong>Schema analysis</strong> - Understand database structure from
          migration files
        </li>
        <li>
          <strong>API discovery</strong> - Extract API contracts from Protobuf
          or GraphQL schemas
        </li>
        <li>
          <strong>Config understanding</strong> - Parse infrastructure-as-code
          files for documentation
        </li>
        <li>
          <strong>Cross-language types</strong> - Generate consistent types
          from schema files
        </li>
      </ul>

      <h2>AI Prompt Template</h2>
      <p>
        The response includes an <code>aiPrompt</code> field with a template for analyzing
        the files. The prompt is tailored to the language and suggests what patterns and
        structures to extract.
      </p>
    </div>
  );
}
