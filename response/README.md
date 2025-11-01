# Response - Consistent API Response Utilities

A lightweight package for creating consistent JSON API responses with a standard format (`code`, `data`, `error`). Provides convenient helper functions for common HTTP status codes and response patterns.

## Installation

```bash
go get github.com/davidsugianto/go-pkgs/response
```

## Quick Start

```go
package main

import (
    "net/http"
    "github.com/davidsugianto/go-pkgs/response"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // Success response
    response.Success(w, map[string]string{"message": "Hello"})
    
    // Error response
    response.BadRequest(w, errors.New("invalid input"))
}
```

## Features

- ✅ **Standard Format** - Consistent `code`, `data`, `error` response structure
- ✅ **Type-Safe** - Helper functions for common HTTP status codes
- ✅ **Simple API** - Clean, intuitive functions for all response types
- ✅ **Flexible** - Support for custom status codes and messages
- ✅ **Zero Dependencies** - Uses only standard library

## Response Format

All responses follow this standard JSON format:

```json
{
  "code": 200,
  "data": { ... },
  "error": "..."
}
```

- `code`: HTTP status code (always present)
- `data`: Response payload (omitted if empty)
- `error`: Error message (omitted if empty)

## Usage

### Success Responses

#### `Success(w, data)` - Status 200 OK

```go
user := User{ID: 1, Name: "John"}
response.Success(w, user)

// Response:
// {
//   "code": 200,
//   "data": {"id": 1, "name": "John"}
// }
```

#### `Created(w, data)` - Status 201 Created

```go
newUser := User{ID: 2, Name: "Jane"}
response.Created(w, newUser)

// Response:
// {
//   "code": 201,
//   "data": {"id": 2, "name": "Jane"}
// }
```

#### `NoContent(w)` - Status 204 No Content

```go
response.NoContent(w)

// Response: Empty body with 204 status
```

### Error Responses

#### `BadRequest(w, err)` - Status 400

```go
response.BadRequest(w, errors.New("invalid input"))

// Response:
// {
//   "code": 400,
//   "error": "invalid input"
// }
```

#### `Unauthorized(w, err)` - Status 401

```go
response.Unauthorized(w, errors.New("authentication required"))

// Response:
// {
//   "code": 401,
//   "error": "authentication required"
// }
```

#### `Forbidden(w, err)` - Status 403

```go
response.Forbidden(w, errors.New("access denied"))

// Response:
// {
//   "code": 403,
//   "error": "access denied"
// }
```

#### `NotFound(w, err)` - Status 404

```go
response.NotFound(w, errors.New("user not found"))

// Response:
// {
//   "code": 404,
//   "error": "user not found"
// }
```

#### `InternalServerError(w, err)` - Status 500

```go
response.InternalServerError(w, errors.New("database error"))

// Response:
// {
//   "code": 500,
//   "error": "database error"
// }
```

### Custom Responses

#### `JSON(w, statusCode, data)` - Custom Status Code

```go
// Any status code with data
response.JSON(w, http.StatusAccepted, map[string]string{
    "message": "Request accepted",
})

// Response:
// {
//   "code": 202,
//   "data": {"message": "Request accepted"}
// }
```

#### `Error(w, statusCode, err)` - Custom Error Status

```go
// Custom error status code
response.Error(w, http.StatusConflict, errors.New("resource exists"))

// Response:
// {
//   "code": 409,
//   "error": "resource exists"
// }
```

#### `StatusCode(w, statusCode, message)` - Custom Status with Message

```go
response.StatusCode(w, http.StatusTeapot, "I'm a teapot")

// Response:
// {
//   "code": 418,
//   "error": "I'm a teapot"
// }
```

## Complete Example

```go
package main

import (
	"errors"
	"net/http"

	"github.com/davidsugianto/go-pkgs/response"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		response.BadRequest(w, errors.New("id parameter required"))
		return
	}

	// Simulate fetching user
	user, err := fetchUser(userID)
	if err != nil {
		response.NotFound(w, err)
		return
	}

	response.Success(w, user)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.BadRequest(w, err)
		return
	}

	// Validate user
	if user.Name == "" || user.Email == "" {
		response.BadRequest(w, errors.New("name and email are required"))
		return
	}

	// Create user
	newUser, err := saveUser(user)
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.Created(w, newUser)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if err := deleteUser(userID); err != nil {
		if err == ErrUserNotFound {
			response.NotFound(w, err)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	response.NoContent(w)
}

func main() {
	http.HandleFunc("/users", getUserHandler)
	http.HandleFunc("/users/create", createUserHandler)
	http.HandleFunc("/users/delete", deleteUserHandler)

	http.ListenAndServe(":8080", nil)
}
```

## Response Structure

The `Response` struct is available if you need to construct responses manually:

```go
type Response struct {
    Code  int         `json:"code"`
    Data  interface{} `json:"data,omitempty"`
    Error string      `json:"error,omitempty"`
}
```

Example:

```go
resp := Response{
    Code: 200,
    Data: map[string]string{"message": "success"},
}

json.NewEncoder(w).Encode(resp)
```

## Error Handling

All response functions return an error (typically from JSON encoding). It's good practice to handle these:

```go
if err := response.Success(w, data); err != nil {
    log.Printf("Failed to write response: %v", err)
    // Response may have been partially written
}
```

However, in most cases, JSON encoding errors are rare and can be logged rather than causing a panic.

## Best Practices

1. **Use Appropriate Status Codes**
   - Use `Success()` for successful GET/PUT/PATCH requests
   - Use `Created()` for successful POST requests
   - Use `NoContent()` for successful DELETE requests

2. **Error Messages**
   - Provide clear, user-friendly error messages
   - Avoid exposing internal implementation details

3. **Data Types**
   - Use structs for structured data
   - Use maps for simple key-value pairs
   - Use slices for arrays

4. **Consistency**
   - Always use the response package for API responses
   - Maintain consistent error message format

## API Reference

### Functions

#### `JSON(w http.ResponseWriter, statusCode int, data interface{}) error`

Writes a JSON response with the given status code and data.

#### `Success(w http.ResponseWriter, data interface{}) error`

Writes a success JSON response with status 200 OK.

#### `Created(w http.ResponseWriter, data interface{}) error`

Writes a success JSON response with status 201 Created.

#### `NoContent(w http.ResponseWriter)`

Writes an empty response with status 204 No Content.

#### `Error(w http.ResponseWriter, statusCode int, err error) error`

Writes an error JSON response with the given status code and error message.

#### `BadRequest(w http.ResponseWriter, err error) error`

Writes a 400 Bad Request error response.

#### `Unauthorized(w http.ResponseWriter, err error) error`

Writes a 401 Unauthorized error response.

#### `Forbidden(w http.ResponseWriter, err error) error`

Writes a 403 Forbidden error response.

#### `NotFound(w http.ResponseWriter, err error) error`

Writes a 404 Not Found error response.

#### `InternalServerError(w http.ResponseWriter, err error) error`

Writes a 500 Internal Server Error response.

#### `StatusCode(w http.ResponseWriter, statusCode int, message string) error`

Writes a JSON response with the given status code and message string.

## Examples

See the `example/` directory for a complete working example.

```bash
cd example
go run main.go
```

Then try the endpoints:
- `GET http://localhost:8080/success`
- `GET http://localhost:8080/created`
- `GET http://localhost:8080/not-found`
- `GET http://localhost:8080/bad-request`
- `GET http://localhost:8080/internal-error`
- `GET http://localhost:8080/no-content`
- `GET http://localhost:8080/custom`

