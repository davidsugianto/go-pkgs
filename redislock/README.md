# Redislock - Distributed Lock with Redis

A distributed lock implementation using Redis and [redsync](https://github.com/go-redsync/redsync). Provides mutual exclusion across multiple application instances.

## Installation

```bash
go get github.com/davidsugianto/go-pkgs/redislock
```

## Quick Start

```go
package main

import (
    "time"
    
    "github.com/davidsugianto/go-pkgs/redislock"
    "github.com/go-redis/redis/v8"
)

func main() {
    // Create Redis client
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    // Create lock manager
    lockManager := redislock.New(redislock.RedisDriver{
        GoRedisClient: []redis.UniversalClient{redisClient},
    })
    
    // Create a mutex
    mutex, _ := lockManager.NewMutexW(redislock.MutexOpt{
        Key:      "my-resource",
        LockTime: 10 * time.Second,
    })
    
    // Acquire lock
    if err := mutex.GetLock(); err != nil {
        // Lock acquisition failed
        return
    }
    defer mutex.Release()
    
    // Critical section - only one instance can execute this
    doWork()
}
```

## Features

- **Distributed Mutual Exclusion** - Lock across multiple instances
- **Auto-Extend** - Automatically extend lock TTL for long-running operations
- **Configurable Retry** - Retry lock acquisition with customizable attempts
- **Manual Extend** - Extend lock TTL manually when needed
- **Multiple Redis Nodes** - Support for Redis clusters

## API Reference

### Creating a Lock Manager

```go
lockManager := redislock.New(redislock.RedisDriver{
    GoRedisClient: []redis.UniversalClient{redisClient},
})
```

For high availability, pass multiple Redis clients:

```go
lockManager := redislock.New(redislock.RedisDriver{
    GoRedisClient: []redis.UniversalClient{redis1, redis2, redis3},
})
```

### Mutex Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `Key` | `string` | (required) | Redis key for the lock |
| `LockTime` | `time.Duration` | 15s | Lock TTL |
| `WaitTimeExtra` | `time.Duration` | 0 | Extra time added to LockTime |
| `RetryCount` | `int` | 1 | Number of retry attempts |
| `AutoExtend` | `bool` | false | Auto-extend lock TTL |
| `CheckTTLTime` | `time.Duration` | 1s | Interval to check and auto-extend |
| `MaxRetryAutoExtend` | `int` | 5 | Max retries for auto-extend |

### Lock Operations

```go
// Create mutex
mutex, err := lockManager.NewMutexW(redislock.MutexOpt{
    Key:       "resource-name",
    LockTime:  10 * time.Second,
    RetryCount: 3,
})

// Acquire lock
err = mutex.GetLock()

// Extend lock (optional)
ok, err := mutex.Extend()

// Release lock
ok, err := mutex.Release()
```

## Usage Patterns

### Basic Lock

```go
mutex, _ := lockManager.NewMutexW(redislock.MutexOpt{
    Key:      "critical-section",
    LockTime: 30 * time.Second,
})

if err := mutex.GetLock(); err != nil {
    return err
}
defer mutex.Release()

// Critical section
processData()
```

### Lock with Auto-Extend

Use auto-extend for operations that may take longer than the initial TTL:

```go
mutex, _ := lockManager.NewMutexW(redislock.MutexOpt{
    Key:             "long-operation",
    LockTime:        5 * time.Second,
    AutoExtend:      true,
    CheckTTLTime:    1 * time.Second,  // Check every second
    MaxRetryAutoExtend: 3,
})

if err := mutex.GetLock(); err != nil {
    return err
}
defer mutex.Release()

// Long-running operation
// Lock will be auto-extended as needed
processLargeDataset()
```

### Lock with Retry

For high-contention scenarios:

```go
mutex, _ := lockManager.NewMutexW(redislock.MutexOpt{
    Key:        "contended-resource",
    LockTime:   10 * time.Second,
    RetryCount: 5,  // Try 5 times
})

if err := mutex.GetLock(); err != nil {
    // Failed after all retries
    return err
}
defer mutex.Release()
```

### Manual Extend

Extend lock TTL during execution:

```go
mutex, _ := lockManager.NewMutexW(redislock.MutexOpt{
    Key:      "flexible-operation",
    LockTime: 10 * time.Second,
})

mutex.GetLock()

// Halfway through, extend the lock
if moreWorkNeeded {
    mutex.Extend()
}

mutex.Release()
```

## Error Handling

```go
import "github.com/davidsugianto/go-pkgs/redislock"

mutex, _ := lockManager.NewMutexW(redislock.MutexOpt{
    Key: "resource",
})

// Errors from lock operations
err := mutex.GetLock()

// Check for specific errors
if err == redislock.ErrKeyEmpty {
    // Key was not provided
}
```

### Error Types

| Error | Description |
|-------|-------------|
| `ErrClientEmpty` | No Redis client provided |
| `ErrKeyEmpty` | Lock key is empty |
| `ErrMutexEmpty` | Mutex not initialized |
| `ErrFailReleaseNotAcquired` | Release called without acquiring |
| `ErrExtendNotAcquired` | Extend called without acquiring |

## Best Practices

1. **Always release locks** - Use `defer mutex.Release()` after acquiring
2. **Set appropriate TTL** - Balance between safety and deadlock risk
3. **Use auto-extend for long operations** - Prevent premature lock expiration
4. **Handle acquisition failures** - Have fallback logic when lock unavailable
5. **Use unique keys** - Different resources should have different lock keys

## Examples

See the `example/` directory for complete working examples:

```bash
# Start Redis
docker run -d -p 6379:6379 redis:latest

# Run example
cd example
go run main.go
```
