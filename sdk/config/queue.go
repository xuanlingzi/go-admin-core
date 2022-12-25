package config

import (
	"github.com/go-admin-team/redisqueue/v2"
	"github.com/go-redis/redis/v9"
	"github.com/xuanlingzi/go-admin-core/storage"
	"github.com/xuanlingzi/go-admin-core/storage/queue"
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
		client := GetRedisClient()
		if client == nil {
			options, err := e.Redis.RedisConnectOptions.GetRedisOptions()
			if err != nil {
				return nil, err
			}
			client = redis.NewClient(options)
			_redis = client
		}
		redisOption := &redisqueue.RedisOptions{
			Network:    client.Options().Network,
			Addr:       client.Options().Addr,
			Username:   client.Options().Username,
			Password:   client.Options().Password,
			DB:         client.Options().DB,
			MaxRetries: client.Options().MaxRetries,
			PoolSize:   client.Options().PoolSize,
		}
		e.Redis.Producer = &redisqueue.ProducerOptions{
			RedisOptions: redisOption,
		} // .RedisClient = client
		e.Redis.Consumer = &redisqueue.ConsumerOptions{
			RedisOptions: redisOption,
		} // .RedisClient = client
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
