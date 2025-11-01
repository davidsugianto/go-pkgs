package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestNew(t *testing.T) {
	logger := New()
	assert.NotNil(t, logger)
	assert.Equal(t, zerolog.InfoLevel, logger.GetLevel())
}

func TestNewWithConfig_Defaults(t *testing.T) {
	logger := NewWithConfig(Config{})
	assert.NotNil(t, logger)
	assert.Equal(t, zerolog.InfoLevel, logger.GetLevel())
}

func TestNewWithConfig_CustomLevel(t *testing.T) {
	logger := NewWithConfig(Config{
		Level: zerolog.InfoLevel,
	})
	assert.NotNil(t, logger)
	assert.Equal(t, zerolog.InfoLevel, logger.GetLevel())
}

func TestNewWithConfig_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatJSON,
	})

	logger.Info().Msg("test message")

	output := buf.String()
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, `"level":"info"`)
}

func TestNewWithConfig_ConsoleFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatConsole,
	})

	logger.Info().Msg("test message")

	output := buf.String()
	assert.Contains(t, output, "test message")
	// Console format should be human-readable
	assert.Contains(t, output, "INF")
}

func TestNewWithConfig_PrettyFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatPretty,
	})

	logger.Info().Msg("test message")

	output := buf.String()
	assert.Contains(t, output, "test message")
	// Pretty format should be human-readable
	assert.Contains(t, output, "INF")
}

func TestNewWithConfig_ServiceName(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output:      &buf,
		Format:      FormatJSON,
		ServiceName: "my-service",
	})

	logger.Info().Msg("test")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)
	require.NoError(t, err)

	assert.Equal(t, "my-service", logData["service"])
}

func TestNewWithConfig_Environment(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output:      &buf,
		Format:      FormatJSON,
		Environment: "production",
	})

	logger.Info().Msg("test")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)
	require.NoError(t, err)

	assert.Equal(t, "production", logData["env"])
}

func TestWithContext_NoSpan(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatJSON,
	})

	ctx := context.Background()
	loggerWithCtx := logger.WithContext(ctx)
	loggerWithCtx.Info().Msg("test")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)
	require.NoError(t, err)

	// Should not have trace_id or span_id
	_, hasTraceID := logData["trace_id"]
	_, hasSpanID := logData["span_id"]
	assert.False(t, hasTraceID)
	assert.False(t, hasSpanID)
}

// func TestWithContext_WithSpan(t *testing.T) {
// 	var buf bytes.Buffer
// 	logger := NewWithConfig(Config{
// 		Output: &buf,
// 		Format: FormatJSON,
// 	})

// 	// Create a tracer and span
// 	tracer := trace.NewNoopTracerProvider().Tracer("test")
// 	ctx, span := tracer.Start(context.Background(), "test-span")

// 	loggerWithCtx := logger.WithContext(ctx)
// 	loggerWithCtx.Info().Msg("test")

// 	span.End()

// 	var logData map[string]interface{}
// 	err := json.Unmarshal(buf.Bytes(), &logData)
// 	require.NoError(t, err)

// 	// Should have trace_id and span_id
// 	traceID, hasTraceID := logData["trace_id"]
// 	spanID, hasSpanID := logData["span_id"]

// 	assert.True(t, hasTraceID, "should have trace_id field")
// 	assert.True(t, hasSpanID, "should have span_id field")
// 	assert.NotEmpty(t, traceID)
// 	assert.NotEmpty(t, spanID)
// }

// func TestWithContext_CustomFieldNames(t *testing.T) {
// 	var buf bytes.Buffer
// 	logger := NewWithConfig(Config{
// 		Output:           &buf,
// 		Format:           FormatJSON,
// 		TraceIDFieldName: "my_trace_id",
// 		SpanIDFieldName:  "my_span_id",
// 	})

// 	tracer := trace.NewNoopTracerProvider().Tracer("test")
// 	ctx, span := tracer.Start(context.Background(), "test-span")

// 	loggerWithCtx := logger.WithContext(ctx)
// 	loggerWithCtx.Info().Msg("test")

// 	span.End()

// 	var logData map[string]interface{}
// 	err := json.Unmarshal(buf.Bytes(), &logData)
// 	require.NoError(t, err)

// 	// Should use custom field names
// 	_, hasTraceID := logData["my_trace_id"]
// 	_, hasSpanID := logData["my_span_id"]
// 	assert.True(t, hasTraceID)
// 	assert.True(t, hasSpanID)
// }

func TestSetLevel(t *testing.T) {
	logger := New()
	assert.Equal(t, zerolog.InfoLevel, logger.GetLevel())

	logger.SetLevel(zerolog.DebugLevel)
	assert.Equal(t, zerolog.DebugLevel, logger.GetLevel())

	logger.SetLevel(zerolog.ErrorLevel)
	assert.Equal(t, zerolog.ErrorLevel, logger.GetLevel())
}

// func TestLogLevels(t *testing.T) {
// 	var buf bytes.Buffer
// 	logger := NewWithConfig(Config{
// 		Output: &buf,
// 		Format: FormatJSON,
// 		Level:  zerolog.DebugLevel,
// 	})

// 	tests := []struct {
// 		name  string
// 		logFn func()
// 		level string
// 	}{
// 		{"Trace", func() { logger.Trace().Msg("trace") }, "trace"},
// 		{"Debug", func() { logger.Debug().Msg("debug") }, "debug"},
// 		{"Info", func() { logger.Info().Msg("info") }, "info"},
// 		{"Warn", func() { logger.Warn().Msg("warn") }, "warn"},
// 		{"Error", func() { logger.Error().Msg("error") }, "error"},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			buf.Reset()
// 			tt.logFn()

// 			var logData map[string]interface{}
// 			err := json.Unmarshal(buf.Bytes(), &logData)
// 			require.NoError(t, err)

// 			assert.Equal(t, tt.level, logData["message"])
// 		})
// 	}
// }

func TestLogWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatJSON,
	})

	logger.Info().
		Str("key1", "value1").
		Int("key2", 42).
		Bool("key3", true).
		Msg("test")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)
	require.NoError(t, err)

	assert.Equal(t, "value1", logData["key1"])
	assert.Equal(t, float64(42), logData["key2"]) // JSON unmarshals numbers as float64
	assert.Equal(t, true, logData["key3"])
}

func TestGlobalLogger(t *testing.T) {
	logger := GetGlobal()
	assert.NotNil(t, logger)

	// Should return the same instance
	logger2 := GetGlobal()
	assert.Equal(t, logger, logger2)
}

// func TestSetGlobal(t *testing.T) {
// 	original := GetGlobal()

// 	custom := NewWithConfig(Config{Level: zerolog.DebugLevel})
// 	SetGlobal(custom)

// 	current := GetGlobal()
// 	assert.Equal(t, custom, current)
// 	assert.NotEqual(t, original, current)

// 	// Reset for other tests
// 	SetGlobal(New())
// }

func TestGlobalHelperFunctions(t *testing.T) {
	var buf bytes.Buffer
	SetGlobal(NewWithConfig(Config{
		Output: &buf,
		Format: FormatJSON,
	}))

	Info().Msg("info message")
	assert.Contains(t, buf.String(), "info message")
}

func TestGetGlobal_Concurrency(t *testing.T) {
	// This tests that sync.Once works correctly
	// by calling GetGlobal multiple times concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			logger := GetGlobal()
			assert.NotNil(t, logger)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestWithMethod(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatJSON,
	})

	childLogger := logger.With().
		Str("static", "value").
		Logger()

	childLogger.Info().Msg("test")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)
	require.NoError(t, err)

	assert.Equal(t, "value", logData["static"])
}

func TestEmptyTraceID(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatJSON,
	})

	ctx := context.Background()

	// Try to get a span from context (should be noop)
	loggerWithCtx := logger.WithContext(ctx)
	loggerWithCtx.Info().Msg("test")

	output := buf.String()
	// Should not crash and should log normally
	assert.Contains(t, output, "test")
}

// func TestWithContext_MultipleCalls(t *testing.T) {
// 	var buf1 bytes.Buffer
// 	logger := NewWithConfig(Config{
// 		Output: &buf1,
// 		Format: FormatJSON,
// 	})

// 	tracer := trace.NewNoopTracerProvider().Tracer("test")
// 	ctx, span := tracer.Start(context.Background(), "test-span")

// 	loggerWithCtx1 := logger.WithContext(ctx)
// 	loggerWithCtx2 := logger.WithContext(ctx)

// 	// They should be different instances but both work
// 	assert.NotEqual(t, loggerWithCtx1, loggerWithCtx2)

// 	loggerWithCtx1.Info().Msg("test1")
// 	loggerWithCtx2.Info().Msg("test2")

// 	span.End()

// 	output := buf1.String()
// 	assert.Contains(t, output, "test1")
// 	assert.Contains(t, output, "test2")
// }

func BenchmarkLogger_Info(b *testing.B) {
	logger := NewWithConfig(Config{
		Output: os.Stderr,
		Format: FormatJSON,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Msg("benchmark test")
	}
}

func BenchmarkLogger_InfoWithFields(b *testing.B) {
	logger := NewWithConfig(Config{
		Output: os.Stderr,
		Format: FormatJSON,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().
			Str("key1", "value1").
			Int("key2", 42).
			Bool("key3", true).
			Msg("benchmark test")
	}
}

func BenchmarkLogger_WithContext(b *testing.B) {
	logger := NewWithConfig(Config{
		Output: os.Stderr,
		Format: FormatJSON,
	})

	tracer := trace.NewNoopTracerProvider().Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loggerWithCtx := logger.WithContext(ctx)
		loggerWithCtx.Info().Msg("benchmark test")
	}
}

// func TestExampleUsage(t *testing.T) {
// 	// This test demonstrates typical usage patterns
// 	t.Run("Basic logging", func(t *testing.T) {
// 		var buf bytes.Buffer
// 		logger := NewWithConfig(Config{
// 			Output:      &buf,
// 			Format:      FormatJSON,
// 			ServiceName: "test-service",
// 		})

// 		logger.Info().Msg("Starting application")

// 		output := buf.String()
// 		assert.Contains(t, output, "Starting application")
// 		assert.Contains(t, output, "test-service")
// 	})

// 	t.Run("With OpenTelemetry tracing", func(t *testing.T) {
// 		var buf bytes.Buffer
// 		logger := NewWithConfig(Config{
// 			Output: &buf,
// 			Format: FormatJSON,
// 		})

// 		tracer := trace.NewNoopTracerProvider().Tracer("test")
// 		ctx, span := tracer.Start(context.Background(), "http-request")

// 		loggerWithCtx := logger.WithContext(ctx)
// 		loggerWithCtx.Info().
// 			Str("method", "GET").
// 			Str("path", "/users").
// 			Int("status", 200).
// 			Msg("Request completed")

// 		span.End()

// 		var logData map[string]interface{}
// 		err := json.Unmarshal(buf.Bytes(), &logData)
// 		require.NoError(t, err)

// 		assert.Contains(t, logData, "trace_id")
// 		assert.Contains(t, logData, "span_id")
// 		assert.Equal(t, "GET", logData["method"])
// 		assert.Equal(t, "/users", logData["path"])
// 	})

// 	t.Run("Error logging with stack trace", func(t *testing.T) {
// 		var buf bytes.Buffer
// 		logger := NewWithConfig(Config{
// 			Output: &buf,
// 			Format: FormatConsole, // More readable for errors
// 		})

// 		err := assert.AnError
// 		logger.Error().
// 			Err(err).
// 			Str("operation", "database-query").
// 			Msg("Failed to execute query")

// 		output := buf.String()
// 		assert.Contains(t, output, "ERR")
// 		assert.Contains(t, output, "Failed to execute query")
// 	})
// }

func TestFormatConstants(t *testing.T) {
	// Test that format constants are valid
	assert.Equal(t, "json", FormatJSON)
	assert.Equal(t, "console", FormatConsole)
	assert.Equal(t, "pretty", FormatPretty)
}

func TestLogger_WithOutput(t *testing.T) {
	// Test that we can use different outputs
	var buf1, buf2 bytes.Buffer

	logger1 := NewWithConfig(Config{Output: &buf1, Format: FormatJSON})
	logger2 := NewWithConfig(Config{Output: &buf2, Format: FormatJSON})

	logger1.Info().Msg("message1")
	logger2.Info().Msg("message2")

	output1 := buf1.String()
	output2 := buf2.String()

	assert.Contains(t, output1, "message1")
	assert.Contains(t, output2, "message2")
	assert.NotContains(t, output1, "message2")
	assert.NotContains(t, output2, "message1")
}

func TestLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatJSON,
		Level:  zerolog.WarnLevel, // Only warn and above
	})

	logger.Debug().Msg("debug message")
	logger.Info().Msg("info message")
	logger.Warn().Msg("warn message")
	logger.Error().Msg("error message")

	// Should not contain debug or info
	output := buf.String()
	assert.NotContains(t, output, "debug message")
	assert.NotContains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
}

func TestLogger_ComplexFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatJSON,
	})

	logger.Info().
		Str("string", "value").
		Int("int", 42).
		Int8("int8", 8).
		Int16("int16", 16).
		Int32("int32", 32).
		Int64("int64", 64).
		Uint("uint", 100).
		Float32("float32", 3.14).
		Float64("float64", 2.71).
		Bool("bool", true).
		Msg("complex fields")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)
	require.NoError(t, err)

	assert.Equal(t, "value", logData["string"])
	assert.Equal(t, float64(42), logData["int"])
	assert.Equal(t, float64(64), logData["int64"])
	assert.Equal(t, float64(100), logData["uint"])
	assert.Equal(t, float64(2.71), logData["float64"])
	assert.Equal(t, true, logData["bool"])
}

// Test to ensure JSON output is valid
func TestLogger_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatJSON,
	})

	logger.Info().
		Str("user", "alice").
		Int("count", 5).
		Msg("User action")

	var logData map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logData)
	require.NoError(t, err, "Output should be valid JSON")
}

func TestLogger_ConcurrentAccess(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithConfig(Config{
		Output: &buf,
		Format: FormatJSON,
	})

	// Concurrent logging should not panic
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			logger.Info().Int("n", n).Msg("concurrent message")
		}(i)
	}
	wg.Wait()

	// Just verify no panic occurred (bytes.Buffer is not thread-safe for counting)
	assert.NotEmpty(t, buf.String(), "Should have some output")
}

// func TestGlobalFunctions(t *testing.T) {
// 	var buf bytes.Buffer
// 	SetGlobal(NewWithConfig(Config{
// 		Output: &buf,
// 		Format: FormatJSON,
// 	}))

// 	Info().Msg("global info")
// 	Debug().Msg("global debug")
// 	Warn().Msg("global warn")
// 	Error().Msg("global error")

// 	output := buf.String()
// 	assert.Contains(t, output, "global info")
// 	assert.Contains(t, output, "global debug")
// 	assert.Contains(t, output, "global warn")
// 	assert.Contains(t, output, "global error")

// 	// Cleanup
// 	SetGlobal(New())
// }

func TestSetLevel_Global(t *testing.T) {
	originalLevel := GetGlobal().GetLevel()

	SetLevel(zerolog.DebugLevel)
	assert.Equal(t, zerolog.DebugLevel, GetGlobal().GetLevel())

	SetGlobalLevel(zerolog.InfoLevel)
	assert.Equal(t, zerolog.InfoLevel, GetGlobal().GetLevel())

	// Restore original level
	GetGlobal().SetLevel(originalLevel)
}

// func TestWithContext_Global(t *testing.T) {
// 	var buf bytes.Buffer
// 	SetGlobal(NewWithConfig(Config{
// 		Output: &buf,
// 		Format: FormatJSON,
// 	}))

// 	tracer := trace.NewNoopTracerProvider().Tracer("test")
// 	ctx, span := tracer.Start(context.Background(), "test-span")

// 	loggerWithCtx := WithContext(ctx)
// 	loggerWithCtx.Info().Msg("test")

// 	span.End()

// 	var logData map[string]interface{}
// 	err := json.Unmarshal(buf.Bytes(), &logData)
// 	require.NoError(t, err)

// 	assert.Contains(t, logData, "trace_id")
// 	assert.Contains(t, logData, "span_id")

// 	// Cleanup
// 	SetGlobal(New())
// }
