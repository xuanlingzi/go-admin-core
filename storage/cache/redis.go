package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var _redis *redis.Client

// GetRedisClient 获取redis客户端
func GetRedisClient() *redis.Client {
	return _redis
}

// SetRedisClient 设置redis客户端
func SetRedisClient(c *redis.Client) {
	if _redis != nil && _redis != c {
		_redis.Shutdown(context.TODO())
	}
	_redis = c
}

// NewRedis redis模式
func NewRedis(client *redis.Client, options *redis.Options) *Redis {
	if client == nil {
		client = redis.NewClient(options)
		_redis = client
	}
	r := &Redis{
		client: client,
	}
	err := r.connect()
	if err != nil {
		panic(fmt.Sprintf("Redis cache init error %s", err.Error()))
	}
	return r
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
	_, err = r.client.Ping(context.TODO()).Result()
	return err
}

// Get from key
func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.TODO(), key).Result()
}

// Set value with key and expire time
func (r *Redis) Set(key string, val interface{}, expire int) error {
	return r.client.Set(context.TODO(), key, val, time.Duration(expire)*time.Second).Err()
}

// Del delete key in redis
func (r *Redis) Del(key ...string) error {
	return r.client.Del(context.TODO(), key...).Err()
}

// DelPattern delete key in redis
func (r *Redis) DelPattern(pattern string) error {
	keys, err := r.client.Keys(context.Background(), pattern).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	return r.client.Del(context.Background(), keys...).Err()
}

// HashKeys from key
func (r *Redis) HashKeys(hk string) ([]string, error) {
	return r.client.HKeys(context.TODO(), hk).Result()
}

// HashGet from key
func (r *Redis) HashGet(hk, key string) (string, error) {
	return r.client.HGet(context.TODO(), hk, key).Result()
}

// HashSet delete key in specify redis's hashtable
func (r *Redis) HashSet(hk, key string, val interface{}, _ int) error {
	return r.client.HSet(context.TODO(), hk, key, val).Err()
}

// HashDel delete key in specify redis's hashtable
func (r *Redis) HashDel(hk string, key ...string) error {
	return r.client.HDel(context.TODO(), hk, key...).Err()
}

// HashDelPattern delete key in specify redis's hashtable
func (r *Redis) HashDelPattern(hk, pattern string) error {
	keys, err := r.client.HKeys(context.Background(), hk).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	var delKeys []string
	for _, key := range keys {
		if strings.Contains(key, pattern) {
			delKeys = append(delKeys, key)
		}
	}
	return r.client.HDel(context.Background(), hk, delKeys...).Err()
}

func (r *Redis) Increase(key string) error {
	return r.client.Incr(context.TODO(), key).Err()
}

func (r *Redis) Decrease(key string) error {
	return r.client.Decr(context.TODO(), key).Err()
}

// Set ttl
func (r *Redis) Expire(key string, dur time.Duration) error {
	return r.client.Expire(context.TODO(), key, dur).Err()
}

// GetClient 暴露原生client
func (r *Redis) GetClient() interface{} {
	return r.client
}
