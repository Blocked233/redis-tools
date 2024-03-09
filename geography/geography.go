package geography

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// GeoIndex 是一个用于地理位置信息存储的封装结构
type GeoIndex struct {
	client *redis.Client
	key    string
}

// NewGeoIndex 创建一个新的 GeoIndex 实例
func NewGeoIndex(client *redis.Client, key string) *GeoIndex {
	return &GeoIndex{
		client: client,
		key:    key,
	}
}

// AddLocation 添加地理位置信息
func (g *GeoIndex) AddLocation(locations ...*redis.GeoLocation) error {
	_, err := g.client.GeoAdd(context.Background(), g.key, locations...).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetLocationsInRange 获取指定范围内的地理位置信息
func (g *GeoIndex) GetLocationsInRange(longitude, latitude float64, radius float64, unit string) ([]string, error) {
	results, err := g.client.GeoRadius(context.Background(), g.key, longitude, latitude, &redis.GeoRadiusQuery{
		Radius: radius,
		Unit:   unit,
	}).Result()
	if err != nil {
		return nil, err
	}

	var locations []string
	for _, result := range results {
		locations = append(locations, result.Name)
	}
	return locations, nil
}
