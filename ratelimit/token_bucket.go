package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// TokenBucket 使用Redis实现的令牌桶
type TokenBucket struct {
	rdb        *redis.Client
	key        string
	capacity   int64     // 令牌桶的容量
	fillRate   int64     // 每秒填充的令牌数
	lastRefill time.Time // 上次填充令牌的时间
}

// NewTokenBucket 创建一个新的TokenBucket实例
func NewTokenBucket(rdb *redis.Client, key string, capacity, fillRate int64) *TokenBucket {
	return &TokenBucket{
		rdb:      rdb,
		key:      key,
		capacity: capacity,
		fillRate: fillRate,
	}
}

// Refill 重新填充令牌桶
func (tb *TokenBucket) Refill() {
	ctx := context.Background()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	// 计算新的令牌数，并更新令牌桶
	tokensToAdd := int64(elapsed) * tb.fillRate
	tb.rdb.Eval(ctx, `
		local key = KEYS[1]
		local tokensToAdd = tonumber(ARGV[1])
		local capacity = tonumber(ARGV[2])

		local currentTokens = tonumber(redis.call("get", key) or "0")
		local newTokenCount = math.min(currentTokens + tokensToAdd, capacity)
		redis.call("set", key, newTokenCount)
	`, []string{tb.key}, tokensToAdd, tb.capacity)
}

// TryAcquire 尝试获取令牌
func (tb *TokenBucket) TryAcquire(count int64) bool {
	ctx := context.Background()

	tb.Refill() // 首先填充令牌桶

	result, err := tb.rdb.Eval(ctx, `
		local key = KEYS[1]
		local count = tonumber(ARGV[1])

		local currentTokens = tonumber(redis.call("get", key) or "0")
		if currentTokens >= count then
			redis.call("decrby", key, count)
			return 1
		else
			return 0
		end
	`, []string{tb.key}, count).Int()

	if err != nil {
		fmt.Println("Error acquiring token:", err)
		return false
	}

	return result == 1
}

// func main() {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	tb := NewTokenBucket(rdb, "myTokenBucket", 100, 10)

// 	// 尝试获取令牌
// 	if tb.TryAcquire(1) {
// 		fmt.Println("Token acquired!")
// 	} else {
// 		fmt.Println("Failed to acquire token.")
// 	}
// }
