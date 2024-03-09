package counter

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// HyperLogLogCounter 是一个用于基数统计的封装结构
type HyperLogLogCounter struct {
	client *redis.Client
	key    string
}

// NewHyperLogLogCounter 创建一个新的 HyperLogLogCounter 实例
func NewHyperLogLogCounter(client *redis.Client, key string) *HyperLogLogCounter {
	return &HyperLogLogCounter{
		client: client,
		key:    key,
	}
}

// Add 将元素添加到 HyperLogLog 中
func (h *HyperLogLogCounter) Add(elements ...interface{}) error {
	_, err := h.client.PFAdd(context.Background(), h.key, elements...).Result()
	if err != nil {
		return err
	}
	return nil
}

// Count 统计 HyperLogLog 中元素的基数
func (h *HyperLogLogCounter) Count() (int64, error) {
	count, err := h.client.PFCount(context.Background(), h.key).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}
