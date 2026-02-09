package otel

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// Config holds the configuration for OpenTelemetry setup.
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string

	// Trace configuration
	EnableTracing bool
	TraceEndpoint string

	// Metrics configuration
	EnableMetrics  bool
	MetricEndpoint string

	// Logging configuration
	EnableLogging bool
	LogLevel      slog.Level

	// Additional resource attributes
	ResourceAttributes map[string]string
}

// Provider manages OpenTelemetry providers and their lifecycle.
type Provider struct {
	config         *Config
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	logger         *slog.Logger
	shutdownFuncs  []func(context.Context) error
}

// NewConfig creates a new Config with defaults.
func NewConfig(serviceName string) *Config {
	return &Config{
		ServiceName:        serviceName,
		ServiceVersion:     "1.0.0",
		Environment:        getEnvOrDefault("ENVIRONMENT", "development"),
		EnableTracing:      true,
		TraceEndpoint:      getEnvOrDefault("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", ""),
		EnableMetrics:      true,
		MetricEndpoint:     getEnvOrDefault("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", ""),
		EnableLogging:      true,
		LogLevel:           slog.LevelInfo,
		ResourceAttributes: make(map[string]string),
	}
}

// WithServiceVersion sets the service version.
func (c *Config) WithServiceVersion(version string) *Config {
	c.ServiceVersion = version
	return c
}

// WithEnvironment sets the environment.
func (c *Config) WithEnvironment(env string) *Config {
	c.Environment = env
	return c
}

// WithTracing enables or disables tracing.
func (c *Config) WithTracing(enabled bool) *Config {
	c.EnableTracing = enabled
	return c
}

// WithTraceEndpoint sets the trace endpoint.
func (c *Config) WithTraceEndpoint(endpoint string) *Config {
	c.TraceEndpoint = endpoint
	return c
}

// WithMetrics enables or disables metrics.
func (c *Config) WithMetrics(enabled bool) *Config {
	c.EnableMetrics = enabled
	return c
}

// WithMetricEndpoint sets the metric endpoint.
func (c *Config) WithMetricEndpoint(endpoint string) *Config {
	c.MetricEndpoint = endpoint
	return c
}

// WithLogging enables or disables logging.
func (c *Config) WithLogging(enabled bool) *Config {
	c.EnableLogging = enabled
	return c
}

// WithLogLevel sets the log level.
func (c *Config) WithLogLevel(level slog.Level) *Config {
	c.LogLevel = level
	return c
}

// WithResourceAttribute adds a resource attribute.
func (c *Config) WithResourceAttribute(key, value string) *Config {
	c.ResourceAttributes[key] = value
	return c
}

// NewProvider initializes OpenTelemetry providers based on the configuration.
func NewProvider(ctx context.Context, config *Config) (*Provider, error) {
	p := &Provider{
		config:        config,
		shutdownFuncs: []func(context.Context) error{},
	}

	// Create resource
	res, err := p.createResource(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Initialize tracing
	if config.EnableTracing && config.TraceEndpoint != "" {
		if err := p.initTracing(ctx, res); err != nil {
			return nil, fmt.Errorf("failed to initialize tracing: %w", err)
		}
	}

	// Initialize metrics
	if config.EnableMetrics && config.MetricEndpoint != "" {
		if err := p.initMetrics(ctx, res); err != nil {
			return nil, fmt.Errorf("failed to initialize metrics: %w", err)
		}
	}

	// Initialize logging
	if config.EnableLogging {
		p.initLogging()
	}

	return p, nil
}

// createResource creates an OpenTelemetry resource with service information.
func (p *Provider) createResource(ctx context.Context) (*resource.Resource, error) {
	attrs := []resource.Option{
		resource.WithAttributes(
			semconv.ServiceNameKey.String(p.config.ServiceName),
			semconv.ServiceVersionKey.String(p.config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(p.config.Environment),
		),
	}

	// Add custom resource attributes
	for key, value := range p.config.ResourceAttributes {
		attrs = append(attrs, resource.WithAttributes(
			attribute.String(key, value),
		))
	}

	return resource.New(ctx, attrs...)
}

// initTracing initializes the trace provider.
func (p *Provider) initTracing(ctx context.Context, res *resource.Resource) error {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(p.config.TraceEndpoint),
		otlptracegrpc.WithInsecure(), // Use WithTLSCredentials in production
	)
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	p.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(p.tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	p.shutdownFuncs = append(p.shutdownFuncs, p.tracerProvider.Shutdown)
	return nil
}

// initMetrics initializes the meter provider.
func (p *Provider) initMetrics(ctx context.Context, res *resource.Resource) error {
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(p.config.MetricEndpoint),
		otlpmetricgrpc.WithInsecure(), // Use WithTLSCredentials in production
	)
	if err != nil {
		return fmt.Errorf("failed to create metric exporter: %w", err)
	}

	p.meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(p.meterProvider)

	p.shutdownFuncs = append(p.shutdownFuncs, p.meterProvider.Shutdown)
	return nil
}

// initLogging initializes structured logging with OpenTelemetry integration.
func (p *Provider) initLogging() {
	opts := &slog.HandlerOptions{
		Level: p.config.LogLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	p.logger = slog.New(handler)
	slog.SetDefault(p.logger)
}

// Tracer returns a tracer for the given name.
func (p *Provider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	if p.tracerProvider == nil {
		return otel.Tracer(name, opts...)
	}
	return p.tracerProvider.Tracer(name, opts...)
}

// Logger returns the configured logger with trace context.
func (p *Provider) Logger(ctx context.Context) *slog.Logger {
	if p.logger == nil {
		return slog.Default()
	}

	logger := p.logger

	// Add trace context if available
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		logger = logger.With(
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	}

	return logger
}

// Shutdown gracefully shuts down all providers.
func (p *Provider) Shutdown(ctx context.Context) error {
	var errs []error

	for _, shutdown := range p.shutdownFuncs {
		if err := shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}

// GetTracerProvider returns the tracer provider (or nil if not initialized).
func (p *Provider) GetTracerProvider() *sdktrace.TracerProvider {
	return p.tracerProvider
}

// GetMeterProvider returns the meter provider (or nil if not initialized).
func (p *Provider) GetMeterProvider() *sdkmetric.MeterProvider {
	return p.meterProvider
}

// getEnvOrDefault returns environment variable value or default.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
