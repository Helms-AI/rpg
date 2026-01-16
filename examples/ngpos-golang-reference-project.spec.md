# NGPOS Go Reference Service

A reference implementation demonstrating standard operational behavior for Go services. This project provides a clear, runnable example of how services should behave operationally, aligned with equivalent Java services. It follows an API-first design approach with REST APIs defined using OpenAPI and async messaging contracts defined using AsyncAPI.

## Target Languages

- go
- java
- typescript
- python
- rust
- csharp

## Meta

- **Version**: 1.0.0
- **License**: Proprietary
- **Author**: Kroger Technology

## Configuration

| Variable | Type | Default | Required | Description |
|----------|------|---------|----------|-------------|
| SERVICE_NAME | string | ngpos-go-reference | No | Name of the service for labeling metrics and logs |
| ENV | string | local | No | Environment identifier (local, dev, staging, prod) |
| PORT | int | 8080 | No | HTTP server listen port |

## Types

### HealthResponse

Health check response indicating service status.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| status | string | No | Health status indicator (e.g., "UP", "DOWN") |

### MetricsResponse

Prometheus-compatible metrics response in text format.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| body | string | Yes | Plain text metrics in Prometheus exposition format |

### ReceiptEvent

Async event payload for receipt-related events.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| receiptId | string | Yes | Unique identifier for the receipt |
| eventType | string | Yes | Type of event (e.g., "created", "updated", "deleted") |

## Functions

### main

Application entry point that initializes configuration and starts the HTTP server.

**Accepts**: None

**Returns**: None (exits on error)

**Logic**:
```
1. Read SERVICE_NAME from environment, default to "ngpos-go-reference" if empty
2. Read ENV from environment, default to "local" if empty
3. Log startup message with service name
4. Log environment information
5. Register HTTP handler for /metrics endpoint
6. Log that metrics are available at /metrics
7. Start HTTP server on port 8080
8. If server fails to start, log fatal error and exit
```

**Errors**:
- Server startup failure logs fatal error and terminates process

---

### handleMetrics

HTTP handler that returns Prometheus-compatible metrics.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| w | http.ResponseWriter | Yes | HTTP response writer |
| r | *http.Request | Yes | HTTP request |

**Returns**: None (writes to response)

**Logic**:
```
1. Atomically increment the global request counter
2. Read current memory statistics from runtime
3. Set response status to 200 OK
4. Write metrics in Prometheus exposition format:
   - service_up{service="<name>",env="<env>"} 1
   - http_requests_total{service="<name>",env="<env>"} <count>
   - go_mem_alloc_bytes{service="<name>",env="<env>"} <bytes>
```

---

### fmtUint

Formats an unsigned 64-bit integer as a string.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| v | uint64 | Yes | Value to format |

**Returns**: string - Decimal string representation of the value

**Logic**:
```
1. Use fmt.Sprintf to convert uint64 to decimal string
2. Return the formatted string
```

---

### NewGetHealth200Response

Constructor that creates a new HealthResponse instance with default values.

**Accepts**: None

**Returns**: *HealthResponse - Pointer to new instance

**Logic**:
```
1. Create new empty HealthResponse struct
2. Return pointer to the struct
```

---

### GetStatus

Returns the status field value if set, otherwise returns empty string.

**Accepts**: None (method receiver: *HealthResponse)

**Returns**: string - Status value or empty string

**Logic**:
```
1. If receiver is nil or Status field is nil, return empty string
2. Otherwise, dereference and return the Status value
```

---

### SetStatus

Sets the status field on a HealthResponse.

**Accepts**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| v | string | Yes | Status value to set |

**Returns**: None

**Logic**:
```
1. Take address of input value
2. Assign pointer to the Status field
```

---

### HasStatus

Checks if the status field has been set.

**Accepts**: None (method receiver: *HealthResponse)

**Returns**: bool - True if status is set, false otherwise

**Logic**:
```
1. If receiver is not nil AND Status is not nil, return true
2. Otherwise return false
```

## API Endpoints

### GET /health

Health check endpoint to verify service availability.

**Summary**: Health Check

**Operation ID**: getHealth

**Request**: None

**Response 200**:
- Content-Type: application/json
- Body: HealthResponse with status field (e.g., {"status": "UP"})

---

### GET /metrics

Prometheus-compatible metrics endpoint for observability.

**Summary**: Metrics Endpoint

**Operation ID**: getMetrics

**Request**: None

**Response 200**:
- Content-Type: text/plain
- Body: Prometheus exposition format metrics

**Example Response**:
```
service_up{service="ngpos-go-reference",env="local"} 1
http_requests_total{service="ngpos-go-reference",env="local"} 42
go_mem_alloc_bytes{service="ngpos-go-reference",env="local"} 1234567
```

## Async Events

### Channel: receipt.events

Publishes receipt-related events for async processing.

**Direction**: Publish

**Payload**: ReceiptEvent

**Example**:
```json
{
  "receiptId": "receipt-123",
  "eventType": "created"
}
```

## Tests

### test_health_endpoint_returns_up

Verifies the health endpoint returns UP status.

**Given**:
- Server is running
- GET request to /health

**Expect**:
- Status code: 200
- Response body contains: {"status": "UP"}

---

### test_metrics_endpoint_returns_valid_format

Verifies the metrics endpoint returns Prometheus format.

**Given**:
- Server is running with SERVICE_NAME="test-service" and ENV="test"
- GET request to /metrics

**Expect**:
- Status code: 200
- Content-Type: text/plain
- Response contains: service_up{service="test-service",env="test"} 1
- Response contains: http_requests_total
- Response contains: go_mem_alloc_bytes

---

### test_request_counter_increments

Verifies request counter increments on each metrics call.

**Given**:
- Server is running
- Make 3 sequential GET requests to /metrics

**Expect**:
- Third response contains http_requests_total value >= 3

---

### test_default_configuration

Verifies default configuration values when env vars not set.

**Given**:
- No SERVICE_NAME environment variable
- No ENV environment variable

**Expect**:
- Service uses "ngpos-go-reference" as service name
- Service uses "local" as environment

---

### test_custom_configuration

Verifies custom configuration via environment variables.

**Given**:
- SERVICE_NAME="custom-service"
- ENV="production"

**Expect**:
- Metrics contain service="custom-service"
- Metrics contain env="production"

## Dependencies

### Go Standard Library

| Package | Purpose |
|---------|---------|
| fmt | String formatting |
| log | Logging |
| net/http | HTTP server and handlers |
| os | Environment variable access |
| runtime | Memory statistics |
| sync/atomic | Thread-safe counter operations |
| encoding/json | JSON marshaling for API responses |
