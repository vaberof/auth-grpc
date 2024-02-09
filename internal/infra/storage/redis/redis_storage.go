package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExpired  = errors.New("key has expired")
)

type RedisStorage struct {
	client *redis.Client

	// TODO: add logger
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func (i *RedisStorage) Set(ctx context.Context, key, value string, exp time.Duration) error {
	err := i.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (i *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	val, err := i.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
