package runtime

import (
	"github.com/go-admin-team/go-admin-core/storage"
)

// NewAmqp 创建对应上下文缓存
func NewAmqp(prefix string, store storage.AdapterAmqp) storage.AdapterAmqp {
	return &Amqp{
		prefix:          prefix,
		store:           store,
	}
}

type Amqp struct {
	prefix          string
	store           storage.AdapterAmqp
}

// String string输出
func (e *Amqp) String() string {
	if e.store == nil {
		return ""
	}
	return e.store.String()
}

// Publish mqssage on queue
func (e *Amqp) PublishOnQueue(queueName string, body []byte) error {
	return e.store.PublishOnQueue(queueName, body)
}

// Subscribe mqssage to queue
func (e *Amqp) SubscribeToQueue(queueName string, consumerName string, f storage.AmqpConsumerFunc) error {
	return e.store.SubscribeToQueue(queueName, consumerName, f)
}