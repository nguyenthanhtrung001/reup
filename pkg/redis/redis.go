package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func Connect(opts ClientOptions) (Client, error) {
	cl := redis.NewClient(opts.clo)
	return &redisClient{cl: cl}, nil
}

type Database interface {
	Client() Client
}

type Client interface {
	Disconnect() error
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value interface{}, expiration int) error
}

type redisClient struct {
	cl *redis.Client
}

func (rc *redisClient) Disconnect() error {
	return rc.cl.Close()
}

func (rc *redisClient) Get(ctx context.Context, key string) ([]byte, error) {
	return rc.cl.Get(ctx, key).Bytes()
}

func (rc *redisClient) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	return rc.cl.Set(ctx, key, value, time.Second*time.Duration(expiration)).Err()
}
