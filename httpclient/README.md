# HTTPClient - Simple & Powerful HTTP Client

A lightweight, idiomatic HTTP client wrapper for Go with automatic JSON serialization, error handling, and flexible configuration options.

## Installation

```bash
go get github.com/davidsugianto/go-pkgs/httpclient
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/davidsugianto/go-pkgs/httpclient"
)

func main() {
    client := httpclient.New("https://api.example.com")
    
    resp, err := client.Get(context.Background(), "/users/1", nil)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    // Use resp as needed
}
```

## Features

- ✅ **Simple API** - Clean, intuitive methods for all HTTP verbs
- ✅ **Auto JSON** - Automatic JSON serialization/deserialization
- ✅ **Context Support** - Full context.Context integration for timeouts and cancellation
- ✅ **Flexible Bodies** - Supports structs, strings, bytes, or nil
- ✅ **Custom Headers** - Easy header configuration
- ✅ **Error Handling** - Automatic error handling for 4xx/5xx responses
- ✅ **Raw Content** - Support for custom content types (XML, plain text, etc.)

## Usage

### Creating a Client

```go
// Basic client with default 10s timeout
client := httpclient.New("https://api.example.com")

// With custom timeout
client := httpclient.New(
    "https://api.example.com",
    httpclient.WithTimeout(30 * time.Second),
)

// With custom headers
client := httpclient.New(
    "https://api.example.com",
    httpclient.WithHeaders(map[string]string{
        "Authorization": "Bearer token123",
        "User-Agent":    "my-app/1.0",
    }),
)

// With both timeout and headers
client := httpclient.New(
    "https://api.example.com",
    httpclient.WithTimeout(30 * time.Second),
    httpclient.WithHeaders(map[string]string{
        "Authorization": "Bearer token123",
    }),
)
```

### GET Request

```go
// Simple GET
resp, err := client.Get(ctx, "/users/1", nil)
if err != nil {
    return err
}
defer resp.Body.Close()

// Decode JSON response
var user User
json.NewDecoder(resp.Body).Decode(&user)
```

### POST Request

```go
// POST with JSON body (struct automatically serialized)
type CreateUser struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

newUser := CreateUser{
    Name:  "John Doe",
    Email: "john@example.com",
}

resp, err := client.Post(ctx, "/users", newUser)
if err != nil {
    return err
}
defer resp.Body.Close()

// POST with raw body
resp, err := client.PostRaw(ctx, "/data", "<xml>...</xml>", "application/xml")
```

### PUT Request

```go
// PUT with JSON body
type UpdateUser struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

updateData := UpdateUser{
    Name:  "Jane Doe",
    Email: "jane@example.com",
}

resp, err := client.Put(ctx, "/users/1", updateData)
if err != nil {
    return err
}
defer resp.Body.Close()

// PUT with raw body
resp, err := client.PutRaw(ctx, "/users/1", "<xml>...</xml>", "application/xml")
```

### DELETE Request

```go
// DELETE without body
resp, err := client.Delete(ctx, "/users/1", nil)
if err != nil {
    return err
}
defer resp.Body.Close()

// DELETE with body (if your API requires it)
deleteReason := map[string]string{
    "reason": "No longer needed",
}
resp, err := client.Delete(ctx, "/users/1", deleteReason)
```

### Using Context

```go
// Context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.Get(ctx, "/slow-endpoint", nil)
if err != nil {
    // Handles timeout automatically
    return err
}

// Context with cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(2 * time.Second)
    cancel() // Cancel request after 2 seconds
}()
resp, err := client.Get(ctx, "/endpoint", nil)
```

### Body Types

The client accepts various body types:

```go
// Struct (automatically JSON encoded)
type Data struct {
    Name string `json:"name"`
}
resp, err := client.Post(ctx, "/endpoint", Data{Name: "test"})

// String
resp, err := client.PostRaw(ctx, "/endpoint", "raw string", "text/plain")

// []byte
resp, err := client.PostRaw(ctx, "/endpoint", []byte("raw bytes"), "application/octet-stream")

// nil (no body)
resp, err := client.Get(ctx, "/endpoint", nil)
```

### Error Handling

The client automatically handles HTTP error responses (4xx, 5xx):

```go
resp, err := client.Get(ctx, "/users/999", nil)
if err != nil {
    // err contains the response body message for 4xx/5xx responses
    fmt.Printf("Error: %v\n", err)
    return
}
// Success (2xx response)
```

### Complete Example

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/davidsugianto/go-pkgs/httpclient"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    client := httpclient.New(
        "https://api.example.com",
        httpclient.WithTimeout(30*time.Second),
        httpclient.WithHeaders(map[string]string{
            "Authorization": "Bearer token123",
        }),
    )

    ctx := context.Background()

    // GET
    resp, err := client.Get(ctx, "/users/1", nil)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    var user User
    json.NewDecoder(resp.Body).Decode(&user)
    fmt.Printf("User: %+v\n", user)

    // POST
    newUser := User{Name: "John", Email: "john@example.com"}
    resp, err = client.Post(ctx, "/users", newUser)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
}
```

## API Reference

### Functions

#### `New(baseURL string, opts ...Option) *Client`

Creates a new HTTP client instance.

- `baseURL`: Base URL for all requests (e.g., `"https://api.example.com"`)
- `opts`: Optional configuration functions

#### `WithTimeout(timeout time.Duration) Option`

Sets a custom timeout for HTTP requests. Default is 10 seconds.

#### `WithHeaders(headers map[string]string) Option`

Sets default headers for all requests.

### Methods

All methods return `(*http.Response, error)` and follow the same pattern.

#### `Get(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)`

Performs a GET request.

#### `Post(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)`

Performs a POST request with JSON body (struct automatically serialized).

#### `PostRaw(ctx context.Context, endpoint string, rawBody string, contentType string) (*http.Response, error)`

Performs a POST request with raw body and custom content type.

#### `Put(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)`

Performs a PUT request with JSON body (struct automatically serialized).

#### `PutRaw(ctx context.Context, endpoint string, rawBody string, contentType string) (*http.Response, error)`

Performs a PUT request with raw body and custom content type.

#### `Delete(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)`

Performs a DELETE request. Body is optional (can be `nil`).

## Error Handling

- Network errors are returned as-is
- HTTP error responses (status code >= 400) return an error containing the response body
- JSON marshaling errors are returned immediately

## Examples

See the `example/` directory for a complete working example.

```bash
cd example
go run main.go
```

