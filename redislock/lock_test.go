package redislock

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMutex(t *testing.T, key string) (IMutexDistLock, *redis.Client) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Test connection
	err := client.Ping(context.Background()).Err()
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	m := New(RedisDriver{
		GoRedisClient: []redis.UniversalClient{client},
	})

	mutex, err := m.NewMutexW(MutexOpt{
		Key:       key,
		LockTime:  5 * time.Second,
		RetryCount: 3,
	})
	require.NoError(t, err)

	return mutex, client
}

func TestLock_GetLock_Release(t *testing.T) {
	mutex, client := newTestMutex(t, "test-lock-1")
	defer client.Close()

	// GetLock
	err := mutex.GetLock()
	require.NoError(t, err)

	// Release
	ok, err := mutex.Release()
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestLock_Release_WithoutAcquire(t *testing.T) {
	mutex, client := newTestMutex(t, "test-lock-2")
	defer client.Close()

	// Try to release without acquiring
	_, err := mutex.Release()
	assert.Equal(t, ErrFailReleaseNotAcquired, err)
}

func TestLock_Extend_WithoutAcquire(t *testing.T) {
	mutex, client := newTestMutex(t, "test-lock-3")
	defer client.Close()

	// Try to extend without acquiring
	_, err := mutex.Extend()
	assert.Equal(t, ErrExtendNotAcquired, err)
}

func TestLock_Extend(t *testing.T) {
	mutex, client := newTestMutex(t, "test-lock-4")
	defer client.Close()

	// GetLock
	err := mutex.GetLock()
	require.NoError(t, err)

	// Extend
	ok, err := mutex.Extend()
	require.NoError(t, err)
	assert.True(t, ok)

	// Release
	_, _ = mutex.Release()
}

func TestLock_AutoExtend(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
	defer client.Close()

	m := New(RedisDriver{
		GoRedisClient: []redis.UniversalClient{client},
	})

	mutex, err := m.NewMutexW(MutexOpt{
		Key:             "test-lock-autoextend",
		LockTime:        3 * time.Second,
		CheckTTLTime:    500 * time.Millisecond,
		AutoExtend:      true,
		MaxRetryAutoExtend: 3,
	})
	require.NoError(t, err)

	// GetLock
	err = mutex.GetLock()
	require.NoError(t, err)

	// Wait for more than the lock time
	// Auto-extend should keep the lock alive
	time.Sleep(4 * time.Second)

	// We should still be able to release (lock wasn't lost)
	ok, err := mutex.Release()
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestLock_SecondAcquireFails(t *testing.T) {
	mutex1, client := newTestMutex(t, "test-lock-exclusive")
	defer client.Close()

	// Create second mutex for same key
	m := New(RedisDriver{
		GoRedisClient: []redis.UniversalClient{client},
	})

	mutex2, err := m.NewMutexW(MutexOpt{
		Key:        "test-lock-exclusive",
		LockTime:   5 * time.Second,
		RetryCount: 1,
	})
	require.NoError(t, err)

	// First mutex acquires lock
	err = mutex1.GetLock()
	require.NoError(t, err)

	// Second mutex should fail to acquire
	err = mutex2.GetLock()
	assert.Error(t, err)

	// Release first
	_, _ = mutex1.Release()
}
