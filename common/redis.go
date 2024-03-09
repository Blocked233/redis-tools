package common

import (
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func init() {
	// 初始化Redis客户端
	RDB = redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Host + ":" + strconv.Itoa(Config.Redis.Port),
		Password: Config.Redis.Password,
		DB:       Config.Redis.DB,
	})
}
