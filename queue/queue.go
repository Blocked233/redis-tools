package queue

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisQueue represents a queue structure based on Redis
type RedisQueue struct {
	client *redis.Client
	key    string
}

// NewRedisQueue creates and returns a new RedisQueue instance
func NewRedisQueue(client *redis.Client, key string) *RedisQueue {
	return &RedisQueue{
		client: client,
		key:    key,
	}
}

// Enqueue adds an element to the queue
func (q *RedisQueue) Enqueue(ctx context.Context, value string) error {
	return q.client.LPush(ctx, q.key, value).Err()
}

// Dequeue removes and returns an element from the queue
func (q *RedisQueue) Dequeue(ctx context.Context) (string, error) {
	return q.client.RPop(ctx, q.key).Result()
}

// DequeueBlocking removes and returns an element from the queue, blocking if the queue is empty
func (q *RedisQueue) DequeueBlocking(ctx context.Context, timeout time.Duration) (string, error) {
	result, err := q.client.BRPop(ctx, timeout, q.key).Result()
	if err != nil {
		return "", err
	}
	if len(result) != 2 {
		return "", redis.Nil
	}
	return result[1], nil
}
