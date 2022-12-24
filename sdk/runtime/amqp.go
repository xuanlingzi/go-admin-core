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
func (e *Amqp) PublishOnQueue(queueName string, body string, tag string) error {
	return e.amqp.PublishOnQueue(queueName, body, tag)
}

// SubscribeToQueue 消费消息
func (e *Amqp) SubscribeToQueue(queueName string, consumerName string, tag string, f message.AmqpConsumerFunc) error {
	return e.amqp.SubscribeToQueue(queueName, consumerName, tag, f)
}
