package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

type Redis struct {
	redis  redis.UniversalClient
	prefix string
}

var (
	ErrKeyNotFound = redis.Nil
)

func New(redis redis.UniversalClient, prefix string) *Redis {
	return &Redis{
		redis:  redis,
		prefix: prefix,
	}
}

// Client returns the underlying redis client
func (r *Redis) Client() redis.UniversalClient {
	return r.redis
}

// Ping tests the connection
func (r *Redis) Ping(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}

func (r *Redis) prefixed(key string) string {
	if r.prefix != "" {
		return r.prefix + ":" + key
	}
	return key
}

// SET command
func (r *Redis) Set(ctx context.Context, key, value string, expire time.Duration) error {
	_, err := r.redis.SetEX(ctx, r.prefixed(key), value, expire).Result()
	return err
}

func (r *Redis) SetBulk(ctx context.Context, keys []string, values []string, expire time.Duration) []error {
	cmd, _ := r.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i, key := range keys {
			pipe.Set(ctx, key, values[i], 7*time.Minute)
		}
		return nil
	})

	errs := []error{}
	for _, cmd := range cmd {
		_, err := cmd.(*redis.StatusCmd).Result()
		errs = append(errs, err)
	}
	return errs
}

func (r *Redis) SetBulkStruct(ctx context.Context, keys []string, values []any, expire time.Duration) []error {
	cmd, _ := r.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i, key := range keys {
			j, err := msgpack.Marshal(values[i])
			if err != nil {
				return err
			}
			pipe.Set(ctx, key, string(j), 7*time.Minute)
		}
		return nil
	})

	errs := []error{}
	for _, cmd := range cmd {
		_, err := cmd.(*redis.StatusCmd).Result()
		errs = append(errs, err)
	}
	return errs
}

func (r *Redis) SetBulkCompound(ctx context.Context, prefix string, keys []string, values []string, expire time.Duration) []error {
	for i, key := range keys {
		keys[i] = prefix + ":" + key
	}
	return r.SetBulk(ctx, keys, values, expire)
}

func (r *Redis) SetBulkCompoundStruct(ctx context.Context, prefix string, keys []string, values []any, expire time.Duration) []error {
	for i, key := range keys {
		keys[i] = prefix + ":" + key
	}
	return r.SetBulkStruct(ctx, keys, values, expire)
}

func (r *Redis) SetStruct(ctx context.Context, key string, value any, expire time.Duration) error {
	j, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, string(j), expire)
}

// GET command
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return val, ErrKeyNotFound
	}
	return val, err
}

func (r *Redis) GetBulk(ctx context.Context, keys []string) ([]string, []error) {
	cmd, _ := r.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(ctx, key)
		}
		return nil
	})

	vals, errs := []string{}, []error{}
	for _, cmd := range cmd {
		val, err := cmd.(*redis.StringCmd).Result()
		if err == redis.Nil {
			err = ErrKeyNotFound
		}
		vals = append(vals, val)
		errs = append(errs, err)
	}

	return vals, errs
}

func (r *Redis) GetBulkCompound(ctx context.Context, prefix string, keys []string) ([]string, []error) {
	for i, key := range keys {
		keys[i] = prefix + ":" + key
	}
	return r.GetBulk(ctx, keys)
}

func (r *Redis) GetBulkCompoundStruct(ctx context.Context, prefix string, keys []string, target []any) ([]error, error) {
	vals, errs := r.GetBulkCompound(ctx, prefix, keys)
	for i, val := range vals {
		if errs[i] != nil {
			continue
		}
		err := msgpack.Unmarshal([]byte(val), target[i])
		if err != nil {
			return nil, err
		}
	}
	return errs, nil
}

func (r *Redis) GetStruct(ctx context.Context, key string, i any) error {
	val, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal([]byte(val), i)
}

func (r *Redis) GetCompound(ctx context.Context, key string, key2 string) (string, error) {
	val, err := r.Get(ctx, key+":"+key2)
	if err == redis.Nil {
		return val, ErrKeyNotFound
	}
	return val, err
}

func (r *Redis) GetCompoundStruct(ctx context.Context, key string, key2 string, i any) error {
	val, err := r.GetCompound(ctx, key, key2)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal([]byte(val), i)
}

// HGET command
func (r *Redis) HGet(ctx context.Context, key string, unique string) (string, error) {
	val, err := r.redis.HGet(ctx, key, unique).Result()
	if err == redis.Nil {
		return val, ErrKeyNotFound
	}
	return val, err
}

func (r *Redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	cmd := r.redis.HGetAll(ctx, key)
	return cmd.Result()
}

func (r *Redis) Del(ctx context.Context, key string) (int64, error) {
	return r.redis.Del(ctx, key).Result()
}

func (r *Redis) DelCompound(ctx context.Context, key, key2 string) (int64, error) {
	cmd := r.redis.Del(ctx, key+":"+key2)
	return cmd.Result()
}

// HDel command
func (r *Redis) HDel(ctx context.Context, key, value string) (int64, error) {
	cmd := r.redis.HDel(ctx, key, value)
	return cmd.Result()
}

func (r *Redis) AddInSet(ctx context.Context, key, value string) error {
	_, err := r.redis.SAdd(ctx, key, value).Result()
	return err
}

// GetSetMembers run SMEMBERS command tp get all members in a set
func (r *Redis) GetSetMembers(ctx context.Context, key string) ([]string, error) {
	return r.redis.SMembers(ctx, key).Result()
}

// HIncrBy increment hash value by n integer
func (r *Redis) HIncrBy(ctx context.Context, key, field string, n int64) (int64, error) {
	cmd := r.redis.HIncrBy(ctx, key, field, n)
	return cmd.Result()
}

// IncrBy increment key by n integer
func (r *Redis) IncrBy(ctx context.Context, key string, n int64) (int64, error) {
	cmd := r.redis.IncrBy(ctx, key, n)
	return cmd.Result()
}

// Expire set time to live of key by duration
func (r *Redis) Expire(ctx context.Context, key string, duration time.Duration) (bool, error) {
	cmd := r.redis.Expire(ctx, key, duration)
	return cmd.Result()
}

// HSetStruct run HSET command
func (r *Redis) HSetStruct(ctx context.Context, key, field string, value any) error {
	j, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}

	cmd := r.redis.HSet(ctx, key, field, string(j))
	_, err = cmd.Result()
	return err
}

// HGetStruct function
func (r *Redis) HGetStruct(ctx context.Context, key, field string, value any) error {
	cmd := r.redis.HGet(ctx, key, field)
	val, err := cmd.Result()
	if err == redis.Nil {
		return ErrKeyNotFound
	}

	return msgpack.Unmarshal([]byte(val), value)
}

// GetExStruct function
func (r *Redis) GetExStruct(ctx context.Context, key string, i any, ttl ...time.Duration) error {
	val, err := r.GetEx(ctx, key, ttl...)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal([]byte(val), i)
}

// GetExCompound run GET redis command with two keys  (concat with ":")
func (r *Redis) GetExCompound(ctx context.Context, key string, key2 string, ttl ...time.Duration) (string, error) {
	dur := 7 * time.Minute
	if len(ttl) > 0 {
		dur = ttl[0]
	}

	return r.GetEx(ctx, key+":"+key2, dur)
}

// GetEx function
func (r *Redis) GetEx(ctx context.Context, key string, ttl ...time.Duration) (string, error) {
	dur := 7 * time.Minute
	if len(ttl) > 0 {
		dur = ttl[0]
	}
	val, err := r.redis.Get(ctx, r.prefixed(key)).Result()
	if err != nil {
		return val, err
	}

	// TODO: temporary workaround for redis proxy
	_, err = r.redis.Expire(ctx, r.prefixed(key), dur).Result()

	return val, err
}

// GetExCompoundStruct run GET redis command with two keys (concat with ":") to struct
func (r *Redis) GetExCompoundStruct(ctx context.Context, key string, key2 string, i any, ttl ...time.Duration) (err error) {
	val, err := r.GetExCompound(ctx, key, key2, ttl...)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal([]byte(val), i)
}

// SetExCompound SET with two keys (concat with ":")
// default ttl 7 minutes
func (r *Redis) SetExCompound(ctx context.Context, key, key2, value string, ttl ...time.Duration) error {
	dur := 7 * time.Minute
	if len(ttl) > 0 {
		dur = ttl[0]
	}
	return r.Set(ctx, key+":"+key2, value, dur)
}

// SetExCompoundStruct SET from struct with two keys (concat with ":")
// default ttl 7 minutes
func (r *Redis) SetExCompoundStruct(ctx context.Context, key, key2 string, value any, ttl ...time.Duration) (err error) {
	j, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	return r.SetExCompound(ctx, key, key2, string(j), ttl...)
}
