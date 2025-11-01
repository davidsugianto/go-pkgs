package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCtx = context.Background()

func TestNew(t *testing.T) {
	client := New("localhost:6379")
	assert.NotNil(t, client)
	assert.NotNil(t, client.Client)
	defer client.Close()
}

func TestNewWithOptions(t *testing.T) {
	client := New("localhost:6379",
		WithPassword("testpass"),
		WithDB(1),
		WithPoolSize(20),
		WithMinIdleConns(10),
		WithTimeout(10*time.Second),
		WithMaxRetries(5),
	)
	assert.NotNil(t, client)
	defer client.Close()
}

func TestPing(t *testing.T) {
	// Skip test if Redis is not available
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
}

func TestSetGetDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Set
	err = client.Set(testCtx, "test:key", "test:value", 10*time.Second)
	require.NoError(t, err)

	// Get
	val, err := client.Get(testCtx, "test:key")
	require.NoError(t, err)
	assert.Equal(t, "test:value", val)

	// Delete
	err = client.Delete(testCtx, "test:key")
	require.NoError(t, err)

	// Get should fail
	_, err = client.Get(testCtx, "test:key")
	assert.Equal(t, ErrKeyNotFound, err)
}

func TestGetKeyNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	_, err = client.Get(testCtx, "test:nonexistent")
	assert.Equal(t, ErrKeyNotFound, err)
}

func TestGetBytes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	data := []byte("test bytes")
	err = client.Set(testCtx, "test:bytes", data, 10*time.Second)
	require.NoError(t, err)

	result, err := client.GetBytes(testCtx, "test:bytes")
	require.NoError(t, err)
	assert.Equal(t, data, result)

	client.Delete(testCtx, "test:bytes")
}

func TestSetJSONGetJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	user := User{ID: 1, Name: "John", Email: "john@example.com"}

	// SetJSON
	err = client.SetJSON(testCtx, "test:user", user, 10*time.Second)
	require.NoError(t, err)

	// GetJSON
	var result User
	err = client.GetJSON(testCtx, "test:user", &result)
	require.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)

	client.Delete(testCtx, "test:user")
}

func TestExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	exists, err := client.Exists(testCtx, "test:nonexistent")
	require.NoError(t, err)
	assert.False(t, exists)

	err = client.Set(testCtx, "test:exists", "value", 10*time.Second)
	require.NoError(t, err)

	exists, err = client.Exists(testCtx, "test:exists")
	require.NoError(t, err)
	assert.True(t, exists)

	client.Delete(testCtx, "test:exists")
}

func TestIncrementDecrement(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Set initial value
	err = client.Set(testCtx, "test:counter", 0, 10*time.Second)
	require.NoError(t, err)

	// Increment by 1
	val, err := client.Increment(testCtx, "test:counter", 1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), val)

	// Increment by 5
	val, err = client.Increment(testCtx, "test:counter", 5)
	require.NoError(t, err)
	assert.Equal(t, int64(6), val)

	// Decrement by 2
	val, err = client.Decrement(testCtx, "test:counter", 2)
	require.NoError(t, err)
	assert.Equal(t, int64(4), val)

	client.Delete(testCtx, "test:counter")
}

func TestExpireTTL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	err = client.Set(testCtx, "test:ttl", "value", 10*time.Second)
	require.NoError(t, err)

	ttl, err := client.TTL(testCtx, "test:ttl")
	require.NoError(t, err)
	assert.InDelta(t, float64(10*time.Second), float64(ttl), float64(2*time.Second))

	client.Delete(testCtx, "test:ttl")
}

func TestSetNX(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// First SetNX should succeed
	ok, err := client.SetNX(testCtx, "test:nx", "value1", 10*time.Second)
	require.NoError(t, err)
	assert.True(t, ok)

	// Second SetNX should fail
	ok, err = client.SetNX(testCtx, "test:nx", "value2", 10*time.Second)
	require.NoError(t, err)
	assert.False(t, ok)

	// Value should still be the first one
	val, err := client.Get(testCtx, "test:nx")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)

	client.Delete(testCtx, "test:nx")
}

func TestSetXX(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// SetXX on non-existent key should fail
	ok, err := client.SetXX(testCtx, "test:xx", "value1", 10*time.Second)
	require.NoError(t, err)
	assert.False(t, ok)

	// Set key first
	err = client.Set(testCtx, "test:xx", "value1", 10*time.Second)
	require.NoError(t, err)

	// SetXX on existing key should succeed
	ok, err = client.SetXX(testCtx, "test:xx", "value2", 10*time.Second)
	require.NoError(t, err)
	assert.True(t, ok)

	// Value should be updated
	val, err := client.Get(testCtx, "test:xx")
	require.NoError(t, err)
	assert.Equal(t, "value2", val)

	client.Delete(testCtx, "test:xx")
}

func TestMGetMSet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// MSet
	err = client.MSet(testCtx, "test:m1", "v1", "test:m2", "v2", "test:m3", "v3")
	require.NoError(t, err)

	// MGet
	values, err := client.MGet(testCtx, "test:m1", "test:m2", "test:m3")
	require.NoError(t, err)
	assert.Len(t, values, 3)
	assert.Equal(t, "v1", values[0])
	assert.Equal(t, "v2", values[1])
	assert.Equal(t, "v3", values[2])

	client.Delete(testCtx, "test:m1", "test:m2", "test:m3")
}

func TestHashOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// HSet
	err = client.HSet(testCtx, "test:hash", "field1", "value1")
	require.NoError(t, err)

	// HGet
	val, err := client.HGet(testCtx, "test:hash", "field1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)

	// HMSet
	err = client.HMSet(testCtx, "test:hash", "field2", "value2", "field3", "value3")
	require.NoError(t, err)

	// HGetAll
	all, err := client.HGetAll(testCtx, "test:hash")
	require.NoError(t, err)
	assert.Len(t, all, 3)
	assert.Equal(t, "value1", all["field1"])
	assert.Equal(t, "value2", all["field2"])
	assert.Equal(t, "value3", all["field3"])

	// HDel
	err = client.HDel(testCtx, "test:hash", "field1")
	require.NoError(t, err)

	all, err = client.HGetAll(testCtx, "test:hash")
	require.NoError(t, err)
	assert.Len(t, all, 2)

	client.Delete(testCtx, "test:hash")
}

func TestListOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// RPush
	err = client.RPush(testCtx, "test:list", "a", "b", "c")
	require.NoError(t, err)

	// LLen
	length, err := client.LLen(testCtx, "test:list")
	require.NoError(t, err)
	assert.Equal(t, int64(3), length)

	// LRange
	items, err := client.LRange(testCtx, "test:list", 0, -1)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, items)

	// LPush
	err = client.LPush(testCtx, "test:list", "x", "y")
	require.NoError(t, err)

	items, err = client.LRange(testCtx, "test:list", 0, -1)
	require.NoError(t, err)
	assert.Equal(t, []string{"y", "x", "a", "b", "c"}, items)

	// RPop
	val, err := client.RPop(testCtx, "test:list")
	require.NoError(t, err)
	assert.Equal(t, "c", val)

	// LPop
	val, err = client.LPop(testCtx, "test:list")
	require.NoError(t, err)
	assert.Equal(t, "y", val)

	client.Delete(testCtx, "test:list")
}

func TestSetOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// SAdd
	err = client.SAdd(testCtx, "test:set", "a", "b", "c")
	require.NoError(t, err)

	// SMembers
	members, err := client.SMembers(testCtx, "test:set")
	require.NoError(t, err)
	assert.Len(t, members, 3)

	// SIsMember
	isMember, err := client.SIsMember(testCtx, "test:set", "a")
	require.NoError(t, err)
	assert.True(t, isMember)

	isMember, err = client.SIsMember(testCtx, "test:set", "x")
	require.NoError(t, err)
	assert.False(t, isMember)

	// SRem
	err = client.SRem(testCtx, "test:set", "a")
	require.NoError(t, err)

	members, err = client.SMembers(testCtx, "test:set")
	require.NoError(t, err)
	assert.Len(t, members, 2)

	client.Delete(testCtx, "test:set")
}

func TestSortedSetOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// ZAdd
	err = client.ZAdd(testCtx, "test:zset",
		redis.Z{Score: 1.0, Member: "a"},
		redis.Z{Score: 2.0, Member: "b"},
		redis.Z{Score: 3.0, Member: "c"},
	)
	require.NoError(t, err)

	// ZRange
	items, err := client.ZRange(testCtx, "test:zset", 0, -1)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, items)

	// ZRangeByScore
	items, err = client.ZRangeByScore(testCtx, "test:zset", "2", "3")
	require.NoError(t, err)
	assert.Equal(t, []string{"b", "c"}, items)

	// ZRem
	err = client.ZRem(testCtx, "test:zset", "a")
	require.NoError(t, err)

	items, err = client.ZRange(testCtx, "test:zset", 0, -1)
	require.NoError(t, err)
	assert.Equal(t, []string{"b", "c"}, items)

	client.Delete(testCtx, "test:zset")
}

func TestStats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := New("localhost:6379")
	defer client.Close()

	err := client.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	stats := client.Stats()
	assert.NotNil(t, stats)
	assert.GreaterOrEqual(t, stats.Hits, uint32(0))
	assert.GreaterOrEqual(t, stats.Misses, uint32(0))
}
