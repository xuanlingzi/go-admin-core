package runtime

import (
	"time"

	"github.com/xuanlingzi/go-admin-core/storage"
)

const (
	intervalTenant = ""
)

// NewCache 创建对应上下文缓存
func NewCache(prefix string, store storage.AdapterCache) storage.AdapterCache {
	return &Cache{
		prefix: prefix,
		store:  store,
	}
}

type Cache struct {
	prefix string
	store  storage.AdapterCache
}

// String string输出
func (e *Cache) String() string {
	if e.store == nil {
		return ""
	}
	return e.store.String()
}

// SetPrefix 设置前缀
func (e *Cache) SetPrefix(prefix string) {
	e.prefix = prefix
}

// Connect 初始化
func (e Cache) Connect() error {
	return nil
	//return e.store.Connect()
}

// Get val in cache
func (e Cache) Get(key string) (string, error) {
	return e.store.Get(e.prefix + intervalTenant + key)
}

// Set val in cache
func (e Cache) Set(key string, val interface{}, expire int) error {
	return e.store.Set(e.prefix+intervalTenant+key, val, expire)
}

// Del delete key in cache
func (e Cache) Del(key ...string) error {
	var keys []string
	for _, k := range key {
		keys = append(keys, e.prefix+intervalTenant+k)
	}
	return e.store.Del(keys...)
}

// DelPattern delete key in cache
func (e Cache) DelPattern(pattern string) error {
	return e.store.DelPattern(e.prefix + intervalTenant + pattern)
}

// HashKeys get val in hashtable cache
func (e Cache) HashKeys(hk string) ([]string, error) {
	return e.store.HashKeys(e.prefix + intervalTenant + hk)
}

// HashGet get val in hashtable cache
func (e Cache) HashGet(hk, key string) (string, error) {
	return e.store.HashGet(e.prefix+intervalTenant+hk, key)
}

// HashSet set val in hashtable cache
func (e Cache) HashSet(hk, key string, val interface{}, expire int) error {
	return e.store.HashSet(e.prefix+intervalTenant+hk, key, val, expire)
}

func (e Cache) HashIncrease(hk, key string, val interface{}) (int64, error) {
	return e.store.HashIncrease(e.prefix+intervalTenant+hk, key, val)
}

// HashDel delete one key:value pair in hashtable cache
func (e Cache) HashDel(hk string, key ...string) error {
	return e.store.HashDel(e.prefix+intervalTenant+hk, key...)
}

// HashDelPattern delete one key:value pair in hashtable cache
func (e Cache) HashDelPattern(hk, pattern string) error {
	return e.store.HashDelPattern(e.prefix+intervalTenant+hk, pattern)
}

// Increase value
func (e Cache) Increase(key string, val interface{}) (int64, error) {
	return e.store.Increase(e.prefix+intervalTenant+key, val)
}

func (e Cache) Decrease(key string, val interface{}) (int64, error) {
	return e.store.Decrease(e.prefix+intervalTenant+key, val)
}

func (e Cache) Expire(key string, dur time.Duration) error {
	return e.store.Expire(e.prefix+intervalTenant+key, dur)
}

func (e Cache) GetClient() interface{} {
	return e.store.GetClient()
}
