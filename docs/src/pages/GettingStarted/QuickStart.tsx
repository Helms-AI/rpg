export default function QuickStart() {
  return (
    <div className="prose-docs">
      <h1>Quick Start</h1>
      <p className="lead">Write your first spec and generate code in 5 minutes.</p>

      <h2>Step 1: Create a Spec File</h2>
      <p>Create a new file called <code>my-function.spec.md</code>:</p>
      <pre><code>{`# slugify

A utility function to convert text to URL-friendly slugs.

## Target Languages
- go
- typescript

## Functions
### slugify
**accepts:** text: Text
**returns:** Text
**logic:**
  - convert to lowercase
  - replace spaces with dashes
  - remove special characters`}</code></pre>

      <h2>Step 2: Generate Code</h2>
      <p>Ask your AI assistant to generate code using RPG:</p>
      <pre><code>{`Use RPG to generate code from my-function.spec.md in Go`}</code></pre>

      <h2>Step 3: Review and Use</h2>
      <p>The generated code will follow Go idioms and conventions.</p>
    </div>
  );
}
