package runtime

import (
	"github.com/go-admin-team/go-admin-core/storage"
)

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

// PublishOnQueue 发送消息
func (e *Amqp) PublishOnQueue(queueName string, body []byte) error {
	return e.store.PublishOnQueue(queueName, body)
}

// SubscribeToQueue 消费消息
func (e *Amqp) SubscribeToQueue(queueName string, consumerName string, f storage.AmqpConsumerFunc) error {
	return e.store.SubscribeToQueue(queueName, consumerName, f)
}