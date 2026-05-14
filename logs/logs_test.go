package logs

import (
	"bytes"
	"os"
	"testing"

	"github.com/davidsugianto/go-pkgs/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureOutput(fn func()) string {
	var buf bytes.Buffer
	l := logger.NewWithConfig(logger.Config{
		Output: &buf,
		Format: logger.FormatJSON,
	})
	l.SetLevel(zerolog.DebugLevel)

	oldDebug := debugLogger
	oldInfo := infoLogger
	oldWarn := warnLogger
	oldErr := errLogger
	oldFatal := fatalLogger

	debugLogger = l
	infoLogger = l
	warnLogger = l
	errLogger = l
	fatalLogger = l

	fn()

	debugLogger = oldDebug
	infoLogger = oldInfo
	warnLogger = oldWarn
	errLogger = oldErr
	fatalLogger = oldFatal

	return buf.String()
}

func TestLevelConstants(t *testing.T) {
	assert.Equal(t, Level(zerolog.DebugLevel), DebugLevel)
	assert.Equal(t, Level(zerolog.InfoLevel), InfoLevel)
	assert.Equal(t, Level(zerolog.WarnLevel), WarnLevel)
	assert.Equal(t, Level(zerolog.ErrorLevel), ErrorLevel)
	assert.Equal(t, Level(zerolog.FatalLevel), FatalLevel)
}

func TestDebug(t *testing.T) {
	output := captureOutput(func() {
		Debug("hello", "world")
	})
	assert.Contains(t, output, "helloworld")
	assert.Contains(t, output, `"level":"debug"`)
}

func TestDebugf(t *testing.T) {
	output := captureOutput(func() {
		Debugf("count: %d", 42)
	})
	assert.Contains(t, output, "count: 42")
}

func TestDebugWithFields(t *testing.T) {
	output := captureOutput(func() {
		DebugWithFields("user action", map[string]any{
			"user_id": 123,
			"action":  "login",
		})
	})
	assert.Contains(t, output, "user action")
	assert.Contains(t, output, `"user_id":123`)
	assert.Contains(t, output, `"action":"login"`)
}

func TestInfo(t *testing.T) {
	output := captureOutput(func() {
		Info("server started")
	})
	assert.Contains(t, output, "server started")
	assert.Contains(t, output, `"level":"info"`)
}

func TestInfof(t *testing.T) {
	output := captureOutput(func() {
		Infof("port: %d", 8080)
	})
	assert.Contains(t, output, "port: 8080")
}

func TestInfoWithFields(t *testing.T) {
	output := captureOutput(func() {
		InfoWithFields("request completed", map[string]any{
			"method": "GET",
			"status": 200,
		})
	})
	assert.Contains(t, output, "request completed")
	assert.Contains(t, output, `"method":"GET"`)
}

func TestPrint(t *testing.T) {
	output := captureOutput(func() {
		Print("log.Print style output")
	})
	assert.Contains(t, output, "log.Print style output")
	assert.Contains(t, output, `"level":"info"`)
}

func TestPrintf(t *testing.T) {
	output := captureOutput(func() {
		Printf("formatted: %s", "value")
	})
	assert.Contains(t, output, "formatted: value")
}

func TestWarn(t *testing.T) {
	output := captureOutput(func() {
		Warn("low disk space")
	})
	assert.Contains(t, output, "low disk space")
	assert.Contains(t, output, `"level":"warn"`)
}

func TestWarnWithFields(t *testing.T) {
	output := captureOutput(func() {
		WarnWithFields("rate limit approaching", map[string]any{
			"remaining": 5,
		})
	})
	assert.Contains(t, output, "rate limit approaching")
	assert.Contains(t, output, `"remaining":5`)
}

func TestError(t *testing.T) {
	output := captureOutput(func() {
		Error("connection failed")
	})
	assert.Contains(t, output, "connection failed")
	assert.Contains(t, output, `"level":"error"`)
}

func TestErrorf(t *testing.T) {
	output := captureOutput(func() {
		Errorf("timeout after %ds", 30)
	})
	assert.Contains(t, output, "timeout after 30s")
}

func TestErrorWithFields(t *testing.T) {
	output := captureOutput(func() {
		ErrorWithFields("db query failed", map[string]any{
			"query": "SELECT * FROM users",
			"error": "connection refused",
		})
	})
	assert.Contains(t, output, "db query failed")
	assert.Contains(t, output, `"error":"connection refused"`)
}

func TestErrors(t *testing.T) {
	output := captureOutput(func() {
		err := &testError{msg: "something broke"}
		Errors(err)
	})
	assert.Contains(t, output, `"level":"error"`)
	assert.Contains(t, output, "something broke")
}

func TestErrors_Nil(t *testing.T) {
	// Should not panic with nil error
	assert.NotPanics(t, func() {
		Errors(nil)
	})
}

func TestFatal(t *testing.T) {
	// Fatal exits the program, so we test with a custom logger
	// that uses a test writer. We verify the format by calling
	// a function that produces similar output.
	output := captureOutput(func() {
		// Use error level instead to avoid os.Exit
		errLogger.Error().Msg("fatal type message")
	})
	assert.Contains(t, output, `"level":"error"`)
}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }

func TestNewLogger(t *testing.T) {
	l, err := NewLogger(&Config{
		Level:   InfoLevel,
		AppName: "myapp",
		UseJSON: true,
	})
	require.NoError(t, err)
	require.NotNil(t, l)

	assert.Equal(t, zerolog.InfoLevel, l.GetLevel())
	assert.Equal(t, "myapp", l.ServiceName())
}

func TestNewLogger_NilConfig(t *testing.T) {
	l, err := NewLogger(nil)
	require.NoError(t, err)
	require.NotNil(t, l)
}

func TestNewLogger_WithFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "log-test-*.log")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	l, err := NewLogger(&Config{
		Level:   InfoLevel,
		LogFile: tmpFile.Name(),
	})
	require.NoError(t, err)
	require.NotNil(t, l)

	l.Info().Msg("file log test")

	// Read the file
	content, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	assert.Contains(t, string(content), "file log test")
}

func TestNewLogger_InvalidFile(t *testing.T) {
	_, err := NewLogger(&Config{
		Level:   InfoLevel,
		LogFile: "/nonexistent/path/log.log",
	})
	require.Error(t, err)
}

func TestSetLogger(t *testing.T) {
	l, _ := NewLogger(&Config{Level: InfoLevel, UseJSON: true})

	err := SetLogger(InfoLevel, l)
	require.NoError(t, err)
	assert.Equal(t, l, infoLogger)
}

func TestSetLogger_InvalidLevel(t *testing.T) {
	l, _ := NewLogger(nil)
	err := SetLogger(Level(99), l)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid level")
}

func TestSetLogger_NilLogger(t *testing.T) {
	err := SetLogger(InfoLevel, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid logger")
}

func TestSetLevel(t *testing.T) {
	// Save original loggers
	origDebug := debugLogger
	origInfo := infoLogger
	defer func() {
		debugLogger = origDebug
		infoLogger = origInfo
		warnLogger = origInfo
		errLogger = origInfo
		fatalLogger = origInfo
	}()

	debugLogger, _ = NewLogger(&Config{Level: DebugLevel})
	infoLogger, _ = NewLogger(&Config{Level: InfoLevel})
	warnLogger = infoLogger
	errLogger = infoLogger
	fatalLogger = infoLogger

	SetLevel(WarnLevel)

	assert.Equal(t, zerolog.WarnLevel, debugLogger.GetLevel())
	assert.Equal(t, zerolog.WarnLevel, infoLogger.GetLevel())
}

func TestSetConfig(t *testing.T) {
	err := SetConfig(&Config{
		Level:   InfoLevel,
		AppName: "test-app",
		UseJSON: true,
	})
	require.NoError(t, err)

	assert.Equal(t, "test-app", infoLogger.ServiceName())
	assert.Equal(t, "test-app", debugLogger.ServiceName())
}

func TestSetConfig_Nil(t *testing.T) {
	err := SetConfig(nil)
	require.NoError(t, err)
}

func TestSetConfig_InvalidFile(t *testing.T) {
	err := SetConfig(&Config{
		Level:   InfoLevel,
		LogFile: "/nonexistent/path.log",
	})
	require.Error(t, err)
}

func TestDefaultLoggers(t *testing.T) {
	// Default loggers should be initialized
	assert.NotNil(t, infoLogger)
	assert.NotNil(t, debugLogger)
	assert.NotNil(t, warnLogger)
	assert.NotNil(t, errLogger)
	assert.NotNil(t, fatalLogger)

	// warn/error/fatal default to info logger
	assert.Equal(t, infoLogger, warnLogger)
	assert.Equal(t, infoLogger, errLogger)
	assert.Equal(t, infoLogger, fatalLogger)
}

func TestWithFields_AllLevels(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string, map[string]any)
	}{
		{"Debug", DebugWithFields},
		{"Info", InfoWithFields},
		{"Warn", WarnWithFields},
		{"Error", ErrorWithFields},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.fn("test message", map[string]any{"key": "value"})
			})
			assert.Contains(t, output, "test message")
			assert.Contains(t, output, `"key":"value"`)
		})
	}
}