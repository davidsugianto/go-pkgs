package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/davidsugianto/go-pkgs/redis"
	redisdriver "github.com/redis/go-redis/v9"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Initialize Redis client
	client := redis.New("localhost:6379",
		redis.WithPassword(""), // Set password if needed
		redis.WithDB(0),
		redis.WithTimeout(5*time.Second),
	)
	defer client.Close()

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("✓ Connected to Redis")
	fmt.Println()

	// Basic string operations
	fmt.Println("=== Basic Operations ===")
	demoStringOperations(ctx, client)
	fmt.Println()

	// JSON operations
	fmt.Println("=== JSON Operations ===")
	demoJSONOperations(ctx, client)
	fmt.Println()

	// Counter operations
	fmt.Println("=== Counter Operations ===")
	demoCounterOperations(ctx, client)
	fmt.Println()

	// Hash operations
	fmt.Println("=== Hash Operations ===")
	demoHashOperations(ctx, client)
	fmt.Println()

	// List operations
	fmt.Println("=== List Operations ===")
	demoListOperations(ctx, client)
	fmt.Println()

	// Set operations
	fmt.Println("=== Set Operations ===")
	demoSetOperations(ctx, client)
	fmt.Println()

	// Sorted set operations
	fmt.Println("=== Sorted Set Operations ===")
	demoSortedSetOperations(ctx, client)
	fmt.Println()

	// Expiration and TTL
	fmt.Println("=== Expiration & TTL ===")
	demoExpiration(ctx, client)
	fmt.Println()

	// Conditional operations
	fmt.Println("=== Conditional Operations ===")
	demoConditionalOperations(ctx, client)
	fmt.Println()
}

func demoStringOperations(ctx context.Context, client *redis.Client) {
	// Set a key
	err := client.Set(ctx, "key1", "value1", time.Hour)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Set key1 = 'value1'")

	// Get a key
	val, err := client.Get(ctx, "key1")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Got key1 = '%s'\n", val)

	// Check if key exists
	exists, err := client.Exists(ctx, "key1")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("key1 exists: %v\n", exists)

	// Delete key
	err = client.Delete(ctx, "key1")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Deleted key1")

	// Try to get deleted key
	_, err = client.Get(ctx, "key1")
	if err == redis.ErrKeyNotFound {
		fmt.Println("key1 not found (expected)")
	}
}

func demoJSONOperations(ctx context.Context, client *redis.Client) {
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Store JSON
	err := client.SetJSON(ctx, "user:1", user, time.Hour)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Stored JSON: %+v\n", user)

	// Retrieve JSON
	var retrievedUser User
	err = client.GetJSON(ctx, "user:1", &retrievedUser)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Retrieved JSON: %+v\n", retrievedUser)

	client.Delete(ctx, "user:1")
}

func demoCounterOperations(ctx context.Context, client *redis.Client) {
	// Set initial value
	err := client.Set(ctx, "counter", 0, time.Hour)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Set counter = 0")

	// Increment
	val, err := client.Increment(ctx, "counter", 1)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Incremented by 1: %d\n", val)

	val, err = client.Increment(ctx, "counter", 5)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Incremented by 5: %d\n", val)

	// Decrement
	val, err = client.Decrement(ctx, "counter", 2)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Decremented by 2: %d\n", val)

	client.Delete(ctx, "counter")
}

func demoHashOperations(ctx context.Context, client *redis.Client) {
	// Set hash fields
	err := client.HSet(ctx, "user:2:profile", "name", "Jane Doe")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	err = client.HSet(ctx, "user:2:profile", "email", "jane@example.com")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Set hash fields")

	// Get hash field
	name, err := client.HGet(ctx, "user:2:profile", "name")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Got name: %s\n", name)

	// Set multiple hash fields at once
	err = client.HMSet(ctx, "user:2:profile", "age", "25", "city", "NYC")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Set multiple hash fields")

	// Get all hash fields
	all, err := client.HGetAll(ctx, "user:2:profile")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("All fields: %+v\n", all)

	client.Delete(ctx, "user:2:profile")
}

func demoListOperations(ctx context.Context, client *redis.Client) {
	// Push to list
	err := client.RPush(ctx, "tasks", "task1", "task2", "task3")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Pushed tasks to list")

	// Get list length
	length, err := client.LLen(ctx, "tasks")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("List length: %d\n", length)

	// Get all items
	items, err := client.LRange(ctx, "tasks", 0, -1)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("All items: %v\n", items)

	// Pop from list
	task, err := client.RPop(ctx, "tasks")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Popped task: %s\n", task)

	client.Delete(ctx, "tasks")
}

func demoSetOperations(ctx context.Context, client *redis.Client) {
	// Add to set
	err := client.SAdd(ctx, "tags", "golang", "redis", "go-pkgs")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Added tags to set")

	// Get all members
	members, err := client.SMembers(ctx, "tags")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("All tags: %v\n", members)

	// Check membership
	isMember, err := client.SIsMember(ctx, "tags", "golang")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Is 'golang' a member? %v\n", isMember)

	client.Delete(ctx, "tags")
}

func demoSortedSetOperations(ctx context.Context, client *redis.Client) {
	// Add to sorted set
	err := client.ZAdd(ctx, "leaderboard",
		redisdriver.Z{Score: 100, Member: "Alice"},
		redisdriver.Z{Score: 200, Member: "Bob"},
		redisdriver.Z{Score: 150, Member: "Charlie"},
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Added scores to leaderboard")

	// Get top players
	top, err := client.ZRange(ctx, "leaderboard", 0, -1)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Leaderboard: %v\n", top)

	client.Delete(ctx, "leaderboard")
}

func demoExpiration(ctx context.Context, client *redis.Client) {
	// Set key with expiration
	err := client.Set(ctx, "temp:key", "temp:value", 10*time.Second)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Set key with 10s expiration")

	// Get TTL
	ttl, err := client.TTL(ctx, "temp:key")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("TTL: %v\n", ttl)

	// Update expiration
	err = client.Expire(ctx, "temp:key", 20*time.Second)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Println("Extended expiration to 20s")

	ttl, err = client.TTL(ctx, "temp:key")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("New TTL: %v\n", ttl)

	client.Delete(ctx, "temp:key")
}

func demoConditionalOperations(ctx context.Context, client *redis.Client) {
	// SetNX - only set if not exists
	ok, err := client.SetNX(ctx, "lock:resource1", "locked", time.Minute)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("SetNX succeeded: %v\n", ok)

	// Try again - should fail
	ok, err = client.SetNX(ctx, "lock:resource1", "locked", time.Minute)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("SetNX succeeded (second time): %v\n", ok)

	// SetXX - only set if exists
	ok, err = client.SetXX(ctx, "lock:resource1", "updated", time.Minute)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("SetXX succeeded: %v\n", ok)

	client.Delete(ctx, "lock:resource1")
}
