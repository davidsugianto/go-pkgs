package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/davidsugianto/go-pkgs/redislock"
	"github.com/go-redis/redis/v8"
)

func main() {
	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	// Test connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")
	fmt.Println()

	// Create distributed lock manager
	lockManager := redislock.New(redislock.RedisDriver{
		GoRedisClient: []redis.UniversalClient{redisClient},
	})

	// Example 1: Basic lock and release
	fmt.Println("=== Basic Lock ===")
	basicLock(lockManager)
	fmt.Println()

	// Example 2: Lock with auto-extend
	fmt.Println("=== Auto-Extend Lock ===")
	autoExtendLock(lockManager)
	fmt.Println()

	// Example 3: Lock with retry
	fmt.Println("=== Lock with Retry ===")
	retryLock(lockManager)
	fmt.Println()

	// Example 4: Multiple locks (mutex pattern)
	fmt.Println("=== Exclusive Lock Pattern ===")
	exclusiveLock(lockManager)
	fmt.Println()
}

func basicLock(lockManager redislock.IMutex) {
	// Create a mutex with 10 second TTL
	mutex, err := lockManager.NewMutexW(redislock.MutexOpt{
		Key:       "my-resource",
		LockTime:  10 * time.Second,
		RetryCount: 1,
	})
	if err != nil {
		log.Printf("Failed to create mutex: %v", err)
		return
	}

	// Acquire the lock
	if err := mutex.GetLock(); err != nil {
		log.Printf("Failed to acquire lock: %v", err)
		return
	}
	fmt.Println("Lock acquired")

	// Do critical section work
	fmt.Println("Doing critical work...")
	time.Sleep(1 * time.Second)

	// Release the lock
	ok, err := mutex.Release()
	if err != nil {
		log.Printf("Failed to release lock: %v", err)
		return
	}
	if ok {
		fmt.Println("Lock released successfully")
	}
}

func autoExtendLock(lockManager redislock.IMutex) {
	// Create a mutex with auto-extend enabled
	// This will automatically extend the lock if the work takes longer than expected
	mutex, err := lockManager.NewMutexW(redislock.MutexOpt{
		Key:             "long-running-task",
		LockTime:        3 * time.Second,   // Initial TTL
		CheckTTLTime:    500 * time.Millisecond, // How often to check and extend
		AutoExtend:      true,               // Enable auto-extend
		MaxRetryAutoExtend: 5,               // Max retries if extend fails
	})
	if err != nil {
		log.Printf("Failed to create mutex: %v", err)
		return
	}

	// Acquire the lock
	if err := mutex.GetLock(); err != nil {
		log.Printf("Failed to acquire lock: %v", err)
		return
	}
	fmt.Println("Lock acquired with auto-extend")

	// Simulate long-running work (longer than TTL)
	fmt.Println("Processing long task (5 seconds)...")
	time.Sleep(5 * time.Second)
	fmt.Println("Long task completed")

	// Release the lock
	ok, err := mutex.Release()
	if err != nil {
		log.Printf("Failed to release lock: %v", err)
		return
	}
	if ok {
		fmt.Println("Lock released successfully")
	}
}

func retryLock(lockManager redislock.IMutex) {
	// Create a mutex with retry
	mutex, err := lockManager.NewMutexW(redislock.MutexOpt{
		Key:        "retry-resource",
		LockTime:   5 * time.Second,
		RetryCount: 3, // Will retry 3 times if lock acquisition fails
	})
	if err != nil {
		log.Printf("Failed to create mutex: %v", err)
		return
	}

	// Acquire the lock
	if err := mutex.GetLock(); err != nil {
		log.Printf("Failed to acquire lock after retries: %v", err)
		return
	}
	fmt.Println("Lock acquired with retry config")

	// Do work
	time.Sleep(500 * time.Millisecond)

	// Release
	_, _ = mutex.Release()
	fmt.Println("Lock released")
}

func exclusiveLock(lockManager redislock.IMutex) {
	// Create two mutex instances for the same resource
	mutex1, _ := lockManager.NewMutexW(redislock.MutexOpt{
		Key:       "exclusive-resource",
		LockTime:  5 * time.Second,
		RetryCount: 1,
	})

	mutex2, _ := lockManager.NewMutexW(redislock.MutexOpt{
		Key:       "exclusive-resource",
		LockTime:  5 * time.Second,
		RetryCount: 1,
	})

	// First instance acquires the lock
	if err := mutex1.GetLock(); err != nil {
		log.Printf("Failed to acquire lock 1: %v", err)
		return
	}
	fmt.Println("Lock 1 acquired")

	// Second instance tries to acquire (should fail)
	if err := mutex2.GetLock(); err != nil {
		fmt.Printf("Lock 2 failed to acquire (expected): %v\n", err)
	}

	// Release first lock
	_, _ = mutex1.Release()
	fmt.Println("Lock 1 released")

	// Now second instance can acquire
	if err := mutex2.GetLock(); err != nil {
		log.Printf("Failed to acquire lock 2: %v", err)
		return
	}
	fmt.Println("Lock 2 acquired after lock 1 released")

	_, _ = mutex2.Release()
	fmt.Println("Lock 2 released")
}
