package redis_storage

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"main/internal/storage"
	"time"
)

var _ storage.Storage = (*RedisStorage)(nil)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr string) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	pong, err := client.Ping(context.Background()).Result()

	if err != nil || pong != "PONG" {
		return nil
	}

	return &RedisStorage{client: client}
}

func (r *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("key not found")
	}
	return value, err
}

func (r *RedisStorage) Set(ctx context.Context, key, value string, lifetimeSeconds int) error {
	var expiration time.Duration

	if lifetimeSeconds > 0 {
		expiration = time.Duration(lifetimeSeconds) * time.Second
	} else {
		expiration = time.Duration(0) * time.Second
	}

	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisStorage) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisStorage) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
