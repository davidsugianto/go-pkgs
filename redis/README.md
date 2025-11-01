# Redis - Go-Redis Wrapper

A lightweight wrapper around `go-redis` with helper methods for caching and connection handling. Provides a simple, idiomatic API for Redis operations including strings, hashes, lists, sets, sorted sets, and JSON serialization.

## Installation

```bash
go get github.com/davidsugianto/go-pkgs/redis
```

You'll also need to add the `go-redis` dependency:

```bash
go get github.com/redis/go-redis/v9
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/davidsugianto/go-pkgs/redis"
)

func main() {
    // Create a new Redis client
    client := redis.New("localhost:6379")
    defer client.Close()

    ctx := context.Background()
    
    // Basic operations
    client.Set(ctx, "key", "value", time.Hour)
    val, _ := client.Get(ctx, "key")
    
    // JSON operations
    user := User{Name: "John"}
    client.SetJSON(ctx, "user:1", user, time.Hour)
    client.GetJSON(ctx, "user:1", &user)
}
```

## Features

- ✅ **Connection Management** - Automatic connection pooling with configurable options
- ✅ **Type Safety** - Helper methods with clear error handling
- ✅ **JSON Support** - Built-in JSON serialization/deserialization
- ✅ **Atomic Operations** - Support for SETNX, SETXX, and other atomic operations
- ✅ **Data Structures** - Full support for strings, hashes, lists, sets, and sorted sets
- ✅ **Expiration** - Easy key expiration and TTL management
- ✅ **Context Support** - All operations support context for cancellation and timeouts
- ✅ **Connection Pool Stats** - Monitor connection pool health

## Client Configuration

### Basic Setup

```go
client := redis.New("localhost:6379")
```

### With Options

```go
client := redis.New("localhost:6379",
    redis.WithPassword("mypassword"),
    redis.WithDB(1),
    redis.WithPoolSize(20),
    redis.WithMinIdleConns(10),
    redis.WithTimeout(5*time.Second),
    redis.WithMaxRetries(3),
)
```

### Available Options

- `WithPassword(password string)` - Set Redis password
- `WithDB(db int)` - Select database number (default: 0)
- `WithPoolSize(size int)` - Connection pool size (default: 10)
- `WithMinIdleConns(conns int)` - Minimum idle connections (default: 5)
- `WithTimeout(timeout time.Duration)` - Dial, read, and write timeout (default: 5s)
- `WithMaxRetries(retries int)` - Maximum retry attempts (default: 3)

## Basic Operations

### Connection Testing

```go
err := client.Ping(ctx)
```

### Set and Get

```go
// Set a key
err := client.Set(ctx, "key", "value", time.Hour)

// Get a key
val, err := client.Get(ctx, "key")
if err == redis.ErrKeyNotFound {
    // Key doesn't exist
}

// Get as bytes
data, err := client.GetBytes(ctx, "key")
```

### Delete and Exists

```go
// Delete one or more keys
err := client.Delete(ctx, "key1", "key2", "key3")

// Check if key exists
exists, err := client.Exists(ctx, "key")
```

### Expiration

```go
// Set expiration
err := client.Expire(ctx, "key", 10*time.Second)

// Get TTL
ttl, err := client.TTL(ctx, "key")
```

## JSON Operations

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Store JSON
user := User{ID: 1, Name: "John"}
err := client.SetJSON(ctx, "user:1", user, time.Hour)

// Retrieve JSON
var user User
err := client.GetJSON(ctx, "user:1", &user)
```

## Atomic Operations

### Conditional Sets

```go
// Set if key doesn't exist (useful for locking)
success, err := client.SetNX(ctx, "lock:resource", "locked", time.Minute)

// Set if key exists (useful for updates)
success, err := client.SetXX(ctx, "user:1", "updated", time.Hour)
```

### Counters

```go
// Increment
newVal, err := client.Increment(ctx, "counter", 1)  // Increment by 1
newVal, err := client.Increment(ctx, "counter", 5)  // Increment by 5

// Decrement
newVal, err := client.Decrement(ctx, "counter", 2)
```

### Batch Operations

```go
// Get multiple keys
values, err := client.MGet(ctx, "key1", "key2", "key3")

// Set multiple keys
err := client.MSet(ctx, "key1", "val1", "key2", "val2", "key3", "val3")
```

## Hash Operations

```go
// Set hash field
err := client.HSet(ctx, "user:1:profile", "name", "John")

// Get hash field
name, err := client.HGet(ctx, "user:1:profile", "name")

// Set multiple fields
err := client.HMSet(ctx, "user:1:profile", "name", "John", "age", "30", "city", "NYC")

// Get all fields
all, err := client.HGetAll(ctx, "user:1:profile")
// all is map[string]string

// Delete fields
err := client.HDel(ctx, "user:1:profile", "name", "age")
```

## List Operations

```go
// Push to list
err := client.LPush(ctx, "tasks", "task1", "task2")  // Push to left
err := client.RPush(ctx, "tasks", "task3", "task4")  // Push to right

// Pop from list
task, err := client.LPop(ctx, "tasks")  // Pop from left
task, err := client.RPop(ctx, "tasks")  // Pop from right

// Get list length
length, err := client.LLen(ctx, "tasks")

// Get range of items
items, err := client.LRange(ctx, "tasks", 0, -1)  // Get all
items, err := client.LRange(ctx, "tasks", 0, 2)   // Get first 3
```

## Set Operations

```go
// Add members
err := client.SAdd(ctx, "tags", "golang", "redis", "go-pkgs")

// Get all members
members, err := client.SMembers(ctx, "tags")

// Check membership
isMember, err := client.SIsMember(ctx, "tags", "golang")

// Remove members
err := client.SRem(ctx, "tags", "golang")
```

## Sorted Set Operations

```go
// Add with scores
err := client.ZAdd(ctx, "leaderboard",
    redis.Z{Score: 100, Member: "Alice"},
    redis.Z{Score: 200, Member: "Bob"},
    redis.Z{Score: 150, Member: "Charlie"},
)

// Get by index range
players, err := client.ZRange(ctx, "leaderboard", 0, -1)  // All players

// Get by score range
top, err := client.ZRangeByScore(ctx, "leaderboard", "150", "300")

// Remove members
err := client.ZRem(ctx, "leaderboard", "Alice")
```

## Pattern Matching

```go
// Get all keys matching pattern
keys, err := client.Keys(ctx, "user:*")

// Scan keys (safer for large datasets)
cursor := uint64(0)
for {
    keys, cursor, err := client.Scan(ctx, cursor, "user:*", 100)
    // Process keys
    if cursor == 0 {
        break
    }
}
```

## Pub/Sub

```go
// Publish message
err := client.Publish(ctx, "channel", "message")

// Subscribe
pubsub := client.Subscribe(ctx, "channel")
defer pubsub.Close()

// Receive messages
for msg := range pubsub.Channel() {
    fmt.Printf("Received: %s\n", msg.Payload)
}
```

## Monitoring

### Connection Pool Stats

```go
stats := client.Stats()
fmt.Printf("Hits: %d\n", stats.Hits)
fmt.Printf("Misses: %d\n", stats.Misses)
fmt.Printf("Timeouts: %d\n", stats.Timeouts)
```

## Error Handling

The package provides two common error types:

```go
if err == redis.ErrKeyNotFound {
    // Key doesn't exist
}

if err == redis.ErrConnectionFailed {
    // Connection failed
}
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/davidsugianto/go-pkgs/redis"
)

func main() {
    // Initialize client
    client := redis.New("localhost:6379")
    defer client.Close()

    ctx := context.Background()
    
    // Test connection
    if err := client.Ping(ctx); err != nil {
        panic(err)
    }
    
    // Basic operations
    client.Set(ctx, "greeting", "Hello, Redis!", time.Hour)
    val, _ := client.Get(ctx, "greeting")
    fmt.Println(val)
    
    // JSON operations
    user := User{ID: 1, Name: "John"}
    client.SetJSON(ctx, "user:1", user, time.Hour)
    var retrievedUser User
    client.GetJSON(ctx, "user:1", &retrievedUser)
    fmt.Printf("Retrieved: %+v\n", retrievedUser)
    
    // Hash operations
    client.HMSet(ctx, "user:1:profile",
        "name", "John Doe",
        "email", "john@example.com",
    )
    profile, _ := client.HGetAll(ctx, "user:1:profile")
    fmt.Printf("Profile: %+v\n", profile)
}
```

## Advanced Usage

### Using Context for Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

val, err := client.Get(ctx, "key")
```

### Accessing Underlying Redis Client

The `Client` embeds `*redis.Client`, so you can access all go-redis functionality:

```go
// Use any go-redis method directly
result := client.Pipeline()
result.Set(ctx, "key1", "val1", 0)
result.Get(ctx, "key1")
_, err := result.Exec(ctx)
```

## Best Practices

1. **Always use context** - Pass context to all operations for cancellation and timeouts
2. **Handle errors** - Check for `ErrKeyNotFound` when getting values
3. **Close connections** - Always call `client.Close()` when done
4. **Use connection pooling** - Configure pool size based on your workload
5. **Set appropriate timeouts** - Configure timeouts to prevent hanging connections
6. **Use JSON for complex data** - Use `SetJSON`/`GetJSON` for structured data
7. **Use atomic operations** - Use `SetNX` for distributed locking
8. **Monitor pool stats** - Track connection pool statistics in production

## Examples

See the `example/` directory for a complete working example:

```bash
# Start Redis (if not already running)
docker run -d -p 6379:6379 redis:latest

# Run example
cd example
go run main.go
```

