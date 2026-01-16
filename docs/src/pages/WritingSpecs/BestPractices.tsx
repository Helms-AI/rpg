import { Check, X, Lightbulb, AlertTriangle } from 'lucide-react';

export default function BestPractices() {
  return (
    <div className="prose-docs">
      <h1>Best Practices</h1>
      <p className="lead">
        Practical tips for writing specs that generate better code.
        These aren't rules—they're patterns that help the AI understand your intent.
      </p>

      <h2>Focus on "What", Not "How"</h2>
      <p>
        Describe what you want the code to do, not how to implement it in a specific language.
        Let the AI choose idiomatic patterns for each target language.
      </p>

      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
        <div className="p-4 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800">
          <div className="flex items-center gap-2 mb-2">
            <X className="h-5 w-5 text-red-500" />
            <span className="font-semibold text-red-900 dark:text-red-100">Too Implementation-Specific</span>
          </div>
          <pre className="text-xs bg-red-100 dark:bg-red-900/40 p-2 rounded text-red-800 dark:text-red-200 overflow-x-auto">
{`Use a HashMap<String, User> to store users.
Iterate with a for loop and filter.
Return an ArrayList of results.`}
          </pre>
        </div>
        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <div className="flex items-center gap-2 mb-2">
            <Check className="h-5 w-5 text-green-500" />
            <span className="font-semibold text-green-900 dark:text-green-100">Behavior-Focused</span>
          </div>
          <pre className="text-xs bg-green-100 dark:bg-green-900/40 p-2 rounded text-green-800 dark:text-green-200 overflow-x-auto">
{`Store users by ID for quick lookup.
Filter users matching criteria.
Return the matching users as a list.`}
          </pre>
        </div>
      </div>

      <h2>Be Specific About Edge Cases</h2>
      <p>
        The AI makes reasonable assumptions, but it can't read your mind about edge cases.
        When behavior matters, be explicit.
      </p>

      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
        <div className="p-4 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800">
          <div className="flex items-center gap-2 mb-2">
            <X className="h-5 w-5 text-red-500" />
            <span className="font-semibold text-red-900 dark:text-red-100">Ambiguous</span>
          </div>
          <pre className="text-xs bg-red-100 dark:bg-red-900/40 p-2 rounded text-red-800 dark:text-red-200 overflow-x-auto">
{`### divide
Divide two numbers.`}
          </pre>
          <p className="text-xs text-red-700 dark:text-red-300 mt-2">
            What happens when dividing by zero? Integer or float division?
          </p>
        </div>
        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <div className="flex items-center gap-2 mb-2">
            <Check className="h-5 w-5 text-green-500" />
            <span className="font-semibold text-green-900 dark:text-green-100">Clear Edge Cases</span>
          </div>
          <pre className="text-xs bg-green-100 dark:bg-green-900/40 p-2 rounded text-green-800 dark:text-green-200 overflow-x-auto">
{`### divide
Divide two numbers (float division).
Return error if divisor is zero.`}
          </pre>
        </div>
      </div>

      <h2>Use Meaningful Names</h2>
      <p>
        Choose descriptive names for types and functions. The AI uses these names
        when generating code, so good names lead to readable output.
      </p>

      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
        <div className="p-4 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800">
          <div className="flex items-center gap-2 mb-2">
            <X className="h-5 w-5 text-red-500" />
            <span className="font-semibold text-red-900 dark:text-red-100">Vague Names</span>
          </div>
          <pre className="text-xs bg-red-100 dark:bg-red-900/40 p-2 rounded text-red-800 dark:text-red-200 overflow-x-auto">
{`### Data
| Field | Type |
|-------|------|
| val1  | string |
| val2  | int |
| flag  | bool |`}
          </pre>
        </div>
        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <div className="flex items-center gap-2 mb-2">
            <Check className="h-5 w-5 text-green-500" />
            <span className="font-semibold text-green-900 dark:text-green-100">Descriptive Names</span>
          </div>
          <pre className="text-xs bg-green-100 dark:bg-green-900/40 p-2 rounded text-green-800 dark:text-green-200 overflow-x-auto">
{`### UserProfile
| Field | Type |
|-------|------|
| email | string |
| loginCount | int |
| isVerified | bool |`}
          </pre>
        </div>
      </div>

      <h2>Include Test Cases for Complex Logic</h2>
      <p>
        Test cases serve two purposes: they define expected behavior, and they help
        verify the generated code works correctly. Include them for anything non-trivial.
      </p>

      <div className="not-prose my-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
        <div className="flex items-start gap-3">
          <Lightbulb className="h-5 w-5 text-blue-500 flex-shrink-0 mt-0.5" />
          <div>
            <h4 className="font-semibold text-blue-900 dark:text-blue-100 mt-0 mb-1">
              Test Case Format
            </h4>
            <p className="text-sm text-blue-800 dark:text-blue-200 mb-3">
              Use the Given/Expect pattern for clear, readable tests:
            </p>
            <pre className="text-xs bg-blue-100 dark:bg-blue-900/40 p-3 rounded text-blue-800 dark:text-blue-200 overflow-x-auto">
{`### test_empty_string_returns_zero
**Given**: input is ""
**Expect**: wordCount returns 0

### test_counts_words_correctly
**Given**: input is "hello world"
**Expect**: wordCount returns 2

### test_handles_multiple_spaces
**Given**: input is "hello    world"
**Expect**: wordCount returns 2 (not 4)`}
            </pre>
          </div>
        </div>
      </div>

      <h2>Keep Specs Focused</h2>
      <p>
        One spec per logical component works better than one giant spec for everything.
        Split large projects into focused specs that can be generated independently.
      </p>

      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
        <div className="p-4 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800">
          <div className="flex items-center gap-2 mb-2">
            <X className="h-5 w-5 text-red-500" />
            <span className="font-semibold text-red-900 dark:text-red-100">Monolithic</span>
          </div>
          <p className="text-xs text-red-700 dark:text-red-300">
            One spec with auth, users, posts, comments, notifications, admin panel,
            analytics, billing, and email service.
          </p>
        </div>
        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <div className="flex items-center gap-2 mb-2">
            <Check className="h-5 w-5 text-green-500" />
            <span className="font-semibold text-green-900 dark:text-green-100">Focused</span>
          </div>
          <p className="text-xs text-green-700 dark:text-green-300">
            Separate specs: auth.spec.md, users.spec.md, posts.spec.md, etc.
            Each focused on one domain.
          </p>
        </div>
      </div>

      <h2>Document Dependencies</h2>
      <p>
        If your code needs external libraries, mention them. The AI will use
        language-appropriate equivalents.
      </p>

      <pre className="code-block">
{`## Dependencies

This service needs:
- HTTP client for making API calls
- JSON parsing
- Logging (structured, with levels)
- Environment variable loading`}
      </pre>

      <p>
        The AI translates these to the right packages: <code>net/http</code> in Go,
        <code>reqwest</code> in Rust, <code>axios</code> in TypeScript, etc.
      </p>

      <h2>Common Patterns</h2>

      <h3>Configuration with Defaults</h3>
      <pre className="code-block">
{`## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | Server listen port |
| TIMEOUT_SECONDS | 30 | Request timeout |
| LOG_LEVEL | info | Logging verbosity |`}
      </pre>

      <h3>Error Handling</h3>
      <pre className="code-block">
{`### fetchUser
Get user by ID from the database.

**Returns**: User or error

**Errors**:
- NotFound: user doesn't exist
- DatabaseError: connection or query failed`}
      </pre>

      <h3>API Endpoints</h3>
      <pre className="code-block">
{`### POST /api/items
Create a new item.

**Auth**: Required (Bearer token)
**Body**: { name: string, quantity: int }
**Response 201**: Created item with ID
**Response 400**: Invalid input
**Response 401**: Not authenticated`}
      </pre>

      <h2>Things to Avoid</h2>

      <div className="not-prose my-6 p-4 bg-amber-50 dark:bg-amber-900/20 rounded-lg border border-amber-200 dark:border-amber-800">
        <div className="flex items-start gap-3">
          <AlertTriangle className="h-5 w-5 text-amber-500 flex-shrink-0 mt-0.5" />
          <div>
            <h4 className="font-semibold text-amber-900 dark:text-amber-100 mt-0 mb-2">
              Watch Out For
            </h4>
            <ul className="text-sm text-amber-800 dark:text-amber-200 space-y-2 mb-0 list-none pl-0">
              <li className="flex items-start gap-2">
                <X className="h-4 w-4 text-amber-600 flex-shrink-0 mt-0.5" />
                <span><strong>Language-specific syntax</strong> — Don't write Java generics or Go channels. Use generic descriptions.</span>
              </li>
              <li className="flex items-start gap-2">
                <X className="h-4 w-4 text-amber-600 flex-shrink-0 mt-0.5" />
                <span><strong>Over-specifying structure</strong> — Don't dictate file names or internal architecture unless necessary.</span>
              </li>
              <li className="flex items-start gap-2">
                <X className="h-4 w-4 text-amber-600 flex-shrink-0 mt-0.5" />
                <span><strong>Contradictory requirements</strong> — Review for consistency. The AI can't resolve "sync and async".</span>
              </li>
              <li className="flex items-start gap-2">
                <X className="h-4 w-4 text-amber-600 flex-shrink-0 mt-0.5" />
                <span><strong>Assuming context</strong> — Don't reference external files or previous specs without explanation.</span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <h2>Multi-Language Tips</h2>
      <p>
        When generating the same spec in multiple languages:
      </p>

      <ul>
        <li>
          <strong>Use portable types</strong> — <code>string</code>, <code>int</code>, <code>bool</code>,
          <code>float</code>, <code>list</code>, <code>map</code>. The AI maps these to native types.
        </li>
        <li>
          <strong>Avoid platform-specific features</strong> — If you need async, say "async" not "goroutines"
          or "CompletableFuture".
        </li>
        <li>
          <strong>Use ensure_parity</strong> — After generating, run the parity check tool to verify
          all implementations match.
        </li>
        <li>
          <strong>Consider language strengths</strong> — Some patterns are more natural in certain languages.
          The AI adapts, but extreme paradigm differences may need separate specs.
        </li>
      </ul>

      <h2>Quick Checklist</h2>
      <div className="not-prose my-6">
        <ul className="space-y-2">
          {[
            'Descriptive title and overview',
            'Meaningful type and function names',
            'Clear input/output descriptions',
            'Edge cases documented',
            'Test cases for complex logic',
            'Dependencies listed',
            'Configuration with defaults',
            'Error conditions specified',
          ].map((item, i) => (
            <li key={i} className="flex items-center gap-2 text-gray-700 dark:text-gray-300">
              <div className="w-5 h-5 rounded border border-gray-300 dark:border-gray-600 flex items-center justify-center">
                <Check className="h-3 w-3 text-gray-400" />
              </div>
              {item}
            </li>
          ))}
        </ul>
      </div>

      <p>
        Remember: these are guidelines, not requirements. A spec that clearly communicates
        intent will generate good code, regardless of format.
      </p>
    </div>
  );
}
