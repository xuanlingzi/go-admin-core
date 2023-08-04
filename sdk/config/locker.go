package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/xuanlingzi/go-admin-core/storage"
	"github.com/xuanlingzi/go-admin-core/storage/cache"
	"github.com/xuanlingzi/go-admin-core/storage/locker"
)

var LockerConfig = new(Locker)

type Locker struct {
	Redis *RedisConnectOptions `json:"redis" yaml:"redis"`
}

// Empty 空设置
func (e Locker) Empty() bool {
	return e.Redis == nil
}

// Setup 启用顺序 redis > 其他 > memory
func (e Locker) Setup() storage.AdapterLocker {
	if e.Redis != nil {
		client := cache.GetRedisClient()
		if client == nil {
			options := e.Redis.GetRedisOptions()
			client = redis.NewClient(options)
		}
		return locker.NewRedis(client)
	}
	return nil
}
