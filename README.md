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

### 🚧 Planned

#### `logger`  
Structured logging with log level support and optional JSON format.  

#### `redis`  
Wrapper around `go-redis` with helper methods for caching and connection handling.  

#### `httpserver`  
Graceful HTTP server with middleware support (logging, recovery, health checks).  

#### `config`  
Configuration loader supporting `.env`, JSON, and YAML with environment overrides.  

#### `db`  
Database connector wrapper for PostgreSQL/MySQL with migration support.  

#### `response`  
Standard API response format (`code, data, error`) with JSON writer helpers.  

#### `workerpool`  
Goroutine worker pool with configurable concurrency and graceful shutdown.  

#### `ratelimiter`  
In-memory or Redis-based rate limiter using token bucket / leaky bucket.  

#### `auth/jwt`  
JWT authentication helpers for token generation, validation, and middleware.  

## 🛠 Getting Started  

### 1. Prerequisites  
- Go 1.21+ installed → [Download Go](https://go.dev/dl/)  
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
  - [ ] Implement `logger` with leveled and structured logging  
  - [ ] Implement `redis` wrapper with connection pool and helper methods  
  - [ ] Implement `httpserver` with graceful shutdown and middleware support  
  - [ ] Implement `config` loader with `.env`, JSON, and YAML support  
  - [ ] Implement `db` connector with migrations support  
  - [ ] Implement `response` utilities for consistent API responses  
  - [ ] Implement `workerpool` with job queue and concurrency control  
  - [ ] Implement `ratelimiter` with in-memory and Redis support  
  - [ ] Implement `auth/jwt` for token generation and validation  

- [x] **Examples**
  - [x] Add usage examples for `grace` package  
  - [x] Add usage examples for `httpclient` package  
  - [ ] Provide a sample microservice using multiple packages  

- [x] **Testing & Quality**
  - [x] Add unit tests for `grace` package  
  - [x] Add unit tests for `httpclient` package  
  - [ ] Add integration tests (Redis, DB, HTTP server)  
  - [ ] Add CI pipeline with GitHub Actions (`go test ./...`, lint, vet)  
  - [ ] Add Go Report Card and Coverage badge  

- [ ] **Enhancements**
  - [ ] Add gRPC server wrapper  
  - [ ] Add metrics exporter with Prometheus  
  - [ ] Add distributed tracing middleware with OpenTelemetry  
  - [ ] Add caching abstraction  

- [x] **Documentation**
  - [x] Write README for `grace` package  
  - [x] Write README for `httpclient` package  
  - [ ] Write package-level docs with `godoc` examples  
  - [ ] Add contribution guide (`CONTRIBUTING.md`)  
  - [ ] Add code of conduct (`CODE_OF_CONDUCT.md`)  

## 📜 License

MIT License – feel free to use and contribute.
