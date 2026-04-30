# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a collection of reusable Go packages for building modern backend services. Each package is independently importable via `go get github.com/davidsugianto/go-pkgs/<package>`.

## Development Commands

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./grace
go test ./logger
go test ./otel

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -v -run TestServeHTTP ./grace
```

## Package Architecture

Each package follows a consistent structure:
- `<package>.go` - Main implementation
- `<package>_test.go` - Unit tests using standard `testing` package and `testify`
- `example/main.go` - Usage example
- `README.md` - Package documentation

### Package Dependencies

- **logger** - Uses `rs/zerolog` with OpenTelemetry `trace` integration for span context correlation
- **otel** - OpenTelemetry SDK with OTLP gRPC exporters for traces and metrics
- **redis** - Wraps `go-redis/v9` with helper methods
- **config** - Supports JSON and YAML via `gopkg.in/yaml.v3`
- **grace**, **httpclient**, **pagination**, **response** - Standard library only

### Common Patterns

**Method chaining for configuration** - Packages like `otel` use a fluent builder pattern:
```go
config := otel.NewConfig("my-service").
    WithTracing(true).
    WithTraceEndpoint("localhost:4317")
```

**Global singleton access** - The `logger` package provides both instance-based and global functions:
```go
// Instance-based
log := logger.NewWithConfig(cfg)
log.Info().Msg("message")

// Global convenience functions
logger.Info().Msg("message")
```

## Go Version

Requires Go 1.23.0+.
