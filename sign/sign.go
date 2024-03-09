package sign

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// BitMapCounter 是一个用于二值状态统计的封装结构
type BitMapCounter struct {
	client *redis.Client
	key    string
}

// NewBitMapCounter 创建一个新的 BitMapCounter 实例
func NewBitMapCounter(client *redis.Client, key string) *BitMapCounter {
	return &BitMapCounter{
		client: client,
		key:    key,
	}
}

// SignIn 将用户的签到状态设置为已签到
func (b *BitMapCounter) SignIn(userID int64) error {
	_, err := b.client.SetBit(context.Background(), b.key, userID, 1).Result()
	if err != nil {
		return err
	}
	return nil
}

// IsSignedIn 检查用户是否已签到
func (b *BitMapCounter) IsSignedIn(userID int64) (bool, error) {
	val, err := b.client.GetBit(context.Background(), b.key, userID).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

// CountSignedUsers 统计已签到用户的总数
func (b *BitMapCounter) CountSignedUsers() (int64, error) {
	count, err := b.client.BitCount(context.Background(), b.key, nil).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}