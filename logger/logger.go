package logger

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Logger wraps zerolog with OpenTelemetry span context integration
type Logger struct {
	zerolog.Logger
	traceIDKey string
	spanIDKey  string
	level      zerolog.Level
	mu         sync.RWMutex

	// OTel integration
	serviceName    string
	environment    string
	resourceAttrs  []attribute.KeyValue
}

// Config holds configuration for the logger
type Config struct {
	// Output specifies the output destination (default: os.Stderr)
	Output io.Writer

	// Level specifies the logging level (default: InfoLevel)
	Level zerolog.Level

	// Format specifies the output format: "json", "console", or "pretty"
	// "json": JSON format for production
	// "console": Human-readable console format
	// "pretty": Colorized pretty format
	Format string

	// ServiceName sets the service name in logs
	ServiceName string

	// Environment sets the environment (dev, staging, prod, etc.)
	Environment string

	// TraceIDFieldName customizes the field name for trace ID in logs (default: "trace_id")
	TraceIDFieldName string

	// SpanIDFieldName customizes the field name for span ID in logs (default: "span_id")
	SpanIDFieldName string

	// PrettyPrint enables pretty JSON formatting (indented) - only affects JSON format
	PrettyPrint bool

	// Caller enables caller information in logs
	Caller bool

	// ResourceAttributes adds OpenTelemetry resource attributes to logs
	ResourceAttributes map[string]string
}

// New creates a new logger with default configuration
func New() *Logger {
	return NewWithConfig(Config{})
}

// NewWithConfig creates a new logger with custom configuration
func NewWithConfig(cfg Config) *Logger {
	// Set defaults
	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}
	if cfg.Level == 0 {
		cfg.Level = zerolog.InfoLevel
	}
	if cfg.Format == "" {
		cfg.Format = "json"
	}
	if cfg.TraceIDFieldName == "" {
		cfg.TraceIDFieldName = "trace_id"
	}
	if cfg.SpanIDFieldName == "" {
		cfg.SpanIDFieldName = "span_id"
	}

	// Configure zerolog
	zerolog.SetGlobalLevel(cfg.Level)
	var logger zerolog.Logger

	switch cfg.Format {
	case "console":
		logger = zerolog.New(cfg.Output).
			With().
			Timestamp().
			Logger().
			Output(zerolog.ConsoleWriter{Out: cfg.Output, NoColor: false})
	case "pretty":
		consoleWriter := zerolog.ConsoleWriter{
			Out:        cfg.Output,
			NoColor:    false,
			TimeFormat: time.RFC3339,
		}
		if cfg.PrettyPrint {
			consoleWriter.PartsOrder = []string{
				zerolog.TimestampFieldName,
				zerolog.LevelFieldName,
				zerolog.CallerFieldName,
				zerolog.MessageFieldName,
			}
		}
		logger = zerolog.New(cfg.Output).
			With().
			Timestamp().
			Logger().
			Output(consoleWriter)
	default: // json
		logger = zerolog.New(cfg.Output).
			With().
			Timestamp().
			Logger()
	}

	// Add context fields
	builder := logger.With()
	if cfg.ServiceName != "" {
		builder = builder.Str("service", cfg.ServiceName)
	}
	if cfg.Environment != "" {
		builder = builder.Str("env", cfg.Environment)
	}
	if cfg.Caller {
		logger = logger.With().Caller().Logger()
	}
	logger = builder.Logger()

	// Build resource attributes
	var resourceAttrs []attribute.KeyValue
	for k, v := range cfg.ResourceAttributes {
		resourceAttrs = append(resourceAttrs, attribute.String(k, v))
	}

	return &Logger{
		Logger:        logger,
		traceIDKey:    cfg.TraceIDFieldName,
		spanIDKey:     cfg.SpanIDFieldName,
		level:         cfg.Level,
		serviceName:   cfg.ServiceName,
		environment:   cfg.Environment,
		resourceAttrs: resourceAttrs,
	}
}

// WithContext adds fields from OpenTelemetry span context to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()

	// Create a child logger with trace context
	builder := l.Logger.With()

	if spanCtx.IsValid() {
		if spanCtx.HasTraceID() {
			builder = builder.Str(l.traceIDKey, spanCtx.TraceID().String())
		}
		if spanCtx.HasSpanID() {
			builder = builder.Str(l.spanIDKey, spanCtx.SpanID().String())
		}
	}

	return &Logger{
		Logger:        builder.Logger(),
		traceIDKey:    l.traceIDKey,
		spanIDKey:     l.spanIDKey,
		level:         l.level,
		serviceName:   l.serviceName,
		environment:   l.environment,
		resourceAttrs: l.resourceAttrs,
	}
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	builder := l.Logger.With()
	for k, v := range fields {
		builder = builder.Interface(k, v)
	}

	return &Logger{
		Logger:        builder.Logger(),
		traceIDKey:    l.traceIDKey,
		spanIDKey:     l.spanIDKey,
		level:         l.level,
		serviceName:   l.serviceName,
		environment:   l.environment,
		resourceAttrs: l.resourceAttrs,
	}
}

// WithError creates a new logger with an error field
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger:        l.Logger.With().Err(err).Logger(),
		traceIDKey:    l.traceIDKey,
		spanIDKey:     l.spanIDKey,
		level:         l.level,
		serviceName:   l.serviceName,
		environment:   l.environment,
		resourceAttrs: l.resourceAttrs,
	}
}

// With creates a zerolog event builder
func (l *Logger) With() zerolog.Context {
	return l.Logger.With()
}

// Info creates an info level log event
func (l *Logger) Info() *zerolog.Event {
	return l.Logger.Info()
}

// Debug creates a debug level log event
func (l *Logger) Debug() *zerolog.Event {
	return l.Logger.Debug()
}

// Error creates an error level log event
func (l *Logger) Error() *zerolog.Event {
	return l.Logger.Error()
}

// Warn creates a warn level log event
func (l *Logger) Warn() *zerolog.Event {
	return l.Logger.Warn()
}

// Fatal creates a fatal level log event (exits the program)
func (l *Logger) Fatal() *zerolog.Event {
	return l.Logger.Fatal()
}

// Panic creates a panic level log event (panics)
func (l *Logger) Panic() *zerolog.Event {
	return l.Logger.Panic()
}

// Trace creates a trace level log event
func (l *Logger) Trace() *zerolog.Event {
	return l.Logger.Trace()
}

// GetLevel returns the current logging level
func (l *Logger) GetLevel() zerolog.Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// SetLevel updates the logging level
func (l *Logger) SetLevel(level zerolog.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
	zerolog.SetGlobalLevel(level)
}

// ServiceName returns the service name
func (l *Logger) ServiceName() string {
	return l.serviceName
}

// Environment returns the environment
func (l *Logger) Environment() string {
	return l.environment
}

// ResourceAttributes returns the resource attributes
func (l *Logger) ResourceAttributes() []attribute.KeyValue {
	return l.resourceAttrs
}

// Global logger instance
var (
	globalLogger *Logger
	globalOnce   sync.Once
	globalMu     sync.RWMutex
)

// GetGlobal returns the global logger instance (singleton)
func GetGlobal() *Logger {
	globalOnce.Do(func() {
		globalLogger = New()
	})
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalLogger
}

// SetGlobal sets the global logger instance
func SetGlobal(logger *Logger) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalLogger = logger
}

// Helper functions for global logger access
var (
	Info  = func() *zerolog.Event { return GetGlobal().Info() }
	Debug = func() *zerolog.Event { return GetGlobal().Debug() }
	Error = func() *zerolog.Event { return GetGlobal().Error() }
	Warn  = func() *zerolog.Event { return GetGlobal().Warn() }
	Fatal = func() *zerolog.Event { return GetGlobal().Fatal() }
	Panic = func() *zerolog.Event { return GetGlobal().Panic() }
	Trace = func() *zerolog.Event { return GetGlobal().Trace() }
)

// WithContext returns a logger with context from the global logger
func WithContext(ctx context.Context) *Logger {
	return GetGlobal().WithContext(ctx)
}

// WithFields returns a logger with additional fields from the global logger
func WithFields(fields map[string]interface{}) *Logger {
	return GetGlobal().WithFields(fields)
}

// WithError returns a logger with an error field from the global logger
func WithError(err error) *Logger {
	return GetGlobal().WithError(err)
}

// SetLevel sets the level for the global logger
func SetLevel(level zerolog.Level) {
	GetGlobal().SetLevel(level)
}

// SetGlobalLevel is an alias for SetLevel
func SetGlobalLevel(level zerolog.Level) {
	SetLevel(level)
}

// Format constant values
const (
	FormatJSON    = "json"
	FormatConsole = "console"
	FormatPretty  = "pretty"
)
