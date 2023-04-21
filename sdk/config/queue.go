package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/xuanlingzi/go-admin-core/storage"
	"github.com/xuanlingzi/go-admin-core/storage/queue"
	"github.com/xuanlingzi/redisqueue/v2"
	"time"
)

type Queue struct {
	Redis  *QueueRedis  `json:"redis" yaml:"redis"`
	Memory *QueueMemory `json:"memory" yaml:"memory"`
	NSQ    *QueueNSQ    `json:"nsq" yaml:"nsq"`
}

type QueueRedis struct {
	RedisConnectOptions
	Producer *redisqueue.ProducerOptions
	Consumer *redisqueue.ConsumerOptions
}

type QueueMemory struct {
	PoolSize uint `json:"pool_size" yaml:"pool_size"`
}

type QueueNSQ struct {
	NSQOptions
	ChannelPrefix string `json:"channel_prefix" yaml:"channel_prefix"`
}

var QueueConfig = new(Queue)

// Empty 空设置
func (e Queue) Empty() bool {
	return e.Memory == nil && e.Redis == nil && e.NSQ == nil
}

// Setup 启用顺序 redis > 其他 > memory
func (e Queue) Setup() (storage.AdapterQueue, error) {
	if e.Redis != nil {
		e.Redis.Consumer.ReclaimInterval = e.Redis.Consumer.ReclaimInterval * time.Second
		e.Redis.Consumer.BlockingTimeout = e.Redis.Consumer.BlockingTimeout * time.Second
		e.Redis.Consumer.VisibilityTimeout = e.Redis.Consumer.VisibilityTimeout * time.Second
		universalOption := redis.UniversalOptions{
			Addrs:      []string{e.Redis.Addr},
			ClientName: ApplicationConfig.Name,
			DB:         e.Redis.DB,
			Username:   e.Redis.Username,
			Password:   e.Redis.Password,
			MaxRetries: e.Redis.MaxRetries,
			PoolSize:   e.Redis.PoolSize,
		}
		client := redis.NewUniversalClient(&universalOption)
		redisOption := &redisqueue.RedisOptions{}
		e.Redis.Producer = &redisqueue.ProducerOptions{
			RedisClient:  client,
			RedisOptions: redisOption,
		}
		e.Redis.Consumer = &redisqueue.ConsumerOptions{
			RedisClient:  client,
			RedisOptions: redisOption,
		}
		return queue.NewRedis(e.Redis.Producer, e.Redis.Consumer)
	}
	if e.NSQ != nil {
		cfg, err := e.NSQ.GetNSQOptions()
		if err != nil {
			return nil, err
		}
		return queue.NewNSQ(e.NSQ.Addresses, cfg, e.NSQ.ChannelPrefix)
	}
	return queue.NewMemory(e.Memory.PoolSize), nil
}
