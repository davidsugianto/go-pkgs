package redislock

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	m := New(RedisDriver{
		GoRedisClient: []redis.UniversalClient{client},
	})

	assert.NotNil(t, m)
}

func TestNewMutexW(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	m := New(RedisDriver{
		GoRedisClient: []redis.UniversalClient{client},
	})

	mutex, err := m.NewMutexW(MutexOpt{
		Key: "test-lock",
	})

	require.NoError(t, err)
	require.NotNil(t, mutex)
}

func TestNewMutexW_EmptyKey(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	m := New(RedisDriver{
		GoRedisClient: []redis.UniversalClient{client},
	})

	_, err := m.NewMutexW(MutexOpt{
		Key: "",
	})

	assert.Equal(t, ErrKeyEmpty, err)
}

func TestNewMutexW_NilClient(t *testing.T) {
	m := New(RedisDriver{
		GoRedisClient: nil,
	})

	_, err := m.NewMutexW(MutexOpt{
		Key: "test-lock",
	})

	assert.Equal(t, ErrClientEmpty, err)
}

func TestNewMutexW_Defaults(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	m := New(RedisDriver{
		GoRedisClient: []redis.UniversalClient{client},
	})

	mutex, err := m.NewMutexW(MutexOpt{
		Key: "test-lock",
	})
	require.NoError(t, err)
	require.NotNil(t, mutex)
}

func TestValidateRequest(t *testing.T) {
	tests := []struct {
		name     string
		input    MutexOpt
		expected MutexOpt
		err      error
	}{
		{
			name:  "empty key",
			input: MutexOpt{Key: ""},
			err:   ErrKeyEmpty,
		},
		{
			name:     "apply defaults",
			input:    MutexOpt{Key: "test"},
			expected: MutexOpt{Key: "test", LockTime: 15 * time.Second, RetryCount: 1, CheckTTLTime: 1 * time.Second, MaxRetryAutoExtend: 5},
		},
		{
			name:     "custom values preserved",
			input:    MutexOpt{Key: "test", LockTime: 30 * time.Second, RetryCount: 5, CheckTTLTime: 2 * time.Second, MaxRetryAutoExtend: 10},
			expected: MutexOpt{Key: "test", LockTime: 30 * time.Second, RetryCount: 5, CheckTTLTime: 2 * time.Second, MaxRetryAutoExtend: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequest(&tt.input)
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, tt.input)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	assert.Equal(t, "no mutex wrapper client found", ErrClientEmpty.Error())
	assert.Equal(t, "key empty", ErrKeyEmpty.Error())
	assert.Equal(t, "release fail because never acquired the lock ", ErrFailReleaseNotAcquired.Error())
	assert.Equal(t, "Extend fail because never acquired the lock ", ErrExtendNotAcquired.Error())
	assert.Equal(t, "mutex client empty", ErrMutexEmpty.Error())
}

func TestConstants(t *testing.T) {
	assert.Equal(t, 15, DefaultLockTime)
	assert.Equal(t, 1, DefaultRetryCount)
	assert.Equal(t, 1, DefaultCheckTTLTime)
	assert.Equal(t, 5, DefaultMaxRetryAutoExtend)
}

func TestInterface(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	m := New(RedisDriver{
		GoRedisClient: []redis.UniversalClient{client},
	})

	// Verify interface implementation
	var _ IMutex = m
	var _ IMutexDistLock = &mutex{}
}
