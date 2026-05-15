package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/davidsugianto/go-pkgs/redis"
	redisdriver "github.com/go-redis/redis/v8"
)

type User struct {
	ID   int    `msgpack:"id"`
	Name string `msgpack:"name"`
}

func main() {
	// Create underlying redis client
	redisClient := redisdriver.NewClient(&redisdriver.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Create wrapper with prefix
	client := redis.New(redisClient, "myapp")
	defer client.Client().Close()

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")
	fmt.Println()

	// Basic string operations
	fmt.Println("=== Basic Operations ===")
	demoStringOperations(ctx, client)
	fmt.Println()

	// Struct operations with msgpack
	fmt.Println("=== Struct Operations ===")
	demoStructOperations(ctx, client)
	fmt.Println()

	// Compound key operations
	fmt.Println("=== Compound Key Operations ===")
	demoCompoundOperations(ctx, client)
	fmt.Println()

	// Hash operations
	fmt.Println("=== Hash Operations ===")
	demoHashOperations(ctx, client)
	fmt.Println()

	// Set operations
	fmt.Println("=== Set Operations ===")
	demoSetOperations(ctx, client)
	fmt.Println()

	// Counter operations
	fmt.Println("=== Counter Operations ===")
	demoCounterOperations(ctx, client)
	fmt.Println()

	// Expiration
	fmt.Println("=== Expiration ===")
	demoExpiration(ctx, client)
	fmt.Println()
}

func demoStringOperations(ctx context.Context, client *redis.Redis) {
	// Set a key with prefix
	err := client.Set(ctx, "greeting", "Hello, Redis!", 10*time.Minute)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Set greeting = 'Hello, Redis!'")

	// Get - note: Get doesn't apply prefix, need to use prefixed key
	val, err := client.Get(ctx, "myapp:greeting")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Got greeting = '%s'\n", val)

	// Delete
	_, err = client.Del(ctx, "myapp:greeting")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Deleted greeting")
}

func demoStructOperations(ctx context.Context, client *redis.Redis) {
	user := User{ID: 1, Name: "John Doe"}

	// SetStruct - serializes with msgpack
	err := client.SetStruct(ctx, "user:1", user, 10*time.Minute)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Stored user: %+v\n", user)

	// GetStruct - deserializes from msgpack
	var retrieved User
	err = client.GetStruct(ctx, "myapp:user:1", &retrieved)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Retrieved user: %+v\n", retrieved)

	// Cleanup
	client.Del(ctx, "myapp:user:1")
}

func demoCompoundOperations(ctx context.Context, client *redis.Redis) {
	// Compound keys are keys joined with ":"
	// Useful for hierarchical data like "user:123:profile"

	// SetExCompound - SET with compound key and expiration
	err := client.SetExCompound(ctx, "user", "123:profile", "profile data", 10*time.Minute)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Set compound key 'user:123:profile'")

	// GetCompound - GET with compound key
	val, err := client.GetCompound(ctx, "user", "123:profile")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Got compound value: '%s'\n", val)

	// GetEx - GET with automatic expiration refresh
	val, err = client.GetEx(ctx, "greeting", 5*time.Minute)
	if err != nil {
		if err == redisdriver.Nil {
			fmt.Println("Key not found")
		} else {
			log.Printf("Error: %v", err)
		}
	} else {
		fmt.Printf("Got with expiration refresh: '%s'\n", val)
	}

	// DelCompound - DELETE with compound key
	_, err = client.DelCompound(ctx, "user", "123:profile")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Deleted compound key")
}

func demoHashOperations(ctx context.Context, client *redis.Redis) {
	// HSetStruct - store struct in hash field
	data := User{ID: 2, Name: "Jane"}
	err := client.HSetStruct(ctx, "myapp:users", "user:2", data)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Stored struct in hash")

	// HGetStruct - retrieve struct from hash field
	var retrieved User
	err = client.HGetStruct(ctx, "myapp:users", "user:2", &retrieved)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Retrieved from hash: %+v\n", retrieved)

	// HGetAll - get all hash fields
	client.Client().HSet(ctx, "myapp:profile", "name", "Alice")
	client.Client().HSet(ctx, "myapp:profile", "email", "alice@example.com")

	all, err := client.HGetAll(ctx, "myapp:profile")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("All hash fields: %+v\n", all)

	// HDel - delete hash field
	_, err = client.HDel(ctx, "myapp:profile", "name")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Deleted hash field")

	// Cleanup
	client.Del(ctx, "myapp:users")
	client.Del(ctx, "myapp:profile")
}

func demoSetOperations(ctx context.Context, client *redis.Redis) {
	// AddInSet - SADD
	err := client.AddInSet(ctx, "myapp:tags", "golang")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	client.AddInSet(ctx, "myapp:tags", "redis")
	client.AddInSet(ctx, "myapp:tags", "backend")
	fmt.Println("Added members to set")

	// GetSetMembers - SMEMBERS
	members, err := client.GetSetMembers(ctx, "myapp:tags")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Set members: %v\n", members)

	// Cleanup
	client.Del(ctx, "myapp:tags")
}

func demoCounterOperations(ctx context.Context, client *redis.Redis) {
	// Setup initial value
	client.Client().Set(ctx, "myapp:counter", "0", 0)

	// IncrBy - increment by value
	val, err := client.IncrBy(ctx, "myapp:counter", 5)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("After increment by 5: %d\n", val)

	// HIncrBy - increment hash field
	client.Client().HSet(ctx, "myapp:stats", "views", "0")
	val, err = client.HIncrBy(ctx, "myapp:stats", "views", 10)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Hash field after increment by 10: %d\n", val)

	// Cleanup
	client.Del(ctx, "myapp:counter")
	client.Del(ctx, "myapp:stats")
}

func demoExpiration(ctx context.Context, client *redis.Redis) {
	// Set a key
	client.Set(ctx, "tempkey", "tempvalue", 5*time.Minute)
	fmt.Println("Set key with 5 minute expiration")

	// Expire - update expiration
	ok, err := client.Expire(ctx, "myapp:tempkey", 10*time.Minute)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Expiration updated: %v\n", ok)

	// Cleanup
	client.Del(ctx, "myapp:tempkey")
}
