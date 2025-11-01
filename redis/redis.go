package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrKeyNotFound indicates a key does not exist in Redis
	ErrKeyNotFound = errors.New("key not found")

	// ErrConnectionFailed indicates a failed connection attempt
	ErrConnectionFailed = errors.New("connection failed")
)

// Client wraps the go-redis client with helper methods
type Client struct {
	*redis.Client
}

// Option configures the Redis client
type Option func(*redis.Options)

// New creates a new Redis client with default options
func New(addr string, opts ...Option) *Client {
	options := &redis.Options{
		Addr:         addr,
		Password:     "",
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
	}

	for _, opt := range opts {
		opt(options)
	}

	return &Client{
		Client: redis.NewClient(options),
	}
}

// WithPassword sets the Redis password
func WithPassword(password string) Option {
	return func(opts *redis.Options) {
		opts.Password = password
	}
}

// WithDB sets the Redis database number
func WithDB(db int) Option {
	return func(opts *redis.Options) {
		opts.DB = db
	}
}

// WithPoolSize sets the connection pool size
func WithPoolSize(size int) Option {
	return func(opts *redis.Options) {
		opts.PoolSize = size
	}
}

// WithMinIdleConns sets the minimum idle connections
func WithMinIdleConns(conns int) Option {
	return func(opts *redis.Options) {
		opts.MinIdleConns = conns
	}
}

// WithTimeout sets dial, read, and write timeouts
func WithTimeout(timeout time.Duration) Option {
	return func(opts *redis.Options) {
		opts.DialTimeout = timeout
		opts.ReadTimeout = timeout
		opts.WriteTimeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(retries int) Option {
	return func(opts *redis.Options) {
		opts.MaxRetries = retries
	}
}

// Ping checks the Redis connection
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.Client.Ping(ctx).Result()
	return err
}

// Set stores a key-value pair with expiration
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key (returns ErrKeyNotFound if key doesn't exist)
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	return val, err
}

// GetBytes retrieves a value as bytes by key
func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	val, err := c.Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrKeyNotFound
	}
	return val, err
}

// Delete removes one or more keys
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.Client.Del(ctx, keys...).Err()
}

// Exists checks if one or more keys exist
func (c *Client) Exists(ctx context.Context, keys ...string) (bool, error) {
	count, err := c.Client.Exists(ctx, keys...).Result()
	return count > 0, err
}

// SetJSON stores a JSON-serialized value with expiration
func (c *Client) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return c.Set(ctx, key, jsonData, expiration)
}

// GetJSON retrieves and unmarshals a JSON value into the provided type
func (c *Client) GetJSON(ctx context.Context, key string, dest interface{}) error {
	jsonData, err := c.GetBytes(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, dest); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// Increment increments the value of a key by the specified amount
func (c *Client) Increment(ctx context.Context, key string, value int64) (int64, error) {
	if value == 1 {
		return c.Client.Incr(ctx, key).Result()
	}
	return c.Client.IncrBy(ctx, key, value).Result()
}

// Decrement decrements the value of a key by the specified amount
func (c *Client) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	if value == 1 {
		return c.Client.Decr(ctx, key).Result()
	}
	return c.Client.DecrBy(ctx, key, value).Result()
}

// Expire sets a key's expiration time
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Client.Expire(ctx, key, expiration).Err()
}

// TTL returns the remaining time to live of a key
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.Client.TTL(ctx, key).Result()
}

// SetNX sets a key only if it doesn't already exist (atomic operation)
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.Client.SetNX(ctx, key, value, expiration).Result()
}

// SetXX sets a key only if it already exists (atomic operation)
func (c *Client) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.Client.SetXX(ctx, key, value, expiration).Result()
}

// MGet retrieves multiple values at once
func (c *Client) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return c.Client.MGet(ctx, keys...).Result()
}

// MSet sets multiple key-value pairs at once
func (c *Client) MSet(ctx context.Context, pairs ...interface{}) error {
	return c.Client.MSet(ctx, pairs...).Err()
}

// Keys finds all keys matching a pattern
func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.Client.Keys(ctx, pattern).Result()
}

// Scan iterates over keys matching a pattern (safer than Keys for large datasets)
func (c *Client) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.Client.Scan(ctx, cursor, match, count).Result()
}

// HSet sets a field in a hash
func (c *Client) HSet(ctx context.Context, key string, field string, value interface{}) error {
	return c.Client.HSet(ctx, key, field, value).Err()
}

// HGet retrieves a field from a hash
func (c *Client) HGet(ctx context.Context, key string, field string) (string, error) {
	return c.Client.HGet(ctx, key, field).Result()
}

// HGetAll retrieves all fields from a hash
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.Client.HGetAll(ctx, key).Result()
}

// HDel deletes one or more fields from a hash
func (c *Client) HDel(ctx context.Context, key string, fields ...string) error {
	return c.Client.HDel(ctx, key, fields...).Err()
}

// HMSet sets multiple fields in a hash at once
func (c *Client) HMSet(ctx context.Context, key string, pairs ...interface{}) error {
	return c.Client.HMSet(ctx, key, pairs...).Err()
}

// LPush prepends one or more values to a list
func (c *Client) LPush(ctx context.Context, key string, values ...interface{}) error {
	return c.Client.LPush(ctx, key, values...).Err()
}

// RPush appends one or more values to a list
func (c *Client) RPush(ctx context.Context, key string, values ...interface{}) error {
	return c.Client.RPush(ctx, key, values...).Err()
}

// LPop removes and returns the first element of a list
func (c *Client) LPop(ctx context.Context, key string) (string, error) {
	return c.Client.LPop(ctx, key).Result()
}

// RPop removes and returns the last element of a list
func (c *Client) RPop(ctx context.Context, key string) (string, error) {
	return c.Client.RPop(ctx, key).Result()
}

// LLen returns the length of a list
func (c *Client) LLen(ctx context.Context, key string) (int64, error) {
	return c.Client.LLen(ctx, key).Result()
}

// LRange returns elements from a list
func (c *Client) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.Client.LRange(ctx, key, start, stop).Result()
}

// SAdd adds one or more members to a set
func (c *Client) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.Client.SAdd(ctx, key, members...).Err()
}

// SMembers returns all members of a set
func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.Client.SMembers(ctx, key).Result()
}

// SIsMember checks if a value is a member of a set
func (c *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.Client.SIsMember(ctx, key, member).Result()
}

// SRem removes one or more members from a set
func (c *Client) SRem(ctx context.Context, key string, members ...interface{}) error {
	return c.Client.SRem(ctx, key, members...).Err()
}

// ZAdd adds one or more members with scores to a sorted set
func (c *Client) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return c.Client.ZAdd(ctx, key, members...).Err()
}

// ZRange returns elements from a sorted set by index range
func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.Client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeByScore returns elements from a sorted set by score range
func (c *Client) ZRangeByScore(ctx context.Context, key string, min, max string) ([]string, error) {
	opt := &redis.ZRangeBy{Min: min, Max: max}
	return c.Client.ZRangeByScore(ctx, key, opt).Result()
}

// ZRem removes one or more members from a sorted set
func (c *Client) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.Client.ZRem(ctx, key, members...).Err()
}

// Publish publishes a message to a channel
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.Client.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to one or more channels
func (c *Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.Client.Subscribe(ctx, channels...)
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.Client.Close()
}

// Stats returns connection pool statistics
func (c *Client) Stats() *redis.PoolStats {
	return c.Client.PoolStats()
}
