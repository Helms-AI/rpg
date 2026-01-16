export default function ParityChecking() {
  return (
    <div className="prose-docs">
      <h1>Feature Parity Checking</h1>
      <p className="lead">
        Ensure all your language implementations have consistent behavior using RPG's parity tools.
      </p>

      <div className="not-prose my-6 p-4 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
        <p className="text-blue-800 dark:text-blue-300 text-sm">
          <strong>Time:</strong> 15 minutes &nbsp;|&nbsp; <strong>Difficulty:</strong> Advanced
        </p>
      </div>

      <h2>Why Parity Matters</h2>
      <p>
        When you generate the same spec in multiple languages, each implementation should
        behave identically. However, differences can creep in due to:
      </p>
      <ul>
        <li>Language-specific edge cases</li>
        <li>Manual modifications after generation</li>
        <li>Incomplete regeneration</li>
        <li>Different interpretations of ambiguous spec language</li>
      </ul>

      <h2>Step 1: Set Up Multiple Implementations</h2>
      <p>
        Assuming you have the URL shortener from the previous tutorial in both Go and TypeScript:
      </p>

      <pre className="code-block">
{`generated/
├── go/
│   └── url-shortener/
│       ├── cmd/main.go
│       └── internal/shortener/...
└── typescript/
    └── url-shortener/
        ├── src/index.ts
        └── src/shortener/...`}
      </pre>

      <h2>Step 2: Run Parity Check</h2>
      <p>Use RPG's <code>ensure_parity</code> tool:</p>

      <pre className="code-block">
{`> Check parity between Go and TypeScript implementations of url-shortener

Checking parity for url-shortener.spec.md...

Reference: go (generated/go/url-shortener)

Comparing: typescript (generated/typescript/url-shortener)

┌─────────────────────┬────────┬─────────────────────────────────┐
│ Feature             │ Status │ Details                         │
├─────────────────────┼────────┼─────────────────────────────────┤
│ ShortURL type       │ ✓      │ All fields present              │
│ CreateRequest type  │ ✓      │ All fields present              │
│ CreateResponse type │ ✓      │ All fields present              │
│ createShortUrl      │ ✓      │ Implemented correctly           │
│ getOriginalUrl      │ ⚠      │ Missing click increment         │
│ handleRedirect      │ ✓      │ Implemented correctly           │
│ POST /api/shorten   │ ✓      │ Endpoint working                │
│ GET /{shortCode}    │ ✓      │ Endpoint working                │
└─────────────────────┴────────┴─────────────────────────────────┘

Overall Parity: 87.5%

Issues Found:
1. getOriginalUrl in TypeScript doesn't increment click counter

Suggested Fix:
In src/shortener/shortener.ts, add click increment:

  export function getOriginalUrl(shortCode: string): string {
    const entry = storage.get(shortCode);
    if (!entry) throw new NotFoundError();
+   entry.clicks++;
    return entry.originalUrl;
  }`}
      </pre>

      <h2>Step 3: Fix Discrepancies</h2>
      <p>Apply the suggested fixes to bring implementations into parity:</p>

      <pre className="code-block">
{`> Apply the suggested fix to the TypeScript implementation

Updated src/shortener/shortener.ts
- Added click increment in getOriginalUrl function`}
      </pre>

      <h2>Step 4: Verify Parity</h2>
      <p>Run the check again:</p>

      <pre className="code-block">
{`> Check parity again

Overall Parity: 100%
All features implemented consistently across both languages.`}
      </pre>

      <h2>What Parity Checks</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Category</th>
            <th>What's Checked</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><strong>Types</strong></td>
            <td>All type definitions present with correct fields</td>
          </tr>
          <tr>
            <td><strong>Functions</strong></td>
            <td>All functions implemented with correct signatures</td>
          </tr>
          <tr>
            <td><strong>Logic</strong></td>
            <td>Core algorithm steps are present</td>
          </tr>
          <tr>
            <td><strong>Error Handling</strong></td>
            <td>All error cases handled consistently</td>
          </tr>
          <tr>
            <td><strong>Tests</strong></td>
            <td>All spec test cases have corresponding tests</td>
          </tr>
          <tr>
            <td><strong>Configuration</strong></td>
            <td>Environment variables handled identically</td>
          </tr>
        </tbody>
      </table>

      <h2>Best Practices</h2>
      <ol>
        <li>
          <strong>Run parity checks after regeneration</strong> - Always verify after
          generating or updating code
        </li>
        <li>
          <strong>Use a reference language</strong> - Pick one implementation as the
          "source of truth" (usually the first generated)
        </li>
        <li>
          <strong>Update spec, not just code</strong> - If you find issues, consider
          whether the spec needs clarification
        </li>
        <li>
          <strong>Automate in CI</strong> - Add parity checks to your continuous
          integration pipeline
        </li>
      </ol>

      <h2>Continuous Parity in CI</h2>
      <pre className="code-block">
{`# .github/workflows/parity.yml
name: Feature Parity Check

on: [push, pull_request]

jobs:
  parity:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build RPG
        run: make build

      - name: Check Parity
        run: |
          # Run parity check via MCP
          ./bin/rpg ensure-parity \\
            --spec specs/url-shortener.spec.md \\
            --projects go:generated/go,typescript:generated/typescript`}
      </pre>

      <h2>Congratulations!</h2>
      <p>
        You've completed all the tutorials! You now know how to:
      </p>
      <ul>
        <li>Write effective specification files</li>
        <li>Generate code in multiple languages</li>
        <li>Build complete applications from specs</li>
        <li>Maintain consistency across implementations</li>
      </ul>
      <p>
        Check out the <a href="/mcp-tools">MCP Tools Reference</a> for complete
        documentation on all available tools.
      </p>
    </div>
  );
}
