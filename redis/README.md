# Redis - Go-Redis Wrapper

A lightweight wrapper around `go-redis/v8` with helper methods for caching, struct serialization with msgpack, and key prefix support.

## Installation

```bash
go get github.com/davidsugianto/go-pkgs/redis
```

## Quick Start

```go
package main

import (
    "context"
    "time"
    
    "github.com/davidsugianto/go-pkgs/redis"
    redisdriver "github.com/go-redis/redis/v8"
)

func main() {
    // Create underlying redis client
    redisClient := redisdriver.NewClient(&redisdriver.Options{
        Addr: "localhost:6379",
    })

    // Create wrapper with prefix
    client := redis.New(redisClient, "myapp")
    defer client.Client().Close()

    ctx := context.Background()
    
    // Basic operations
    client.Set(ctx, "key", "value", time.Hour)
    val, _ := client.Get(ctx, "myapp:key")
    
    // Struct operations (uses msgpack)
    user := User{Name: "John"}
    client.SetStruct(ctx, "user:1", user, time.Hour)
    client.GetStruct(ctx, "myapp:user:1", &user)
}
```

## Features

- **Key Prefix Support** - Automatically prefix all keys
- **Struct Serialization** - Msgpack-based struct storage/retrieval
- **Compound Keys** - Helper methods for hierarchical keys like `user:123:profile`
- **Bulk Operations** - Pipeline-based bulk get/set
- **Hash Operations** - HSet/HGet with struct support
- **Set Operations** - SADD/SMEMBERS helpers
- **Counter Operations** - INCRBY/HINCRBY wrappers
- **Expiration Helpers** - GetEx with automatic TTL refresh

## API Reference

### Creating a Client

```go
import (
    "github.com/davidsugianto/go-pkgs/redis"
    redisdriver "github.com/go-redis/redis/v8"
)

redisClient := redisdriver.NewClient(&redisdriver.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

client := redis.New(redisClient, "myapp")
```

### Basic Operations

| Method | Description |
|--------|-------------|
| `New(client, prefix)` | Create wrapper with key prefix |
| `Client()` | Get underlying redis client |
| `Ping(ctx)` | Test connection |
| `Set(ctx, key, value, expire)` | Set string value |
| `Get(ctx, key)` | Get string value |
| `Del(ctx, key)` | Delete a key |
| `Expire(ctx, key, duration)` | Set key expiration |

### Struct Operations (Msgpack)

```go
type User struct {
    ID   int    `msgpack:"id"`
    Name string `msgpack:"name"`
}

user := User{ID: 1, Name: "John"}

// Store struct
client.SetStruct(ctx, "user:1", user, time.Hour)

// Retrieve struct
var result User
client.GetStruct(ctx, "myapp:user:1", &result)
```

### Compound Key Operations

Compound keys join parts with `:` separator for hierarchical data.

```go
// Set with compound key
client.SetExCompound(ctx, "user", "123:profile", "data", time.Hour)
// Creates key: "user:123:profile"

// Get with compound key
val, err := client.GetCompound(ctx, "user", "123:profile")

// Delete with compound key
client.DelCompound(ctx, "user", "123:profile")
```

### Bulk Operations

```go
// Bulk set
keys := []string{"k1", "k2", "k3"}
values := []string{"v1", "v2", "v3"}
errs := client.SetBulk(ctx, keys, values, time.Hour)

// Bulk get
vals, errs := client.GetBulk(ctx, keys)
```

### Hash Operations

```go
// Set hash field
client.HSetStruct(ctx, "users", "user:1", user)

// Get hash field
var user User
client.HGetStruct(ctx, "users", "user:1", &user)

// Get all fields
all, _ := client.HGetAll(ctx, "users")

// Delete field
client.HDel(ctx, "users", "user:1")
```

### Set Operations

```go
// Add to set
client.AddInSet(ctx, "tags", "golang")

// Get all members
members, _ := client.GetSetMembers(ctx, "tags")
```

### Counter Operations

```go
// Increment key
val, _ := client.IncrBy(ctx, "counter", 1)

// Increment hash field
val, _ := client.HIncrBy(ctx, "stats", "views", 1)
```

### GetEx - Get with TTL Refresh

```go
// Get and extend TTL atomically
val, err := client.GetEx(ctx, "cache:key", 5*time.Minute)
```

## Error Handling

```go
val, err := client.Get(ctx, "key")
if err == redis.ErrKeyNotFound {
    // Key doesn't exist
}
```

## Prefix Behavior

| Method | Applies Prefix |
|--------|----------------|
| `Set` | Yes |
| `Get` | No |
| `SetStruct` | Yes |
| `GetStruct` | No (uses Get) |
| `GetEx` | Yes |
| `Del` | No |
| `Expire` | No |
| `IncrBy` | No |
| `HGet/HSet/HDel` | No |

**Note:** When in doubt, check if the method uses `prefixed(key)` internally. Methods that don't apply prefix expect the full key including prefix.
