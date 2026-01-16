import { useState } from 'react';
import { Zap, Server, Wrench, Database } from 'lucide-react';

type ExampleKey = 'minimal' | 'utility' | 'api' | 'detailed';

interface Example {
  id: ExampleKey;
  title: string;
  description: string;
  icon: typeof Zap;
  color: string;
  spec: string;
}

const examples: Example[] = [
  {
    id: 'minimal',
    title: 'Minimal Style',
    description: 'Brief, conversational description. Let the AI infer details.',
    icon: Zap,
    color: 'yellow',
    spec: `# String Utilities

A collection of common string manipulation functions.

## Functions

### capitalize
Capitalize the first letter of a string, lowercase the rest.

### slugify
Convert a string to a URL-safe slug (lowercase, hyphens for spaces,
remove special characters).

### truncate
Truncate a string to a max length, adding "..." if truncated.
Should not break in the middle of a word.

### wordCount
Count the number of words in a string. Words are separated by whitespace.`,
  },
  {
    id: 'utility',
    title: 'Utility Library',
    description: 'Moderate detail with types and clear function signatures.',
    icon: Wrench,
    color: 'blue',
    spec: `# JSON Schema Validator

Validate JSON data against schemas with detailed error reporting.

## Types

### ValidationError
| Field | Type | Description |
|-------|------|-------------|
| path | string | JSON path to the error (e.g., "user.email") |
| message | string | Human-readable error description |
| expected | string | What was expected |
| actual | string | What was received |

### ValidationResult
| Field | Type | Description |
|-------|------|-------------|
| valid | bool | Whether validation passed |
| errors | []ValidationError | List of validation errors (empty if valid) |

## Functions

### validate
Validates a JSON object against a schema definition.

**Accepts**:
- data: any JSON value to validate
- schema: schema definition object

**Returns**: ValidationResult

**Logic**:
- Check required fields exist
- Validate types match schema
- Check string patterns if specified
- Validate number ranges if specified
- Recurse into nested objects

### formatErrors
Formats validation errors as a human-readable string.

**Accepts**: ValidationResult
**Returns**: string with one error per line, empty string if valid`,
  },
  {
    id: 'api',
    title: 'REST API Service',
    description: 'Full API specification with endpoints and behavior.',
    icon: Server,
    color: 'green',
    spec: `# Task Manager API

A RESTful API for managing tasks with user authentication.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | HTTP server port |
| DATABASE_URL | - | PostgreSQL connection string |
| JWT_SECRET | - | Secret for JWT tokens |

## Types

### Task
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | uuid | Yes | Unique identifier |
| title | string | Yes | Task title (max 200 chars) |
| description | string | No | Detailed description |
| status | enum | Yes | "pending", "in_progress", "completed" |
| createdAt | timestamp | Yes | Creation timestamp |
| userId | uuid | Yes | Owner's user ID |

### CreateTaskRequest
| Field | Type | Required |
|-------|------|----------|
| title | string | Yes |
| description | string | No |

## API Endpoints

### POST /api/tasks
Create a new task for the authenticated user.

**Auth**: Required (JWT Bearer token)
**Body**: CreateTaskRequest
**Response 201**: Created Task object
**Response 401**: Unauthorized

### GET /api/tasks
List all tasks for the authenticated user.

**Auth**: Required
**Query params**:
- status: filter by status (optional)
- limit: max results, default 50
- offset: pagination offset

**Response 200**: Array of Task objects

### PATCH /api/tasks/:id
Update a task's status or details.

**Auth**: Required (must own the task)
**Body**: Partial Task (only fields to update)
**Response 200**: Updated Task
**Response 403**: Not your task
**Response 404**: Task not found

### DELETE /api/tasks/:id
Delete a task.

**Auth**: Required (must own the task)
**Response 204**: No content
**Response 403**: Not your task`,
  },
  {
    id: 'detailed',
    title: 'Detailed Implementation',
    description: 'Complete specification with logic, tests, and error handling.',
    icon: Database,
    color: 'purple',
    spec: `# Rate Limiter

A sliding window rate limiter for API protection.

## Overview

Implement a rate limiter using the sliding window algorithm.
Track requests per client (by IP or API key) and reject requests
that exceed the configured threshold.

## Types

### RateLimitConfig
| Field | Type | Default | Description |
|-------|------|---------|-------------|
| maxRequests | int | 100 | Max requests per window |
| windowSeconds | int | 60 | Window size in seconds |
| keyPrefix | string | "rl:" | Redis key prefix |

### RateLimitResult
| Field | Type | Description |
|-------|------|-------------|
| allowed | bool | Whether the request is allowed |
| remaining | int | Remaining requests in window |
| resetAt | timestamp | When the window resets |
| retryAfter | int | Seconds until retry (if blocked) |

## Functions

### NewRateLimiter
Create a new rate limiter instance.

**Accepts**: RateLimitConfig, Redis client
**Returns**: *RateLimiter

### Check
Check if a request should be allowed.

**Accepts**:
- key: string (client identifier like IP or API key)

**Returns**: RateLimitResult, error

**Logic**:
1. Get current timestamp in seconds
2. Calculate window start (now - windowSeconds)
3. Remove entries older than window start (ZREMRANGEBYSCORE)
4. Count current entries (ZCARD)
5. If count >= maxRequests:
   - Get oldest entry timestamp
   - Calculate retryAfter = oldest + windowSeconds - now
   - Return blocked result
6. Add current timestamp to sorted set (ZADD)
7. Set TTL on key to windowSeconds (EXPIRE)
8. Return allowed result with remaining count

## Tests

### test_allows_under_limit
**Given**: Limiter with max 5 requests per 60 seconds
**When**: Make 3 requests with same key
**Expect**: All return allowed=true, remaining decreases

### test_blocks_over_limit
**Given**: Limiter with max 5 requests per 60 seconds
**When**: Make 6 requests with same key
**Expect**: First 5 allowed, 6th blocked with retryAfter > 0

### test_window_slides
**Given**: Limiter with max 5 requests per 2 seconds
**When**: Make 5 requests, wait 3 seconds, make 1 more
**Expect**: All 6 requests allowed (window slid)

### test_different_keys_independent
**Given**: Limiter with max 2 requests
**When**: Make 2 requests with key "a", 2 with key "b"
**Expect**: All 4 allowed (separate counters)

## Dependencies

- Redis client library for sorted set operations
- Time utilities for timestamp handling`,
  },
];

const colorClasses: Record<string, { bg: string; border: string; icon: string }> = {
  yellow: {
    bg: 'bg-yellow-50 dark:bg-yellow-900/20',
    border: 'border-yellow-200 dark:border-yellow-800 hover:border-yellow-300 dark:hover:border-yellow-700',
    icon: 'bg-yellow-500',
  },
  blue: {
    bg: 'bg-blue-50 dark:bg-blue-900/20',
    border: 'border-blue-200 dark:border-blue-800 hover:border-blue-300 dark:hover:border-blue-700',
    icon: 'bg-blue-500',
  },
  green: {
    bg: 'bg-green-50 dark:bg-green-900/20',
    border: 'border-green-200 dark:border-green-800 hover:border-green-300 dark:hover:border-green-700',
    icon: 'bg-green-500',
  },
  purple: {
    bg: 'bg-purple-50 dark:bg-purple-900/20',
    border: 'border-purple-200 dark:border-purple-800 hover:border-purple-300 dark:hover:border-purple-700',
    icon: 'bg-purple-500',
  },
};

export default function Examples() {
  const [selected, setSelected] = useState<ExampleKey>('minimal');
  const selectedExample = examples.find((e) => e.id === selected)!;

  return (
    <div className="prose-docs">
      <h1>Spec Examples</h1>
      <p className="lead">
        Real spec examples showing different styles—from minimal descriptions
        to detailed implementations. All styles work with RPG.
      </p>

      <h2>Choose a Style</h2>
      <p>
        Click an example to see the full spec. Notice how each style conveys
        different levels of detail while remaining readable markdown.
      </p>

      {/* Example selector */}
      <div className="not-prose grid grid-cols-1 sm:grid-cols-2 gap-3 my-6">
        {examples.map((example) => {
          const colors = colorClasses[example.color];
          const Icon = example.icon;
          const isSelected = selected === example.id;

          return (
            <button
              key={example.id}
              onClick={() => setSelected(example.id)}
              className={`p-4 rounded-lg border-2 text-left transition-all ${colors.bg} ${
                isSelected
                  ? `${colors.border} ring-2 ring-offset-2 ring-${example.color}-500`
                  : `border-transparent hover:${colors.border}`
              }`}
            >
              <div className="flex items-center gap-3">
                <div className={`p-2 ${colors.icon} rounded-lg`}>
                  <Icon className="h-5 w-5 text-white" />
                </div>
                <div>
                  <h4 className="font-semibold text-gray-900 dark:text-white text-sm">
                    {example.title}
                  </h4>
                  <p className="text-xs text-gray-600 dark:text-gray-400 mt-0.5">
                    {example.description}
                  </p>
                </div>
              </div>
            </button>
          );
        })}
      </div>

      {/* Selected example display */}
      <div className="not-prose my-8">
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            {selectedExample.title}
          </h3>
          <span className="text-xs text-gray-500 dark:text-gray-400 font-mono">
            {selectedExample.id}.spec.md
          </span>
        </div>
        <pre className="code-block text-sm overflow-x-auto max-h-[600px] overflow-y-auto">
          {selectedExample.spec}
        </pre>
      </div>

      <h2>When to Use Each Style</h2>

      <div className="not-prose space-y-4">
        <div className="p-4 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg border border-yellow-200 dark:border-yellow-800">
          <h4 className="font-semibold text-yellow-900 dark:text-yellow-100 mb-2">
            Minimal Style
          </h4>
          <p className="text-sm text-yellow-800 dark:text-yellow-200">
            Best for simple utilities, prototypes, or when you trust the AI to make
            reasonable implementation choices. Write less, let AI fill in standard patterns.
          </p>
        </div>

        <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <h4 className="font-semibold text-blue-900 dark:text-blue-100 mb-2">
            Utility Library Style
          </h4>
          <p className="text-sm text-blue-800 dark:text-blue-200">
            Good balance for libraries and shared code. Define types clearly so the AI
            generates consistent interfaces. Describe function behavior without over-specifying.
          </p>
        </div>

        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <h4 className="font-semibold text-green-900 dark:text-green-100 mb-2">
            API Service Style
          </h4>
          <p className="text-sm text-green-800 dark:text-green-200">
            Ideal for REST APIs where endpoint contracts matter. Specify routes, auth requirements,
            request/response shapes. The AI handles the implementation details.
          </p>
        </div>

        <div className="p-4 bg-purple-50 dark:bg-purple-900/20 rounded-lg border border-purple-200 dark:border-purple-800">
          <h4 className="font-semibold text-purple-900 dark:text-purple-100 mb-2">
            Detailed Implementation Style
          </h4>
          <p className="text-sm text-purple-800 dark:text-purple-200">
            Use when behavior must be precise—algorithms, business logic, security features.
            Include step-by-step logic and test cases to ensure exact implementation.
          </p>
        </div>
      </div>

      <h2>Key Observations</h2>
      <ul>
        <li>
          <strong>All examples are valid specs</strong> — RPG doesn't require a specific format.
          The AI adapts to your writing style.
        </li>
        <li>
          <strong>Tables are optional</strong> — You can use prose, bullets, or tables
          for type definitions. Use what's clearest for your content.
        </li>
        <li>
          <strong>Logic blocks guide implementation</strong> — When you need specific behavior,
          describe it step-by-step. Otherwise, let the AI choose idiomatic patterns.
        </li>
        <li>
          <strong>Tests define expectations</strong> — Including test cases helps the AI
          understand edge cases and expected behavior.
        </li>
      </ul>

      <h2>Mixing Styles</h2>
      <p>
        You can mix styles within a single spec. For example, use minimal descriptions
        for simple CRUD functions but detailed logic for complex algorithms:
      </p>

      <pre className="code-block">
{`# User Service

## Functions

### getUser
Get a user by ID. Return null if not found.

### listUsers
List all users with optional pagination.

### calculateEngagementScore
Calculate a user's engagement score based on activity.

**Logic**:
1. Get login count from last 30 days
2. Get content interactions (views, likes, comments)
3. Weight recent activity higher (exponential decay)
4. Normalize to 0-100 scale
5. Cache result for 1 hour`}
      </pre>

      <p>
        The simple functions get minimal specs while the complex algorithm gets detailed logic.
        RPG handles both appropriately.
      </p>
    </div>
  );
}
