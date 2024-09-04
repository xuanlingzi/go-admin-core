package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	rabbitmq "github.com/rabbitmq/amqp091-go"
	"github.com/xuanlingzi/go-admin-core/logger"
	"github.com/xuanlingzi/go-admin-core/message"
	"strings"
	"time"
)

type Rabbit struct {
	config   rabbitmq.Config
	endpoint string

	conn       *rabbitmq.Connection
	connNotify chan *rabbitmq.Error
}

// NewRabbit redis模式
func NewRabbit(addr, accessKey, secretKey, vhost string) *Rabbit {
	config := rabbitmq.Config{
		Vhost:      vhost,
		Properties: rabbitmq.NewConnectionProperties(),
		Heartbeat:  30 * time.Second,
	}

	endpoint := fmt.Sprintf("amqp://%v:%v@%v", accessKey, secretKey, addr)

	r := &Rabbit{
		config:     config,
		endpoint:   endpoint,
		connNotify: make(chan *rabbitmq.Error),
	}

	r.reconnect()
	go r.keepAlive()

	return r
}

func (m *Rabbit) reconnect() {
	var err error
	m.conn, err = rabbitmq.DialConfig(m.endpoint, m.config)
	if err != nil {
		panic(fmt.Sprintf("Error to reconnect to RabbitMQ: %v", err.Error()))
	}
	m.connNotify = m.conn.NotifyClose(make(chan *rabbitmq.Error))
}

func (m *Rabbit) keepAlive() {
	for {
		select {
		case errNotify := <-m.connNotify:

			logger.Errorf("RabbitMQ connection closed: %v", errNotify.Error())
			m.reconnect()
		}

		for err := range m.connNotify {
			logger.Errorf("RabbitMQ connection closed: %v", err.Error())
		}
	}
}

// Close 关闭连接
func (m *Rabbit) Close() {
	_ = m.conn.Close()
}

// PublishOnQueue 发布消息
func (m *Rabbit) PublishOnQueue(exchangeName, exchangeType, queueName, key, tag string, durableQueue bool, body interface{}) error {
	var err error

	if strings.EqualFold(exchangeType, "topic") && queueName == "" {
		queueName = exchangeName
	}
	if key == "" {
		key = exchangeName
	}

	if m.conn.IsClosed() {
		m.reconnect()
	}

	var channel *rabbitmq.Channel
	channel, err = m.conn.Channel()
	if err != nil {
		if channel != nil {
			err = channel.Close()
			if err != nil {
				logger.Errorf("RabbitMQ channel close error: %v", err.Error())
			}
		}
		return err
	}
	defer func() {
		// 关闭 AMQP 通道，注意在关闭前检查是否为 nil
		if channel != nil {
			err = channel.Close()
			if err != nil {
				logger.Errorf("RabbitMQ channel close error: %v", err.Error())
			}
		}
	}()

	err = channel.ExchangeDeclare(exchangeName, exchangeType, durableQueue, false, false, false, nil)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var b []byte
	switch body.(type) {
	case string:
		b = []byte(body.(string))
	default:
		b, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	err = channel.PublishWithContext(ctx, exchangeName, key, false, false, rabbitmq.Publishing{
		DeliveryMode: rabbitmq.Persistent,
		ContentType:  "application/json",
		Body:         b,
	})

	return err
}

func (m *Rabbit) SubscribeToQueue(exchangeName, exchangeType, queueName, key, tag string, durableQueue bool, consumerExclusive bool, handlerFunc message.AmqpConsumerFunc) error {

	if queueName == "" {
		queueName = exchangeName
	}
	if key == "" {
		key = exchangeName
	}

	conn, err := rabbitmq.DialConfig(m.endpoint, m.config)
	if err != nil {
		panic(fmt.Sprintf("Error to reconnect to RabbitMQ: %v", err.Error()))
	}

	channel, err := conn.Channel()
	if err != nil {
		if channel != nil {
			err = channel.Close()
			if err != nil {
				logger.Errorf("RabbitMQ channel close error: %v", err.Error())
			}
		}
		if conn != nil {
			err = conn.Close()
			if err != nil {
				logger.Errorf("RabbitMQ connection close error: %v", err.Error())
			}
		}
		return err
	}
	defer func() {
		if channel != nil {
			err = channel.Close()
			if err != nil {
				logger.Errorf("RabbitMQ channel close error: %v", err.Error())
			}
		}
		if conn != nil {
			err = conn.Close()
			if err != nil {
				logger.Errorf("RabbitMQ connection close error: %v", err.Error())
			}
		}
	}()

	/*
		Exchange 的 Durable：
		持久化 exchange 在 RabbitMQ 重启后仍然存在。
		持久化 exchange 本身不会保存消息，它们只是消息路由的定义。
		autoDelete 自动删除，当没有消费者连接到队列时，队列会被自动删除
	*/
	err = channel.ExchangeDeclare(exchangeName, exchangeType, durableQueue, false, false, false, nil)
	if err != nil {
		return err
	}

	/*
		Queue 的 Exclusive 属性指定队列是否是独占队列，独占队列有以下特性：
		仅限当前连接：队列仅对声明它的连接可见。其他连接无法访问该队列。
		自动删除：当声明该队列的连接关闭时，队列会被自动删除。
		使用场景：独占队列通常用于临时性的、单个客户端使用的场景，比如临时RPC响应队列。

		Queue 的 Durable：
		持久化 queue 在 RabbitMQ 重启后仍然存在。
		持久化 queue 可以保存持久化的消息，使这些消息在 RabbitMQ 重启后也能继续存在。
		注意：要确保消息在重启后不丢失，消息本身也需要被标记为持久化的（将 deliveryMode 设置为 2）。
		autoDelete 自动删除，当没有消费者连接到队列时，队列会被自动删除
	*/
	queue, err := channel.QueueDeclare(queueName, durableQueue, false, false, false, nil)
	if err != nil {
		return err
	}

	err = channel.QueueBind(queue.Name, key, exchangeName, false, nil)
	if err != nil {
		return err
	}

	/*
		Consume 的 Exclusive 属性指定消费者是否是独占消费者，独占消费者有以下特性：
		独占访问：队列仅能有一个独占消费者。如果已经有一个独占消费者，其他消费者（无论是否独占）都无法消费该队列。
		优先级高：独占消费者比普通消费者优先级高，确保其唯一的消息消费权。
		使用场景：独占消费者通常用于确保某个消费者独享队列消息的场景，比如在主备模式中，主节点需要独占消息消费权。
	*/
	deliver, err := channel.Consume(queue.Name, "", false, consumerExclusive, false, false, nil)
	if err != nil {
		return err
	}

	for d := range deliver {
		err = handlerFunc(d.Body)
		if err != nil {
			err = d.Nack(true, true)
		} else {
			err = d.Ack(true)
		}
		if err != nil {
			logger.Errorf("RabbitMQ Nack error: %v", err.Error())
		}
	}

	return nil
}

// String 字符
func (m *Rabbit) String() string {
	return m.endpoint
}
