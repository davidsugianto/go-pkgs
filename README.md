# go-pkgs

A collection of reusable Golang packages for building modern backend services.

## 🚀 Features

- Lightweight & idiomatic Go packages
- Ready-to-use building blocks for microservices
- Examples included for quick start
- Modular design — use only what you need

## 📂 Modules

### ✅ Implemented

#### `grace` - Graceful HTTP Server Shutdown
Drop-in replacement for `http.ListenAndServe` with graceful shutdown handling. Works with any framework (Gin, Echo, standard library, etc.).

```bash
go get github.com/davidsugianto/go-pkgs/grace
```

[📖 Documentation](./grace/README.md) | [💡 Example](./grace/example/)

#### `httpclient` - Simple & Powerful HTTP Client
Lightweight HTTP client wrapper with automatic JSON serialization, context support, and flexible configuration.

```bash
go get github.com/davidsugianto/go-pkgs/httpclient
```

[📖 Documentation](./httpclient/README.md) | [💡 Example](./httpclient/example/)

#### `pagination` - Pagination Helper
A lightweight pagination package that helps you handle pagination logic with default values, offset/limit calculation, and total pages.

```bash
go get github.com/davidsugianto/go-pkgs/pagination
```

[📖 Documentation](./pagination/README.md) | [💡 Example](./pagination/example/)

#### `config` - Configuration Loader
Type-safe configuration loader with JSON and YAML support. Features automatic format detection, generics for type safety, and support for nested structures.

```bash
go get github.com/davidsugianto/go-pkgs/config
```

[📖 Documentation](./config/README.md) | [💡 Example](./config/example/)

#### `response` - Consistent API Response Utilities
A lightweight package for creating consistent JSON API responses with a standard format (`code`, `data`, `error`). Provides convenient helper functions for common HTTP status codes. Includes built-in support for Gin framework.

```bash
go get github.com/davidsugianto/go-pkgs/response
```

[📖 Documentation](./response/README.md) | [💡 Example](./response/example/)

#### `redis` - Go-Redis Wrapper
A lightweight wrapper around `go-redis` with helper methods for caching and connection handling. Provides a simple, idiomatic API for Redis operations including strings, hashes, lists, sets, sorted sets, and JSON serialization.

```bash
go get github.com/davidsugianto/go-pkgs/redis
```

[📖 Documentation](./redis/README.md) | [💡 Example](./redis/example/)

#### `logger` - Structured Logging with OpenTelemetry Integration
A structured logging package built on zerolog with automatic OpenTelemetry span context correlation. Provides high-performance, zero-allocation JSON logging with seamless integration for observability platforms like Grafana and Loki.

```bash
go get github.com/davidsugianto/go-pkgs/logger
```

[📖 Documentation](./logger/README.md) | [💡 Example](./logger/example/)

#### `otel` - OpenTelemetry Integration
A comprehensive OpenTelemetry integration package providing unified configuration for distributed tracing, metrics, and logging. Features selective enablement, method chaining API, trace-aware logging, and graceful shutdown.

```bash
go get github.com/davidsugianto/go-pkgs/otel
```

[📖 Documentation](./otel/README.md) | [💡 Example](./otel/example/)

#### `db` - Database Package
Multi-database support with GORM, automated migrations, and OpenTelemetry instrumentation. Supports PostgreSQL, MySQL, and MSSQL with connection pool metrics and embedded schema migrations.

```bash
go get github.com/davidsugianto/go-pkgs/db
```

[📖 Documentation](./db/README.md) | [💡 Example](./db/example/)

#### `logs` - Unified Logging Package
A unified logging layer built on top of the logger package with a simple, package-level API. Supports structured logging with fields, multiple log levels, and flexible configuration.

```bash
go get github.com/davidsugianto/go-pkgs/logs
```

[📖 Documentation](./logs/README.md) | [💡 Example](./logs/example/)

#### `redislock` - Distributed Lock with Redis
A distributed lock implementation using Redis and redsync. Provides mutual exclusion across multiple application instances with auto-extend support for long-running operations.

```bash
go get github.com/davidsugianto/go-pkgs/redislock
```

[📖 Documentation](./redislock/README.md) | [💡 Example](./redislock/example/)

### 🚧 Planned

#### `config` (`.env` support)
Add `.env` file support and environment variable overrides to the config package.

#### `workerpool`
Goroutine worker pool with configurable concurrency and graceful shutdown.

#### `ratelimiter`
In-memory or Redis-based rate limiter using token bucket / leaky bucket.

#### `auth/jwt`
JWT authentication helpers for token generation, validation, and middleware.

## 🛠 Getting Started

### 1. Prerequisites
- Go 1.25+ installed → [Download Go](https://go.dev/dl/)
- (Optional) Docker & Docker Compose for running Redis/DB examples

### 2. Clone the Repository
```bash
git clone https://github.com/davidsugianto/go-pkgs.git
cd go-pkgs
```

## 📌 Roadmap

- [x] **Core Packages**
  - [x] Implement `grace` with graceful HTTP server shutdown
  - [x] Implement `httpclient` with automatic JSON serialization and context support
  - [x] Implement `pagination` with offset/limit calculation and total pages
  - [x] Implement `config` loader with JSON and YAML support
  - [x] Implement `response` utilities for consistent API responses
  - [x] Implement `redis` wrapper with connection pool and helper methods
  - [x] Implement `logger` with OpenTelemetry integration and structured logging
  - [x] Implement `otel` with unified OpenTelemetry configuration for tracing, metrics, and logging
  - [x] Implement `db` connector with GORM, migrations, and OpenTelemetry support
  - [x] Implement `logs` unified logging layer with package-level API
  - [x] Implement `redislock` distributed lock with Redis
  - [ ] Implement `httpserver` with graceful shutdown and middleware support
  - [ ] Add `.env` file support to `config` package
  - [ ] Implement `workerpool` with job queue and concurrency control
  - [ ] Implement `ratelimiter` with in-memory and Redis support
  - [ ] Implement `auth/jwt` for token generation and validation

- [x] **Examples**
  - [x] Add usage examples for `grace` package
  - [x] Add usage examples for `httpclient` package
  - [x] Add usage examples for `pagination` package
  - [x] Add usage examples for `config` package
  - [x] Add usage examples for `response` package
  - [x] Add usage examples for `redis` package
  - [x] Add usage examples for `logger` package
  - [x] Add usage examples for `otel` package
  - [x] Add usage examples for `db` package
  - [x] Add usage examples for `logs` package
  - [x] Add usage examples for `redislock` package
  - [ ] Provide a sample microservice using multiple packages

- [x] **Testing & Quality**
  - [x] Add unit tests for `grace` package
  - [x] Add unit tests for `httpclient` package
  - [x] Add unit tests for `pagination` package
  - [x] Add unit tests for `config` package
  - [x] Add unit tests for `response` package
  - [x] Add unit tests for `redis` package
  - [x] Add unit tests for `logger` package
  - [x] Add unit tests for `otel` package
  - [x] Add unit tests for `db` package
  - [x] Add unit tests for `logs` package
  - [x] Add unit tests for `redislock` package
  - [ ] Add integration tests (Redis, DB, HTTP server)
  - [ ] Add CI pipeline with GitHub Actions (`go test ./...`, lint, vet)
  - [ ] Add Go Report Card and Coverage badge

- [ ] **Enhancements**
  - [ ] Add gRPC server wrapper
  - [ ] Add metrics exporter with Prometheus
  - [x] Add distributed tracing middleware with OpenTelemetry
  - [ ] Add caching abstraction

- [x] **Documentation**
  - [x] Write README for `grace` package
  - [x] Write README for `httpclient` package
  - [x] Write README for `pagination` package
  - [x] Write README for `config` package
  - [x] Write README for `response` package
  - [x] Write README for `redis` package
  - [x] Write README for `logger` package
  - [x] Write README for `otel` package
  - [x] Write README for `db` package
  - [x] Write README for `logs` package
  - [x] Write README for `redislock` package
  - [ ] Write package-level docs with `godoc` examples
  - [ ] Add contribution guide (`CONTRIBUTING.md`)
  - [ ] Add code of conduct (`CODE_OF_CONDUCT.md`)

## 📜 License

MIT License – feel free to use and contribute.
