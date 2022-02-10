package storage

import (
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
	Del(key string) error
	HashKeys(hk string) ([]string, error)
	HashGet(hk, key string) (string, error)
	HashSet(hk, key string, val interface{}, expire int) error
	HashDel(hk, key string) error
	Increase(key string) error
	Decrease(key string) error
	Expire(key string, dur time.Duration) error
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
}

type ConsumerFunc func(Messager) error

type AdapterLocker interface {
	String() string
	Lock(key string, ttl int64, options *redislock.Options) (*redislock.Lock, error)
}

type AdapterFileStore interface {
	String() string
	Upload(name, location string) (string, error)
}
