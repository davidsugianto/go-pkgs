package main

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/davidsugianto/go-pkgs/logger"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	// Example 1: Basic logger configuration
	basicExample()

	// Example 2: Logger with OpenTelemetry integration
	otelExample()

	// Example 3: Global logger usage
	globalExample()

	// Example 4: Different output formats
	formatExample()
}

func basicExample() {
	printTitle("Example 1: Basic Logger Usage")

	// Create a logger with custom configuration
	log := logger.NewWithConfig(logger.Config{
		ServiceName: "my-service",
		Environment: "development",
		Format:      logger.FormatJSON,
		Level:       zerolog.DebugLevel,
	})

	// Log at different levels
	log.Debug().Msg("Debug message - only visible in debug mode")
	log.Info().Msg("Application started")
	log.Warn().Str("warning_type", "deprecated_api").Msg("Using deprecated API")
	log.Error().Err(errors.New("failed to connect")).Msg("Connection error")

	// Log with structured fields
	log.Info().
		Str("user_id", "12345").
		Str("action", "login").
		Dur("duration_ms", 150*time.Millisecond).
		Msg("User logged in")

	// Log with multiple field types
	log.Debug().
		Str("string_field", "value").
		Int("int_field", 42).
		Float64("float_field", 3.14159).
		Bool("bool_field", true).
		Msg("Complex structured log")
}

func otelExample() {
	printTitle("Example 2: OpenTelemetry Integration")

	// Create a logger
	log := logger.NewWithConfig(logger.Config{
		ServiceName:      "api-server",
		Environment:      "production",
		Format:           logger.FormatJSON,
		TraceIDFieldName: "trace_id",
		SpanIDFieldName:  "span_id",
	})

	// Create a tracer (normally you'd use a real OpenTelemetry tracer)
	tracer := otel.Tracer("example-service")

	// Simulate an HTTP request handler
	handler := func(ctx context.Context, method, path string) {
		// Start a span for this operation
		ctx, span := tracer.Start(ctx, "http.request",
			trace.WithAttributes(
				attribute.String("http.method", method),
				attribute.String("http.path", path),
			),
		)
		defer span.End()

		// Create a logger with context to automatically add trace/span IDs
		requestLogger := log.WithContext(ctx)

		// Log request started
		requestLogger.Info().
			Str("method", method).
			Str("path", path).
			Msg("Request started")

		// Simulate processing
		time.Sleep(10 * time.Millisecond)

		// Log with span context automatically included
		requestLogger.Info().
			Int("status_code", 200).
			Dur("duration_ms", 10*time.Millisecond).
			Msg("Request completed")

		// Nested operation example
		dbOperation(ctx, log)
	}

	// Simulate multiple requests
	ctx := context.Background()
	handler(ctx, "GET", "/users")
	handler(ctx, "POST", "/orders")
}

func dbOperation(ctx context.Context, parentLog *logger.Logger) {
	tracer := otel.Tracer("example-service")

	// Create a child span
	ctx, span := tracer.Start(ctx, "db.query")
	defer span.End()

	// Log with new span context
	dbLog := parentLog.WithContext(ctx)
	dbLog.Info().
		Str("query", "SELECT * FROM users").
		Int("rows_returned", 42).
		Msg("Database query executed")
}

func globalExample() {
	printTitle("Example 3: Global Logger Usage")

	// Set up global logger
	logger.SetGlobal(logger.NewWithConfig(logger.Config{
		ServiceName: "global-service",
		Format:      logger.FormatConsole,
		Level:       zerolog.InfoLevel,
	}))

	// Use global helper functions
	logger.Info().Msg("Global logger info message")
	logger.Warn().Msg("Global logger warning message")
	logger.Error().Err(errors.New("global error")).Msg("Global error message")

	// Change log level at runtime
	logger.SetLevel(zerolog.DebugLevel)
	logger.Debug().Msg("This debug message is now visible")
	logger.SetLevel(zerolog.InfoLevel)

	// Use WithContext with global logger
	tracer := otel.Tracer("example-service")
	ctx, span := tracer.Start(context.Background(), "global.operation")
	defer span.End()

	logger.WithContext(ctx).Info().Msg("Global logger with tracing")
}

func formatExample() {
	printTitle("Example 4: Different Output Formats")

	// JSON format (for production/log aggregation)
	jsonLog := logger.NewWithConfig(logger.Config{
		Output: os.Stderr,
		Format: logger.FormatJSON,
	})
	jsonLog.Info().
		Str("service", "api").
		Int("port", 8080).
		Msg("JSON format log")

	// Console format (for development)
	consoleLog := logger.NewWithConfig(logger.Config{
		Output: os.Stderr,
		Format: logger.FormatConsole,
	})
	consoleLog.Info().
		Str("service", "api").
		Int("port", 8080).
		Msg("Console format log")

	// Pretty format (for development with colors)
	prettyLog := logger.NewWithConfig(logger.Config{
		Output: os.Stderr,
		Format: logger.FormatPretty,
	})
	prettyLog.Info().
		Str("service", "api").
		Int("port", 8080).
		Msg("Pretty format log")
}

func printTitle(title string) {
	println("\n" + strings.Repeat("=", 60))
	println(title)
	println(strings.Repeat("=", 60))
	time.Sleep(500 * time.Millisecond)
}
