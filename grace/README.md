# Grace - Graceful HTTP Server Shutdown

## Installation

```bash
go get github.com/davidsugianto/go-pkgs/grace
```

## Usage

Replace your existing server start call:

```go
// BEFORE
http.ListenAndServe(":8080", handler)

// AFTER
grace.ServeHTTP(":8080", handler)
```

That's it! Works with any framework.

### Examples

**Standard Library**
```go
import "github.com/davidsugianto/go-pkgs/grace"

http.HandleFunc("/", handler)
grace.ServeHTTP(":8080", nil)
```

**Gin**
```go
r := gin.Default()
grace.ServeHTTP(":8080", r)
```

**Echo**
```go
e := echo.New()
grace.ServeHTTP(":8080", e)
```

**Any Framework**
```go
grace.ServeHTTP(":8080", yourHandler)
```

### HTTPS

```go
grace.ServeHTTPS(":8443", "cert.pem", "key.pem", handler)
```

### Custom Server

```go
server := &http.Server{...}
grace.ServeServer(server)
```

## What It Does

- Starts your HTTP server normally
- Listens for SIGINT/SIGTERM signals
- Stops accepting new connections
- Waits up to 30 seconds for active requests to complete
- Gracefully shuts down
