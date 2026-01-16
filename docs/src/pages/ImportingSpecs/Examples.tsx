import { useState } from 'react';
import { Zap, Server, Package, Building2, ArrowRight } from 'lucide-react';

type ExampleKey = 'utility' | 'library' | 'api' | 'legacy';

interface Example {
  id: ExampleKey;
  title: string;
  description: string;
  icon: typeof Zap;
  color: string;
  sourceCode: string;
  sourceLanguage: string;
  generatedSpec: string;
}

const examples: Example[] = [
  {
    id: 'utility',
    title: 'Single Utility File',
    description: 'Import a small utility and generate in other languages.',
    icon: Zap,
    color: 'yellow',
    sourceLanguage: 'Go',
    sourceCode: `// stringutil.go
package stringutil

import (
    "strings"
    "unicode"
)

// Slugify converts a string to URL-safe slug
func Slugify(s string) string {
    s = strings.ToLower(s)
    s = strings.Map(func(r rune) rune {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            return r
        }
        if r == ' ' || r == '-' {
            return '-'
        }
        return -1
    }, s)
    return strings.Trim(s, "-")
}

// WordCount returns number of words in string
func WordCount(s string) int {
    return len(strings.Fields(s))
}

// Truncate shortens string to max length with ellipsis
func Truncate(s string, max int) string {
    if len(s) <= max {
        return s
    }
    return s[:max-3] + "..."
}`,
    generatedSpec: `# String Utilities

A collection of common string manipulation functions.

## Functions

### Slugify
Converts a string to a URL-safe slug.

**Accepts**: string to convert
**Returns**: slug string (lowercase, hyphens, no special chars)

**Logic**:
1. Convert to lowercase
2. Replace spaces with hyphens
3. Remove non-alphanumeric characters (except hyphens)
4. Trim leading/trailing hyphens

### WordCount
Counts the number of words in a string.

**Accepts**: input string
**Returns**: int (word count)

Words are separated by whitespace (spaces, tabs, newlines).

### Truncate
Shortens a string to a maximum length with ellipsis.

**Accepts**:
- s: string to truncate
- max: maximum length (including ellipsis)

**Returns**: truncated string

**Logic**:
- If string length <= max, return unchanged
- Otherwise return first (max-3) characters + "..."`,
  },
  {
    id: 'library',
    title: 'Library Migration',
    description: 'Import a multi-file library with tests.',
    icon: Package,
    color: 'blue',
    sourceLanguage: 'Java',
    sourceCode: `// Result.java
public sealed interface Result<T, E> {
    record Ok<T, E>(T value) implements Result<T, E> {}
    record Err<T, E>(E error) implements Result<T, E> {}

    static <T, E> Result<T, E> ok(T value) {
        return new Ok<>(value);
    }

    static <T, E> Result<T, E> err(E error) {
        return new Err<>(error);
    }

    default <U> Result<U, E> map(Function<T, U> f) {
        return switch (this) {
            case Ok<T, E> ok -> Result.ok(f.apply(ok.value()));
            case Err<T, E> err -> Result.err(err.error());
        };
    }

    default T unwrapOr(T defaultValue) {
        return switch (this) {
            case Ok<T, E> ok -> ok.value();
            case Err<T, E> err -> defaultValue;
        };
    }
}

// ResultTest.java
class ResultTest {
    @Test void okValueCanBeUnwrapped() {
        Result<Integer, String> r = Result.ok(42);
        assertEquals(42, r.unwrapOr(0));
    }

    @Test void errReturnsDefault() {
        Result<Integer, String> r = Result.err("failed");
        assertEquals(0, r.unwrapOr(0));
    }
}`,
    generatedSpec: `# Result Type

A functional Result type for explicit error handling without exceptions.

## Types

### Result<T, E>
A discriminated union representing either success (Ok) or failure (Err).

| Variant | Fields | Description |
|---------|--------|-------------|
| Ok | value: T | Successful result containing value |
| Err | error: E | Failed result containing error |

## Functions

### ok
Creates a successful Result containing a value.

**Accepts**: value of type T
**Returns**: Result<T, E> in Ok state

### err
Creates a failed Result containing an error.

**Accepts**: error of type E
**Returns**: Result<T, E> in Err state

### map
Transforms the success value using a function.

**Accepts**: function (T → U)
**Returns**: Result<U, E>

**Logic**:
- If Ok: apply function to value, return new Ok
- If Err: return Err unchanged (error passes through)

### unwrapOr
Extracts the value or returns a default.

**Accepts**: default value of type T
**Returns**: T (the value if Ok, default if Err)

## Tests

### test_ok_value_can_be_unwrapped
**Given**: Result.ok(42)
**Expect**: unwrapOr(0) returns 42

### test_err_returns_default
**Given**: Result.err("failed")
**Expect**: unwrapOr(0) returns 0`,
  },
  {
    id: 'api',
    title: 'REST API Service',
    description: 'Import an API service with endpoints and models.',
    icon: Server,
    color: 'green',
    sourceLanguage: 'TypeScript',
    sourceCode: `// types.ts
interface Todo {
  id: string;
  title: string;
  completed: boolean;
  createdAt: Date;
}

// routes.ts
app.get('/api/todos', async (req, res) => {
  const userId = req.user.id;
  const todos = await db.todos.findMany({ userId });
  res.json(todos);
});

app.post('/api/todos', async (req, res) => {
  const { title } = req.body;
  if (!title?.trim()) {
    return res.status(400).json({ error: 'Title required' });
  }
  const todo = await db.todos.create({
    id: uuid(),
    title: title.trim(),
    completed: false,
    userId: req.user.id,
    createdAt: new Date()
  });
  res.status(201).json(todo);
});

app.patch('/api/todos/:id', async (req, res) => {
  const todo = await db.todos.findUnique({ id: req.params.id });
  if (!todo) return res.status(404).json({ error: 'Not found' });
  if (todo.userId !== req.user.id) {
    return res.status(403).json({ error: 'Forbidden' });
  }
  const updated = await db.todos.update({
    id: req.params.id,
    ...req.body
  });
  res.json(updated);
});`,
    generatedSpec: `# Todo API

A RESTful API for managing user todo items.

## Types

### Todo
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| id | uuid | Yes | Unique identifier |
| title | string | Yes | Todo title |
| completed | bool | Yes | Completion status |
| createdAt | timestamp | Yes | Creation time |

## API Endpoints

### GET /api/todos
List all todos for the authenticated user.

**Auth**: Required
**Response 200**: Array of Todo objects

### POST /api/todos
Create a new todo.

**Auth**: Required
**Body**: { title: string }
**Response 201**: Created Todo object
**Response 400**: Title is required or empty

**Logic**:
1. Validate title is non-empty
2. Trim whitespace from title
3. Generate UUID for id
4. Set completed to false
5. Set createdAt to current time
6. Associate with authenticated user

### PATCH /api/todos/:id
Update an existing todo.

**Auth**: Required (must own todo)
**Params**: id - Todo ID
**Body**: Partial Todo (fields to update)
**Response 200**: Updated Todo
**Response 403**: Not your todo
**Response 404**: Todo not found

**Logic**:
1. Find todo by ID
2. Verify ownership matches authenticated user
3. Apply updates
4. Return updated todo`,
  },
  {
    id: 'legacy',
    title: 'Legacy Modernization',
    description: 'Import complex legacy code and extract intent.',
    icon: Building2,
    color: 'purple',
    sourceLanguage: 'PHP',
    sourceCode: `<?php
// OrderProcessor.php - Legacy e-commerce code
class OrderProcessor {
    private $db;
    private $mailer;

    function processOrder($cart, $user) {
        // Calculate totals
        $subtotal = 0;
        foreach ($cart->items as $item) {
            $subtotal += $item->price * $item->qty;
        }

        // Apply discount
        $discount = 0;
        if ($user->isPremium && $subtotal > 100) {
            $discount = $subtotal * 0.1; // 10% premium discount
        }

        // Calculate tax (8.5%)
        $tax = ($subtotal - $discount) * 0.085;
        $total = $subtotal - $discount + $tax;

        // Create order
        $order = new Order();
        $order->user_id = $user->id;
        $order->subtotal = $subtotal;
        $order->discount = $discount;
        $order->tax = $tax;
        $order->total = $total;
        $order->status = 'pending';

        $this->db->save($order);

        // Send confirmation
        $this->mailer->send($user->email,
            'Order Confirmed',
            "Your order total: $" . number_format($total, 2)
        );

        return $order;
    }
}`,
    generatedSpec: `# Order Processing

Processes shopping cart orders with discounts, tax, and notifications.

## Types

### Order
| Field | Type | Description |
|-------|------|-------------|
| id | uuid | Order identifier |
| userId | uuid | Customer ID |
| subtotal | decimal | Sum of item prices × quantities |
| discount | decimal | Applied discount amount |
| tax | decimal | Calculated tax |
| total | decimal | Final amount (subtotal - discount + tax) |
| status | enum | "pending", "confirmed", "shipped", "delivered" |

### CartItem
| Field | Type | Description |
|-------|------|-------------|
| productId | uuid | Product reference |
| price | decimal | Unit price |
| quantity | int | Number of items |

## Functions

### processOrder
Processes a cart into a finalized order.

**Accepts**:
- cart: Cart with items
- user: User placing order

**Returns**: Created Order

**Logic**:
1. Calculate subtotal (sum of price × quantity for all items)
2. Apply discount:
   - Premium users with subtotal > $100 get 10% off
   - Otherwise no discount
3. Calculate tax: 8.5% of (subtotal - discount)
4. Calculate total: subtotal - discount + tax
5. Create order with status "pending"
6. Save to database
7. Send confirmation email with total
8. Return order

## Configuration

| Variable | Value | Description |
|----------|-------|-------------|
| TAX_RATE | 0.085 | Sales tax rate (8.5%) |
| PREMIUM_DISCOUNT | 0.10 | Premium member discount (10%) |
| PREMIUM_THRESHOLD | 100 | Minimum subtotal for premium discount |

## Dependencies

- Database for order persistence
- Email service for confirmations`,
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
  const [selected, setSelected] = useState<ExampleKey>('utility');
  const [showSpec, setShowSpec] = useState(false);
  const selectedExample = examples.find((e) => e.id === selected)!;

  return (
    <div className="prose-docs">
      <h1>Import Examples</h1>
      <p className="lead">
        See how existing code transforms into natural language specs.
        From simple utilities to complex legacy systems.
      </p>

      <h2>Choose a Scenario</h2>
      <p>
        Click an example to see the source code and generated spec side by side.
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
              onClick={() => {
                setSelected(example.id);
                setShowSpec(false);
              }}
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

      {/* Source/Spec toggle and display */}
      <div className="not-prose my-8">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-4">
            <button
              onClick={() => setShowSpec(false)}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                !showSpec
                  ? 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
              }`}
            >
              Source ({selectedExample.sourceLanguage})
            </button>
            <ArrowRight className="h-4 w-4 text-gray-400" />
            <button
              onClick={() => setShowSpec(true)}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                showSpec
                  ? 'bg-purple-200 dark:bg-purple-700 text-purple-900 dark:text-white'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
              }`}
            >
              Generated Spec
            </button>
          </div>
          <span className="text-xs text-gray-500 dark:text-gray-400">
            {showSpec ? 'spec.md' : `original.${selectedExample.sourceLanguage.toLowerCase()}`}
          </span>
        </div>
        <pre className="code-block text-sm overflow-x-auto max-h-[500px] overflow-y-auto">
          {showSpec ? selectedExample.generatedSpec : selectedExample.sourceCode}
        </pre>
      </div>

      <h2>What to Notice</h2>

      <div className="not-prose space-y-4">
        <div className="p-4 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg border border-yellow-200 dark:border-yellow-800">
          <h4 className="font-semibold text-yellow-900 dark:text-yellow-100 mb-2">
            Simple Utility
          </h4>
          <p className="text-sm text-yellow-800 dark:text-yellow-200">
            The AI extracts function purposes, parameter meanings, and implementation logic.
            Even simple code benefits from natural language documentation.
          </p>
        </div>

        <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <h4 className="font-semibold text-blue-900 dark:text-blue-100 mb-2">
            Library Migration
          </h4>
          <p className="text-sm text-blue-800 dark:text-blue-200">
            Tests reveal expected behavior—they're gold for spec generation. The AI uses
            test cases to understand edge cases and create spec-level tests.
          </p>
        </div>

        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <h4 className="font-semibold text-green-900 dark:text-green-100 mb-2">
            API Service
          </h4>
          <p className="text-sm text-green-800 dark:text-green-200">
            Route handlers become API endpoint specs. Request validation, response codes,
            and auth requirements are extracted from the implementation.
          </p>
        </div>

        <div className="p-4 bg-purple-50 dark:bg-purple-900/20 rounded-lg border border-purple-200 dark:border-purple-800">
          <h4 className="font-semibold text-purple-900 dark:text-purple-100 mb-2">
            Legacy Modernization
          </h4>
          <p className="text-sm text-purple-800 dark:text-purple-200">
            Complex business logic gets documented. Magic numbers become named configuration.
            Implicit behavior becomes explicit spec that can be reviewed and ported.
          </p>
        </div>
      </div>

      <h2>The Transformation</h2>
      <p>
        Notice how the generated specs:
      </p>
      <ul>
        <li>
          <strong>Abstract implementation details</strong> — Language-specific syntax becomes
          portable descriptions
        </li>
        <li>
          <strong>Capture business rules</strong> — Magic numbers like 0.085 become named
          constants with explanations
        </li>
        <li>
          <strong>Document behavior, not mechanics</strong> — "Calculate tax" not
          "multiply by 0.085"
        </li>
        <li>
          <strong>Preserve tests as specs</strong> — Unit tests become Given/Expect test cases
        </li>
        <li>
          <strong>Extract types cleanly</strong> — Classes and interfaces become type definitions
        </li>
      </ul>

      <h2>After Import</h2>
      <p>
        Once you have a generated spec, you can:
      </p>
      <ol>
        <li>
          <strong>Review for accuracy</strong> — Does the spec match what the code actually does?
        </li>
        <li>
          <strong>Enhance descriptions</strong> — Add the "why" that wasn't in the code
        </li>
        <li>
          <strong>Generate in new languages</strong> — Use <code>get_generation_context</code> to
          produce TypeScript, Rust, Go, etc.
        </li>
        <li>
          <strong>Run parity checks</strong> — Use <code>ensure_parity</code> to verify
          implementations match
        </li>
      </ol>
    </div>
  );
}
