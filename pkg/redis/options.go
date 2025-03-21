package redis

import (
	"time"

	"book-store/config"

	"github.com/redis/go-redis/v9"
)

type ClientOptions struct {
	clo *redis.Options
}

// NewClientOptions creates a new ClientOptions instance.
func NewClientOptions() ClientOptions {
	return ClientOptions{
		clo: &redis.Options{},
	}
}

func (co ClientOptions) SetOptions(opts config.RedisConfig) ClientOptions {
	co.clo.Addr = opts.RedisAddr
	co.clo.MinIdleConns = opts.MinIdleConns
	co.clo.PoolSize = opts.PoolSize
	co.clo.PoolTimeout = time.Duration(opts.PoolTimeout) * time.Second
	co.clo.Password = opts.Password
	co.clo.DB = opts.DB
	return co
}
