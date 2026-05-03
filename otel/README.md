# OpenTelemetry Integration Package

A comprehensive, reusable Go package for OpenTelemetry integration providing centralized configuration for distributed tracing, metrics, and logging.

## Overview

The `otel` package provides centralized OpenTelemetry configuration that enables observability across all library packages. It supports the three pillars of observability:

- **Traces** - Distributed tracing for request flows
- **Metrics** - Performance and health measurements
- **Logs** - Structured logging via OpenTelemetry standard

## Features

- **Unified Configuration**: Single config object for all telemetry pillars
- **Selective Enablement**: Enable only the telemetry you need
- **No-op by Default**: Zero overhead when providers are not configured
- **Method Chaining**: Fluent API for easy configuration
- **Trace-Aware Logging**: Automatic trace correlation in logs
- **Graceful Shutdown**: Proper resource cleanup
- **Production Ready**: Thread-safe and tested

## Installation

```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
go get go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc
go get go.opentelemetry.io/otel/sdk
```

## Quick Start

### Basic Setup

```go
package main

import (
    "context"
    "log"
    "your-module/otel"
)

func main() {
    ctx := context.Background()

    // Configure OpenTelemetry
    config := otel.NewConfig("my-service").
        WithServiceVersion("1.0.0").
        WithEnvironment("production").
        WithTraceEndpoint("localhost:4317").
        WithMetricEndpoint("localhost:4317")

    // Initialize provider
    provider, err := otel.NewProvider(ctx, config)
    if err != nil {
        log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
    }
    defer provider.Shutdown(context.Background())

    // Your application code here
}
```

## Configuration

### Environment Variables

The package respects standard OpenTelemetry environment variables:

- `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` - Trace collector endpoint
- `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT` - Metrics collector endpoint
- `ENVIRONMENT` - Deployment environment

### Configuration Options

```go
config := otel.NewConfig("my-service").
    WithServiceVersion("1.2.3").              // Service version
    WithEnvironment("staging").                // Environment (dev/staging/prod)
    WithTracing(true).                        // Enable/disable tracing
    WithTraceEndpoint("localhost:4317").      // OTLP trace endpoint
    WithMetrics(true).                        // Enable/disable metrics
    WithMetricEndpoint("localhost:4317").     // OTLP metrics endpoint
    WithLogging(true).                        // Enable/disable logging
    WithLogLevel(slog.LevelDebug).           // Log level
    WithResourceAttribute("team", "backend"). // Custom resource attributes
    WithResourceAttribute("region", "us-west")
```

## Usage Examples

### 1. Distributed Tracing

```go
// Get a tracer
tracer := provider.Tracer("my-component")

// Start a span
ctx, span := tracer.Start(ctx, "processOrder")
defer span.End()

// Add attributes
span.SetAttributes(
    attribute.String("user.id", userID),
    attribute.Int("order.items", itemCount),
    attribute.Float64("order.total", total),
)

// Add events
span.AddEvent("validation.started")
span.AddEvent("payment.completed", trace.WithAttributes(
    attribute.String("payment.method", "credit_card"),
))

// Record errors
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
}
```

### 2. Metrics Collection

```go
// Get a meter
meter := provider.GetMeterProvider().Meter("my-component")

// Create a counter
requestCounter, _ := meter.Int64Counter(
    "requests.total",
    metric.WithDescription("Total number of requests"),
    metric.WithUnit("{request}"),
)

// Increment counter with labels
requestCounter.Add(ctx, 1, metric.WithAttributes(
    attribute.String("method", "GET"),
    attribute.String("status", "200"),
    attribute.String("endpoint", "/api/users"),
))

// Create a histogram
responseTime, _ := meter.Float64Histogram(
    "request.duration",
    metric.WithDescription("Request duration in milliseconds"),
    metric.WithUnit("ms"),
)

// Record values
responseTime.Record(ctx, duration, metric.WithAttributes(
    attribute.String("endpoint", "/api/users"),
))

// Create a gauge (async)
meter.Int64ObservableGauge(
    "memory.usage",
    metric.WithDescription("Current memory usage"),
    metric.WithUnit("By"),
    metric.WithInt64Callback(func(ctx context.Context, o metric.Int64Observer) error {
        o.Observe(getCurrentMemoryUsage())
        return nil
    }),
)
```

### 3. Structured Logging with Trace Context

```go
// Get logger with automatic trace correlation
logger := provider.Logger(ctx)

// Log with structured fields
logger.Info("Processing request",
    slog.String("user_id", userID),
    slog.Int("item_count", count),
    slog.Duration("elapsed", elapsed),
)

// Log levels
logger.Debug("Debug information")
logger.Info("Information message")
logger.Warn("Warning message")
logger.Error("Error occurred", slog.Any("error", err))

// When a span is active, logs automatically include trace_id and span_id
```

### 4. HTTP Middleware

```go
func TracingMiddleware(provider *otel.Provider) func(http.Handler) http.Handler {
    tracer := provider.Tracer("http-server")
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, span := tracer.Start(r.Context(), 
                r.Method+" "+r.URL.Path,
                trace.WithAttributes(
                    attribute.String("http.method", r.Method),
                    attribute.String("http.url", r.URL.String()),
                    attribute.String("http.user_agent", r.UserAgent()),
                ),
            )
            defer span.End()

            // Add logger with trace context
            logger := provider.Logger(ctx)
            logger.Info("Incoming request")

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// Usage
http.Handle("/", TracingMiddleware(provider)(myHandler))
```

### 5. Database Query Tracing

```go
func QueryDatabase(ctx context.Context, provider *otel.Provider, query string) error {
    tracer := provider.Tracer("database")
    ctx, span := tracer.Start(ctx, "db.query",
        trace.WithAttributes(
            attribute.String("db.system", "postgresql"),
            attribute.String("db.statement", query),
        ),
    )
    defer span.End()

    logger := provider.Logger(ctx)
    logger.Debug("Executing query", slog.String("query", query))

    // Execute query
    result, err := db.QueryContext(ctx, query)
    if err != nil {
        span.RecordError(err)
        logger.Error("Query failed", slog.Any("error", err))
        return err
    }

    span.SetAttributes(attribute.Int("db.rows_affected", rowsAffected))
    return nil
}
```

## Best Practices

### 1. Always Propagate Context

```go
// Good ✅
func ProcessOrder(ctx context.Context, order Order) error {
    ctx, span := tracer.Start(ctx, "processOrder")
    defer span.End()
    
    // Pass context to child functions
    if err := validateOrder(ctx, order); err != nil {
        return err
    }
    
    return chargePayment(ctx, order)
}

// Bad ❌
func ProcessOrder(order Order) error {
    ctx := context.Background() // Creates disconnected trace
    ctx, span := tracer.Start(ctx, "processOrder")
    defer span.End()
    // ...
}
```

### 2. Defer Span.End()

Always defer `span.End()` immediately after creating a span to ensure it's closed even if panics occur.

```go
ctx, span := tracer.Start(ctx, "operation")
defer span.End() // Called even if function panics
```

### 3. Use Semantic Conventions

Follow OpenTelemetry semantic conventions for attributes:

```go
// HTTP
attribute.String("http.method", "GET")
attribute.Int("http.status_code", 200)

// Database
attribute.String("db.system", "postgresql")
attribute.String("db.statement", query)

// RPC
attribute.String("rpc.service", "UserService")
attribute.String("rpc.method", "GetUser")
```

### 4. Handle Errors Properly

```go
result, err := someOperation(ctx)
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
    logger.Error("Operation failed", slog.Any("error", err))
    return err
}
```

### 5. Graceful Shutdown

```go
func main() {
    provider, err := otel.NewProvider(ctx, config)
    if err != nil {
        log.Fatal(err)
    }

    // Ensure cleanup happens
    defer func() {
        shutdownCtx, cancel := context.WithTimeout(
            context.Background(), 
            5*time.Second,
        )
        defer cancel()
        
        if err := provider.Shutdown(shutdownCtx); err != nil {
            log.Printf("Error shutting down: %v", err)
        }
    }()

    // Application code
}
```

## Testing

### Running Tests

```bash
go test -v ./otel
```

### Running Benchmarks

```bash
go test -bench=. -benchmem ./otel
```

### Example Test

```go
func TestTracing(t *testing.T) {
    ctx := context.Background()
    config := otel.NewConfig("test-service")
    provider, err := otel.NewProvider(ctx, config)
    require.NoError(t, err)
    defer provider.Shutdown(ctx)

    tracer := provider.Tracer("test")
    ctx, span := tracer.Start(ctx, "test-operation")
    defer span.End()

    // Your test assertions
}
```

## Integration with Popular Frameworks

### Chi Router

```go
import "github.com/go-chi/chi/v5"

r := chi.NewRouter()
r.Use(TracingMiddleware(provider))
```

### Gin

```go
import "github.com/gin-gonic/gin"

func GinTracingMiddleware(provider *otel.Provider) gin.HandlerFunc {
    tracer := provider.Tracer("gin")
    return func(c *gin.Context) {
        ctx, span := tracer.Start(c.Request.Context(), c.Request.Method+" "+c.FullPath())
        defer span.End()
        
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}

router := gin.New()
router.Use(GinTracingMiddleware(provider))
```

### gRPC

```go
import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

server := grpc.NewServer(
    grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
    grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
)
```

## Observability Stack Setup

### Jaeger (Tracing)

```bash
docker run -d --name jaeger \
  -p 4317:4317 \
  -p 16686:16686 \
  jaegertracing/all-in-one:latest
```

### Prometheus (Metrics)

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'otel-collector'
    static_configs:
      - targets: ['localhost:9090']
```

### OpenTelemetry Collector

```yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

exporters:
  jaeger:
    endpoint: localhost:14250
  prometheus:
    endpoint: 0.0.0.0:9090

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [jaeger]
    metrics:
      receivers: [otlp]
      exporters: [prometheus]
```

## Performance Considerations

- **Sampling**: Use appropriate sampling strategies in production
- **Batch Processing**: Exporters use batching by default
- **Resource Usage**: Minimal overhead when no endpoints configured
- **Context Propagation**: Efficiently passes trace context

## Troubleshooting

### Traces Not Appearing

1. Verify endpoint is reachable: `telnet localhost 4317`
2. Check if tracing is enabled: `config.EnableTracing = true`
3. Ensure span is ended: `defer span.End()`
4. Verify collector is running

### High Memory Usage

1. Adjust batch size in exporter configuration
2. Implement sampling strategy
3. Check for span leaks (spans not ended)

### Logs Missing Trace Context

1. Ensure context is passed to logger: `provider.Logger(ctx)`
2. Verify span is active when logging
3. Check if logging is enabled

## Contributing

Contributions are welcome! Please ensure:

- All tests pass
- Code is formatted with `gofmt`
- Documentation is updated
- Examples are provided for new features

## Resources

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [OpenTelemetry Go SDK](https://github.com/open-telemetry/opentelemetry-go)
- [Semantic Conventions](https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/)
- [Best Practices](https://opentelemetry.io/docs/instrumentation/go/manual/)