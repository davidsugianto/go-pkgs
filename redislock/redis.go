package redislock

import (
	goredislib "github.com/go-redis/redis/v8"
	redsynclib "github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type RedisDriver struct {
	GoRedisClient []goredislib.UniversalClient
}

type redsyncWrap struct {
	rs          *redsynclib.Redsync
	hasClients  bool
}

type IMutex interface {
	NewMutexW(mutexOpt MutexOpt) (IMutexDistLock, error)
}

// New : init Redis Client to store the lock data
func New(opt RedisDriver) IMutex {
	var rs []redis.Pool
	for _, val := range opt.GoRedisClient {
		if val != nil {
			rs = append(rs, goredis.NewPool(val))
		}
	}

	return &redsyncWrap{
		rs:         redsynclib.New(rs...),
		hasClients: len(rs) > 0,
	}
}

// GetRedsync returns the underlying redsync instance
func (rw *redsyncWrap) GetRedsync() *redsynclib.Redsync {
	return rw.rs
}
