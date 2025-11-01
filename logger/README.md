# Logger - Structured Logging with OpenTelemetry Integration

A structured logging package built on [zerolog](https://github.com/rs/zerolog) with automatic OpenTelemetry span context correlation. This package provides high-performance, zero-allocation JSON logging with seamless integration for observability platforms like Grafana, Loki, and other log aggregation systems.

## Installation

```bash
go get github.com/davidsugianto/go-pkgs/logger
```

You'll also need to add the required dependencies:

```bash
go get github.com/rs/zerolog
go get go.opentelemetry.io/otel/trace
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/davidsugianto/go-pkgs/logger"
    "github.com/rs/zerolog"
)

func main() {
    // Create a logger
    log := logger.NewWithConfig(logger.Config{
        ServiceName: "my-service",
        Environment: "production",
        Format:      logger.FormatJSON,
        Level:       zerolog.InfoLevel,
    })

    // Basic logging
    log.Info().Msg("Application started")
    
    // Structured logging
    log.Info().
        Str("user_id", "12345").
        Int("status_code", 200).
        Msg("User logged in")
}
```

## Features

- ✅ **Zero Allocation** - High-performance logging with zero allocations for most operations
- ✅ **Structured Logging** - JSON-first design for log aggregation systems
- ✅ **OpenTelemetry Integration** - Automatic trace/span ID correlation
- ✅ **Multiple Formats** - JSON, Console, and Pretty output formats
- ✅ **Context Support** - Automatic span context extraction from Go context
- ✅ **Configurable** - Flexible configuration for different environments
- ✅ **Global Logger** - Convenient global logger for application-wide logging
- ✅ **Level Filtering** - Runtime log level configuration
- ✅ **Observability Ready** - Perfect for Grafana, Loki, and other observability platforms

## Configuration

### Basic Configuration

```go
log := logger.New() // Uses defaults: JSON format, Info level, stderr output
```

### Advanced Configuration

```go
log := logger.NewWithConfig(logger.Config{
    Output:      os.Stderr,              // Output destination
    Level:       zerolog.DebugLevel,     // Log level
    Format:      logger.FormatJSON,      // Output format
    ServiceName: "api-server",           // Service name
    Environment: "production",           // Environment
    TraceIDFieldName: "trace_id",        // Custom trace ID field name
    SpanIDFieldName:  "span_id",         // Custom span ID field name
    PrettyPrint: false,                  // Pretty print JSON
})
```

### Configuration Options

- `Output` (`io.Writer`) - Output destination (default: `os.Stderr`)
- `Level` (`zerolog.Level`) - Minimum log level (default: `InfoLevel`)
- `Format` (`string`) - Output format: `"json"`, `"console"`, or `"pretty"`
- `ServiceName` (`string`) - Service name to include in logs
- `Environment` (`string`) - Environment (e.g., `"production"`, `"staging"`, `"dev"`)
- `TraceIDFieldName` (`string`) - Field name for trace ID (default: `"trace_id"`)
- `SpanIDFieldName` (`string`) - Field name for span ID (default: `"span_id"`)
- `PrettyPrint` (`bool`) - Enable pretty JSON formatting (indented)

## Output Formats

### JSON Format (Production)

JSON format is recommended for production use and log aggregation systems:

```go
log := logger.NewWithConfig(logger.Config{
    Format: logger.FormatJSON,
})

log.Info().
    Str("user_id", "12345").
    Int("status_code", 200).
    Msg("Request completed")
```

Output:
```json
{"level":"info","time":"2024-01-15T10:30:00Z","user_id":"12345","status_code":200,"message":"Request completed"}
```

### Console Format (Development)

Human-readable console format for development:

```go
log := logger.NewWithConfig(logger.Config{
    Format: logger.FormatConsole,
})
```

Output:
```
INF 10:30:00 user_id=12345 status_code=200 Request completed
```

### Pretty Format (Development)

Colorized pretty format with better readability:

```go
log := logger.NewWithConfig(logger.Config{
    Format: logger.FormatPretty,
})
```

Output:
```
10:30:00 | INF | user_id=12345 status_code=200 Request completed
```

## Log Levels

```go
log.Trace().Msg("Trace level message")   // Most verbose
log.Debug().Msg("Debug level message")   // Debug information
log.Info().Msg("Info level message")     // General information (default)
log.Warn().Msg("Warning message")        // Warning
log.Error().Msg("Error message")         // Error
log.Fatal().Msg("Fatal message")         // Fatal error (exits)
log.Panic().Msg("Panic message")         // Panic (panics)
```

### Setting Log Level

```go
// Set level when creating logger
log := logger.NewWithConfig(logger.Config{
    Level: zerolog.DebugLevel,
})

// Change level at runtime
log.SetLevel(zerolog.WarnLevel)
```

## Structured Logging

### Adding Fields

```go
log.Info().
    Str("key1", "value1").        // String field
    Int("key2", 42).              // Integer field
    Float64("key3", 3.14).        // Float field
    Bool("key4", true).           // Boolean field
    Dur("duration_ms", 150*time.Millisecond).  // Duration field
    Msg("Message")
```

### Field Types

- `Str(key, value)` - String
- `Bool(key, value)` - Boolean
- `Int(key, value)` - Integer (32-bit)
- `Int8/Int16/Int64(key, value)` - Integer variants
- `Uint/Uint8/Uint16/Uint32/Uint64(key, value)` - Unsigned integers
- `Float32/Float64(key, value)` - Floating point
- `Dur(key, value)` - Duration
- `Err(error)` - Error
- `Interface(key, value)` - Any Go interface

### Error Logging

```go
err := errors.New("something went wrong")
log.Error().
    Err(err).
    Str("operation", "user-creation").
    Msg("Failed to create user")
```

## OpenTelemetry Integration

### Automatic Span Correlation

The logger automatically extracts trace and span IDs from OpenTelemetry spans in the context:

```go
import (
    "context"
    "github.com/davidsugianto/go-pkgs/logger"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func handler(ctx context.Context) {
    log := logger.New()
    
    // Start a span
    tracer := otel.Tracer("my-service")
    ctx, span := tracer.Start(ctx, "http.request")
    defer span.End()
    
    // Create logger with context - automatically adds trace_id and span_id
    ctxLogger := log.WithContext(ctx)
    ctxLogger.Info().Msg("Request started")
    
    // Nested operations also get trace context
    dbOperation(ctx, ctxLogger)
}

func dbOperation(ctx context.Context, log *logger.Logger) {
    tracer := otel.Tracer("my-service")
    ctx, span := tracer.Start(ctx, "db.query")
    defer span.End()
    
    // New span context is automatically captured
    dbLog := log.WithContext(ctx)
    dbLog.Info().
        Str("query", "SELECT * FROM users").
        Msg("Query executed")
}
```

The logs will automatically include trace and span IDs:

```json
{"level":"info","time":"2024-01-15T10:30:00Z","trace_id":"abc123...","span_id":"def456...","message":"Query executed","query":"SELECT * FROM users"}
```

### Custom Field Names

You can customize trace and span ID field names:

```go
log := logger.NewWithConfig(logger.Config{
    TraceIDFieldName: "otel_trace_id",
    SpanIDFieldName:  "otel_span_id",
})
```

## Global Logger

For application-wide logging convenience:

```go
// Set up global logger
logger.SetGlobal(logger.NewWithConfig(logger.Config{
    ServiceName: "my-service",
    Format:      logger.FormatJSON,
}))

// Use global helpers
logger.Info().Msg("Global info message")
logger.Warn().Msg("Global warning")
logger.Error().Err(err).Msg("Global error")

// With context
ctxLogger := logger.WithContext(ctx)
ctxLogger.Info().Msg("Global logger with trace context")

// Change level
logger.SetLevel(zerolog.DebugLevel)
```

## Real-World Examples

### HTTP Server Logging

```go
func httpHandler(w http.ResponseWriter, r *http.Request) {
    log := logger.NewWithConfig(logger.Config{
        ServiceName: "http-server",
        Format:      logger.FormatJSON,
    })
    
    start := time.Now()
    
    // Your handler logic
    // ...
    
    log.Info().
        Str("method", r.Method).
        Str("path", r.URL.Path).
        Int("status", 200).
        Dur("duration_ms", time.Since(start)).
        Msg("Request completed")
}
```

### With OpenTelemetry Tracing

```go
func httpHandlerWithTracing(w http.ResponseWriter, r *http.Request) {
    log := logger.New()
    
    // Start OpenTelemetry span
    tracer := otel.Tracer("http-server")
    ctx, span := tracer.Start(r.Context(), "http.request",
        trace.WithAttributes(
            attribute.String("http.method", r.Method),
            attribute.String("http.path", r.URL.Path),
        ),
    )
    defer span.End()
    
    // Create logger with trace context
    ctxLog := log.WithContext(ctx)
    
    start := time.Now()
    
    // Your handler logic
    // ...
    
    ctxLog.Info().
        Str("method", r.Method).
        Str("path", r.URL.Path).
        Int("status", 200).
        Dur("duration_ms", time.Since(start)).
        Msg("Request completed")
}
```

### Request Middleware

```go
func loggingMiddleware(next http.Handler) http.Handler {
    log := logger.NewWithConfig(logger.Config{
        ServiceName: "api",
        Format:      logger.FormatJSON,
    })
    
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        lw := &loggingResponseWriter{ResponseWriter: w, statusCode: 200}
        
        // Continue to next handler
        next.ServeHTTP(lw, r)
        
        // Log request
        log.Info().
            Str("method", r.Method).
            Str("path", r.URL.Path).
            Str("remote_addr", r.RemoteAddr).
            Int("status_code", lw.statusCode).
            Dur("duration_ms", time.Since(start)).
            Msg("HTTP request")
    })
}

type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
    lw.statusCode = code
    lw.ResponseWriter.WriteHeader(code)
}
```

### Database Operations

```go
func QueryUsers(ctx context.Context, log *logger.Logger) ([]User, error) {
    tracer := otel.Tracer("database")
    ctx, span := tracer.Start(ctx, "db.query_users")
    defer span.End()
    
    queryLog := log.WithContext(ctx)
    queryLog.Debug().Msg("Executing user query")
    
    // Execute query
    users, err := db.Query("SELECT * FROM users")
    if err != nil {
        queryLog.Error().
            Err(err).
            Msg("Failed to query users")
        return nil, err
    }
    
    queryLog.Info().
        Int("count", len(users)).
        Msg("Users queried successfully")
    
    return users, nil
}
```

## Observability Platforms

### Grafana Loki

Logs from this logger work seamlessly with Grafana Loki. JSON format is recommended:

```go
log := logger.NewWithConfig(logger.Config{
    ServiceName: "my-service",
    Format:      logger.FormatJSON,
})
```

The structured JSON logs with trace/span IDs enable powerful correlation in Grafana:
- Correlate logs with traces
- Filter by service, environment, or trace ID
- Create dashboards with log metrics

### CloudWatch Logs / Logz.io / Datadog

All major log aggregation platforms support JSON logs with trace correlation:

```go
log := logger.NewWithConfig(logger.Config{
    ServiceName: "my-service",
    Environment: "production",
    Format:      logger.FormatJSON,
})
```

## Performance

zerolog is designed for high performance:
- **Zero allocation** for most logging operations
- **Fast JSON encoding**
- **Sampling support** for high-traffic scenarios
- **Efficient field builders**

Benchmark results (from zerolog documentation):

```
BenchmarkLogEmpty-8        100000000    19.1 ns/op     0 B/op       0 allocs/op
BenchmarkDisabled-8        500000000    4.07 ns/op     0 B/op       0 allocs/op
BenchmarkInfo-8            30000000     42.5 ns/op     0 B/op       0 allocs/op
BenchmarkContextFields-8   30000000     44.9 ns/op     0 B/op       0 allocs/op
BenchmarkLogFields-8       10000000     184 ns/op      0 B/op       0 allocs/op
```

## Best Practices

1. **Use JSON format in production** - Better for log aggregation systems
2. **Set appropriate log levels** - Use Debug in development, Info in production
3. **Include context** - Add relevant fields to every log message
4. **Use structured fields** - Prefer structured fields over string formatting
5. **Correlate with traces** - Always use `WithContext()` when you have a span context
6. **Add service metadata** - Set ServiceName and Environment in configuration
7. **Monitor performance** - Use sampling for high-volume logs if needed
8. **Handle errors gracefully** - Use Error level with proper error context

## Troubleshooting

### Logs not showing trace/span IDs

Make sure you're using `WithContext()` with a context that contains an active OpenTelemetry span:

```go
// ❌ Wrong - no span context
log.Info().Msg("No trace IDs")

// ✅ Correct - with span context
ctxLogger := log.WithContext(ctx)
ctxLogger.Info().Msg("Has trace IDs")
```

### Too many logs in production

Set appropriate log level:

```go
log := logger.NewWithConfig(logger.Config{
    Level: zerolog.InfoLevel,  // Only info and above
})

// Or change at runtime
log.SetLevel(zerolog.WarnLevel)  // Only warnings and errors
```

### JSON not pretty printed

Use Console or Pretty format for development:

```go
log := logger.NewWithConfig(logger.Config{
    Format: logger.FormatConsole,  // Human-readable
})
```

## Migration from Other Loggers

### From logrus

Replace:
```go
logrus.WithFields(logrus.Fields{
    "user": "alice",
}).Info("User logged in")
```

With:
```go
log.Info().
    Str("user", "alice").
    Msg("User logged in")
```

### From standard library log

Replace:
```go
log.Printf("User %s logged in with status %d", user, status)
```

With:
```go
log.Info().
    Str("user", user).
    Int("status", status).
    Msg("User logged in")
```

### From zap

Replace:
```go
zapLogger.Info("User logged in",
    zap.String("user", user),
    zap.Int("status", status),
)
```

With:
```go
log.Info().
    Str("user", user).
    Int("status", status).
    Msg("User logged in")
```

## API Reference

### Types

- `Logger` - Main logger struct
- `Config` - Logger configuration

### Functions

- `New()` - Create logger with defaults
- `NewWithConfig(cfg Config)` - Create logger with custom config
- `GetGlobal()` - Get global logger instance
- `SetGlobal(logger *Logger)` - Set global logger
- `WithContext(ctx context.Context)` - Get logger with context
- `SetLevel(level zerolog.Level)` - Set global log level

### Logger Methods

- `WithContext(ctx)` - Add span context from context
- `With()` - Create event builder with fields
- `Info()`, `Debug()`, `Warn()`, `Error()`, `Fatal()`, `Panic()`, `Trace()` - Create log events
- `GetLevel()` - Get current log level
- `SetLevel(level)` - Set log level

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## References

- [zerolog](https://github.com/rs/zerolog) - The underlying logging library
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/) - OpenTelemetry integration
- [Grafana Loki](https://grafana.com/docs/loki/latest/) - Log aggregation system
- [OpenTelemetry Logging](https://opentelemetry.io/docs/specs/otel/logs/) - OpenTelemetry logging specification

