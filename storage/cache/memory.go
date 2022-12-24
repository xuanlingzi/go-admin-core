package cache

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type item struct {
	Value   string
	Expired time.Time
}

// NewMemory memory模式
func NewMemory() *Memory {
	return &Memory{
		items: new(sync.Map),
	}
}

type Memory struct {
	items *sync.Map
	mutex sync.RWMutex
}

func (*Memory) String() string {
	return "memory"
}

func (m *Memory) connect() {
}

func (m *Memory) Get(key string) (string, error) {
	item, err := m.getItem(key)
	if err != nil || item == nil {
		return "", err
	}
	return item.Value, nil
}

func (m *Memory) getItem(key string) (*item, error) {
	var err error
	i, ok := m.items.Load(key)
	if !ok {
		return nil, nil
	}
	switch i.(type) {
	case *item:
		item := i.(*item)
		if item.Expired.Before(time.Now()) {
			//过期
			_ = m.del(key)
			//过期后删除
			return nil, nil
		}
		return item, nil
	default:
		err = fmt.Errorf("value of %s type error", key)
		return nil, err
	}
}

func (m *Memory) Set(key string, val interface{}, expire int) error {
	s, err := cast.ToStringE(val)
	if err != nil {
		return err
	}
	item := &item{
		Value:   s,
		Expired: time.Now().Add(time.Duration(expire) * time.Second),
	}
	return m.setItem(key, item)
}

func (m *Memory) setItem(key string, item *item) error {
	m.items.Store(key, item)
	return nil
}

func (m *Memory) Del(key ...string) error {
	return m.del(key...)
}

func (m *Memory) del(key ...string) error {
	for _, k := range key {
		m.items.Delete(k)
	}
	return nil
}

func (m *Memory) DelPattern(pattern string) error {
	//var keys []string
	//m.items.Range(func(key, value any) bool {
	//	keyString := cast.ToString(key)
	//	if strings.Contains(strings.ToLower(keyString), strings.ToLower(pattern)) {
	//		keys = append(keys, keyString)
	//	}
	//	return true
	//})
	//for _, key := range keys {
	//	if err := m.del(key); err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (m *Memory) HashKeys(hk string) ([]string, error) {
	// TODO 返回Keys
	return []string{}, nil
}

func (m *Memory) HashGet(hk, key string) (string, error) {
	item, err := m.getItem(hk + key)
	if err != nil || item == nil {
		return "", err
	}
	return item.Value, err
}

func (m *Memory) HashSet(hk, key string, val interface{}, expire int) error {
	s, err := cast.ToStringE(val)
	if err != nil {
		return err
	}
	item := &item{
		Value:   s,
		Expired: time.Now().Add(time.Duration(expire) * time.Second),
	}
	return m.setItem(hk+key, item)
}

func (m *Memory) HashDel(hk string, key ...string) error {
	var err error
	for _, k := range key {
		err = m.del(hk + k)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Memory) HashDelPattern(hk, pattern string) error {
	return m.DelPattern(hk + pattern)
}

func (m *Memory) Increase(key string) error {
	return m.calculate(key, 1)
}

func (m *Memory) Decrease(key string) error {
	return m.calculate(key, -1)
}

func (m *Memory) calculate(key string, num int) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	item, err := m.getItem(key)
	if err != nil {
		return err
	}

	if item == nil {
		err = fmt.Errorf("%s not exist", key)
		return err
	}
	var n int
	n, err = cast.ToIntE(item.Value)
	if err != nil {
		return err
	}
	n += num
	item.Value = strconv.Itoa(n)
	return m.setItem(key, item)
}

func (m *Memory) Expire(key string, dur time.Duration) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	item, err := m.getItem(key)
	if err != nil {
		return err
	}
	if item == nil {
		err = fmt.Errorf("%s not exist", key)
		return err
	}
	item.Expired = time.Now().Add(dur)
	return m.setItem(key, item)
}

func (m *Memory) GetClient() interface{} {
	return m
}
