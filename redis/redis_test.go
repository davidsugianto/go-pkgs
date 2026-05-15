package redis

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCtx = context.Background()

func newTestClient(t *testing.T) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return New(client, "test")
}

func TestNew(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r := New(client, "myprefix")
	assert.NotNil(t, r)
}

func TestClient(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r := New(client, "test")
	assert.Equal(t, client, r.Client())
}

func TestPing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}
}

func TestPrefixed(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r := New(client, "app")
	assert.Equal(t, "app:key", r.prefixed("key"))

	r2 := New(client, "")
	assert.Equal(t, "key", r2.prefixed("key"))
}

func TestSetGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Set
	err = r.Set(testCtx, "key1", "value1", 10*time.Second)
	require.NoError(t, err)

	// Get - note: Get doesn't use prefix, so we need to use the prefixed key
	val, err := r.Get(testCtx, "test:key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)

	// Cleanup
	r.Del(testCtx, "test:key1")
}

func TestSetStructGetStruct(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	type User struct {
		ID   int    `msgpack:"id"`
		Name string `msgpack:"name"`
	}

	user := User{ID: 1, Name: "John"}

	// SetStruct
	err = r.SetStruct(testCtx, "user:1", user, 10*time.Second)
	require.NoError(t, err)

	// GetStruct - note: GetStruct uses Get which doesn't apply prefix
	var retrieved User
	err = r.GetStruct(testCtx, "test:user:1", &retrieved)
	require.NoError(t, err)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, user.Name, retrieved.Name)

	// Cleanup
	r.Del(testCtx, "test:user:1")
}

func TestGetKeyNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	_, err = r.Get(testCtx, "test:nonexistent")
	assert.Equal(t, ErrKeyNotFound, err)
}

func TestHSetHGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// HSet via client directly
	err = r.Client().HSet(testCtx, "test:hash", "field1", "value1").Err()
	require.NoError(t, err)

	// HGet
	val, err := r.HGet(testCtx, "test:hash", "field1")
	require.NoError(t, err)
	assert.Equal(t, "value1", val)

	// Cleanup
	r.Del(testCtx, "test:hash")
}

func TestHGetAll(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Setup
	r.Client().HSet(testCtx, "test:hash2", "field1", "value1")
	r.Client().HSet(testCtx, "test:hash2", "field2", "value2")

	// HGetAll
	all, err := r.HGetAll(testCtx, "test:hash2")
	require.NoError(t, err)
	assert.Len(t, all, 2)
	assert.Equal(t, "value1", all["field1"])
	assert.Equal(t, "value2", all["field2"])

	// Cleanup
	r.Del(testCtx, "test:hash2")
}

func TestHDel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Setup
	r.Client().HSet(testCtx, "test:hash3", "field1", "value1")

	// HDel
	n, err := r.HDel(testCtx, "test:hash3", "field1")
	require.NoError(t, err)
	assert.Equal(t, int64(1), n)

	// Cleanup
	r.Del(testCtx, "test:hash3")
}

func TestDel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Setup
	r.Set(testCtx, "delkey", "delvalue", 10*time.Second)

	// Del
	n, err := r.Del(testCtx, "test:delkey")
	require.NoError(t, err)
	assert.Equal(t, int64(1), n)
}

func TestDelCompound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Setup
	r.Set(testCtx, "part1:part2", "value", 10*time.Second)

	// DelCompound
	n, err := r.DelCompound(testCtx, "part1", "part2")
	require.NoError(t, err)
	assert.Equal(t, int64(1), n)
}

func TestAddInSetGetSetMembers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// AddInSet
	err = r.AddInSet(testCtx, "test:set", "member1")
	require.NoError(t, err)

	err = r.AddInSet(testCtx, "test:set", "member2")
	require.NoError(t, err)

	// GetSetMembers
	members, err := r.GetSetMembers(testCtx, "test:set")
	require.NoError(t, err)
	assert.Len(t, members, 2)
	assert.Contains(t, members, "member1")
	assert.Contains(t, members, "member2")

	// Cleanup
	r.Del(testCtx, "test:set")
}

func TestHIncrBy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Setup
	r.Client().HSet(testCtx, "test:counter", "count", "0")

	// HIncrBy
	n, err := r.HIncrBy(testCtx, "test:counter", "count", 5)
	require.NoError(t, err)
	assert.Equal(t, int64(5), n)

	// Cleanup
	r.Del(testCtx, "test:counter")
}

func TestIncrBy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Setup
	r.Client().Set(testCtx, "test:incrkey", "0", 0)

	// IncrBy
	n, err := r.IncrBy(testCtx, "test:incrkey", 10)
	require.NoError(t, err)
	assert.Equal(t, int64(10), n)

	// Cleanup
	r.Del(testCtx, "test:incrkey")
}

func TestExpire(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Setup
	r.Set(testCtx, "expirekey", "value", 10*time.Second)

	// Expire
	ok, err := r.Expire(testCtx, "test:expirekey", 30*time.Second)
	require.NoError(t, err)
	assert.True(t, ok)

	// Cleanup
	r.Del(testCtx, "test:expirekey")
}

func TestHSetStructHGetStruct(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	type Data struct {
		Value string `msgpack:"value"`
	}

	data := Data{Value: "test"}

	// HSetStruct
	err = r.HSetStruct(testCtx, "test:hstruct", "field1", data)
	require.NoError(t, err)

	// HGetStruct
	var retrieved Data
	err = r.HGetStruct(testCtx, "test:hstruct", "field1", &retrieved)
	require.NoError(t, err)
	assert.Equal(t, data.Value, retrieved.Value)

	// Cleanup
	r.Del(testCtx, "test:hstruct")
}

func TestGetEx(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Setup
	r.Set(testCtx, "getexkey", "value", 10*time.Second)

	// GetEx - applies prefix
	val, err := r.GetEx(testCtx, "getexkey", 30*time.Second)
	require.NoError(t, err)
	assert.Equal(t, "value", val)

	// Cleanup
	r.Del(testCtx, "test:getexkey")
}

func TestGetCompound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Setup
	r.Set(testCtx, "key1:key2", "value", 10*time.Second)

	// GetCompound
	val, err := r.GetCompound(testCtx, "key1", "key2")
	require.NoError(t, err)
	assert.Equal(t, "value", val)

	// Cleanup
	r.Del(testCtx, "key1:key2")
}

func TestSetExCompound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// SetExCompound
	err = r.SetExCompound(testCtx, "part1", "part2", "value", 10*time.Second)
	require.NoError(t, err)

	// Verify
	val, err := r.Get(testCtx, "part1:part2")
	require.NoError(t, err)
	assert.Equal(t, "value", val)

	// Cleanup
	r.Del(testCtx, "part1:part2")
}

func TestSetBulk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	keys := []string{"bulk1", "bulk2", "bulk3"}
	values := []string{"val1", "val2", "val3"}

	// SetBulk
	errs := r.SetBulk(testCtx, keys, values, 10*time.Second)
	for _, e := range errs {
		assert.NoError(t, e)
	}

	// Cleanup
	r.Del(testCtx, "bulk1")
	r.Del(testCtx, "bulk2")
	r.Del(testCtx, "bulk3")
}

func TestGetBulk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	r := newTestClient(t)
	defer r.Client().Close()

	err := r.Ping(testCtx)
	if err != nil {
		t.Skip("Redis not available, skipping test")
	}

	// Setup
	r.Client().Set(testCtx, "gb1", "val1", 0)
	r.Client().Set(testCtx, "gb2", "val2", 0)

	keys := []string{"gb1", "gb2", "gb3"}

	// GetBulk
	vals, errs := r.GetBulk(testCtx, keys)
	assert.Len(t, vals, 3)
	assert.Len(t, errs, 3)
	assert.Equal(t, "val1", vals[0])
	assert.Equal(t, "val2", vals[1])
	assert.Equal(t, ErrKeyNotFound, errs[2])

	// Cleanup
	r.Del(testCtx, "gb1")
	r.Del(testCtx, "gb2")
}

func TestErrKeyNotFound(t *testing.T) {
	assert.Equal(t, redis.Nil, ErrKeyNotFound)
}
