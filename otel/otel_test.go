package otel

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig("test-service")

	if config.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", config.ServiceName)
	}

	if config.ServiceVersion != "1.0.0" {
		t.Errorf("Expected ServiceVersion '1.0.0', got '%s'", config.ServiceVersion)
	}

	if !config.EnableTracing {
		t.Error("Expected EnableTracing to be true by default")
	}

	if !config.EnableMetrics {
		t.Error("Expected EnableMetrics to be true by default")
	}

	if !config.EnableLogging {
		t.Error("Expected EnableLogging to be true by default")
	}
}

func TestConfigChaining(t *testing.T) {
	config := NewConfig("test-service").
		WithServiceVersion("2.0.0").
		WithEnvironment("staging").
		WithTracing(false).
		WithMetrics(true).
		WithLogging(false).
		WithResourceAttribute("team", "platform")

	if config.ServiceVersion != "2.0.0" {
		t.Errorf("Expected ServiceVersion '2.0.0', got '%s'", config.ServiceVersion)
	}

	if config.Environment != "staging" {
		t.Errorf("Expected Environment 'staging', got '%s'", config.Environment)
	}

	if config.EnableTracing {
		t.Error("Expected EnableTracing to be false")
	}

	if !config.EnableMetrics {
		t.Error("Expected EnableMetrics to be true")
	}

	if config.EnableLogging {
		t.Error("Expected EnableLogging to be false")
	}

	if config.ResourceAttributes["team"] != "platform" {
		t.Errorf("Expected ResourceAttribute 'team' to be 'platform', got '%s'", config.ResourceAttributes["team"])
	}
}

func TestProviderWithoutEndpoints(t *testing.T) {
	ctx := context.Background()

	config := NewConfig("test-service").
		WithTracing(true).
		WithMetrics(true).
		WithLogging(true)

	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}
	defer provider.Shutdown(ctx)

	// Provider should be created successfully even without endpoints
	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	// Tracer should still be available (no-op)
	tracer := provider.Tracer("test")
	if tracer == nil {
		t.Error("Expected tracer to be available")
	}
}

func TestTracingOperations(t *testing.T) {
	ctx := context.Background()

	config := NewConfig("test-service")
	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}
	defer provider.Shutdown(ctx)

	tracer := provider.Tracer("test-component")

	// Test span creation
	ctx, span := tracer.Start(ctx, "test-operation")
	defer span.End()

	// Add attributes
	span.SetAttributes(
		attribute.String("test.key", "test-value"),
		attribute.Int("test.count", 42),
	)

	// Add event
	span.AddEvent("test-event", trace.WithAttributes(
		attribute.String("event.type", "test"),
	))

	// Nested span
	_, childSpan := tracer.Start(ctx, "child-operation")
	time.Sleep(10 * time.Millisecond)
	childSpan.End()
}

func TestProviderShutdown(t *testing.T) {
	ctx := context.Background()

	config := NewConfig("test-service")
	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Test shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = provider.Shutdown(shutdownCtx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}

func TestResourceAttributes(t *testing.T) {
	ctx := context.Background()

	config := NewConfig("test-service").
		WithResourceAttribute("custom.key1", "value1").
		WithResourceAttribute("custom.key2", "value2")

	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}
	defer provider.Shutdown(ctx)

	if len(config.ResourceAttributes) != 2 {
		t.Errorf("Expected 2 resource attributes, got %d", len(config.ResourceAttributes))
	}
}

func TestGetProviders(t *testing.T) {
	ctx := context.Background()

	config := NewConfig("test-service")
	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}
	defer provider.Shutdown(ctx)

	// Without endpoints, providers should be nil
	if provider.GetTracerProvider() != nil {
		t.Log("TracerProvider is available (has endpoint configured)")
	}

	if provider.GetMeterProvider() != nil {
		t.Log("MeterProvider is available (has endpoint configured)")
	}
}

func TestProviderMetadata(t *testing.T) {
	ctx := context.Background()

	config := NewConfig("test-service").
		WithServiceVersion("2.0.0").
		WithEnvironment("production")

	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}
	defer provider.Shutdown(ctx)

	if provider.ServiceName() != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", provider.ServiceName())
	}

	if provider.ServiceVersion() != "2.0.0" {
		t.Errorf("Expected ServiceVersion '2.0.0', got '%s'", provider.ServiceVersion())
	}

	if provider.Environment() != "production" {
		t.Errorf("Expected Environment 'production', got '%s'", provider.Environment())
	}
}

func TestLoggerConfig(t *testing.T) {
	ctx := context.Background()

	config := NewConfig("test-service").
		WithServiceVersion("1.5.0").
		WithEnvironment("staging").
		WithResourceAttribute("team", "backend")

	provider, err := NewProvider(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}
	defer provider.Shutdown(ctx)

	loggerCfg := provider.LoggerConfig()
	if loggerCfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", loggerCfg.ServiceName)
	}

	if loggerCfg.Environment != "staging" {
		t.Errorf("Expected Environment 'staging', got '%s'", loggerCfg.Environment)
	}
}

func TestLoggerOutput(t *testing.T) {
	config := NewConfig("test-service").
		WithLoggerOutput(LoggerOutputStdout)

	if config.LoggerOutput != LoggerOutputStdout {
		t.Errorf("Expected LoggerOutputStdout, got '%s'", config.LoggerOutput)
	}
}

// Benchmark tests
func BenchmarkSpanCreation(b *testing.B) {
	ctx := context.Background()
	config := NewConfig("bench-service")
	provider, _ := NewProvider(ctx, config)
	defer provider.Shutdown(ctx)

	tracer := provider.Tracer("bench")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, span := tracer.Start(ctx, "operation")
		span.End()
	}
}
