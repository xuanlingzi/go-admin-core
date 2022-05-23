package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// NewRedis redis模式
func NewRedis(client *redis.Client, options *redis.Options) (*Redis, error) {
	if client == nil {
		client = redis.NewClient(options)
	}
	r := &Redis{
		client: client,
	}
	err := r.connect()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Redis cache implement
type Redis struct {
	client *redis.Client
}

func (*Redis) String() string {
	return "redis"
}

// connect connect test
func (r *Redis) connect() error {
	var err error
	_, err = r.client.Ping(context.Background()).Result()
	return err
}

// Get from key
func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

// Set value with key and expire time
func (r *Redis) Set(key string, val interface{}, expire int) error {
	return r.client.Set(context.Background(), key, val, time.Duration(expire)*time.Second).Err()
}

// Del delete key in redis
func (r *Redis) Del(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

// HashKeys from key
func (r *Redis) HashKeys(hk string) ([]string, error) {
	return r.client.HKeys(context.Background(), hk).Result()
}

// HashGet from key
func (r *Redis) HashGet(hk, key string) (string, error) {
	return r.client.HGet(context.Background(), hk, key).Result()
}

// HashSet delete key in specify redis's hashtable
func (r *Redis) HashSet(hk, key string, val interface{}, _ int) error {
	return r.client.HSet(context.Background(), hk, key, val).Err()
}

// HashDel delete key in specify redis's hashtable
func (r *Redis) HashDel(hk, key string) error {
	return r.client.HDel(context.Background(), hk, key).Err()
}

// Increase
func (r *Redis) Increase(key string) error {
	return r.client.Incr(context.Background(), key).Err()
}

func (r *Redis) Decrease(key string) error {
	return r.client.Decr(context.Background(), key).Err()
}

// Set ttl
func (r *Redis) Expire(key string, dur time.Duration) error {
	return r.client.Expire(context.Background(), key, dur).Err()
}

// GetClient 暴露原生client
func (r *Redis) GetClient() *redis.Client {
	return r.client
}
