export default function Rust() {
  return (
    <div className="prose-docs">
      <h1>Rust</h1>
      <p className="lead">
        RPG generates safe, idiomatic Rust code with proper ownership semantics, Result/Option types,
        and zero-cost abstractions.
      </p>

      <div className="not-prose my-6 p-4 rounded-lg bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-800">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-orange-800 dark:text-orange-300 font-semibold">Philosophy</span>
        </div>
        <p className="text-orange-700 dark:text-orange-400 text-sm">
          Rust prioritizes safety, performance, and fearless concurrency. RPG generates code that
          leverages Rust's type system to eliminate entire classes of bugs at compile time.
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
            <td>Crate/Module</td>
            <td>snake_case</td>
            <td><code>url_shortener</code>, <code>auth_service</code></td>
          </tr>
          <tr>
            <td>Functions</td>
            <td>snake_case</td>
            <td><code>create_user</code>, <code>validate_token</code></td>
          </tr>
          <tr>
            <td>Types/Structs</td>
            <td>PascalCase</td>
            <td><code>ShortUrl</code>, <code>UserAccount</code></td>
          </tr>
          <tr>
            <td>Enums</td>
            <td>PascalCase (type and variants)</td>
            <td><code>Status::Active</code>, <code>Error::NotFound</code></td>
          </tr>
          <tr>
            <td>Traits</td>
            <td>PascalCase</td>
            <td><code>Validator</code>, <code>Serializable</code></td>
          </tr>
          <tr>
            <td>Constants</td>
            <td>SCREAMING_SNAKE_CASE</td>
            <td><code>MAX_RETRIES</code>, <code>DEFAULT_PORT</code></td>
          </tr>
          <tr>
            <td>Variables</td>
            <td>snake_case</td>
            <td><code>user_id</code>, <code>request_count</code></td>
          </tr>
        </tbody>
      </table>

      <h2>Type Mappings</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Spec Type</th>
            <th>Rust Type</th>
            <th>Notes</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>string</code></td>
            <td><code>String</code> / <code>&str</code></td>
            <td>Owned vs borrowed</td>
          </tr>
          <tr>
            <td><code>number</code></td>
            <td><code>f64</code></td>
            <td>64-bit floating point</td>
          </tr>
          <tr>
            <td><code>integer</code></td>
            <td><code>i64</code></td>
            <td>Signed 64-bit integer</td>
          </tr>
          <tr>
            <td><code>boolean</code></td>
            <td><code>bool</code></td>
            <td>true/false</td>
          </tr>
          <tr>
            <td><code>datetime</code></td>
            <td><code>chrono::DateTime&lt;Utc&gt;</code></td>
            <td>From chrono crate</td>
          </tr>
          <tr>
            <td><code>array</code></td>
            <td><code>Vec&lt;T&gt;</code></td>
            <td>Growable vector</td>
          </tr>
          <tr>
            <td><code>map</code></td>
            <td><code>HashMap&lt;K, V&gt;</code></td>
            <td>From std::collections</td>
          </tr>
          <tr>
            <td><code>optional</code></td>
            <td><code>Option&lt;T&gt;</code></td>
            <td>Some(value) or None</td>
          </tr>
          <tr>
            <td><code>error</code></td>
            <td><code>Result&lt;T, E&gt;</code></td>
            <td>Ok(value) or Err(error)</td>
          </tr>
        </tbody>
      </table>

      <h2>Error Handling</h2>
      <p>
        Rust uses the <code>Result</code> type for recoverable errors. RPG generates custom error
        types with the <code>thiserror</code> crate:
      </p>

      <pre className="code-block">
{`use thiserror::Error;

#[derive(Error, Debug)]
pub enum CalculatorError {
    #[error("division by zero")]
    DivisionByZero,

    #[error("invalid input: {0}")]
    InvalidInput(String),

    #[error("overflow occurred")]
    Overflow,
}

pub fn divide(a: f64, b: f64) -> Result<f64, CalculatorError> {
    if b == 0.0 {
        return Err(CalculatorError::DivisionByZero);
    }
    Ok(a / b)
}

// Using the ? operator for error propagation
pub fn calculate(input: &str) -> Result<f64, CalculatorError> {
    let values = parse_input(input)?;
    let result = divide(values.0, values.1)?;
    Ok(result)
}`}
      </pre>

      <h2>Project Structure</h2>
      <pre className="code-block">
{`project-name/
├── src/
│   ├── lib.rs              # Library entry point
│   ├── main.rs             # Binary entry point (optional)
│   ├── error.rs            # Error types
│   ├── models/
│   │   ├── mod.rs
│   │   └── user.rs
│   └── services/
│       ├── mod.rs
│       └── auth.rs
├── tests/
│   └── integration_test.rs # Integration tests
├── benches/
│   └── benchmark.rs        # Benchmarks
├── Cargo.toml              # Package manifest
├── Cargo.lock              # Dependency lock file
└── README.md`}
      </pre>

      <h2>Testing Patterns</h2>
      <pre className="code-block">
{`#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_slugify_simple() {
        assert_eq!(slugify("Hello World"), "hello-world");
    }

    #[test]
    fn test_slugify_special_chars() {
        assert_eq!(slugify("Hello, World!"), "hello-world");
    }

    #[test]
    fn test_divide_success() {
        let result = divide(10.0, 2.0);
        assert!(result.is_ok());
        assert_eq!(result.unwrap(), 5.0);
    }

    #[test]
    fn test_divide_by_zero() {
        let result = divide(10.0, 0.0);
        assert!(result.is_err());
        assert!(matches!(result.unwrap_err(), CalculatorError::DivisionByZero));
    }

    // Property-based testing with proptest
    #[cfg(feature = "proptest")]
    proptest! {
        #[test]
        fn test_slugify_never_panics(s in ".*") {
            let _ = slugify(&s);
        }
    }
}`}
      </pre>

      <h2>Generated Code Example</h2>
      <p>Given this spec:</p>
      <pre className="code-block">
{`### create_user
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
{`use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use thiserror::Error;
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    pub id: String,
    pub email: String,
    #[serde(skip_serializing)]
    password_hash: String,
    pub created_at: DateTime<Utc>,
}

#[derive(Error, Debug)]
pub enum UserError {
    #[error("invalid email format: {0}")]
    InvalidEmail(String),

    #[error("password must be at least 8 characters")]
    PasswordTooShort,

    #[error("failed to hash password: {0}")]
    HashError(String),
}

pub fn create_user(email: &str, password: &str) -> Result<User, UserError> {
    if !is_valid_email(email) {
        return Err(UserError::InvalidEmail(email.to_string()));
    }

    if password.len() < 8 {
        return Err(UserError::PasswordTooShort);
    }

    let password_hash = hash_password(password)
        .map_err(|e| UserError::HashError(e.to_string()))?;

    Ok(User {
        id: Uuid::new_v4().to_string(),
        email: email.to_string(),
        password_hash,
        created_at: Utc::now(),
    })
}`}
      </pre>

      <h2>Best Practices</h2>
      <ul>
        <li><strong>Ownership</strong> — Prefer borrowing over cloning when possible</li>
        <li><strong>Result over panic</strong> — Use Result for recoverable errors, reserve panic for unrecoverable bugs</li>
        <li><strong>Derive macros</strong> — Use #[derive(...)] for common traits like Debug, Clone, Serialize</li>
        <li><strong>Documentation</strong> — Use /// doc comments for public APIs</li>
        <li><strong>Formatting</strong> — All code passes <code>cargo fmt</code></li>
        <li><strong>Linting</strong> — All code passes <code>cargo clippy</code></li>
      </ul>

      <h2>Common Dependencies</h2>
      <ul>
        <li><code>serde</code> / <code>serde_json</code> — Serialization</li>
        <li><code>thiserror</code> — Error type derivation</li>
        <li><code>anyhow</code> — Application error handling</li>
        <li><code>tokio</code> — Async runtime</li>
        <li><code>chrono</code> — Date/time handling</li>
        <li><code>uuid</code> — UUID generation</li>
      </ul>
    </div>
  );
}
