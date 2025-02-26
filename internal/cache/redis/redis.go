package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tclemos/go-expert-rate-limiter/config"
)

type RedisCache struct {
	client redis.Client
}

func NewRedisCache(ctx context.Context, config config.WebServerConfig) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPass,
		DB:       config.RedisDB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic("failed to connect to redis:" + err.Error())
	}

	return &RedisCache{
		client: *client,
	}
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, duration time.Duration) error {
	return r.client.Set(ctx, key, value, duration).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
