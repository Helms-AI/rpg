export default function Java() {
  return (
    <div className="prose-docs">
      <h1>Java</h1>
      <p className="lead">
        RPG generates modern Java 17+ code with records, sealed types, pattern matching,
        and enterprise-grade patterns following Oracle and community conventions.
      </p>

      <div className="not-prose my-6 p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-red-800 dark:text-red-300 font-semibold">Philosophy</span>
        </div>
        <p className="text-red-700 dark:text-red-400 text-sm">
          Java emphasizes readability, maintainability, and enterprise-grade reliability.
          RPG generates code that leverages modern Java features while maintaining compatibility
          with established patterns and frameworks.
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
            <td>Packages</td>
            <td>Lowercase, dot-separated</td>
            <td><code>com.example.userservice</code></td>
          </tr>
          <tr>
            <td>Classes</td>
            <td>PascalCase</td>
            <td><code>UserService</code>, <code>ValidationException</code></td>
          </tr>
          <tr>
            <td>Interfaces</td>
            <td>PascalCase (no I prefix)</td>
            <td><code>Repository</code>, <code>Validator</code></td>
          </tr>
          <tr>
            <td>Methods</td>
            <td>camelCase</td>
            <td><code>createUser</code>, <code>validateEmail</code></td>
          </tr>
          <tr>
            <td>Variables</td>
            <td>camelCase</td>
            <td><code>userId</code>, <code>requestCount</code></td>
          </tr>
          <tr>
            <td>Constants</td>
            <td>SCREAMING_SNAKE_CASE</td>
            <td><code>MAX_RETRIES</code>, <code>DEFAULT_PORT</code></td>
          </tr>
          <tr>
            <td>Enums</td>
            <td>PascalCase (type), SCREAMING_SNAKE_CASE (values)</td>
            <td><code>Status.ACTIVE</code>, <code>Role.ADMIN</code></td>
          </tr>
          <tr>
            <td>Generics</td>
            <td>Single uppercase letter</td>
            <td><code>T</code>, <code>K</code>, <code>V</code>, <code>E</code></td>
          </tr>
        </tbody>
      </table>

      <h2>Type Mappings</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Spec Type</th>
            <th>Java Type</th>
            <th>Notes</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>string</code></td>
            <td><code>String</code></td>
            <td>Immutable string class</td>
          </tr>
          <tr>
            <td><code>number</code></td>
            <td><code>double</code></td>
            <td>64-bit floating point</td>
          </tr>
          <tr>
            <td><code>integer</code></td>
            <td><code>int</code> or <code>long</code></td>
            <td>32-bit or 64-bit signed</td>
          </tr>
          <tr>
            <td><code>boolean</code></td>
            <td><code>boolean</code></td>
            <td>Primitive boolean</td>
          </tr>
          <tr>
            <td><code>datetime</code></td>
            <td><code>Instant</code> or <code>LocalDateTime</code></td>
            <td>From java.time package</td>
          </tr>
          <tr>
            <td><code>array</code></td>
            <td><code>List&lt;T&gt;</code></td>
            <td>Usually ArrayList</td>
          </tr>
          <tr>
            <td><code>map</code></td>
            <td><code>Map&lt;K, V&gt;</code></td>
            <td>Usually HashMap</td>
          </tr>
          <tr>
            <td><code>optional</code></td>
            <td><code>Optional&lt;T&gt;</code></td>
            <td>Container for nullable values</td>
          </tr>
        </tbody>
      </table>

      <h2>Error Handling</h2>
      <p>
        Java uses exceptions for error handling. RPG generates custom exception hierarchies
        with checked exceptions for recoverable errors:
      </p>

      <pre className="code-block">
{`// Base exception for application errors
public class AppException extends Exception {
    private final String code;

    public AppException(String message, String code) {
        super(message);
        this.code = code;
    }

    public AppException(String message, String code, Throwable cause) {
        super(message, cause);
        this.code = code;
    }

    public String getCode() {
        return code;
    }
}

// Validation exception with field information
public class ValidationException extends AppException {
    private final String field;

    public ValidationException(String message, String field) {
        super(message, "VALIDATION_ERROR");
        this.field = field;
    }

    public String getField() {
        return field;
    }
}

// Not found exception
public class NotFoundException extends AppException {
    public NotFoundException(String resource) {
        super(resource + " not found", "NOT_FOUND");
    }
}

// Function that throws checked exception
public double divide(double a, double b) throws ValidationException {
    if (b == 0) {
        throw new ValidationException("Cannot divide by zero", "divisor");
    }
    return a / b;
}`}
      </pre>

      <h2>Project Structure</h2>
      <pre className="code-block">
{`project-name/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/
│   │   │       └── example/
│   │   │           └── projectname/
│   │   │               ├── Application.java     # Entry point
│   │   │               ├── model/
│   │   │               │   └── User.java
│   │   │               ├── service/
│   │   │               │   └── UserService.java
│   │   │               ├── repository/
│   │   │               │   └── UserRepository.java
│   │   │               ├── exception/
│   │   │               │   ├── AppException.java
│   │   │               │   └── ValidationException.java
│   │   │               └── util/
│   │   │                   └── ValidationUtils.java
│   │   └── resources/
│   │       └── application.properties
│   └── test/
│       └── java/
│           └── com/
│               └── example/
│                   └── projectname/
│                       └── service/
│                           └── UserServiceTest.java
├── pom.xml                                      # Maven config
├── build.gradle                                 # Gradle config (alternative)
└── README.md`}
      </pre>

      <h2>Testing Patterns</h2>
      <pre className="code-block">
{`import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import static org.assertj.core.api.Assertions.*;
import static org.mockito.Mockito.*;

class UserServiceTest {

    @Mock
    private UserRepository userRepository;

    private UserService userService;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
        userService = new UserService(userRepository);
    }

    @Nested
    @DisplayName("createUser")
    class CreateUserTests {

        @Test
        @DisplayName("creates user with valid input")
        void createsUserWithValidInput() throws ValidationException {
            // Arrange
            when(userRepository.save(any())).thenAnswer(i -> {
                User user = i.getArgument(0);
                return user.withId("123");
            });

            // Act
            User user = userService.createUser("test@example.com", "password123");

            // Assert
            assertThat(user.getId()).isEqualTo("123");
            assertThat(user.getEmail()).isEqualTo("test@example.com");
            assertThat(user.getCreatedAt()).isNotNull();
        }

        @Test
        @DisplayName("throws ValidationException for invalid email")
        void throwsValidationExceptionForInvalidEmail() {
            assertThatThrownBy(() -> userService.createUser("invalid", "password123"))
                .isInstanceOf(ValidationException.class)
                .hasFieldOrPropertyWithValue("field", "email");
        }

        @Test
        @DisplayName("throws ValidationException for short password")
        void throwsValidationExceptionForShortPassword() {
            assertThatThrownBy(() -> userService.createUser("test@example.com", "123"))
                .isInstanceOf(ValidationException.class)
                .hasMessageContaining("at least 8 characters");
        }
    }

    @ParameterizedTest
    @CsvSource({
        "Hello World, hello-world",
        "'Hello, World!', hello-world",
        "Hello    World, hello-world"
    })
    @DisplayName("slugify converts text correctly")
    void slugifyConvertsTextCorrectly(String input, String expected) {
        assertThat(StringUtils.slugify(input)).isEqualTo(expected);
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
{`package com.example.user;

import java.time.Instant;
import java.util.UUID;
import java.util.regex.Pattern;

/**
 * Represents a user account.
 */
public record User(
    String id,
    String email,
    Instant createdAt
) {
    public User withId(String newId) {
        return new User(newId, email, createdAt);
    }
}

/**
 * Exception thrown when validation fails.
 */
public class ValidationException extends Exception {
    private final String field;

    public ValidationException(String message, String field) {
        super(message);
        this.field = field;
    }

    public String getField() {
        return field;
    }
}

/**
 * Service for managing user accounts.
 */
public class UserService {
    private static final Pattern EMAIL_PATTERN =
        Pattern.compile("^[^\\\\s@]+@[^\\\\s@]+\\\\.[^\\\\s@]+$");

    private final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    /**
     * Creates a new user account.
     *
     * @param email User's email address
     * @param password User's password
     * @return The created user
     * @throws ValidationException If email or password is invalid
     */
    public User createUser(String email, String password) throws ValidationException {
        // Validate email
        if (!EMAIL_PATTERN.matcher(email).matches()) {
            throw new ValidationException("Invalid email format", "email");
        }

        // Validate password
        if (password.length() < 8) {
            throw new ValidationException(
                "Password must be at least 8 characters",
                "password"
            );
        }

        // Create user (password would be hashed in real implementation)
        User user = new User(
            UUID.randomUUID().toString(),
            email,
            Instant.now()
        );

        return userRepository.save(user);
    }
}`}
      </pre>

      <h2>Modern Java Features</h2>
      <p>
        RPG leverages modern Java 17+ features when generating code:
      </p>

      <pre className="code-block">
{`// Records for immutable data classes
public record User(String id, String email, Instant createdAt) {}

// Sealed types for restricted hierarchies
public sealed interface Shape
    permits Circle, Rectangle, Triangle {}

public record Circle(double radius) implements Shape {}
public record Rectangle(double width, double height) implements Shape {}
public record Triangle(double base, double height) implements Shape {}

// Pattern matching in switch expressions
public double area(Shape shape) {
    return switch (shape) {
        case Circle c -> Math.PI * c.radius() * c.radius();
        case Rectangle r -> r.width() * r.height();
        case Triangle t -> 0.5 * t.base() * t.height();
    };
}

// Text blocks for multi-line strings
String json = """
    {
        "name": "%s",
        "email": "%s"
    }
    """.formatted(name, email);

// var for local type inference
var users = new ArrayList<User>();
var result = userService.createUser(email, password);`}
      </pre>

      <h2>Best Practices</h2>
      <ul>
        <li><strong>Immutability</strong> — Prefer records and final fields for data classes</li>
        <li><strong>Null safety</strong> — Use Optional for nullable return values</li>
        <li><strong>Checked exceptions</strong> — Use for recoverable errors in public APIs</li>
        <li><strong>Dependency injection</strong> — Constructor injection for dependencies</li>
        <li><strong>Documentation</strong> — JavaDoc for all public classes and methods</li>
        <li><strong>Formatting</strong> — Code follows Google Java Style Guide</li>
        <li><strong>Testing</strong> — JUnit 5 with AssertJ and Mockito</li>
      </ul>

      <h2>Common Dependencies</h2>
      <ul>
        <li><code>org.junit.jupiter</code> — JUnit 5 testing framework</li>
        <li><code>org.assertj</code> — Fluent assertions library</li>
        <li><code>org.mockito</code> — Mocking framework</li>
        <li><code>com.fasterxml.jackson</code> — JSON processing</li>
        <li><code>org.slf4j</code> — Logging facade</li>
        <li><code>com.google.guava</code> — Common utilities</li>
      </ul>
    </div>
  );
}
