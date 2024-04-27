package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	rabbitmq "github.com/rabbitmq/amqp091-go"
	"github.com/xuanlingzi/go-admin-core/logger"
	"github.com/xuanlingzi/go-admin-core/message"
	"sync"
	"time"
)

type Rabbit struct {
	config   rabbitmq.Config
	endpoint string

	conn       *rabbitmq.Connection
	connNotify chan *rabbitmq.Error

	channels sync.Map
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

		if m.conn.IsClosed() == false {
			m.channels.Range(func(k, v interface{}) bool {
				ch, ok := v.(*rabbitmq.Channel)
				if ok {
					err := ch.Close()
					if err != nil {
						logger.Errorf("Error to close channel: %v", err.Error())
						return false
					}
				}
				return true
			})
		}

		for err := range m.connNotify {
			logger.Errorf("RabbitMQ connection closed: %v", err.Error())
		}
	}
}

// Close 关闭连接
func (m *Rabbit) Close() {
	m.channels.Range(func(k, v interface{}) bool {
		ch, ok := v.(*rabbitmq.Channel)
		if ok {
			err := ch.Close()
			if err != nil {
				logger.Errorf("Error to close channel: %v", err.Error())
				return false
			}
		}
		return true
	})

	_ = m.conn.Close()
}

// PublishOnQueue 发布消息
func (m *Rabbit) PublishOnQueue(exchangeName, exchangeType, queueName, key, tag string, body interface{}) error {
	var err error

	if m.conn.IsClosed() {
		m.reconnect()
	}

	var channel *rabbitmq.Channel
	ch, ok := m.channels.Load(exchangeName)
	if ok {
		channel = ch.(*rabbitmq.Channel)
		if channel == nil || channel.IsClosed() {
			channel = nil
		}
	}
	if channel == nil {
		channel, err = m.conn.Channel()
		if err != nil {
			channel.Close()
			return err
		}

		m.channels.Store(exchangeName, channel)
	}

	err = channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
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
		ContentType: "application/json",
		Body:        b,
	})

	return err
}

func (m *Rabbit) SubscribeToQueue(exchangeName, exchangeType, queueName, key, tag string, handlerFunc message.AmqpConsumerFunc) error {

	conn, err := rabbitmq.DialConfig(m.endpoint, m.config)
	if err != nil {
		panic(fmt.Sprintf("Error to reconnect to RabbitMQ: %v", err.Error()))
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}
	defer func() {
		_ = channel.Close()
		_ = conn.Close()
	}()

	err = channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	if err != nil {
		return err
	}

	queue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return err
	}

	err = channel.QueueBind(queue.Name, key, exchangeName, false, nil)
	if err != nil {
		return err
	}

	deliver, err := channel.Consume(queue.Name, "", false, true, false, false, nil)
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
