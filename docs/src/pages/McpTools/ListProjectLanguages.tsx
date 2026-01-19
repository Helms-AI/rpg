import Tabs from '../../components/ui/Tabs';

export default function ListProjectLanguages() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              What languages are used in this monorepo? I need to know before analyzing it.
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll scan the project to detect all programming languages.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using list_project_languages with sourcePath: "./monorepo"
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Found 5 languages in your monorepo:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li><strong>Go</strong> - 45 files (backend services) - semantic parser available</li>
              <li><strong>TypeScript</strong> - 120 files (frontend) - semantic parser available</li>
              <li><strong>Python</strong> - 23 files (scripts) - semantic parser available</li>
              <li><strong>SQL</strong> - 15 files (migrations) - use AI interpretation</li>
              <li><strong>Protobuf</strong> - 8 files (API schemas) - use AI interpretation</li>
            </ul>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              For Go, TypeScript, and Python, I can use <code>deep_analyze_source</code>.
              For SQL and Protobuf, use <code>get_files_for_language</code> for AI interpretation.
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
{`// Scan project for languages
const languages = await mcp.callTool('rpg', 'list_project_languages', {
  sourcePath: './my-project'
});

// Use results to plan analysis strategy
for (const lang of languages.detected) {
  if (lang.hasSemanticParser) {
    // Use deep_analyze_source for semantic analysis
    await mcp.callTool('rpg', 'deep_analyze_source', {
      sourcePath: './my-project',
      language: lang.id
    });
  } else {
    // Use get_files_for_language for AI interpretation
    await mcp.callTool('rpg', 'get_files_for_language', {
      sourcePath: './my-project',
      language: lang.id
    });
  }
}`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>list_project_languages</h1>
      <p className="lead">Detect all programming languages in a project with analysis recommendations.</p>

      <h2>Description</h2>
      <p>
        Scans a source directory to identify all programming languages used, with file counts
        and metadata. Indicates which languages have semantic parsers available (for
        <code>deep_analyze_source</code>) versus which need AI interpretation (via
        <code>get_files_for_language</code>). Use this first to plan analysis strategy.
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
            <td>Path to the source code directory to scan</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns detected languages with metadata:</p>

      <pre className="code-block">
{`{
  "sourcePath": "./monorepo",
  "detected": [
    {
      "id": "go",
      "name": "Go",
      "fileCount": 45,
      "totalLines": 12500,
      "extensions": [".go"],
      "hasSemanticParser": true,
      "analysisRecommendation": "deep_analyze_source"
    },
    {
      "id": "typescript",
      "name": "TypeScript",
      "fileCount": 120,
      "totalLines": 28000,
      "extensions": [".ts", ".tsx"],
      "hasSemanticParser": true,
      "analysisRecommendation": "deep_analyze_source"
    },
    {
      "id": "sql",
      "name": "SQL",
      "fileCount": 15,
      "totalLines": 800,
      "extensions": [".sql"],
      "hasSemanticParser": false,
      "analysisRecommendation": "get_files_for_language"
    },
    {
      "id": "protobuf",
      "name": "Protocol Buffers",
      "fileCount": 8,
      "totalLines": 450,
      "extensions": [".proto"],
      "hasSemanticParser": false,
      "analysisRecommendation": "get_files_for_language"
    }
  ],
  "primaryLanguage": "typescript",
  "summary": {
    "totalFiles": 188,
    "totalLines": 41750,
    "languageCount": 4
  }
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Language Detection</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Language</th>
            <th>Extensions</th>
            <th>Semantic Parser</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Go</td>
            <td>.go</td>
            <td>Yes (native go/ast)</td>
          </tr>
          <tr>
            <td>TypeScript</td>
            <td>.ts, .tsx</td>
            <td>Yes</td>
          </tr>
          <tr>
            <td>JavaScript</td>
            <td>.js, .jsx, .mjs</td>
            <td>Yes</td>
          </tr>
          <tr>
            <td>Python</td>
            <td>.py</td>
            <td>Yes</td>
          </tr>
          <tr>
            <td>Java</td>
            <td>.java</td>
            <td>Yes</td>
          </tr>
          <tr>
            <td>Rust</td>
            <td>.rs</td>
            <td>Yes</td>
          </tr>
          <tr>
            <td>C#</td>
            <td>.cs</td>
            <td>Yes</td>
          </tr>
          <tr>
            <td>SQL</td>
            <td>.sql</td>
            <td>No (AI interpretation)</td>
          </tr>
          <tr>
            <td>Protobuf</td>
            <td>.proto</td>
            <td>No (AI interpretation)</td>
          </tr>
          <tr>
            <td>GraphQL</td>
            <td>.graphql, .gql</td>
            <td>No (AI interpretation)</td>
          </tr>
        </tbody>
      </table>

      <h2>Use Cases</h2>
      <ul>
        <li>
          <strong>Project discovery</strong> - Understand the technology stack
          of an unfamiliar codebase
        </li>
        <li>
          <strong>Analysis planning</strong> - Choose the right analysis tools
          for each language
        </li>
        <li>
          <strong>Monorepo mapping</strong> - Identify all components in a
          multi-language monorepo
        </li>
        <li>
          <strong>Migration planning</strong> - Inventory languages before
          porting or consolidating
        </li>
      </ul>
    </div>
  );
}
