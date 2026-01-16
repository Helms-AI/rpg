export default function TypeScript() {
  return (
    <div className="prose-docs">
      <h1>TypeScript</h1>
      <p className="lead">
        RPG generates strict TypeScript code with comprehensive type safety, modern ES features,
        and industry-standard patterns.
      </p>

      <div className="not-prose my-6 p-4 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-blue-800 dark:text-blue-300 font-semibold">Philosophy</span>
        </div>
        <p className="text-blue-700 dark:text-blue-400 text-sm">
          TypeScript adds static typing to JavaScript while maintaining full interoperability.
          RPG generates code with strict mode enabled, catching errors at compile time rather than runtime.
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
            <td>Files</td>
            <td>kebab-case or camelCase</td>
            <td><code>user-service.ts</code>, <code>userService.ts</code></td>
          </tr>
          <tr>
            <td>Functions</td>
            <td>camelCase</td>
            <td><code>createUser</code>, <code>validateToken</code></td>
          </tr>
          <tr>
            <td>Classes</td>
            <td>PascalCase</td>
            <td><code>UserService</code>, <code>AuthController</code></td>
          </tr>
          <tr>
            <td>Interfaces</td>
            <td>PascalCase (no I prefix)</td>
            <td><code>User</code>, <code>CreateUserRequest</code></td>
          </tr>
          <tr>
            <td>Type Aliases</td>
            <td>PascalCase</td>
            <td><code>UserId</code>, <code>ValidationResult</code></td>
          </tr>
          <tr>
            <td>Enums</td>
            <td>PascalCase (type and members)</td>
            <td><code>Status.Active</code>, <code>Role.Admin</code></td>
          </tr>
          <tr>
            <td>Constants</td>
            <td>SCREAMING_SNAKE_CASE or camelCase</td>
            <td><code>MAX_RETRIES</code>, <code>defaultConfig</code></td>
          </tr>
        </tbody>
      </table>

      <h2>Type Mappings</h2>
      <table className="w-full">
        <thead>
          <tr>
            <th>Spec Type</th>
            <th>TypeScript Type</th>
            <th>Notes</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>string</code></td>
            <td><code>string</code></td>
            <td>Primitive string type</td>
          </tr>
          <tr>
            <td><code>number</code></td>
            <td><code>number</code></td>
            <td>All numbers are floating point</td>
          </tr>
          <tr>
            <td><code>integer</code></td>
            <td><code>number</code></td>
            <td>No separate integer type</td>
          </tr>
          <tr>
            <td><code>boolean</code></td>
            <td><code>boolean</code></td>
            <td>true/false</td>
          </tr>
          <tr>
            <td><code>datetime</code></td>
            <td><code>Date</code></td>
            <td>Built-in Date object</td>
          </tr>
          <tr>
            <td><code>array</code></td>
            <td><code>T[]</code> or <code>Array&lt;T&gt;</code></td>
            <td>Generic array type</td>
          </tr>
          <tr>
            <td><code>map</code></td>
            <td><code>Record&lt;K, V&gt;</code> or <code>Map&lt;K, V&gt;</code></td>
            <td>Object or Map class</td>
          </tr>
          <tr>
            <td><code>optional</code></td>
            <td><code>T | undefined</code></td>
            <td>Union with undefined</td>
          </tr>
          <tr>
            <td><code>nullable</code></td>
            <td><code>T | null</code></td>
            <td>Union with null</td>
          </tr>
        </tbody>
      </table>

      <h2>Error Handling</h2>
      <p>
        TypeScript uses exceptions for error handling. RPG generates custom error classes
        with proper inheritance:
      </p>

      <pre className="code-block">
{`// Custom error classes
export class AppError extends Error {
  constructor(
    message: string,
    public readonly code: string,
    public readonly statusCode: number = 500
  ) {
    super(message);
    this.name = 'AppError';
    Error.captureStackTrace(this, this.constructor);
  }
}

export class ValidationError extends AppError {
  constructor(
    message: string,
    public readonly field: string
  ) {
    super(message, 'VALIDATION_ERROR', 400);
    this.name = 'ValidationError';
  }
}

export class NotFoundError extends AppError {
  constructor(resource: string) {
    super(\`\${resource} not found\`, 'NOT_FOUND', 404);
    this.name = 'NotFoundError';
  }
}

// Function that throws
export function divide(a: number, b: number): number {
  if (b === 0) {
    throw new ValidationError('Cannot divide by zero', 'divisor');
  }
  return a / b;
}

// Result pattern for functional approach
type Result<T, E> = { ok: true; value: T } | { ok: false; error: E };

export function safeDivide(a: number, b: number): Result<number, string> {
  if (b === 0) {
    return { ok: false, error: 'Division by zero' };
  }
  return { ok: true, value: a / b };
}`}
      </pre>

      <h2>Project Structure</h2>
      <pre className="code-block">
{`project-name/
├── src/
│   ├── index.ts            # Main entry point
│   ├── types/
│   │   └── index.ts        # Type definitions
│   ├── models/
│   │   └── user.ts
│   ├── services/
│   │   └── user-service.ts
│   ├── utils/
│   │   └── validation.ts
│   └── errors/
│       └── index.ts
├── tests/
│   ├── unit/
│   │   └── user-service.test.ts
│   └── integration/
│       └── api.test.ts
├── package.json
├── tsconfig.json
├── jest.config.js          # Or vitest.config.ts
└── README.md`}
      </pre>

      <h2>Testing Patterns</h2>
      <pre className="code-block">
{`import { describe, it, expect, beforeEach, vi } from 'vitest';
import { createUser, UserService } from './user-service';

describe('createUser', () => {
  it('should create a user with valid input', () => {
    const user = createUser('test@example.com', 'password123');

    expect(user).toMatchObject({
      email: 'test@example.com',
    });
    expect(user.id).toBeDefined();
    expect(user.createdAt).toBeInstanceOf(Date);
  });

  it('should throw ValidationError for invalid email', () => {
    expect(() => createUser('invalid', 'password123'))
      .toThrow(ValidationError);
  });

  it('should throw ValidationError for short password', () => {
    expect(() => createUser('test@example.com', '123'))
      .toThrow('Password must be at least 8 characters');
  });
});

describe('UserService', () => {
  let service: UserService;
  let mockRepository: MockRepository;

  beforeEach(() => {
    mockRepository = {
      save: vi.fn(),
      findById: vi.fn(),
    };
    service = new UserService(mockRepository);
  });

  it('should save user to repository', async () => {
    mockRepository.save.mockResolvedValue({ id: '123' });

    const result = await service.create('test@example.com', 'password123');

    expect(mockRepository.save).toHaveBeenCalled();
    expect(result.id).toBe('123');
  });
});`}
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
{`import { v4 as uuidv4 } from 'uuid';

export interface User {
  id: string;
  email: string;
  createdAt: Date;
}

export class ValidationError extends Error {
  constructor(
    message: string,
    public readonly field: string
  ) {
    super(message);
    this.name = 'ValidationError';
  }
}

const EMAIL_REGEX = /^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$/;

export function createUser(email: string, password: string): User {
  // Validate email
  if (!EMAIL_REGEX.test(email)) {
    throw new ValidationError('Invalid email format', 'email');
  }

  // Validate password
  if (password.length < 8) {
    throw new ValidationError(
      'Password must be at least 8 characters',
      'password'
    );
  }

  // Create user (password would be hashed in real implementation)
  return {
    id: uuidv4(),
    email,
    createdAt: new Date(),
  };
}`}
      </pre>

      <h2>Best Practices</h2>
      <ul>
        <li><strong>Strict mode</strong> — Always enable strict: true in tsconfig.json</li>
        <li><strong>Explicit types</strong> — Prefer explicit return types for public APIs</li>
        <li><strong>Interfaces over types</strong> — Use interfaces for object shapes, types for unions</li>
        <li><strong>Avoid any</strong> — Use unknown instead of any when type is truly unknown</li>
        <li><strong>Null checks</strong> — Enable strictNullChecks and handle null explicitly</li>
        <li><strong>Formatting</strong> — Code passes Prettier with default config</li>
        <li><strong>Linting</strong> — Code passes ESLint with TypeScript rules</li>
      </ul>

      <h2>Common Dependencies</h2>
      <ul>
        <li><code>zod</code> — Runtime type validation</li>
        <li><code>uuid</code> — UUID generation</li>
        <li><code>date-fns</code> — Date manipulation</li>
        <li><code>axios</code> — HTTP client</li>
        <li><code>vitest</code> or <code>jest</code> — Testing framework</li>
      </ul>
    </div>
  );
}
