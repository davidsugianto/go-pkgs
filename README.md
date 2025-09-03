# go-pkgs

A collection of reusable Golang packages for building modern backend services.

## 🚀 Features

- Lightweight & idiomatic Go packages
- Ready-to-use building blocks for microservices
- Examples included for quick start
- Modular design — use only what you need

## 📂 Modules

### 1. `pkg/logger`  
Structured logging with log level support and optional JSON format.  

### 2. `pkg/redis`  
Wrapper around `go-redis` with helper methods for caching and connection handling.  

### 3. `pkg/httpserver`  
Graceful HTTP server with middleware support (logging, recovery, health checks).  

### 4. `pkg/config`  
Configuration loader supporting `.env`, JSON, and YAML with environment overrides.  

### 5. `pkg/db`  
Database connector wrapper for PostgreSQL/MySQL with migration support.  

### 6. `pkg/response`  
Standard API response format (`code, data, error`) with JSON writer helpers.  

### 7. `pkg/workerpool`  
Goroutine worker pool with configurable concurrency and graceful shutdown.  

### 8. `pkg/ratelimiter`  
In-memory or Redis-based rate limiter using token bucket / leaky bucket.  

### 9. `pkg/auth/jwt`  
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

- [ ] **Core Packages**
  - [ ] Implement `logger` with leveled and structured logging  
  - [ ] Implement `redis` wrapper with connection pool and helper methods  
  - [ ] Implement `httpserver` with graceful shutdown and middleware support  
  - [ ] Implement `config` loader with `.env`, JSON, and YAML support  
  - [ ] Implement `db` connector with migrations support  
  - [ ] Implement `response` utilities for consistent API responses  
  - [ ] Implement `workerpool` with job queue and concurrency control  
  - [ ] Implement `ratelimiter` with in-memory and Redis support  
  - [ ] Implement `auth/jwt` for token generation and validation  

- [ ] **Examples**
  - [ ] Add usage examples for each package under `examples/`  
  - [ ] Provide a sample microservice using multiple packages  

- [ ] **Testing & Quality**
  - [ ] Add unit tests for all modules  
  - [ ] Add integration tests (Redis, DB, HTTP server)  
  - [ ] Add CI pipeline with GitHub Actions (`go test ./...`, lint, vet)  
  - [ ] Add Go Report Card and Coverage badge  

- [ ] **Enhancements**
  - [ ] Add gRPC server wrapper (`pkg/grpcserver`)  
  - [ ] Add metrics exporter with Prometheus (`pkg/metrics`)  
  - [ ] Add distributed tracing middleware with OpenTelemetry  
  - [ ] Add caching abstraction (`pkg/cache`)  

- [ ] **Documentation**
  - [ ] Write package-level docs with `godoc` examples  
  - [ ] Add contribution guide (`CONTRIBUTING.md`)  
  - [ ] Add code of conduct (`CODE_OF_CONDUCT.md`)  

## 📜 License

MIT License – feel free to use and contribute.
