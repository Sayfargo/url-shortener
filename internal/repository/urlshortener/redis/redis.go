package repository_urlshortener_redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

var (
	ErrNotExists = errors.New("url doesn't exists")
)

func (r *RedisRepository) Set(ctx context.Context, originalURL, shortCode string) error {

	if err := r.client.Set(ctx, shortCode, originalURL, time.Hour).Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return context.DeadlineExceeded
		}
		return fmt.Errorf("failed to set : %w", err)
	}

	return nil
}

func (r *RedisRepository) Get(ctx context.Context, shortCode string) (string, error) {

	url, err := r.client.Get(ctx, shortCode).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrNotExists
		} else if errors.Is(err, context.DeadlineExceeded) {
			return "", context.DeadlineExceeded
		}
		return "", fmt.Errorf("failed to get : %w", err)
	}

	return url, nil
}
