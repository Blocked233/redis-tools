package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// LeakyBucket 使用Redis实现的漏桶
type LeakyBucket struct {
	rdb          *redis.Client
	key          string
	capacity     int64     // 桶的容量
	leakRate     int64     // 每秒漏水速率
	lastLeakTime time.Time // 上次漏水时间
}

// NewLeakyBucket 创建一个新的LeakyBucket实例
func NewLeakyBucket(rdb *redis.Client, key string, capacity, leakRate int64) *LeakyBucket {
	return &LeakyBucket{
		rdb:          rdb,
		key:          key,
		capacity:     capacity,
		leakRate:     leakRate,
		lastLeakTime: time.Now(),
	}
}

// Leak 漏水操作
func (lb *LeakyBucket) Leak() {
	ctx := context.Background()

	now := time.Now()
	elapsed := now.Sub(lb.lastLeakTime).Seconds()
	lb.lastLeakTime = now

	waterToLeak := int64(elapsed) * lb.leakRate
	lb.rdb.Eval(ctx, `
		local key = KEYS[1]
		local waterToLeak = tonumber(ARGV[1])

		local currentWaterLevel = tonumber(redis.call("get", key) or "0")
		local newWaterLevel = math.max(0, currentWaterLevel - waterToLeak)
		redis.call("set", key, newWaterLevel)
	`, []string{lb.key}, waterToLeak)
}

// TryFill 尝试向桶中填充水
func (lb *LeakyBucket) TryFill(count int64) bool {
	ctx := context.Background()

	lb.Leak() // 首先执行漏水操作

	result, err := lb.rdb.Eval(ctx, `
		local key = KEYS[1]
		local count = tonumber(ARGV[1])
		local capacity = tonumber(ARGV[2])

		local currentWaterLevel = tonumber(redis.call("get", key) or "0")
		if currentWaterLevel + count <= capacity then
			redis.call("incrby", key, count)
			return 1
		else
			return 0
		end
	`, []string{lb.key}, count, lb.capacity).Int()

	if err != nil {
		fmt.Println("Error filling bucket:", err)
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

// 	lb := NewLeakyBucket(rdb, "myLeakyBucket", 100, 10)

// 	// 尝试向桶中填充水
// 	if lb.TryFill(1) {
// 		fmt.Println("Water filled!")
// 	} else {
// 		fmt.Println("Failed to fill water.")
// 	}
// }
