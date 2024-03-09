package counter

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// StringCounter 结构体包含Redis客户端和计数器的key
type StringCounter struct {
	client *redis.Client
	key    string
}

// NewStringCounter 创建一个新的StringCounter实例
func NewStringCounter(client *redis.Client, key string) *StringCounter {
	return &StringCounter{
		client: client,
		key:    key,
	}
}

// Increment 递增计数器的值
func (c *StringCounter) Increment(ctx context.Context) (int64, error) {
	return c.client.Incr(ctx, c.key).Result()
}

// Get 获取当前计数器的值
func (c *StringCounter) Get(ctx context.Context) (int64, error) {
	return c.client.Get(ctx, c.key).Int64()
}
