export default function MultiLanguage() {
  return (
    <div className="prose-docs">
      <h1>Multi-Language Generation</h1>
      <p className="lead">
        Generate the same project in multiple programming languages from a single spec.
      </p>

      <div className="not-prose my-6 p-4 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
        <p className="text-blue-800 dark:text-blue-300 text-sm">
          <strong>Time:</strong> 15 minutes &nbsp;|&nbsp; <strong>Difficulty:</strong> Intermediate
        </p>
      </div>

      <h2>What You'll Learn</h2>
      <ul>
        <li>How to specify multiple target languages</li>
        <li>Language-specific conventions RPG applies</li>
        <li>Getting the recommended project structure for each language</li>
        <li>Generating implementations in all languages</li>
      </ul>

      <h2>Step 1: Check Available Languages</h2>
      <p>First, see what languages RPG supports:</p>

      <pre className="code-block">
{`> Use RPG to list available languages

Languages:
- go: Go (error returns, internal packages)
- rust: Rust (Result/Option, Cargo)
- java: Java (checked exceptions, Maven/Gradle)
- csharp: C# (nullable refs, .NET 6+)
- python: Python (type hints, 3.10+)
- typescript: TypeScript (strict mode, npm)`}
      </pre>

      <h2>Step 2: Create a Multi-Language Spec</h2>
      <p>
        Update your spec's Target Languages section to include all languages
        you want to generate:
      </p>

      <pre className="code-block">
{`# Calculator

A basic calculator library with common operations.

## Target Languages

- go
- rust
- typescript
- python

## Functions

### add
Adds two numbers together.

**Accepts:**
- \`a\` (number): First operand
- \`b\` (number): Second operand

**Returns:**
- number: The sum of a and b

### divide
Divides the first number by the second.

**Accepts:**
- \`a\` (number): Dividend
- \`b\` (number): Divisor

**Returns:**
- number: The quotient
- error: If divisor is zero

## Tests

### Addition
- **Given**: add(2, 3)
- **Expect**: 5

### Division
- **Given**: divide(10, 2)
- **Expect**: 5

### Division by zero
- **Given**: divide(10, 0)
- **Expect**: error "division by zero"`}
      </pre>

      <h2>Step 3: Get Project Structure</h2>
      <p>For each language, get the recommended file structure:</p>

      <pre className="code-block">
{`> Get the project structure for calculator in Go

generated/go/calculator/
├── cmd/
│   └── main.go
├── calculator.go
├── calculator_test.go
├── go.mod
└── README.md

> Get the project structure for calculator in Rust

generated/rust/calculator/
├── src/
│   ├── lib.rs
│   └── main.rs
├── tests/
│   └── calculator_test.rs
├── Cargo.toml
└── README.md`}
      </pre>

      <h2>Step 4: Generate Each Implementation</h2>
      <p>Generate code for each target language:</p>

      <pre className="code-block">
{`> Generate Go implementation from calculator.spec.md
> Generate Rust implementation from calculator.spec.md
> Generate TypeScript implementation from calculator.spec.md
> Generate Python implementation from calculator.spec.md`}
      </pre>

      <h2>Language-Specific Differences</h2>
      <p>
        Notice how RPG adapts to each language's idioms:
      </p>

      <h3>Go - Multiple Returns</h3>
      <pre className="code-block">
{`func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}`}
      </pre>

      <h3>Rust - Result Type</h3>
      <pre className="code-block">
{`pub fn divide(a: f64, b: f64) -> Result<f64, CalculatorError> {
    if b == 0.0 {
        return Err(CalculatorError::DivisionByZero);
    }
    Ok(a / b)
}`}
      </pre>

      <h3>TypeScript - Exceptions</h3>
      <pre className="code-block">
{`export function divide(a: number, b: number): number {
  if (b === 0) {
    throw new Error('division by zero');
  }
  return a / b;
}`}
      </pre>

      <h3>Python - Type Hints</h3>
      <pre className="code-block">
{`def divide(a: float, b: float) -> float:
    if b == 0:
        raise ValueError("division by zero")
    return a / b`}
      </pre>

      <h2>Next Steps</h2>
      <p>
        Ready to build something more complex? Try the{' '}
        <a href="/tutorials/rest-api">Building a REST API</a> tutorial.
      </p>
    </div>
  );
}
