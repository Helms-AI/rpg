# ngpos-golang-reference-project

A reference implementation demonstrating standard operational behavior for Go services. Provides health checks and Prometheus-style metrics exposure with environment-based configuration.

## Meta

- version: 1.0.0
- license: Internal
- description: Minimal reference service for operational alignment with Java services

## Target Languages

- go
- java
- typescript
- csharp

## Configuration

### Environment Variables

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| SERVICE_NAME | Text | "ngpos-go-reference" | Name of the service for metrics labeling |
| ENV | Text | "local" | Environment identifier (local, dev, prod) |

### Server Configuration

- Port: 8080 (hardcoded)

## Types

### GetHealth200Response

Health check response payload.

contains:
- status: Optional Text (JSON: "status")

### NullableGetHealth200Response

Nullable wrapper for GetHealth200Response with explicit set/unset tracking.

contains:
- value: Optional GetHealth200Response
- isSet: Boolean

### MemStats

Runtime memory statistics (from Go runtime).

contains:
- Alloc: Number (currently allocated bytes)

### MetricsResponse

Prometheus-style metrics output.

contains:
- service_up: Gauge with service and env labels
- http_requests_total: Counter with service and env labels
- go_mem_alloc_bytes: Gauge with service and env labels

## Functions

### main

Application entry point that initializes configuration and starts the HTTP server.

**accepts:** none

**returns:** none (runs indefinitely)

**logic:**
```
read SERVICE_NAME from environment variable
if SERVICE_NAME is empty:
    set SERVICE_NAME to "ngpos-go-reference"

read ENV from environment variable
if ENV is empty:
    set ENV to "local"

log "Starting service:" with SERVICE_NAME
log "Environment:" with ENV

register handler for "/metrics" endpoint
log "Metrics available at /metrics"

start HTTP server on port 8080
if server fails to start:
    log fatal error and exit
```

### metricsHandler [inline]

Handles requests to the /metrics endpoint, returning Prometheus-style metrics.

**accepts:**
- w: ResponseWriter
- r: Request

**returns:** none (writes to ResponseWriter)

**logic:**
```
atomically increment requestCount by 1

read memory statistics from runtime

set response status to 200 OK

write metrics in Prometheus format:
    service_up{service="<SERVICE_NAME>",env="<ENV>"} 1
    http_requests_total{service="<SERVICE_NAME>",env="<ENV>"} <requestCount>
    go_mem_alloc_bytes{service="<SERVICE_NAME>",env="<ENV>"} <memory.Alloc>
```

### fmtUint

Formats an unsigned 64-bit integer as a decimal string.

**accepts:**
- v: UInt64

**returns:** Text

**logic:**
```
format v as decimal string using sprintf("%d", v)
return the formatted string
```

### NewGetHealth200Response

Constructor that creates a new GetHealth200Response with default values.

**accepts:** none

**returns:** GetHealth200Response

**logic:**
```
create new GetHealth200Response instance
return the instance
```

### NewGetHealth200ResponseWithDefaults

Constructor that creates a new GetHealth200Response with only default values set.

**accepts:** none

**returns:** GetHealth200Response

**logic:**
```
create new GetHealth200Response instance
return the instance
```

### GetHealth200Response.GetStatus

Returns the Status field value if set, otherwise returns empty string.

**accepts:** none (method on GetHealth200Response)

**returns:** Text

**logic:**
```
if receiver is nil or Status is nil:
    return empty string
return dereferenced Status value
```

### GetHealth200Response.GetStatusOk

Returns the Status field value and a boolean indicating if it was set.

**accepts:** none (method on GetHealth200Response)

**returns:** Tuple of (Optional Text, Boolean)

**logic:**
```
if receiver is nil or Status is nil:
    return (nil, false)
return (Status, true)
```

### GetHealth200Response.HasStatus

Checks if the Status field has been set.

**accepts:** none (method on GetHealth200Response)

**returns:** Boolean

**logic:**
```
if receiver is not nil and Status is not nil:
    return true
return false
```

### GetHealth200Response.SetStatus

Sets the Status field to the given value.

**accepts:**
- v: Text

**returns:** none

**logic:**
```
set receiver.Status to pointer of v
```

### GetHealth200Response.MarshalJSON

Serializes the response to JSON bytes.

**accepts:** none (method on GetHealth200Response)

**returns:** Result of Bytes

**logic:**
```
call ToMap() to get serializable map
if error:
    return empty bytes and error
marshal map to JSON and return
```

### GetHealth200Response.ToMap

Converts the response to a map for JSON serialization.

**accepts:** none (method on GetHealth200Response)

**returns:** Map of Text to Any

**logic:**
```
create empty map
if Status is not nil:
    add "status" key with Status value to map
return map
```

### NullableGetHealth200Response.Get

Returns the wrapped value.

**accepts:** none (method on NullableGetHealth200Response)

**returns:** Optional GetHealth200Response

**logic:**
```
return receiver.value
```

### NullableGetHealth200Response.Set

Sets the wrapped value and marks as set.

**accepts:**
- val: GetHealth200Response

**returns:** none

**logic:**
```
set receiver.value to val
set receiver.isSet to true
```

### NullableGetHealth200Response.IsSet

Checks if a value has been set.

**accepts:** none (method on NullableGetHealth200Response)

**returns:** Boolean

**logic:**
```
return receiver.isSet
```

### NullableGetHealth200Response.Unset

Clears the wrapped value and marks as unset.

**accepts:** none (method on NullableGetHealth200Response)

**returns:** none

**logic:**
```
set receiver.value to nil
set receiver.isSet to false
```

### NewNullableGetHealth200Response

Constructor that creates a NullableGetHealth200Response with an initial value.

**accepts:**
- val: GetHealth200Response

**returns:** NullableGetHealth200Response

**logic:**
```
create NullableGetHealth200Response with:
    value: val
    isSet: true
return the instance
```

### NullableGetHealth200Response.MarshalJSON

Serializes the nullable wrapper to JSON (serializes the inner value).

**accepts:** none (method on NullableGetHealth200Response)

**returns:** Result of Bytes

**logic:**
```
marshal receiver.value to JSON and return
```

### NullableGetHealth200Response.UnmarshalJSON

Deserializes JSON bytes into the nullable wrapper.

**accepts:**
- src: Bytes

**returns:** Optional Error

**logic:**
```
set receiver.isSet to true
unmarshal src into receiver.value
return any error from unmarshal
```

## API Endpoints

### GET /health

Health check endpoint to verify service is running.

**operationId:** getHealth

**responses:**
- 200: Service is healthy
  - Content-Type: application/json
  - Body: GetHealth200Response with status: "UP"

### GET /metrics

Prometheus-style metrics endpoint.

**operationId:** getMetrics

**responses:**
- 200: Metrics response
  - Content-Type: text/plain
  - Body: Prometheus metrics format

**example response:**
```
service_up{service="ngpos-go-reference",env="local"} 1
http_requests_total{service="ngpos-go-reference",env="local"} 42
go_mem_alloc_bytes{service="ngpos-go-reference",env="local"} 1234567
```

## Async Events

### Channel: receipt.events

Receipt event publishing channel.

**publish:**
- summary: Receipt event published
- payload:
  - receiptId: Text
  - eventType: Text

## Tests

No tests defined in source project.

## Dependencies

### Required
- Standard library HTTP server (net/http)
- JSON encoding (encoding/json)
- Runtime memory stats (runtime)
- Atomic operations (sync/atomic)
- Logging (log)
- OS environment (os)
- String formatting (fmt)

### Framework
- None (uses Go standard library only)

## Project Structure

```
/
├── main.go                                    # Application entry point with HTTP server
├── internal/
│   └── models/
│       └── model_get_health_200_response.go   # OpenAPI generated response model
├── spec/
│   ├── openapi.yaml                           # REST API specification
│   └── asyncapi.yaml                          # Async messaging specification
├── go.mod                                     # Go module definition
└── README.md                                  # Project documentation
```
