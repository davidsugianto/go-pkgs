package main

import (
	"context"
	"log"
	"time"

	"github.com/davidsugianto/go-pkgs/logger"
	"github.com/davidsugianto/go-pkgs/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry
	otelConfig := otel.NewConfig("my-service").
		WithServiceVersion("1.2.3").
		WithEnvironment("production").
		WithTraceEndpoint("localhost:4317").
		WithMetricEndpoint("localhost:4317").
		WithResourceAttribute("team", "backend")

	// Initialize OTel provider
	provider, err := otel.NewProvider(ctx, otelConfig)
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down provider: %v", err)
		}
	}()

	// Create logger with OTel integration
	loggerConfig := logger.NewWithConfig(logger.Config{
		ServiceName:       provider.ServiceName(),
		Environment:       provider.Environment(),
		Format:            logger.FormatJSON,
		ResourceAttributes: provider.ResourceAttributes(),
	})
	logger.SetGlobal(loggerConfig)

	// Example 1: Using Tracing
	tracer := provider.Tracer("my-component")
	ctx, span := tracer.Start(ctx, "operation")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", "123"),
		attribute.Int("items.count", 5),
	)

	// Example 2: Using Logging with trace context
	log := logger.WithContext(ctx)
	log.Info().
		Str("user", "john").
		Int("items", 5).
		Msg("Processing request")

	// Example 3: Using Metrics
	meter := provider.GetMeterProvider()
	if meter != nil {
		m := meter.Meter("my-component")
		counter, err := m.Int64Counter(
			"requests.total",
			metric.WithDescription("Total number of requests"),
		)
		if err != nil {
			log.Printf("Failed to create counter: %v", err)
		} else {
			counter.Add(ctx, 1, metric.WithAttributes(
				attribute.String("method", "GET"),
				attribute.String("status", "200"),
			))
		}
	}

	// Simulate work
	time.Sleep(100 * time.Millisecond)

	log.Info().Msg("Request completed successfully")
}
