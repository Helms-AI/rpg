import Tabs from '../../components/ui/Tabs';

export default function IterativeRefinementLoop() {
  const exampleTabs = [
    {
      id: 'conversation',
      label: 'Conversation',
      content: (
        <div className="space-y-4">
          <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <div className="text-xs text-blue-600 dark:text-blue-400 font-semibold mb-1">User</div>
            <p className="text-gray-800 dark:text-gray-200">
              Port my Go API to TypeScript and keep refining until it's at least 95% parity.
            </p>
          </div>
          <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="text-xs text-gray-600 dark:text-gray-400 font-semibold mb-1">Assistant</div>
            <p className="text-gray-800 dark:text-gray-200 mb-3">
              I'll start an iterative refinement loop to achieve 95% parity.
            </p>
            <div className="text-sm bg-gray-100 dark:bg-gray-900 p-3 rounded font-mono text-gray-700 dark:text-gray-300">
              Using iterative_refinement_loop with convergenceThreshold: 0.95, targetLanguages: ["typescript"]
            </div>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Refinement loop progress:
            </p>
            <ul className="list-disc list-inside text-gray-800 dark:text-gray-200 mt-2 space-y-1">
              <li>Iteration 1: 78% parity - Generated initial TypeScript code</li>
              <li>Iteration 2: 89% parity - Added missing error types</li>
              <li>Iteration 3: 93% parity - Fixed async handling</li>
              <li>Iteration 4: <strong>96% parity</strong> - Added edge case tests</li>
            </ul>
            <p className="text-gray-800 dark:text-gray-200 mt-3">
              Reached 96% parity after 4 iterations. The TypeScript implementation is ready
              at <code>./output/my-api/typescript/</code>.
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
{`// Basic refinement loop
const result = await mcp.callTool('rpg', 'iterative_refinement_loop', {
  sourcePath: './my-go-project',
  targetLanguages: ['typescript'],
  outputDir: './output',
  convergenceThreshold: 0.95
});

// Multi-language refinement with custom settings
const multiLang = await mcp.callTool('rpg', 'iterative_refinement_loop', {
  sourcePath: './source',
  sourceLanguage: 'go',
  targetLanguages: ['typescript', 'rust', 'python'],
  outputDir: './generated',
  convergenceThreshold: 0.90,
  maxIterations: 10,
  refinementStrategy: 'aggressive'
});

// With spec file
const specBased = await mcp.callTool('rpg', 'iterative_refinement_loop', {
  sourcePath: './reference-impl',
  specPath: './my-api.spec.md',
  targetLanguages: ['java'],
  outputDir: './output'
});`}
        </pre>
      ),
    },
  ];

  return (
    <div className="prose-docs">
      <h1>iterative_refinement_loop</h1>
      <p className="lead">Automated refinement loop to achieve maximum parity between source and generated code.</p>

      <h2>Description</h2>
      <p>
        Orchestrates an iterative process that analyzes source code, generates implementations,
        compares parity, and refines until a convergence threshold is reached. Automates the
        generate-compare-fix cycle for production-quality code generation.
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
            <td>Path to source code to use as reference</td>
          </tr>
          <tr>
            <td><code>targetLanguages</code></td>
            <td>string[]</td>
            <td>Yes</td>
            <td>Languages to generate (e.g., ["typescript", "rust"])</td>
          </tr>
          <tr>
            <td><code>outputDir</code></td>
            <td>string</td>
            <td>Yes</td>
            <td>Directory for generated code</td>
          </tr>
          <tr>
            <td><code>sourceLanguage</code></td>
            <td>string</td>
            <td>No</td>
            <td>Source language (auto-detected if not specified)</td>
          </tr>
          <tr>
            <td><code>specPath</code></td>
            <td>string</td>
            <td>No</td>
            <td>Optional spec file for additional context</td>
          </tr>
          <tr>
            <td><code>convergenceThreshold</code></td>
            <td>number</td>
            <td>No</td>
            <td>Target parity score 0.0-1.0 (default: 0.90)</td>
          </tr>
          <tr>
            <td><code>maxIterations</code></td>
            <td>number</td>
            <td>No</td>
            <td>Maximum iterations before stopping (default: 5)</td>
          </tr>
          <tr>
            <td><code>refinementStrategy</code></td>
            <td>string</td>
            <td>No</td>
            <td>"conservative", "balanced", or "aggressive"</td>
          </tr>
        </tbody>
      </table>

      <h2>Response</h2>
      <p>Returns refinement results:</p>

      <pre className="code-block">
{`{
  "converged": true,
  "iterations": 4,
  "results": [
    {
      "language": "typescript",
      "finalParity": 0.96,
      "outputPath": "./output/my-api/typescript",
      "iterationHistory": [
        { "iteration": 1, "parity": 0.78, "gapsFixed": 0 },
        { "iteration": 2, "parity": 0.89, "gapsFixed": 5 },
        { "iteration": 3, "parity": 0.93, "gapsFixed": 3 },
        { "iteration": 4, "parity": 0.96, "gapsFixed": 2 }
      ],
      "remainingGaps": [
        { "description": "Optional: Add JSDoc comments", "severity": "low" }
      ]
    }
  ],
  "totalDuration": "45s"
}`}
      </pre>

      <h2>Example Usage</h2>
      <Tabs tabs={exampleTabs} defaultTab="conversation" />

      <h2>Refinement Strategies</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Strategy</th>
            <th>Behavior</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><strong>conservative</strong></td>
            <td>Fix high-severity gaps first, minimal changes per iteration</td>
          </tr>
          <tr>
            <td><strong>balanced</strong></td>
            <td>Fix gaps by priority, moderate changes per iteration (default)</td>
          </tr>
          <tr>
            <td><strong>aggressive</strong></td>
            <td>Fix all identified gaps each iteration, faster convergence</td>
          </tr>
        </tbody>
      </table>

      <h2>Use Cases</h2>
      <ul>
        <li>
          <strong>Automated porting</strong> - Port entire codebases with minimal
          manual intervention
        </li>
        <li>
          <strong>CI/CD integration</strong> - Run as part of build pipeline to
          maintain multi-language parity
        </li>
        <li>
          <strong>SDK generation</strong> - Generate and refine SDKs in multiple
          languages from a reference implementation
        </li>
        <li>
          <strong>Quality gates</strong> - Ensure generated code meets minimum
          parity threshold before release
        </li>
      </ul>

      <h2>Stopping Conditions</h2>
      <p>The loop stops when any of these conditions are met:</p>
      <ul>
        <li>All targets reach <code>convergenceThreshold</code></li>
        <li><code>maxIterations</code> is reached</li>
        <li>No parity improvement for 2 consecutive iterations</li>
        <li>An unrecoverable error occurs</li>
      </ul>
    </div>
  );
}
