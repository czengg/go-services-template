# API Client Package

A flexible and extensible HTTP client package that provides optional authentication handling, request/response processing, and error management. This package is designed to simplify API integrations while providing robust functionality and configuration options.

## Features 

- Configurable HTTP client with timeout control
- Token-based authentication with automatic token refresh
- Request middleware support
- Flexible request options
- Comprehensive error handling
- Thread-safe token management

## Installation

```
import "template/packages/api-client-go"
```

## Core Components

### Client

The main client struct provides the foundation for making HTTP requests:

```
client, err := NewClient(
    "https://api.example.com",
    authenticator,
    WithTimeout(30 * time.Second),
)
```

#### Client Configuration Options

The client can be configured using various options when creating a new instance:

```
// Configuration options for NewClient
WithTimeout(timeout time.Duration)        // Set custom timeout duration
WithHTTPClient(httpClient *http.Client)   // Use custom HTTP client

// Example
client, err := NewClient(
    "https://api.example.com",
    authenticator,
    WithTimeout(30 * time.Second),
    WithHTTPClient(customHTTPClient),
)
```

#### Request Options

Each request can be customized using these options:

```
// Request configuration options
WithMethod(method string)                 // Set HTTP method (GET, POST, PUT, DELETE, PATCH)
WithHeaders(headers map[string]string)    // Add custom headers to the request
WithBody(body interface{})                // Set request body (will be JSON-encoded)
WithQueryParams(params url.Values)        // Add URL query parameters
WithLogger(logger *zap.Logger)            // Configure request-specific logging

// Example
resp, err := client.Request(
    ctx,
    "/endpoint",
    WithMethod(MethodPost),
    WithHeaders(map[string]string{
        "Custom-Header": "value",
    }),
    WithBody(requestBody),
    WithQueryParams(url.Values{
        "filter": []string{"active"},
    }),
    WithLogger(logger),
)
```

### Authentication

The package provides flexible authentication support through interfaces:

```
// Core authentication interface
type Authenticator interface {
    Authenticate(*http.Request) error
}

// Token provider interface
type TokenProvider interface {
    GetToken(context.Context) (string, error)
}
```

#### Authentication Options

1. TokenAuthenticator - For authenticated endpoints

```
// Create a token authenticator
auth := NewTokenAuthenticator(tokenProvider)
```

2. No Authentication - For public endpoints

```
// Pass nil as the authenticator
client, err := NewClient(baseURL, nil)
```

3. Basic Authenticator with username and password

```
// Create a basic authenticator
auth := client.NewBasicAuthenticator("username", "password")
```

## Making Requests

The client provides a flexible Request method with various options:

```
resp, err := client.Request(
    ctx,
    "/api/endpoint",
    WithMethod("POST"),
    WithHeaders(map[string]string{"Custom-Header": "value"}),
    WithBody(requestBody),
    WithQueryParams(url.Values{"key": {"value"}}),
)
```

### Available Request Options

```
WithMethod(method string)           // Set HTTP method
WithHeaders(headers map[string]string) // Add custom headers
WithBody(body interface{})          // Set request body
WithQueryParams(params url.Values)  // Add query parameters
WithLogger(logger *zap.Logger)      // Configure logging
```

## Error Handling

The package provides structured error types:

```
// Client-specific errors
type ClientError struct {
    Code    string
    Message string
    Details string
}

// HTTP response errors
type HTTPError struct {
    StatusCode int
    Body       []byte
}
```

## Usage Examples

### Basic Request

```
client, _ := api_client.NewClient(
	config.BaseURL,
	nil,
	api_client.WithTimeout(30*time.Second),
)
```

### Authenticated Request

```
tokenProvider := NewPartnerTokenProvider(TokenTypePartner, Config{
	BaseURL:       baseURL,
	PartnerID:     partnerID,
	PartnerSecret: partnerSecret,
})
client, err := api_client.NewClient(baseURL, api_client.NewTokenAuthenticator(tokenProvider))
```

## Best Practices

1. Token Management
 - Use TokenAuthenticator for automatic token refresh
 - Implement proper token expiration handling
 - Use thread-safe token storage
2. Request Handling
 - Always provide a context for timeouts and cancellation
 - Use appropriate request options for different endpoints
 - Handle response bodies properly
3. Error Handling
 - Check for specific error types
 - Log appropriate error details
 - Implement proper error recovery strategies
4. Configuration
 - Set appropriate timeouts
 - Configure custom HTTP clients when needed
 - Use proper base URLs

## Dependencies

- `go.uber.org/zap`: Structured logging
- Standard library packages:
 - `net/http`
 - `encoding/json`
 - `context`
 - `sync`
 - `time`

## Thread Safety

The package is designed to be thread-safe:

- Token management uses proper synchronization
- HTTP client can be safely used concurrently
- Request options are created per-request

## Logging

The package supports structured logging through `zap.Logger`:

- Request details can be logged
- Error conditions are captured
- Authentication events can be tracked

## Constants

```
const (
    DefaultTimeout = 30 * time.Second

    MethodGet    = http.MethodGet
    MethodPost   = http.MethodPost
    MethodPut    = http.MethodPut
    MethodDelete = http.MethodDelete
    MethodPatch  = http.MethodPatch
)
```

## Error Codes

```
var (
    ErrInvalidBaseURL = &ClientError{
        Code:    "INVALID_BASE_URL",
        Message: "invalid base url",
        Details: "base url cannot be empty",
    }
)
```