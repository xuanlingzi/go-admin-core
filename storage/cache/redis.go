package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
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
func NewRedis(client *redis.Client, options *redis.Options) (*Redis, error) {
	if client == nil {
		client = redis.NewClient(options)
		_redis = client
	}
	r := &Redis{
		client: client,
	}
	if err := r.connect(); err != nil {
		return nil, fmt.Errorf("redis cache init: %w", err)
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

// DelPattern delete key in redis using SCAN to avoid blocking
func (r *Redis) DelPattern(pattern string) error {
	ctx := context.Background()
	var cursor uint64
	var deletedCount int64

	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("scan keys: %w", err)
		}

		if len(keys) > 0 {
			// 使用 Unlink 而不是 Del，更快（异步删除）
			deleted, err := r.client.Unlink(ctx, keys...).Result()
			if err != nil {
				return fmt.Errorf("delete keys: %w", err)
			}
			deletedCount += deleted
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

// HashKeys from key
func (r *Redis) HashKeys(hk string) ([]string, error) {
	return r.client.HKeys(context.TODO(), hk).Result()
}

// HashGet from key
func (r *Redis) HashGet(hk, key string) (string, error) {
	return r.client.HGet(context.TODO(), hk, key).Result()
}

// HashSet set key in specify redis's hashtable
func (r *Redis) HashSet(hk, key string, val interface{}, _ int) error {
	return r.client.HSet(context.TODO(), hk, key, val).Err()
}

// HashIncrease increase key in specify redis's hashtable
func (r *Redis) HashIncrease(hk, key string, val interface{}) (int64, error) {
	return r.client.HIncrBy(context.TODO(), hk, key, cast.ToInt64(val)).Result()
}

// HashDel delete key in specify redis's hashtable
func (r *Redis) HashDel(hk string, key ...string) error {
	return r.client.HDel(context.TODO(), hk, key...).Err()
}

// HashDelPattern delete key in specify redis's hashtable using HSCAN
func (r *Redis) HashDelPattern(hk, pattern string) error {
	ctx := context.Background()
	var cursor uint64
	var delKeys []string

	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.HScan(ctx, hk, cursor, pattern, 100).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil
			}
			return fmt.Errorf("hscan keys: %w", err)
		}

		// HScan returns key-value pairs, so we need to extract keys only
		for i := 0; i < len(keys); i += 2 {
			delKeys = append(delKeys, keys[i])
		}

		if cursor == 0 {
			break
		}
	}

	if len(delKeys) > 0 {
		return r.client.HDel(ctx, hk, delKeys...).Err()
	}
	return nil
}

func (r *Redis) Increase(key string, val interface{}) (int64, error) {
	return r.client.IncrBy(context.TODO(), key, cast.ToInt64(val)).Result()
}

func (r *Redis) Decrease(key string, val interface{}) (int64, error) {
	return r.client.DecrBy(context.TODO(), key, cast.ToInt64(val)).Result()
}

// Expire Set ttl
func (r *Redis) Expire(key string, dur time.Duration) error {
	return r.client.Expire(context.TODO(), key, dur).Err()
}

// GetClient 暴露原生client
func (r *Redis) GetClient() interface{} {
	return r.client
}
