export default function Go() {
  return (
    <div className="prose-docs">
      <h1>Go</h1>
      <p className="lead">
        RPG generates idiomatic Go code following community standards with explicit error handling,
        clear package structure, and comprehensive testing.
      </p>

      <div className="not-prose my-6 p-4 rounded-lg bg-cyan-50 dark:bg-cyan-900/20 border border-cyan-200 dark:border-cyan-800">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-cyan-800 dark:text-cyan-300 font-semibold">Philosophy</span>
        </div>
        <p className="text-cyan-700 dark:text-cyan-400 text-sm">
          Go values simplicity, clarity, and explicit behavior. RPG generates code that feels
          natural to Go developers: no magic, no hidden complexity, just clear and maintainable code.
        </p>
      </div>

      <h2>Naming Conventions</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Element</th>
            <th>Convention</th>
            <th>Example</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Package</td>
            <td>Lowercase, single word</td>
            <td><code>shortener</code>, <code>auth</code></td>
          </tr>
          <tr>
            <td>Public Function</td>
            <td>PascalCase</td>
            <td><code>CreateShortURL</code>, <code>ValidateToken</code></td>
          </tr>
          <tr>
            <td>Private Function</td>
            <td>camelCase</td>
            <td><code>generateCode</code>, <code>hashPassword</code></td>
          </tr>
          <tr>
            <td>Types/Structs</td>
            <td>PascalCase</td>
            <td><code>ShortURL</code>, <code>UserAccount</code></td>
          </tr>
          <tr>
            <td>Interfaces</td>
            <td>PascalCase, often -er suffix</td>
            <td><code>Storage</code>, <code>Validator</code>, <code>Reader</code></td>
          </tr>
          <tr>
            <td>Constants</td>
            <td>PascalCase or ALL_CAPS</td>
            <td><code>MaxRetries</code>, <code>DEFAULT_PORT</code></td>
          </tr>
          <tr>
            <td>Variables</td>
            <td>camelCase</td>
            <td><code>userID</code>, <code>requestCount</code></td>
          </tr>
        </tbody>
      </table>

      <h2>Type Mappings</h2>
      <p>RPG maps spec pseudo-types to Go types:</p>

      <table className="w-full">
        <thead>
          <tr>
            <th>Spec Type</th>
            <th>Go Type</th>
            <th>Notes</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>string</code></td>
            <td><code>string</code></td>
            <td>UTF-8 encoded</td>
          </tr>
          <tr>
            <td><code>number</code></td>
            <td><code>float64</code></td>
            <td>IEEE 754 double precision</td>
          </tr>
          <tr>
            <td><code>integer</code></td>
            <td><code>int</code></td>
            <td>Platform-dependent size (32/64 bit)</td>
          </tr>
          <tr>
            <td><code>boolean</code></td>
            <td><code>bool</code></td>
            <td>true/false</td>
          </tr>
          <tr>
            <td><code>datetime</code></td>
            <td><code>time.Time</code></td>
            <td>From standard library</td>
          </tr>
          <tr>
            <td><code>array</code></td>
            <td><code>[]T</code></td>
            <td>Slice of element type T</td>
          </tr>
          <tr>
            <td><code>map</code></td>
            <td><code>map[K]V</code></td>
            <td>Key type K, value type V</td>
          </tr>
          <tr>
            <td><code>optional</code></td>
            <td><code>*T</code></td>
            <td>Pointer for nullable types</td>
          </tr>
        </tbody>
      </table>

      <h2>Error Handling</h2>
      <p>
        Go uses explicit error returns as the last return value. RPG generates consistent
        error handling patterns:
      </p>

      <pre className="code-block">
{`// Function that can fail returns (result, error)
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Custom error types for rich error handling
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// Error wrapping for context preservation
func ProcessRequest(data []byte) error {
    user, err := parseUser(data)
    if err != nil {
        return fmt.Errorf("processing request: %w", err)
    }
    return nil
}`}
      </pre>

      <h2>Project Structure</h2>
      <p>RPG generates Go projects following the standard layout:</p>

      <pre className="code-block">
{`project-name/
├── cmd/
│   └── main.go              # Application entry point
├── internal/                # Private application code
│   ├── handler/            # HTTP handlers
│   │   └── handler.go
│   ├── service/            # Business logic
│   │   ├── service.go
│   │   └── service_test.go
│   └── storage/            # Data persistence
│       └── storage.go
├── pkg/                    # Public libraries (optional)
│   └── utils/
│       └── utils.go
├── go.mod                  # Module definition
├── go.sum                  # Dependency checksums
└── README.md`}
      </pre>

      <h2>Testing Patterns</h2>
      <p>
        Go uses table-driven tests with the built-in <code>testing</code> package:
      </p>

      <pre className="code-block">
{`func TestSlugify(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "simple text",
            input:    "Hello World",
            expected: "hello-world",
        },
        {
            name:     "special characters",
            input:    "Hello, World!",
            expected: "hello-world",
        },
        {
            name:     "multiple spaces",
            input:    "Hello    World",
            expected: "hello-world",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Slugify(tt.input)
            if got != tt.expected {
                t.Errorf("Slugify(%q) = %q, want %q",
                    tt.input, got, tt.expected)
            }
        })
    }
}`}
      </pre>

      <h2>Generated Code Example</h2>
      <p>Given this spec:</p>

      <pre className="code-block">
{`### createUser
Creates a new user account.

**Accepts:**
- email (string): User's email address
- password (string): User's password

**Returns:**
- User: The created user
- error: If validation fails`}
      </pre>

      <p>RPG generates:</p>

      <pre className="code-block">
{`package user

import (
    "errors"
    "fmt"
    "regexp"
    "time"
)

// User represents a user account
type User struct {
    ID        string    \`json:"id"\`
    Email     string    \`json:"email"\`
    Password  string    \`json:"-"\`
    CreatedAt time.Time \`json:"created_at"\`
}

// CreateUser creates a new user account
func CreateUser(email, password string) (*User, error) {
    if err := validateEmail(email); err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }

    if len(password) < 8 {
        return nil, errors.New("password must be at least 8 characters")
    }

    hashedPassword, err := hashPassword(password)
    if err != nil {
        return nil, fmt.Errorf("hashing password: %w", err)
    }

    return &User{
        ID:        generateID(),
        Email:     email,
        Password:  hashedPassword,
        CreatedAt: time.Now(),
    }, nil
}`}
      </pre>

      <h2>Best Practices</h2>
      <ul>
        <li><strong>Explicit over implicit</strong> — Always handle errors explicitly, never panic in library code</li>
        <li><strong>Simplicity</strong> — Prefer simple solutions over clever ones</li>
        <li><strong>Composition</strong> — Use interfaces and composition over inheritance</li>
        <li><strong>Documentation</strong> — Document all exported functions with godoc comments</li>
        <li><strong>Formatting</strong> — All generated code passes <code>go fmt</code></li>
        <li><strong>Linting</strong> — Code passes <code>golangci-lint</code> with standard config</li>
      </ul>

      <h2>Common Dependencies</h2>
      <p>RPG may include these common Go dependencies:</p>
      <ul>
        <li><code>github.com/gorilla/mux</code> — HTTP routing</li>
        <li><code>github.com/stretchr/testify</code> — Enhanced testing assertions</li>
        <li><code>github.com/spf13/cobra</code> — CLI applications</li>
        <li><code>go.uber.org/zap</code> — Structured logging</li>
      </ul>
    </div>
  );
}
