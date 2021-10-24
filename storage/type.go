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
	HashGet(hk, key string) (string, error)
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

type AdapterSms interface {
	String() string
	Send(phones []string, templateId string, params []string) error
}

type AdapterCos interface {
	String() string
	PutFromFile(fileLocation string) error
}

type AdapterAmqp interface {
	String() string
	PublishOnQueue(queueName string, body []byte) error
	SubscribeToQueue(queueName string, consumerName string, handlerFunc AmqpConsumerFunc) error
}

type AmqpConsumerFunc func([]byte) error
