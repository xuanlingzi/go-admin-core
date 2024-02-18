package runtime

import (
	"github.com/xuanlingzi/go-admin-core/message"
)

type Amqp struct {
	prefix string
	amqp   message.AdapterAmqp
}

// String string输出
func (e *Amqp) String() string {
	if e.amqp == nil {
		return ""
	}
	return e.amqp.String()
}

// PublishOnQueue 发送消息
func (e *Amqp) PublishOnQueue(exchangeName, exchangeType, queueName, key, body string, tag string) error {
	return e.amqp.PublishOnQueue(exchangeName, exchangeType, queueName, key, tag, body)
}

// SubscribeToQueue 消费消息
func (e *Amqp) SubscribeToQueue(exchangeName, exchangeType, queueName, key, tag string, f message.AmqpConsumerFunc) error {
	return e.amqp.SubscribeToQueue(exchangeName, exchangeType, queueName, key, tag, f)
}
