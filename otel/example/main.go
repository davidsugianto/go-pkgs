package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/davidsugianto/go-pkgs/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry
	config := otel.NewConfig("my-service").
		WithServiceVersion("1.2.3").
		WithEnvironment("production").
		WithTraceEndpoint("localhost:4317").
		WithMetricEndpoint("localhost:4317").
		WithLogLevel(slog.LevelDebug).
		WithResourceAttribute("team", "backend")

	// Initialize provider
	provider, err := otel.NewProvider(ctx, config)
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down provider: %v", err)
		}
	}()

	// Example 1: Using Tracing
	tracer := provider.Tracer("my-component")
	ctx, span := tracer.Start(ctx, "operation")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", "123"),
		attribute.Int("items.count", 5),
	)

	// Example 2: Using Logging with trace context
	logger := provider.Logger(ctx)
	logger.Info("Processing request",
		slog.String("user", "john"),
		slog.Int("items", 5),
	)

	// Example 3: Using Metrics
	meter := provider.GetMeterProvider().Meter("my-component")

	counter, err := meter.Int64Counter(
		"requests.total",
		metric.WithDescription("Total number of requests"),
	)
	if err != nil {
		log.Printf("Failed to create counter: %v", err)
	}

	counter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("method", "GET"),
		attribute.String("status", "200"),
	))

	// Simulate work
	time.Sleep(100 * time.Millisecond)

	logger.Info("Request completed successfully")
}
