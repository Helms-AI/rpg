import { Link } from 'react-router-dom';
import { Code2, FileSearch, ArrowRight, Languages, RefreshCw, FileText } from 'lucide-react';

export default function ImportingSpecs() {
  return (
    <div className="prose-docs">
      <h1>Importing Specs</h1>
      <p className="lead">
        Transform existing code into natural language specs. Extract the "what" from the "how"
        to enable porting, documentation, and modernization.
      </p>

      {/* Philosophy box */}
      <div className="not-prose my-8 p-6 bg-gradient-to-br from-purple-50 to-indigo-50 dark:from-purple-900/20 dark:to-indigo-900/20 border border-purple-200 dark:border-purple-800 rounded-xl">
        <div className="flex items-start gap-4">
          <div className="p-3 bg-purple-500 rounded-lg">
            <Code2 className="h-6 w-6 text-white" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-purple-900 dark:text-purple-100 mt-0 mb-2">
              Code Already Written?
            </h3>
            <p className="text-purple-800 dark:text-purple-200 mb-0">
              Don't rewrite specs from scratch. Point RPG at your existing codebase and let AI
              extract the specification. It's reverse-engineering for the AI age—extracting
              intent from implementation.
            </p>
          </div>
        </div>
      </div>

      <h2>What is Spec Importing?</h2>
      <p>
        Spec importing analyzes your existing source code and generates a natural language
        specification that describes what the code does. This isn't line-by-line translation—it's
        <strong> intent extraction</strong>. The AI reads your implementation and captures the
        underlying behavior, data structures, and business logic.
      </p>

      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
        <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-2 mb-2">
            <FileSearch className="h-5 w-5 text-purple-500" />
            <span className="font-semibold text-gray-900 dark:text-white">RPG Analyzes</span>
          </div>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            Your source files, tests, configs, and documentation. It understands structure,
            patterns, and relationships.
          </p>
        </div>
        <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-2 mb-2">
            <FileText className="h-5 w-5 text-green-500" />
            <span className="font-semibold text-gray-900 dark:text-white">AI Generates</span>
          </div>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            A natural language spec describing types, functions, behavior, and tests—ready
            for review, refinement, or multi-language generation.
          </p>
        </div>
      </div>

      <h2>The Import Workflow</h2>
      <div className="not-prose my-8">
        {/* Visual workflow */}
        <div className="flex flex-col md:flex-row items-center justify-center gap-4 p-6 bg-gray-50 dark:bg-gray-800/30 rounded-xl">
          <div className="flex flex-col items-center text-center">
            <div className="w-16 h-16 rounded-full bg-purple-100 dark:bg-purple-900/40 flex items-center justify-center mb-2">
              <Code2 className="h-8 w-8 text-purple-600 dark:text-purple-400" />
            </div>
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Existing Code</span>
          </div>
          <ArrowRight className="h-6 w-6 text-gray-400 rotate-90 md:rotate-0" />
          <div className="flex flex-col items-center text-center">
            <div className="w-16 h-16 rounded-full bg-blue-100 dark:bg-blue-900/40 flex items-center justify-center mb-2">
              <FileSearch className="h-8 w-8 text-blue-600 dark:text-blue-400" />
            </div>
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Analysis</span>
          </div>
          <ArrowRight className="h-6 w-6 text-gray-400 rotate-90 md:rotate-0" />
          <div className="flex flex-col items-center text-center">
            <div className="w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/40 flex items-center justify-center mb-2">
              <FileText className="h-8 w-8 text-green-600 dark:text-green-400" />
            </div>
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Spec</span>
          </div>
          <ArrowRight className="h-6 w-6 text-gray-400 rotate-90 md:rotate-0" />
          <div className="flex flex-col items-center text-center">
            <div className="w-16 h-16 rounded-full bg-orange-100 dark:bg-orange-900/40 flex items-center justify-center mb-2">
              <Languages className="h-8 w-8 text-orange-600 dark:text-orange-400" />
            </div>
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">New Languages</span>
          </div>
        </div>
      </div>

      <div className="not-prose space-y-4 my-6">
        <div className="flex gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-purple-500 text-white flex items-center justify-center font-semibold">
            1
          </div>
          <div>
            <h4 className="font-semibold text-gray-900 dark:text-white mt-0">Point to Source</h4>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              Provide the path to your existing codebase. RPG identifies source files,
              tests, configs, and documentation.
            </p>
          </div>
        </div>
        <div className="flex gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-purple-500 text-white flex items-center justify-center font-semibold">
            2
          </div>
          <div>
            <h4 className="font-semibold text-gray-900 dark:text-white mt-0">AI Analysis</h4>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              The tool collects relevant files and creates an analysis prompt. The AI
              examines structure, patterns, and behavior to understand intent.
            </p>
          </div>
        </div>
        <div className="flex gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-purple-500 text-white flex items-center justify-center font-semibold">
            3
          </div>
          <div>
            <h4 className="font-semibold text-gray-900 dark:text-white mt-0">Generate Spec</h4>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              The AI produces a natural language specification capturing types,
              functions, behavior, and test cases.
            </p>
          </div>
        </div>
        <div className="flex gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-purple-500 text-white flex items-center justify-center font-semibold">
            4
          </div>
          <div>
            <h4 className="font-semibold text-gray-900 dark:text-white mt-0">Review & Refine</h4>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              Review the generated spec for accuracy. Add context, clarify behavior,
              or enhance descriptions where needed.
            </p>
          </div>
        </div>
        <div className="flex gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-purple-500 text-white flex items-center justify-center font-semibold">
            5
          </div>
          <div>
            <h4 className="font-semibold text-gray-900 dark:text-white mt-0">Generate Code</h4>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              Use the spec to generate idiomatic code in any supported language.
              Validate with parity checking.
            </p>
          </div>
        </div>
      </div>

      <h2>When to Import vs Write</h2>
      <p>
        Importing isn't always the right choice. Here's a quick guide:
      </p>

      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
        <div className="p-4 bg-purple-50 dark:bg-purple-900/20 rounded-lg border border-purple-200 dark:border-purple-800">
          <h4 className="font-semibold text-purple-900 dark:text-purple-100 mt-0 mb-3 flex items-center gap-2">
            <FileSearch className="h-5 w-5" />
            Import When...
          </h4>
          <ul className="text-sm text-purple-800 dark:text-purple-200 space-y-2 mb-0 list-none pl-0">
            <li>✓ Code already exists and works</li>
            <li>✓ Porting to a new language</li>
            <li>✓ Documentation is missing or outdated</li>
            <li>✓ Creating multi-language versions</li>
            <li>✓ Understanding legacy behavior</li>
          </ul>
        </div>
        <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <h4 className="font-semibold text-blue-900 dark:text-blue-100 mt-0 mb-3 flex items-center gap-2">
            <FileText className="h-5 w-5" />
            Write When...
          </h4>
          <ul className="text-sm text-blue-800 dark:text-blue-200 space-y-2 mb-0 list-none pl-0">
            <li>✓ Starting a new project</li>
            <li>✓ Code needs major restructuring</li>
            <li>✓ Existing code is too messy to extract</li>
            <li>✓ You want to redesign behavior</li>
            <li>✓ Simpler to describe than reverse-engineer</li>
          </ul>
        </div>
      </div>

      <h2>What Gets Analyzed</h2>
      <p>
        The import tool examines multiple file types to build a complete picture:
      </p>

      <div className="not-prose overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-gray-200 dark:border-gray-700">
              <th className="text-left py-2 pr-4 font-semibold text-gray-900 dark:text-white">File Type</th>
              <th className="text-left py-2 font-semibold text-gray-900 dark:text-white">What It Reveals</th>
            </tr>
          </thead>
          <tbody className="text-gray-600 dark:text-gray-400">
            <tr className="border-b border-gray-100 dark:border-gray-800">
              <td className="py-2 pr-4">Source files</td>
              <td className="py-2">Functions, types, classes, modules, logic</td>
            </tr>
            <tr className="border-b border-gray-100 dark:border-gray-800">
              <td className="py-2 pr-4">Test files</td>
              <td className="py-2">Expected behavior, edge cases, usage examples</td>
            </tr>
            <tr className="border-b border-gray-100 dark:border-gray-800">
              <td className="py-2 pr-4">Config files</td>
              <td className="py-2">Dependencies, environment variables, settings</td>
            </tr>
            <tr className="border-b border-gray-100 dark:border-gray-800">
              <td className="py-2 pr-4">API specs</td>
              <td className="py-2">Endpoints, request/response formats (OpenAPI, GraphQL)</td>
            </tr>
            <tr>
              <td className="py-2 pr-4">Documentation</td>
              <td className="py-2">Intent, context, usage instructions (README, comments)</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="not-prose my-8 p-4 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg">
        <div className="flex gap-3">
          <RefreshCw className="h-5 w-5 text-amber-600 dark:text-amber-400 flex-shrink-0 mt-0.5" />
          <div>
            <p className="text-amber-800 dark:text-amber-200 text-sm mb-0">
              <strong>AI-Assisted, Not Automatic:</strong> Importing is a collaborative process.
              The AI provides a strong starting point, but you should review and refine the
              generated spec to ensure accuracy and completeness.
            </p>
          </div>
        </div>
      </div>

      <h2>Use Cases</h2>
      <ul>
        <li>
          <strong>Legacy modernization</strong> — Port a working Java monolith to Go microservices
          without losing functionality
        </li>
        <li>
          <strong>Multi-language libraries</strong> — Create TypeScript, Python, and Rust versions
          of your utility library
        </li>
        <li>
          <strong>Documentation recovery</strong> — Generate specs for undocumented code
          before the original author leaves
        </li>
        <li>
          <strong>Refactoring confidence</strong> — Understand exactly what code does before
          making changes
        </li>
        <li>
          <strong>Knowledge transfer</strong> — Create readable specs that explain complex
          systems to new team members
        </li>
      </ul>

      <h2>Next Steps</h2>
      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4">
        <Link
          to="/importing-specs/examples"
          className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-purple-300 dark:hover:border-purple-700 transition-colors no-underline"
        >
          <h4 className="font-semibold text-gray-900 dark:text-white mt-0 mb-1">Examples →</h4>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-0">
            See real import scenarios from utility functions to full applications.
          </p>
        </Link>
        <Link
          to="/importing-specs/best-practices"
          className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-purple-300 dark:hover:border-purple-700 transition-colors no-underline"
        >
          <h4 className="font-semibold text-gray-900 dark:text-white mt-0 mb-1">Best Practices →</h4>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-0">
            Tips for successful imports and common pitfalls to avoid.
          </p>
        </Link>
      </div>
    </div>
  );
}
