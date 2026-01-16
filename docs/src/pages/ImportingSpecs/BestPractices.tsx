import { Check, X, Lightbulb, AlertTriangle, FolderTree, RefreshCw } from 'lucide-react';

export default function BestPractices() {
  return (
    <div className="prose-docs">
      <h1>Import Best Practices</h1>
      <p className="lead">
        Tips for successful spec imports. Prepare your codebase, review generated specs,
        and iterate toward accurate documentation.
      </p>

      <h2>Prepare Your Codebase</h2>
      <p>
        A little preparation goes a long way. The cleaner your input, the better your spec.
      </p>

      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <div className="flex items-center gap-2 mb-3">
            <Check className="h-5 w-5 text-green-500" />
            <span className="font-semibold text-green-900 dark:text-green-100">Do Include</span>
          </div>
          <ul className="text-sm text-green-800 dark:text-green-200 space-y-2 mb-0 list-none pl-0">
            <li>✓ Source files with business logic</li>
            <li>✓ Test files (they reveal behavior)</li>
            <li>✓ README or documentation</li>
            <li>✓ Configuration files</li>
            <li>✓ API specs (OpenAPI, GraphQL)</li>
            <li>✓ Type definitions</li>
          </ul>
        </div>
        <div className="p-4 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800">
          <div className="flex items-center gap-2 mb-3">
            <X className="h-5 w-5 text-red-500" />
            <span className="font-semibold text-red-900 dark:text-red-100">Don't Include</span>
          </div>
          <ul className="text-sm text-red-800 dark:text-red-200 space-y-2 mb-0 list-none pl-0">
            <li>✗ Generated code (build output)</li>
            <li>✗ Vendor/node_modules directories</li>
            <li>✗ Binary files or assets</li>
            <li>✗ Lock files (package-lock.json)</li>
            <li>✗ IDE configuration files</li>
            <li>✗ Unrelated projects in same repo</li>
          </ul>
        </div>
      </div>

      <h2>Keep It Focused</h2>
      <p>
        Import modules, not monoliths. Smaller, focused imports produce better specs.
      </p>

      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
        <div className="p-4 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800">
          <div className="flex items-center gap-2 mb-2">
            <X className="h-5 w-5 text-red-500" />
            <span className="font-semibold text-red-900 dark:text-red-100">Too Broad</span>
          </div>
          <pre className="text-xs bg-red-100 dark:bg-red-900/40 p-2 rounded text-red-800 dark:text-red-200 overflow-x-auto">
{`import_spec_from_source
  inputPath: "./entire-monolith"

# 500+ files, multiple domains
# AI can't focus on what matters`}
          </pre>
        </div>
        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <div className="flex items-center gap-2 mb-2">
            <Check className="h-5 w-5 text-green-500" />
            <span className="font-semibold text-green-900 dark:text-green-100">Focused</span>
          </div>
          <pre className="text-xs bg-green-100 dark:bg-green-900/40 p-2 rounded text-green-800 dark:text-green-200 overflow-x-auto">
{`import_spec_from_source
  inputPath: "./src/auth"
  name: "auth-service"

# Just the auth module
# Clear boundaries, better spec`}
          </pre>
        </div>
      </div>

      <div className="not-prose my-6 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
        <div className="flex items-start gap-3">
          <FolderTree className="h-5 w-5 text-blue-500 flex-shrink-0 mt-0.5" />
          <div>
            <h4 className="font-semibold text-blue-900 dark:text-blue-100 mt-0 mb-1">
              Large Codebase Strategy
            </h4>
            <p className="text-sm text-blue-800 dark:text-blue-200 mb-2">
              For large projects, break down by domain:
            </p>
            <ul className="text-sm text-blue-800 dark:text-blue-200 mb-0 space-y-1">
              <li>• <code>./src/auth</code> → auth.spec.md</li>
              <li>• <code>./src/users</code> → users.spec.md</li>
              <li>• <code>./src/billing</code> → billing.spec.md</li>
              <li>• <code>./src/shared</code> → shared-utils.spec.md</li>
            </ul>
          </div>
        </div>
      </div>

      <h2>Review Generated Specs</h2>
      <p>
        AI-generated specs are a starting point, not a final product. Always review for:
      </p>

      <ul>
        <li>
          <strong>Accuracy</strong> — Does the spec match what the code actually does?
          The AI might misinterpret complex logic.
        </li>
        <li>
          <strong>Completeness</strong> — Are all important functions captured? Did any
          edge cases get missed?
        </li>
        <li>
          <strong>Clarity</strong> — Would someone unfamiliar with the code understand
          the spec? Add context where needed.
        </li>
        <li>
          <strong>Intent vs Implementation</strong> — Does the spec describe "what" or
          accidentally describe "how"?
        </li>
      </ul>

      <h2>Enhance Your Specs</h2>
      <p>
        After review, improve the generated spec:
      </p>

      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
        <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
          <h4 className="font-semibold text-gray-900 dark:text-white mt-0 mb-2">
            Generated (Raw)
          </h4>
          <pre className="text-xs bg-gray-100 dark:bg-gray-900 p-2 rounded text-gray-700 dark:text-gray-300 overflow-x-auto">
{`### calculateDiscount
Calculates discount amount.

**Logic**:
1. If premium and total > 100
2. Return total * 0.1
3. Else return 0`}
          </pre>
        </div>
        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <h4 className="font-semibold text-green-900 dark:text-green-100 mt-0 mb-2">
            Enhanced (After Review)
          </h4>
          <pre className="text-xs bg-green-100 dark:bg-green-900/40 p-2 rounded text-green-800 dark:text-green-200 overflow-x-auto">
{`### calculateDiscount
Applies loyalty discount for qualifying orders.

**Business Rule**: Premium members get 10% off
orders over $100 (pre-tax).

**Accepts**: order total, user membership
**Returns**: discount amount (0 if not eligible)

**Config**: PREMIUM_THRESHOLD = 100`}
          </pre>
        </div>
      </div>

      <h2>Iterate and Validate</h2>
      <p>
        The import process is iterative. Use this cycle:
      </p>

      <div className="not-prose my-8 p-6 bg-gray-50 dark:bg-gray-800/30 rounded-xl">
        <div className="flex flex-col md:flex-row items-center justify-center gap-4">
          <div className="text-center">
            <div className="w-12 h-12 rounded-full bg-purple-100 dark:bg-purple-900/40 flex items-center justify-center mx-auto mb-2">
              <span className="text-purple-600 dark:text-purple-400 font-semibold">1</span>
            </div>
            <span className="text-sm text-gray-600 dark:text-gray-400">Import</span>
          </div>
          <RefreshCw className="h-4 w-4 text-gray-400 rotate-90 md:rotate-0" />
          <div className="text-center">
            <div className="w-12 h-12 rounded-full bg-blue-100 dark:bg-blue-900/40 flex items-center justify-center mx-auto mb-2">
              <span className="text-blue-600 dark:text-blue-400 font-semibold">2</span>
            </div>
            <span className="text-sm text-gray-600 dark:text-gray-400">Review</span>
          </div>
          <RefreshCw className="h-4 w-4 text-gray-400 rotate-90 md:rotate-0" />
          <div className="text-center">
            <div className="w-12 h-12 rounded-full bg-green-100 dark:bg-green-900/40 flex items-center justify-center mx-auto mb-2">
              <span className="text-green-600 dark:text-green-400 font-semibold">3</span>
            </div>
            <span className="text-sm text-gray-600 dark:text-gray-400">Refine</span>
          </div>
          <RefreshCw className="h-4 w-4 text-gray-400 rotate-90 md:rotate-0" />
          <div className="text-center">
            <div className="w-12 h-12 rounded-full bg-orange-100 dark:bg-orange-900/40 flex items-center justify-center mx-auto mb-2">
              <span className="text-orange-600 dark:text-orange-400 font-semibold">4</span>
            </div>
            <span className="text-sm text-gray-600 dark:text-gray-400">Generate</span>
          </div>
          <RefreshCw className="h-4 w-4 text-gray-400 rotate-90 md:rotate-0" />
          <div className="text-center">
            <div className="w-12 h-12 rounded-full bg-red-100 dark:bg-red-900/40 flex items-center justify-center mx-auto mb-2">
              <span className="text-red-600 dark:text-red-400 font-semibold">5</span>
            </div>
            <span className="text-sm text-gray-600 dark:text-gray-400">Test</span>
          </div>
        </div>
      </div>

      <ol>
        <li>
          <strong>Import</strong> — Run <code>import_spec_from_source</code> on your code
        </li>
        <li>
          <strong>Review</strong> — Check accuracy against original behavior
        </li>
        <li>
          <strong>Refine</strong> — Add context, fix errors, clarify intent
        </li>
        <li>
          <strong>Generate</strong> — Create code in target language
        </li>
        <li>
          <strong>Test</strong> — Run original tests against new implementation
        </li>
      </ol>

      <p>
        If tests fail, refine the spec and regenerate. The spec is your source of truth.
      </p>

      <h2>Language-Specific Considerations</h2>

      <div className="not-prose space-y-4">
        <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
          <h4 className="font-semibold text-gray-900 dark:text-white mt-0 mb-2">
            Dynamic Languages (Python, JavaScript, PHP)
          </h4>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
            Runtime behavior might not be visible in static analysis. Include:
          </p>
          <ul className="text-sm text-gray-600 dark:text-gray-400 space-y-1 mb-0">
            <li>• Type hints/annotations where available</li>
            <li>• Tests that exercise dynamic behavior</li>
            <li>• Documentation of duck typing expectations</li>
          </ul>
        </div>

        <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
          <h4 className="font-semibold text-gray-900 dark:text-white mt-0 mb-2">
            Strongly Typed Languages (Java, C#, Go, Rust)
          </h4>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
            Type information helps spec generation. Leverage:
          </p>
          <ul className="text-sm text-gray-600 dark:text-gray-400 space-y-1 mb-0">
            <li>• Interface definitions for contracts</li>
            <li>• Generic type parameters</li>
            <li>• Nullable annotations</li>
          </ul>
        </div>

        <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
          <h4 className="font-semibold text-gray-900 dark:text-white mt-0 mb-2">
            Framework-Heavy Code
          </h4>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
            Separate business logic from framework specifics:
          </p>
          <ul className="text-sm text-gray-600 dark:text-gray-400 space-y-1 mb-0">
            <li>• Focus on domain logic, not framework boilerplate</li>
            <li>• Document framework behaviors separately</li>
            <li>• Note which parts are framework-specific</li>
          </ul>
        </div>
      </div>

      <h2>Common Pitfalls</h2>

      <div className="not-prose my-6 p-4 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg">
        <div className="flex items-start gap-3">
          <AlertTriangle className="h-5 w-5 text-amber-500 flex-shrink-0 mt-0.5" />
          <div>
            <h4 className="font-semibold text-amber-900 dark:text-amber-100 mt-0 mb-2">
              Watch Out For
            </h4>
            <ul className="text-sm text-amber-800 dark:text-amber-200 space-y-2 mb-0 list-none pl-0">
              <li className="flex items-start gap-2">
                <X className="h-4 w-4 text-amber-600 flex-shrink-0 mt-0.5" />
                <span>
                  <strong>Assuming perfection</strong> — Generated specs need review.
                  Don't blindly trust the output.
                </span>
              </li>
              <li className="flex items-start gap-2">
                <X className="h-4 w-4 text-amber-600 flex-shrink-0 mt-0.5" />
                <span>
                  <strong>Skipping tests</strong> — Tests are your best validation.
                  Run original tests against generated code.
                </span>
              </li>
              <li className="flex items-start gap-2">
                <X className="h-4 w-4 text-amber-600 flex-shrink-0 mt-0.5" />
                <span>
                  <strong>Importing too much</strong> — Large imports dilute focus.
                  Break into modules.
                </span>
              </li>
              <li className="flex items-start gap-2">
                <X className="h-4 w-4 text-amber-600 flex-shrink-0 mt-0.5" />
                <span>
                  <strong>Ignoring context</strong> — Include README, docs, and comments.
                  They provide intent.
                </span>
              </li>
              <li className="flex items-start gap-2">
                <X className="h-4 w-4 text-amber-600 flex-shrink-0 mt-0.5" />
                <span>
                  <strong>Not iterating</strong> — First import is rarely perfect.
                  Refine and regenerate.
                </span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <h2>When NOT to Import</h2>
      <p>
        Sometimes writing a spec from scratch is better:
      </p>

      <ul>
        <li>
          <strong>Severe technical debt</strong> — If the code is a mess, the spec will be too.
          Consider redesigning instead.
        </li>
        <li>
          <strong>Major behavior changes</strong> — If you're changing how things work,
          write the new spec directly.
        </li>
        <li>
          <strong>Very simple code</strong> — A 10-line utility might be faster to spec by hand.
        </li>
        <li>
          <strong>Highly framework-specific</strong> — Rails/Django/Spring magic doesn't
          translate well to generic specs.
        </li>
        <li>
          <strong>Unclear boundaries</strong> — If you don't know where to draw the line,
          clean up first.
        </li>
      </ul>

      <h2>Quick Checklist</h2>
      <div className="not-prose my-6">
        <ul className="space-y-2">
          {[
            'Focused scope (module, not monolith)',
            'Tests included in import',
            'Generated spec reviewed for accuracy',
            'Business logic enhanced with context',
            'Magic numbers converted to named config',
            'Edge cases documented',
            'New code tested against original tests',
            'Parity check run on generated code',
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

      <div className="not-prose my-8 p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
        <div className="flex gap-3">
          <Lightbulb className="h-5 w-5 text-blue-500 flex-shrink-0 mt-0.5" />
          <div>
            <p className="text-blue-800 dark:text-blue-200 text-sm mb-0">
              <strong>Pro Tip:</strong> Keep both the original code and generated spec
              version-controlled. This creates a clear audit trail and makes iteration easier.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
