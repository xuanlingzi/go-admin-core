package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	rabbitmq "github.com/rabbitmq/amqp091-go"
	"github.com/xuanlingzi/go-admin-core/logger"
	"github.com/xuanlingzi/go-admin-core/message"
	"time"
)

type Rabbit struct {
	publishConn   *rabbitmq.Connection
	subscribeConn *rabbitmq.Connection
	endpoint      string
}

// NewRabbit redis模式
func NewRabbit(addr, accessKey, secretKey, vhost string) *Rabbit {
	config := rabbitmq.Config{
		Vhost:      vhost,
		Properties: rabbitmq.NewConnectionProperties(),
		Heartbeat:  30 * time.Second,
	}
	config.Vhost = "/"

	endpoint := fmt.Sprintf("amqp://%v:%v@%v", accessKey, secretKey, addr)
	publishConn, err := rabbitmq.DialConfig(endpoint, config)
	if err != nil {
		logger.Errorf("Error to connect to RabbitMQ: %v", err.Error())
	}

	subscribeConn, err := rabbitmq.DialConfig(endpoint, config)
	if err != nil {
		logger.Errorf("Error to connect to RabbitMQ: %v", err.Error())
	}

	r := &Rabbit{
		publishConn:   publishConn,
		subscribeConn: subscribeConn,
		endpoint:      endpoint,
	}
	return r
}

// Close 关闭连接
func (m *Rabbit) Close() {
	if m.publishConn != nil {
		ch, err := m.publishConn.Channel()
		if err != nil {
			logger.Errorf("Error to open a channel: %v", err.Error())
		}
		if ch != nil {
			err = ch.Close()
			if err != nil {
				logger.Errorf("Error to close channel: %v", err.Error())
			}
		}

		err = m.publishConn.Close()
		if err != nil {
			logger.Errorf("Error to close publish connection: %v", err.Error())
		}
	}
	if m.subscribeConn != nil {
		ch, err := m.subscribeConn.Channel()
		if err != nil {
			logger.Errorf("Error to open a channel: %v", err.Error())
		}
		if ch != nil {
			err = ch.Close()
			if err != nil {
				logger.Errorf("Error to close channel: %v", err.Error())
			}
		}

		err = m.subscribeConn.Close()
		if err != nil {
			logger.Errorf("Error to close subscribe connection: %v", err.Error())
		}
	}
}

// PublishOnQueue 发布消息
func (m *Rabbit) PublishOnQueue(exchangeName, exchangeType, queueName, key, tag string, body interface{}) error {
	var err error

	channel, err := m.publishConn.Channel()
	if err != nil {
		logger.Errorf("Error to open a channel: %v", err.Error())
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	if err != nil {
		logger.Errorf("Error to declare a exchange: %v", err.Error())
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
	var err error
	channel, err := m.subscribeConn.Channel()
	if err != nil {
		logger.Errorf("Error to open a channel: %v", err.Error())
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	if err != nil {
		logger.Errorf("Error to declare a exchange: %v", err.Error())
	}

	queue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		logger.Errorf("Error to declare a queue: %v", err.Error())
	}

	err = channel.QueueBind(queue.Name, key, exchangeName, false, nil)
	if err != nil {
		logger.Errorf("Error to bind a queue: %v", err.Error())
	}

	deliver, err := channel.Consume(queue.Name, "", false, true, false, false, nil)
	if err != nil {
		logger.Errorf("Error to consumer a queue: %v", err.Error())
	}

	var forever chan struct{}

	go func() {
		for d := range deliver {
			err = handlerFunc(d.Body)
			if err != nil {
				logger.Errorf("Error to handle message: %v", err.Error())
				err = d.Ack(false)
			} else {
				err = d.Ack(true)
				if err != nil {
					logger.Errorf("Error to ack message: %v", err.Error())
				}
			}
		}
	}()

	<-forever

	return err
}

// String 字符
func (m *Rabbit) String() string {
	return m.endpoint
}
