export default function CSharp() {
  return (
    <div className="prose-docs">
      <h1>C#</h1>
      <p className="lead">
        RPG generates modern C# 12+ code with records, pattern matching, nullable reference types,
        and .NET conventions following Microsoft's official guidelines.
      </p>

      <div className="not-prose my-6 p-4 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-800">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-purple-800 dark:text-purple-300 font-semibold">Philosophy</span>
        </div>
        <p className="text-purple-700 dark:text-purple-400 text-sm">
          C# balances power and productivity with a strong type system and modern language features.
          RPG generates code that leverages the latest C# capabilities while maintaining .NET ecosystem
          compatibility and idiomatic patterns.
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
            <td>Namespaces</td>
            <td>PascalCase, dot-separated</td>
            <td><code>MyCompany.UserService</code></td>
          </tr>
          <tr>
            <td>Classes</td>
            <td>PascalCase</td>
            <td><code>UserService</code>, <code>ValidationException</code></td>
          </tr>
          <tr>
            <td>Interfaces</td>
            <td>PascalCase with I prefix</td>
            <td><code>IRepository</code>, <code>IValidator</code></td>
          </tr>
          <tr>
            <td>Methods</td>
            <td>PascalCase</td>
            <td><code>CreateUser</code>, <code>ValidateEmail</code></td>
          </tr>
          <tr>
            <td>Properties</td>
            <td>PascalCase</td>
            <td><code>UserId</code>, <code>CreatedAt</code></td>
          </tr>
          <tr>
            <td>Private fields</td>
            <td>_camelCase with underscore</td>
            <td><code>_userRepository</code>, <code>_logger</code></td>
          </tr>
          <tr>
            <td>Parameters</td>
            <td>camelCase</td>
            <td><code>userId</code>, <code>requestCount</code></td>
          </tr>
          <tr>
            <td>Constants</td>
            <td>PascalCase</td>
            <td><code>MaxRetries</code>, <code>DefaultPort</code></td>
          </tr>
          <tr>
            <td>Enums</td>
            <td>PascalCase (type and values)</td>
            <td><code>Status.Active</code>, <code>Role.Admin</code></td>
          </tr>
        </tbody>
      </table>

      <h2>Type Mappings</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Spec Type</th>
            <th>C# Type</th>
            <th>Notes</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>string</code></td>
            <td><code>string</code></td>
            <td>Immutable reference type</td>
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
            <td><code>bool</code></td>
            <td>Boolean value type</td>
          </tr>
          <tr>
            <td><code>datetime</code></td>
            <td><code>DateTime</code> or <code>DateTimeOffset</code></td>
            <td>System.DateTime</td>
          </tr>
          <tr>
            <td><code>array</code></td>
            <td><code>List&lt;T&gt;</code> or <code>IEnumerable&lt;T&gt;</code></td>
            <td>Generic collections</td>
          </tr>
          <tr>
            <td><code>map</code></td>
            <td><code>Dictionary&lt;K, V&gt;</code></td>
            <td>Generic dictionary</td>
          </tr>
          <tr>
            <td><code>optional</code></td>
            <td><code>T?</code></td>
            <td>Nullable reference/value types</td>
          </tr>
        </tbody>
      </table>

      <h2>Error Handling</h2>
      <p>
        C# uses exceptions for error handling. RPG generates custom exception hierarchies
        following .NET conventions:
      </p>

      <pre className="code-block">
{`// Base exception for application errors
public class AppException : Exception
{
    public string Code { get; }

    public AppException(string message, string code)
        : base(message)
    {
        Code = code;
    }

    public AppException(string message, string code, Exception innerException)
        : base(message, innerException)
    {
        Code = code;
    }
}

// Validation exception with field information
public class ValidationException : AppException
{
    public string Field { get; }

    public ValidationException(string message, string field)
        : base(message, "VALIDATION_ERROR")
    {
        Field = field;
    }
}

// Not found exception
public class NotFoundException : AppException
{
    public NotFoundException(string resource)
        : base($"{resource} not found", "NOT_FOUND")
    {
    }
}

// Function that throws
public double Divide(double a, double b)
{
    if (b == 0)
    {
        throw new ValidationException("Cannot divide by zero", "divisor");
    }
    return a / b;
}

// Result pattern for functional approach
public readonly record struct Result<T>(T? Value, string? Error)
{
    public bool IsSuccess => Error is null;
    public static Result<T> Success(T value) => new(value, null);
    public static Result<T> Failure(string error) => new(default, error);
}`}
      </pre>

      <h2>Project Structure</h2>
      <pre className="code-block">
{`ProjectName/
├── src/
│   └── ProjectName/
│       ├── ProjectName.csproj      # Project file
│       ├── Program.cs              # Entry point
│       ├── Models/
│       │   └── User.cs
│       ├── Services/
│       │   ├── IUserService.cs
│       │   └── UserService.cs
│       ├── Repositories/
│       │   ├── IUserRepository.cs
│       │   └── UserRepository.cs
│       ├── Exceptions/
│       │   ├── AppException.cs
│       │   └── ValidationException.cs
│       └── Extensions/
│           └── StringExtensions.cs
├── tests/
│   └── ProjectName.Tests/
│       ├── ProjectName.Tests.csproj
│       ├── Services/
│       │   └── UserServiceTests.cs
│       └── Usings.cs               # Global usings
├── ProjectName.sln                 # Solution file
└── README.md`}
      </pre>

      <h2>Testing Patterns</h2>
      <pre className="code-block">
{`using FluentAssertions;
using Moq;
using Xunit;

namespace ProjectName.Tests.Services;

public class UserServiceTests
{
    private readonly Mock<IUserRepository> _mockRepository;
    private readonly UserService _userService;

    public UserServiceTests()
    {
        _mockRepository = new Mock<IUserRepository>();
        _userService = new UserService(_mockRepository.Object);
    }

    [Fact]
    public async Task CreateUser_WithValidInput_ReturnsUser()
    {
        // Arrange
        _mockRepository
            .Setup(r => r.SaveAsync(It.IsAny<User>()))
            .ReturnsAsync((User u) => u with { Id = "123" });

        // Act
        var user = await _userService.CreateUserAsync("test@example.com", "password123");

        // Assert
        user.Id.Should().Be("123");
        user.Email.Should().Be("test@example.com");
        user.CreatedAt.Should().BeCloseTo(DateTime.UtcNow, TimeSpan.FromSeconds(1));
    }

    [Fact]
    public async Task CreateUser_WithInvalidEmail_ThrowsValidationException()
    {
        // Act
        var act = () => _userService.CreateUserAsync("invalid", "password123");

        // Assert
        await act.Should()
            .ThrowAsync<ValidationException>()
            .Where(e => e.Field == "email");
    }

    [Fact]
    public async Task CreateUser_WithShortPassword_ThrowsValidationException()
    {
        // Act
        var act = () => _userService.CreateUserAsync("test@example.com", "123");

        // Assert
        await act.Should()
            .ThrowAsync<ValidationException>()
            .WithMessage("*at least 8 characters*");
    }

    [Theory]
    [InlineData("Hello World", "hello-world")]
    [InlineData("Hello, World!", "hello-world")]
    [InlineData("Hello    World", "hello-world")]
    public void Slugify_ConvertsTextCorrectly(string input, string expected)
    {
        // Act
        var result = input.Slugify();

        // Assert
        result.Should().Be(expected);
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
{`using System.Text.RegularExpressions;

namespace ProjectName.Models;

/// <summary>
/// Represents a user account.
/// </summary>
public record User(
    string Id,
    string Email,
    DateTime CreatedAt
);

namespace ProjectName.Exceptions;

/// <summary>
/// Exception thrown when validation fails.
/// </summary>
public class ValidationException : Exception
{
    public string Field { get; }

    public ValidationException(string message, string field)
        : base(message)
    {
        Field = field;
    }
}

namespace ProjectName.Services;

/// <summary>
/// Service for managing user accounts.
/// </summary>
public class UserService : IUserService
{
    private static readonly Regex EmailRegex =
        new(@"^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$", RegexOptions.Compiled);

    private readonly IUserRepository _userRepository;

    public UserService(IUserRepository userRepository)
    {
        _userRepository = userRepository;
    }

    /// <summary>
    /// Creates a new user account.
    /// </summary>
    /// <param name="email">User's email address</param>
    /// <param name="password">User's password</param>
    /// <returns>The created user</returns>
    /// <exception cref="ValidationException">Thrown if email or password is invalid</exception>
    public async Task<User> CreateUserAsync(string email, string password)
    {
        // Validate email
        if (!EmailRegex.IsMatch(email))
        {
            throw new ValidationException("Invalid email format", "email");
        }

        // Validate password
        if (password.Length < 8)
        {
            throw new ValidationException(
                "Password must be at least 8 characters",
                "password"
            );
        }

        // Create user (password would be hashed in real implementation)
        var user = new User(
            Id: Guid.NewGuid().ToString(),
            Email: email,
            CreatedAt: DateTime.UtcNow
        );

        return await _userRepository.SaveAsync(user);
    }
}`}
      </pre>

      <h2>Modern C# Features</h2>
      <p>
        RPG leverages modern C# 12+ features when generating code:
      </p>

      <pre className="code-block">
{`// Records for immutable data classes
public record User(string Id, string Email, DateTime CreatedAt);

// Record structs for value semantics
public readonly record struct Point(double X, double Y);

// Primary constructors (C# 12)
public class UserService(IUserRepository repository, ILogger<UserService> logger)
{
    public async Task<User> GetUserAsync(string id) =>
        await repository.FindByIdAsync(id)
            ?? throw new NotFoundException("User");
}

// Pattern matching with switch expressions
public string GetStatusMessage(UserStatus status) => status switch
{
    UserStatus.Active => "User is active",
    UserStatus.Pending => "User is pending verification",
    UserStatus.Suspended => "User is suspended",
    _ => throw new ArgumentOutOfRangeException(nameof(status))
};

// Nullable reference types
public async Task<User?> FindUserByEmailAsync(string email)
{
    return await _context.Users.FirstOrDefaultAsync(u => u.Email == email);
}

// Collection expressions (C# 12)
List<int> numbers = [1, 2, 3, 4, 5];
int[] moreNumbers = [..numbers, 6, 7, 8];

// Raw string literals
var json = """
    {
        "name": "{name}",
        "email": "{email}"
    }
    """;

// Required properties
public class CreateUserRequest
{
    public required string Email { get; init; }
    public required string Password { get; init; }
}

// File-scoped namespaces
namespace ProjectName.Models;

public record User(string Id, string Email, DateTime CreatedAt);`}
      </pre>

      <h2>Async/Await Patterns</h2>
      <pre className="code-block">
{`// Async method naming convention (Async suffix)
public async Task<User> CreateUserAsync(string email, string password)
{
    // Async operations
    var hashedPassword = await HashPasswordAsync(password);
    var user = new User(Guid.NewGuid().ToString(), email, DateTime.UtcNow);
    return await _repository.SaveAsync(user);
}

// Async LINQ
public async Task<List<User>> GetActiveUsersAsync()
{
    return await _context.Users
        .Where(u => u.Status == UserStatus.Active)
        .OrderBy(u => u.CreatedAt)
        .ToListAsync();
}

// Cancellation token support
public async Task<User> GetUserAsync(
    string id,
    CancellationToken cancellationToken = default)
{
    return await _repository.FindByIdAsync(id, cancellationToken)
        ?? throw new NotFoundException("User");
}

// Async disposal
public class DbConnection : IAsyncDisposable
{
    public async ValueTask DisposeAsync()
    {
        await CloseAsync();
        GC.SuppressFinalize(this);
    }
}`}
      </pre>

      <h2>Best Practices</h2>
      <ul>
        <li><strong>Nullable reference types</strong> — Enable nullable context for null safety</li>
        <li><strong>Records</strong> — Use records for immutable data transfer objects</li>
        <li><strong>Async/await</strong> — Use async methods for I/O operations</li>
        <li><strong>Dependency injection</strong> — Constructor injection via interfaces</li>
        <li><strong>XML documentation</strong> — Document all public APIs with XML comments</li>
        <li><strong>Formatting</strong> — Code follows .NET coding conventions</li>
        <li><strong>Testing</strong> — xUnit with FluentAssertions and Moq</li>
      </ul>

      <h2>Common Dependencies</h2>
      <ul>
        <li><code>xunit</code> — Testing framework</li>
        <li><code>FluentAssertions</code> — Fluent assertion library</li>
        <li><code>Moq</code> — Mocking framework</li>
        <li><code>System.Text.Json</code> — JSON serialization (built-in)</li>
        <li><code>Microsoft.Extensions.Logging</code> — Logging abstractions</li>
        <li><code>Microsoft.Extensions.DependencyInjection</code> — DI container</li>
      </ul>
    </div>
  );
}
