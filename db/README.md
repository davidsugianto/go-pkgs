# Database Package

Multi-database support with GORM, automated migrations, and OpenTelemetry instrumentation.

## Overview

The `db` package provides a unified interface for connecting to multiple database systems with automatic OpenTelemetry tracing and metrics collection. Built on GORM and golang-migrate, it simplifies database operations while providing production-ready observability.

## Features

- **Multi-Database Support**: PostgreSQL, MySQL, MSSQL
- **GORM Integration**: Full ORM capabilities with GORM v2
- **Automatic Tracing**: Query-level distributed tracing
- **Connection Pool Metrics**: Real-time pool health monitoring
- **Schema Migrations**: Embedded migrations with golang-migrate
- **Type-Safe Configuration**: Validation with struct tags
- **Zero Configuration OTel**: Optional but seamless observability
- **Fluent Builder Pattern**: Chain configuration methods for clean setup

## Installation

```bash
go get github.com/davidsugianto/go-pkgs/db
```

## Quick Start

### Basic Connection

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/davidsugianto/go-pkgs/db"
)

func main() {
    ctx := context.Background()

    // Create config with fluent builder pattern
    config := db.NewConfig(db.Postgres, "localhost", "myapp").
        WithPort(5432).
        WithCredentials("postgres", "password").
        WithPool(25, 5, 5*time.Minute, 10*time.Minute)

    client, err := db.New(ctx, config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Verify connection
    if err := client.Ping(ctx); err != nil {
        log.Fatal(err)
    }
}
```

### Using Options Pattern

```go
client, err := db.New(ctx,
    db.NewConfig(db.MySQL, "localhost", "myapp"),
    db.WithPort(3306),
    db.WithCredentials("root", "password"),
    db.WithTracing(true),
    db.WithMetricsOption(true),
)
```

## Configuration

### Fluent Builder Methods

```go
config := db.NewConfig(db.Postgres, "localhost", "mydb").
    WithPort(5432).                              // Set custom port
    WithCredentials("user", "pass").             // Set username/password
    WithSSL("require", "cert.pem", "key.pem", "root.pem").  // SSL config
    WithPool(50, 10, 10*time.Minute, 20*time.Minute).       // Pool settings
    WithTracing(true).                           // Enable OTel tracing
    WithMetrics(true)                            // Enable OTel metrics
```

### Configuration Struct

```go
type Config struct {
    Type     string // postgres, mysql, mssql
    Host     string
    Port     int
    Username string
    Password string
    Database string

    // SSL
    SSLMode     string
    SSLCert     string
    SSLKey      string
    SSLRootCert string

    // Pool
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration

    // OpenTelemetry
    EnableTracing bool
    EnableMetrics bool
}
```

## Database Operations

### Auto Migration

```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:100;not null"`
    Email     string    `gorm:"size:255;uniqueIndex"`
    CreatedAt time.Time
}

if err := client.DB.AutoMigrate(&User{}); err != nil {
    log.Fatal(err)
}
```

### Transactions

```go
err := client.Transaction(ctx, func(tx *gorm.DB) error {
    user := User{Name: "John", Email: "john@example.com"}
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    return nil
})
```

### Context Propagation

```go
// Use WithContext for request-scoped operations
var users []User
if err := client.WithContext(ctx).Find(&users).Error; err != nil {
    log.Fatal(err)
}
```

### Health Check

```go
if client.IsHealthy(ctx) {
    fmt.Println("Database is healthy!")
}
```

### Connection Pool Stats

```go
stats, err := client.Stats()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Open: %d, InUse: %d, Idle: %d\n",
    stats.OpenConnections, stats.InUse, stats.Idle)
```

## Schema Migrations

### Embedded Migrations

```go
//go:embed migrations/*.sql
var migrationsFS embed.FS

// Run migrations
if err := client.Migrate(ctx, migrationsFS, "migrations"); err != nil {
    log.Fatal(err)
}
```

### Migration File Structure

```
migrations/
├── 000001_init.up.sql
├── 000001_init.down.sql
├── 000002_add_users.up.sql
└── 000002_add_users.down.sql
```

## OpenTelemetry Integration

### Automatic Tracing

When `EnableTracing` is true, all database operations are automatically traced:

```go
config := db.NewConfig(db.Postgres, "localhost", "myapp").
    WithTracing(true)
```

### Connection Pool Metrics

When `EnableMetrics` is true, pool metrics are collected:

```go
config := db.NewConfig(db.Postgres, "localhost", "myapp").
    WithMetrics(true)
```

Metrics exposed:
- `db.connection.pool.size` - Connection counts by state (max, open, in_use, idle, wait_count)

## Supported Databases

| Database | Type Constant | Default Port |
|----------|--------------|--------------|
| PostgreSQL | `db.Postgres` | 5432 |
| MySQL | `db.MySQL` | 3306 |
| MSSQL | `db.MSSQL` | 1433 |

## Error Handling

```go
import "errors"

client, err := db.New(ctx, config)
if errors.Is(err, db.ErrInvalidDatabaseType) {
    // Handle invalid database type
}
if errors.Is(err, db.ErrConnectionFailed) {
    // Handle connection failure
}
if errors.Is(err, db.ErrMigrationFailed) {
    // Handle migration failure
}
```

## Example

See [example/main.go](./example/main.go) for a complete working example.

## Dependencies

- `gorm.io/gorm` - ORM
- `gorm.io/driver/postgres` - PostgreSQL driver
- `gorm.io/driver/mysql` - MySQL driver
- `gorm.io/driver/sqlserver` - MSSQL driver
- `gorm.io/plugin/opentelemetry/tracing` - OTel tracing plugin
- `github.com/golang-migrate/migrate/v4` - Schema migrations
- `go.opentelemetry.io/otel` - OpenTelemetry SDK

## License

MIT
