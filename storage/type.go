package storage

import (
	"context"
	"time"

	"github.com/bsm/redislock"
)

const (
	PrefixKey = "__host"
)

type AdapterCache interface {
	String() string
	Get(key string) (string, error)
	Set(key string, val interface{}, expire int) error
	Del(key ...string) error
	DelPattern(pattern string) error
	HashKeys(hk string) ([]string, error)
	HashGet(hk, key string) (string, error)
	HashSet(hk, key string, val interface{}, expire int) error
	HashIncrease(hk, key string, val interface{}) (int64, error)
	HashDel(hk string, key ...string) error
	HashDelPattern(hk, pattern string) error
	Increase(key string, val interface{}) (int64, error)
	Decrease(key string, val interface{}) (int64, error)
	Expire(key string, dur time.Duration) error
	GetClient() interface{}
}

type AdapterQueue interface {
	String() string
	Append(message Messager) error
	Register(name string, f ConsumerFunc)
	Run()
	Shutdown()
}

type Messager interface {
	SetID(string)
	SetStream(string)
	SetValues(map[string]interface{})
	GetID() string
	GetStream() string
	GetValues() map[string]interface{}
	GetPrefix() string
	SetPrefix(string)
	SetErrorCount(count int)
	GetErrorCount() int
}

type ConsumerFunc func(Messager) error

type AdapterLocker interface {
	String() string
	Lock(key string, ttl int64, options *redislock.Options) (*redislock.Lock, error)
}

type AdapterFileStore interface {
	String() string
	Upload(ctx context.Context, name, location string) (string, error)
	GetClient() interface{}
}
