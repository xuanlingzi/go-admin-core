package config

import (
	"github.com/xuanlingzi/go-admin-core/storage"
	"github.com/xuanlingzi/go-admin-core/storage/cache"
)

type Cache struct {
	Redis  *RedisConnectOptions `json:"redis" yaml:"redis"`
	Memory interface{}
}

// CacheConfig cache配置
var CacheConfig = new(Cache)

// Setup 构造cache 顺序 redis > 其他 > memory
func (e Cache) Setup() storage.AdapterCache {
	if e.Redis != nil {
		options := e.Redis.GetRedisOptions()
		r := cache.NewRedis(cache.GetRedisClient(), options)
		return r
	}
	return cache.NewMemory()
}
