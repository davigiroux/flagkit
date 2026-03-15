package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/davigiroux/flagkit/api/internal/model"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func New(redisURL string, ttl time.Duration) (*Cache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return &Cache{client: client, ttl: ttl}, nil
}

func (c *Cache) GetFlag(ctx context.Context, key string) (*model.Flag, error) {
	data, err := c.client.Get(ctx, flagCacheKey(key)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var flag model.Flag
	if err := json.Unmarshal(data, &flag); err != nil {
		return nil, err
	}
	return &flag, nil
}

func (c *Cache) SetFlag(ctx context.Context, flag *model.Flag) error {
	data, err := json.Marshal(flag)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, flagCacheKey(flag.Key), data, c.ttl).Err()
}

func (c *Cache) InvalidateFlag(ctx context.Context, key string) error {
	return c.client.Del(ctx, flagCacheKey(key)).Err()
}

func flagCacheKey(key string) string {
	return "flag:" + key
}
