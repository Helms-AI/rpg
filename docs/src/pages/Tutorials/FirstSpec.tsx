export default function FirstSpec() {
  return (
    <div className="prose-docs">
      <h1>Your First Spec</h1>
      <p className="lead">
        Learn the basics of writing specification files that RPG can transform into working code.
      </p>

      <div className="not-prose my-6 p-4 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
        <p className="text-blue-800 dark:text-blue-300 text-sm">
          <strong>Time:</strong> 10 minutes &nbsp;|&nbsp; <strong>Difficulty:</strong> Beginner
        </p>
      </div>

      <h2>What You'll Build</h2>
      <p>
        A simple <code>slugify</code> function that converts strings into URL-friendly slugs.
        For example: "Hello World!" becomes "hello-world".
      </p>

      <h2>Step 1: Create the Spec File</h2>
      <p>Create a new file called <code>slugify.spec.md</code>:</p>

      <pre className="code-block">
{`# Slugify

A utility function that converts strings into URL-friendly slugs.

## Meta

- **Version**: 1.0.0
- **Author**: Your Name

## Target Languages

- go
- typescript
- python

## Functions

### slugify

Converts a string into a URL-friendly slug.

**Accepts:**
- \`input\` (string): The string to convert

**Returns:**
- string: The slugified string

**Logic:**
1. Convert to lowercase
2. Replace spaces with hyphens
3. Remove non-alphanumeric characters (except hyphens)
4. Remove consecutive hyphens
5. Trim hyphens from start and end

## Tests

### Basic conversion
- **Given**: "Hello World"
- **Expect**: "hello-world"

### Special characters
- **Given**: "Hello, World!"
- **Expect**: "hello-world"

### Multiple spaces
- **Given**: "Hello    World"
- **Expect**: "hello-world"

### Unicode
- **Given**: "Café Latté"
- **Expect**: "caf-latt"`}
      </pre>

      <h2>Step 2: Validate the Spec</h2>
      <p>Ask your AI assistant to validate the spec:</p>

      <pre className="code-block">
{`> Use RPG to validate specs/slugify.spec.md

✓ Spec file is valid and contains content`}
      </pre>

      <h2>Step 3: Generate Code</h2>
      <p>Now generate code in your preferred language:</p>

      <pre className="code-block">
{`> Generate Go code from specs/slugify.spec.md using RPG`}
      </pre>

      <p>The AI will use RPG's <code>get_generation_context</code> tool and produce idiomatic Go:</p>

      <pre className="code-block">
{`package slugify

import (
    "regexp"
    "strings"
)

// Slugify converts a string into a URL-friendly slug.
func Slugify(input string) string {
    // Convert to lowercase
    s := strings.ToLower(input)

    // Replace spaces with hyphens
    s = strings.ReplaceAll(s, " ", "-")

    // Remove non-alphanumeric characters except hyphens
    reg := regexp.MustCompile("[^a-z0-9-]")
    s = reg.ReplaceAllString(s, "")

    // Remove consecutive hyphens
    reg = regexp.MustCompile("-+")
    s = reg.ReplaceAllString(s, "-")

    // Trim hyphens from start and end
    return strings.Trim(s, "-")
}`}
      </pre>

      <h2>Step 4: Generate Tests</h2>
      <p>Ask for the test file as well:</p>

      <pre className="code-block">
{`> Also generate the Go tests for slugify`}
      </pre>

      <h2>Key Takeaways</h2>
      <ul>
        <li>Specs use simple markdown with predictable sections</li>
        <li>The AI interprets your natural language descriptions</li>
        <li>Test cases become real unit tests in the target language</li>
        <li>The same spec generates idiomatic code in any supported language</li>
      </ul>

      <h2>Next Steps</h2>
      <p>
        Ready to generate in multiple languages? Continue to the{' '}
        <a href="/tutorials/multi-language">Multi-Language Generation</a> tutorial.
      </p>
    </div>
  );
}
