export default function RestApi() {
  return (
    <div className="prose-docs">
      <h1>Building a REST API</h1>
      <p className="lead">
        Create a complete URL shortener REST API from spec to working implementation.
      </p>

      <div className="not-prose my-6 p-4 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
        <p className="text-blue-800 dark:text-blue-300 text-sm">
          <strong>Time:</strong> 20 minutes &nbsp;|&nbsp; <strong>Difficulty:</strong> Intermediate
        </p>
      </div>

      <h2>What You'll Build</h2>
      <p>
        A URL shortener service with endpoints to create short URLs and redirect to originals.
        This demonstrates how RPG handles more complex specifications with types, dependencies,
        and configuration.
      </p>

      <h2>Step 1: Design the Spec</h2>
      <p>Create <code>url-shortener.spec.md</code>:</p>

      <pre className="code-block">
{`# URL Shortener

A service that creates short URLs and redirects to original URLs.

## Meta

- **Version**: 1.0.0
- **Author**: Your Team

## Target Languages

- go
- typescript

## Types

### ShortURL
Represents a shortened URL mapping.

- \`id\` (string): Unique short code
- \`originalUrl\` (string): The original long URL
- \`createdAt\` (datetime): When the short URL was created
- \`clicks\` (number): Number of times the short URL was accessed

### CreateRequest
Request body for creating a short URL.

- \`url\` (string): The URL to shorten

### CreateResponse
Response after creating a short URL.

- \`shortUrl\` (string): The full short URL
- \`shortCode\` (string): Just the short code

## Dependencies

- HTTP server framework
- In-memory storage (or database)
- URL validation library

## Configuration

- \`PORT\`: Server port (default: 8080)
- \`BASE_URL\`: Base URL for short links (default: http://localhost:8080)

## Functions

### createShortUrl
Creates a new short URL.

**Accepts:**
- \`request\` (CreateRequest): The creation request

**Returns:**
- CreateResponse: The created short URL info
- error: If URL is invalid

**Logic:**
1. Validate the URL format
2. Generate a unique 6-character code
3. Store the mapping
4. Return the short URL

### getOriginalUrl
Retrieves the original URL for a short code.

**Accepts:**
- \`shortCode\` (string): The short code to look up

**Returns:**
- string: The original URL
- error: If short code not found

**Logic:**
1. Look up the short code in storage
2. Increment click counter
3. Return the original URL

### handleRedirect
HTTP handler that redirects to the original URL.

**Accepts:**
- \`shortCode\` (string): From URL path

**Returns:**
- HTTP 302 redirect to original URL
- HTTP 404 if not found

## API Endpoints

### POST /api/shorten
Creates a new short URL.

Request:
\`\`\`json
{ "url": "https://example.com/very/long/url" }
\`\`\`

Response:
\`\`\`json
{
  "shortUrl": "http://localhost:8080/abc123",
  "shortCode": "abc123"
}
\`\`\`

### GET /{shortCode}
Redirects to the original URL.

## Tests

### Create short URL
- **Given**: POST /api/shorten with url "https://example.com"
- **Expect**: 200 with shortCode in response

### Redirect works
- **Given**: GET /abc123 (existing short code)
- **Expect**: 302 redirect to original URL

### Invalid URL
- **Given**: POST /api/shorten with url "not-a-url"
- **Expect**: 400 error

### Not found
- **Given**: GET /nonexistent
- **Expect**: 404 error`}
      </pre>

      <h2>Step 2: Get Generation Context</h2>
      <p>Let the AI understand the full context:</p>

      <pre className="code-block">
{`> Get generation context for url-shortener.spec.md in Go`}
      </pre>

      <h2>Step 3: Generate the Implementation</h2>
      <pre className="code-block">
{`> Generate a complete Go implementation of the URL shortener`}
      </pre>

      <p>The AI will create:</p>
      <ul>
        <li><code>cmd/main.go</code> - Server entry point</li>
        <li><code>internal/shortener/shortener.go</code> - Core logic</li>
        <li><code>internal/shortener/handlers.go</code> - HTTP handlers</li>
        <li><code>internal/shortener/storage.go</code> - In-memory storage</li>
        <li><code>internal/shortener/shortener_test.go</code> - Unit tests</li>
        <li><code>go.mod</code> - Module definition</li>
      </ul>

      <h2>Step 4: Generate TypeScript Version</h2>
      <pre className="code-block">
{`> Now generate a TypeScript/Express version`}
      </pre>

      <h2>Step 5: Verify Parity</h2>
      <p>Ensure both implementations behave identically:</p>

      <pre className="code-block">
{`> Use RPG to check parity between the Go and TypeScript implementations

Parity Report:
- Overall: 100%
- Functions: All implemented
- Types: All defined
- Tests: All passing
- API Endpoints: Identical behavior`}
      </pre>

      <h2>Key Concepts Demonstrated</h2>
      <ul>
        <li><strong>Types</strong> - Structured data definitions become language-appropriate types</li>
        <li><strong>Dependencies</strong> - RPG selects idiomatic libraries per language</li>
        <li><strong>Configuration</strong> - Environment variables handled appropriately</li>
        <li><strong>API Endpoints</strong> - HTTP routing follows language conventions</li>
        <li><strong>Error Handling</strong> - Errors adapt to language patterns</li>
      </ul>

      <h2>Next Steps</h2>
      <p>
        Learn to maintain consistency across implementations in the{' '}
        <a href="/tutorials/parity-checking">Feature Parity Checking</a> tutorial.
      </p>
    </div>
  );
}
