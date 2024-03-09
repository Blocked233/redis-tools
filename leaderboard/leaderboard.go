package leaderboard

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// MixedLeaderboard 表示一个混合排行榜, 前32位是freshness，后32位是hotness
type MixedLeaderboard struct {
	rdb *redis.Client
	key string
}

func NewMixedLeaderboard(client *redis.Client, key string) *MixedLeaderboard {
	return &MixedLeaderboard{
		rdb: client,
		key: key,
	}
}

func (lb *MixedLeaderboard) AddItem(ctx context.Context, item string, timestamp time.Time, hotness int64) error {
	freshness := timestamp.Unix()
	// 计算分数，将freshness和hotness组合成一个64位的整数后，转换为float64
	score := float64((uint64(freshness) << 32) | uint64(hotness&0xFFFFFFFF))
	return lb.rdb.ZAdd(ctx, lb.key, &redis.Z{
		Score:  score,
		Member: item,
	}).Err()
}

func (lb *MixedLeaderboard) GetItems(ctx context.Context, startTime, endTime time.Time, offset, count int64) ([]redis.Z, error) {
	minFreshness := startTime.Unix()
	maxFreshness := endTime.Unix()

	minScore := float64(uint64(minFreshness) << 32)       // 最小分数，热度为0
	maxScore := float64((uint64(maxFreshness)+1)<<32) - 1 // 最大分数，热度为最大可能值

	return lb.rdb.ZRevRangeByScoreWithScores(ctx, lb.key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%f", minScore),
		Max:    fmt.Sprintf("%f", maxScore),
		Offset: offset,
		Count:  count,
	}).Result()
}

type Leaderboard struct {
	rdb    *redis.Client
	key    string
	lbType string
}

type LeaderboardType string

const (
	Latest  LeaderboardType = "l"
	Hottest LeaderboardType = "h"
)

func NewLeaderboard(client *redis.Client, key string, lbType LeaderboardType) *Leaderboard {
	return &Leaderboard{
		rdb:    client,
		key:    key,
		lbType: string(lbType),
	}
}

func (lb *Leaderboard) AddItem(ctx context.Context, item string, score int64) error {
	return lb.rdb.ZAdd(ctx, lb.key+":"+lb.lbType, &redis.Z{
		Score:  float64(score),
		Member: item,
	}).Err()
}

func (lb *Leaderboard) GetItems(ctx context.Context, start, stop int64) ([]string, error) {
	return lb.rdb.ZRevRange(ctx, lb.key+":"+lb.lbType, start, stop).Result()
}
