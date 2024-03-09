package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// RedisCounter 结构体包含Redis客户端和计数器的key
type RedisCounter struct {
	client *redis.Client
	key    string
}

// NewRedisCounter 创建一个新的RedisCounter实例
func NewRedisCounter(client *redis.Client, key string) *RedisCounter {
	return &RedisCounter{
		client: client,
		key:    key,
	}
}

// Increment 递增计数器的值
func (c *RedisCounter) Increment(ctx context.Context) (int64, error) {
	return c.client.Incr(ctx, c.key).Result()
}

// Get 获取当前计数器的值
func (c *RedisCounter) Get(ctx context.Context) (int64, error) {
	return c.client.Get(ctx, c.key).Int64()
}
