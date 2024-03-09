package redislock

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisLock 表示一个Redis分布式锁
type RedisLock struct {
	client *redis.Client
	key    string
	value  string
}

// NewRedisLock 创建一个RedisLock实例
func NewRedisLock(client *redis.Client, key string, value string) *RedisLock {
	return &RedisLock{
		client: client,
		key:    key,
		value:  value,
	}
}

// TryLock 尝试获取锁
func (lock *RedisLock) TryLock(expiration time.Duration) (bool, error) {
	ok, err := lock.client.SetNX(ctx, lock.key, lock.value, expiration).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// Unlock 释放锁
func (lock *RedisLock) Unlock() error {
	// 使用Lua脚本来保证检查和删除是原子操作
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`
	res, err := lock.client.Eval(ctx, script, []string{lock.key}, lock.value).Result()
	if err != nil {
		return err
	}
	if res == int64(0) {
		return errors.New("unlock failed: lock value does not match")
	}
	return nil
}
