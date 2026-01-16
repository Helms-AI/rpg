import { Link } from 'react-router-dom';
import { FileText, Sparkles, Languages, MessageSquare } from 'lucide-react';

export default function WritingSpecs() {
  return (
    <div className="prose-docs">
      <h1>Writing Specs</h1>
      <p className="lead">
        Describe what you want to build in plain markdown. RPG uses AI to interpret your intent
        and generate idiomatic code in any supported language.
      </p>

      {/* Philosophy box */}
      <div className="not-prose my-8 p-6 bg-gradient-to-br from-violet-50 to-purple-50 dark:from-violet-900/20 dark:to-purple-900/20 border border-violet-200 dark:border-violet-800 rounded-xl">
        <div className="flex items-start gap-4">
          <div className="p-3 bg-violet-500 rounded-lg">
            <Sparkles className="h-6 w-6 text-white" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-violet-900 dark:text-violet-100 mt-0 mb-2">
              Natural Language First
            </h3>
            <p className="text-violet-800 dark:text-violet-200 mb-0">
              There's no rigid format to memorize. Write specs the way you'd explain your project
              to another developer. The AI understands context, infers relationships, and adapts
              to your writing style.
            </p>
          </div>
        </div>
      </div>

      <h2>What is a Spec?</h2>
      <p>
        A spec is simply a markdown document that describes the software you want to build.
        It can be as simple as a few paragraphs or as detailed as a full API specification.
        The AI reads your description and generates production-quality code.
      </p>

      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4 my-6">
        <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-2 mb-2">
            <FileText className="h-5 w-5 text-blue-500" />
            <span className="font-semibold text-gray-900 dark:text-white">You Write</span>
          </div>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            A markdown file describing what your code should do, what data it handles,
            and how it should behave.
          </p>
        </div>
        <div className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-2 mb-2">
            <Languages className="h-5 w-5 text-green-500" />
            <span className="font-semibold text-gray-900 dark:text-white">RPG Generates</span>
          </div>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            Idiomatic, production-ready code in Go, Rust, Java, C#, Python, or TypeScript
            with proper project structure.
          </p>
        </div>
      </div>

      <h2>A Simple Example</h2>
      <p>
        Here's a minimal spec that generates a complete URL shortener:
      </p>

      <pre className="code-block">
{`# URL Shortener

A simple service that creates short URLs and redirects to the original.

## Functions

### shorten
Takes a long URL and returns a short code.

### resolve
Takes a short code and returns the original URL, or an error if not found.

## Storage
Store mappings in memory with a map from code to URL.`}
      </pre>

      <p>
        That's it. From this description, RPG can generate a complete implementation
        with proper error handling, data structures, and idiomatic patterns for your
        chosen language.
      </p>

      <h2>How It Works</h2>
      <div className="not-prose space-y-4 my-6">
        <div className="flex gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-blue-500 text-white flex items-center justify-center font-semibold">
            1
          </div>
          <div>
            <h4 className="font-semibold text-gray-900 dark:text-white mt-0">Write Your Spec</h4>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              Describe your project in markdown. Include what it does, the data it uses,
              and the behavior you expect.
            </p>
          </div>
        </div>
        <div className="flex gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-blue-500 text-white flex items-center justify-center font-semibold">
            2
          </div>
          <div>
            <h4 className="font-semibold text-gray-900 dark:text-white mt-0">AI Interprets</h4>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              RPG parses your markdown and combines it with language-specific conventions.
              The AI understands your intent, even if you don't follow a strict format.
            </p>
          </div>
        </div>
        <div className="flex gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-blue-500 text-white flex items-center justify-center font-semibold">
            3
          </div>
          <div>
            <h4 className="font-semibold text-gray-900 dark:text-white mt-0">Code Generation</h4>
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              The AI generates idiomatic code that follows each language's conventions,
              best practices, and standard project structure.
            </p>
          </div>
        </div>
      </div>

      <h2>Flexibility by Design</h2>
      <p>
        Unlike traditional code generators that require exact syntax, RPG embraces flexibility:
      </p>

      <ul>
        <li>
          <strong>Use any structure</strong> — Organize your spec however makes sense for your project.
          Headers, bullets, tables, prose—it all works.
        </li>
        <li>
          <strong>Be as detailed as needed</strong> — A utility function might need one sentence.
          A complex API might need detailed endpoint descriptions.
        </li>
        <li>
          <strong>Include context</strong> — The AI can use background information, examples,
          and constraints to generate better code.
        </li>
        <li>
          <strong>Write naturally</strong> — Don't worry about keywords or formatting rules.
          Explain it like you would to a teammate.
        </li>
      </ul>

      <h2>Common Sections</h2>
      <p>
        While there's no required format, these sections are commonly useful:
      </p>

      <div className="not-prose overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-gray-200 dark:border-gray-700">
              <th className="text-left py-2 pr-4 font-semibold text-gray-900 dark:text-white">Section</th>
              <th className="text-left py-2 font-semibold text-gray-900 dark:text-white">Purpose</th>
            </tr>
          </thead>
          <tbody className="text-gray-600 dark:text-gray-400">
            <tr className="border-b border-gray-100 dark:border-gray-800">
              <td className="py-2 pr-4 font-mono text-xs">## Types</td>
              <td className="py-2">Define data structures, models, or DTOs</td>
            </tr>
            <tr className="border-b border-gray-100 dark:border-gray-800">
              <td className="py-2 pr-4 font-mono text-xs">## Functions</td>
              <td className="py-2">Describe functions with their inputs, outputs, and logic</td>
            </tr>
            <tr className="border-b border-gray-100 dark:border-gray-800">
              <td className="py-2 pr-4 font-mono text-xs">## Configuration</td>
              <td className="py-2">Environment variables and settings</td>
            </tr>
            <tr className="border-b border-gray-100 dark:border-gray-800">
              <td className="py-2 pr-4 font-mono text-xs">## API Endpoints</td>
              <td className="py-2">REST API routes and their behavior</td>
            </tr>
            <tr className="border-b border-gray-100 dark:border-gray-800">
              <td className="py-2 pr-4 font-mono text-xs">## Tests</td>
              <td className="py-2">Test cases with given/expect format</td>
            </tr>
            <tr>
              <td className="py-2 pr-4 font-mono text-xs">## Dependencies</td>
              <td className="py-2">External libraries needed</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="not-prose my-8 p-4 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg">
        <div className="flex gap-3">
          <MessageSquare className="h-5 w-5 text-amber-600 dark:text-amber-400 flex-shrink-0 mt-0.5" />
          <div>
            <p className="text-amber-800 dark:text-amber-200 text-sm mb-0">
              <strong>Remember:</strong> These sections are suggestions, not requirements.
              The AI adapts to your style. Focus on clearly describing what you want to build.
            </p>
          </div>
        </div>
      </div>

      <h2>Next Steps</h2>
      <div className="not-prose grid grid-cols-1 md:grid-cols-2 gap-4">
        <Link
          to="/writing-specs/examples"
          className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-700 transition-colors no-underline"
        >
          <h4 className="font-semibold text-gray-900 dark:text-white mt-0 mb-1">Examples →</h4>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-0">
            See real spec examples in different styles and complexity levels.
          </p>
        </Link>
        <Link
          to="/writing-specs/best-practices"
          className="p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-700 transition-colors no-underline"
        >
          <h4 className="font-semibold text-gray-900 dark:text-white mt-0 mb-1">Best Practices →</h4>
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-0">
            Tips for writing specs that generate better code.
          </p>
        </Link>
      </div>
    </div>
  );
}
