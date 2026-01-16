export default function Python() {
  return (
    <div className="prose-docs">
      <h1>Python</h1>
      <p className="lead">
        RPG generates modern Python 3.10+ code with comprehensive type hints, dataclasses,
        and PEP 8 compliance.
      </p>

      <div className="not-prose my-6 p-4 rounded-lg bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-yellow-800 dark:text-yellow-300 font-semibold">Philosophy</span>
        </div>
        <p className="text-yellow-700 dark:text-yellow-400 text-sm">
          Python emphasizes readability, simplicity, and "there should be one obvious way to do it."
          RPG generates clean, Pythonic code that experienced developers will find natural.
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
            <td>Modules/Files</td>
            <td>snake_case</td>
            <td><code>user_service.py</code>, <code>auth_utils.py</code></td>
          </tr>
          <tr>
            <td>Functions</td>
            <td>snake_case</td>
            <td><code>create_user</code>, <code>validate_token</code></td>
          </tr>
          <tr>
            <td>Classes</td>
            <td>PascalCase</td>
            <td><code>UserService</code>, <code>ValidationError</code></td>
          </tr>
          <tr>
            <td>Variables</td>
            <td>snake_case</td>
            <td><code>user_id</code>, <code>request_count</code></td>
          </tr>
          <tr>
            <td>Constants</td>
            <td>SCREAMING_SNAKE_CASE</td>
            <td><code>MAX_RETRIES</code>, <code>DEFAULT_PORT</code></td>
          </tr>
          <tr>
            <td>Private</td>
            <td>Leading underscore</td>
            <td><code>_internal_method</code>, <code>_private_var</code></td>
          </tr>
          <tr>
            <td>Type Variables</td>
            <td>PascalCase, short</td>
            <td><code>T</code>, <code>K</code>, <code>V</code></td>
          </tr>
        </tbody>
      </table>

      <h2>Type Mappings</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Spec Type</th>
            <th>Python Type</th>
            <th>Notes</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>string</code></td>
            <td><code>str</code></td>
            <td>Unicode string</td>
          </tr>
          <tr>
            <td><code>number</code></td>
            <td><code>float</code></td>
            <td>64-bit floating point</td>
          </tr>
          <tr>
            <td><code>integer</code></td>
            <td><code>int</code></td>
            <td>Arbitrary precision</td>
          </tr>
          <tr>
            <td><code>boolean</code></td>
            <td><code>bool</code></td>
            <td>True/False</td>
          </tr>
          <tr>
            <td><code>datetime</code></td>
            <td><code>datetime</code></td>
            <td>From datetime module</td>
          </tr>
          <tr>
            <td><code>array</code></td>
            <td><code>list[T]</code></td>
            <td>Generic list (3.9+)</td>
          </tr>
          <tr>
            <td><code>map</code></td>
            <td><code>dict[K, V]</code></td>
            <td>Generic dict (3.9+)</td>
          </tr>
          <tr>
            <td><code>optional</code></td>
            <td><code>T | None</code></td>
            <td>Union syntax (3.10+)</td>
          </tr>
        </tbody>
      </table>

      <h2>Error Handling</h2>
      <p>
        Python uses exceptions for error handling. RPG generates custom exception hierarchies:
      </p>

      <pre className="code-block">
{`from dataclasses import dataclass


class AppError(Exception):
    """Base exception for application errors."""

    def __init__(self, message: str, code: str = "UNKNOWN_ERROR"):
        super().__init__(message)
        self.message = message
        self.code = code


class ValidationError(AppError):
    """Raised when validation fails."""

    def __init__(self, message: str, field: str):
        super().__init__(message, "VALIDATION_ERROR")
        self.field = field


class NotFoundError(AppError):
    """Raised when a resource is not found."""

    def __init__(self, resource: str):
        super().__init__(f"{resource} not found", "NOT_FOUND")
        self.resource = resource


def divide(a: float, b: float) -> float:
    """Divide a by b.

    Args:
        a: The dividend
        b: The divisor

    Returns:
        The quotient

    Raises:
        ValidationError: If b is zero
    """
    if b == 0:
        raise ValidationError("Cannot divide by zero", "divisor")
    return a / b`}
      </pre>

      <h2>Project Structure</h2>
      <pre className="code-block">
{`project-name/
├── src/
│   └── project_name/
│       ├── __init__.py
│       ├── main.py           # Entry point
│       ├── models/
│       │   ├── __init__.py
│       │   └── user.py
│       ├── services/
│       │   ├── __init__.py
│       │   └── user_service.py
│       ├── utils/
│       │   ├── __init__.py
│       │   └── validation.py
│       └── exceptions.py
├── tests/
│   ├── __init__.py
│   ├── conftest.py           # Pytest fixtures
│   ├── test_user_service.py
│   └── test_validation.py
├── pyproject.toml            # Project config (PEP 518)
├── requirements.txt          # Dependencies (optional)
└── README.md`}
      </pre>

      <h2>Testing Patterns</h2>
      <pre className="code-block">
{`import pytest
from datetime import datetime
from project_name.services.user_service import create_user, UserService
from project_name.exceptions import ValidationError


class TestCreateUser:
    """Tests for the create_user function."""

    def test_creates_user_with_valid_input(self):
        user = create_user("test@example.com", "password123")

        assert user.email == "test@example.com"
        assert user.id is not None
        assert isinstance(user.created_at, datetime)

    def test_raises_validation_error_for_invalid_email(self):
        with pytest.raises(ValidationError) as exc_info:
            create_user("invalid", "password123")

        assert exc_info.value.field == "email"

    def test_raises_validation_error_for_short_password(self):
        with pytest.raises(ValidationError, match="at least 8 characters"):
            create_user("test@example.com", "123")


class TestUserService:
    """Tests for the UserService class."""

    @pytest.fixture
    def mock_repository(self, mocker):
        return mocker.Mock()

    @pytest.fixture
    def service(self, mock_repository):
        return UserService(mock_repository)

    def test_saves_user_to_repository(self, service, mock_repository):
        mock_repository.save.return_value = {"id": "123"}

        result = service.create("test@example.com", "password123")

        mock_repository.save.assert_called_once()
        assert result.id == "123"


# Parametrized tests
@pytest.mark.parametrize("input_text,expected", [
    ("Hello World", "hello-world"),
    ("Hello, World!", "hello-world"),
    ("Hello    World", "hello-world"),
])
def test_slugify(input_text: str, expected: str):
    from project_name.utils import slugify
    assert slugify(input_text) == expected`}
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
{`import re
import uuid
from dataclasses import dataclass, field
from datetime import datetime


@dataclass
class User:
    """Represents a user account."""

    id: str
    email: str
    created_at: datetime = field(default_factory=datetime.now)


class ValidationError(Exception):
    """Raised when validation fails."""

    def __init__(self, message: str, field_name: str):
        super().__init__(message)
        self.field = field_name


EMAIL_REGEX = re.compile(r"^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$")


def create_user(email: str, password: str) -> User:
    """Create a new user account.

    Args:
        email: User's email address
        password: User's password

    Returns:
        The created user

    Raises:
        ValidationError: If email or password is invalid
    """
    # Validate email
    if not EMAIL_REGEX.match(email):
        raise ValidationError("Invalid email format", "email")

    # Validate password
    if len(password) < 8:
        raise ValidationError(
            "Password must be at least 8 characters",
            "password"
        )

    return User(
        id=str(uuid.uuid4()),
        email=email,
    )`}
      </pre>

      <h2>Best Practices</h2>
      <ul>
        <li><strong>Type hints</strong> — Always include type hints for function signatures</li>
        <li><strong>Dataclasses</strong> — Use @dataclass for data containers</li>
        <li><strong>Docstrings</strong> — Use Google-style docstrings for public functions</li>
        <li><strong>PEP 8</strong> — All code follows PEP 8 style guidelines</li>
        <li><strong>Formatting</strong> — Code passes Black formatter</li>
        <li><strong>Linting</strong> — Code passes Ruff or flake8</li>
        <li><strong>Type checking</strong> — Code passes mypy in strict mode</li>
      </ul>

      <h2>Common Dependencies</h2>
      <ul>
        <li><code>pydantic</code> — Data validation with type hints</li>
        <li><code>httpx</code> — Async HTTP client</li>
        <li><code>pytest</code> — Testing framework</li>
        <li><code>pytest-asyncio</code> — Async test support</li>
        <li><code>python-dotenv</code> — Environment variable management</li>
      </ul>
    </div>
  );
}
