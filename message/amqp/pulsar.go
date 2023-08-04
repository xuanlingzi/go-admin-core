package amqp

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/xuanlingzi/go-admin-core/message"
	"sync"
	"time"
)

var _pulsar pulsar.Client
var pulsarLock sync.Mutex

type Pulsar struct {
	conn      pulsar.Client
	appId     string
	namespace string
	producers map[string]pulsar.Producer
}

// GetPulsarClient 获取sms客户端
func GetPulsarClient() pulsar.Client {
	return _pulsar
}

// NewPulsar redis模式
func NewPulsar(client pulsar.Client, appId, secretKey, addr, namespace string) *Pulsar {
	var err error
	if client == nil {
		options := &pulsar.ClientOptions{
			URL:               addr,
			OperationTimeout:  30 * time.Second,
			ConnectionTimeout: 30 * time.Second,
			Authentication:    pulsar.NewAuthenticationToken(secretKey),
		}

		client, err = pulsar.NewClient(*options)
		if err != nil {
			panic(fmt.Sprintf("Pulsar init error: %v", err))
		}
		_pulsar = client
	}
	r := &Pulsar{
		conn:      client,
		producers: make(map[string]pulsar.Producer),
		appId:     appId,
		namespace: namespace,
	}
	return r
}

// Close 关闭连接
func (m *Pulsar) Close() {
	if m.conn != nil {
		m.conn.Close()
		m.producers = make(map[string]pulsar.Producer)
	}
}

// String 字符
func (m *Pulsar) String() string {
	return m.appId
}

func (m *Pulsar) InitProducer(queueName string) (pulsar.Producer, error) {

	var err error
	var ok bool
	var producer pulsar.Producer
	if producer, ok = m.producers[queueName]; ok {
		return producer, nil
	}

	defer pulsarLock.Unlock()
	pulsarLock.Lock()

	producer, err = m.conn.CreateProducer(pulsar.ProducerOptions{
		Topic: fmt.Sprintf("%v/%v/%v", m.appId, m.namespace, queueName),
	})
	if err != nil {
		return producer, err
	}
	m.producers[queueName] = producer

	return producer, nil
}

// PublishOnQueue 发布消息
func (m *Pulsar) PublishOnQueue(queueName string, body string, tag string) error {

	var err error
	var producer pulsar.Producer
	if producer, err = m.InitProducer(queueName); err != nil {
		panic(fmt.Sprintf("Pulsar producer init error: %v", err))
	}

	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: []byte(body),
	})
	if err != nil {
		return err
	}

	return err
}

func (m *Pulsar) InitSubscribe(queueName, consumerName string) (pulsar.Consumer, error) {

	defer pulsarLock.Unlock()
	pulsarLock.Lock()

	var err error
	var consumer pulsar.Consumer
	channel := make(chan pulsar.ConsumerMessage, 100)
	consumer, err = m.conn.Subscribe(pulsar.ConsumerOptions{
		Topic:            fmt.Sprintf("%v/%v/%v", m.appId, m.namespace, queueName),
		SubscriptionName: consumerName,
		Type:             pulsar.Shared,
		MessageChannel:   channel,
	})
	if err != nil {
		return consumer, err
	}

	return consumer, nil
}

func (m *Pulsar) SubscribeToQueue(queueName string, consumerName string, tag string, handlerFunc message.AmqpConsumerFunc) error {
	defer func() {
		panic("Pulsar consumer stop")
	}()

	var err error
	var consumer pulsar.Consumer
	var message pulsar.Message
	timeWait := 3
	timer := time.After(time.Duration(timeWait) * time.Minute)

	if consumer, err = m.InitSubscribe(queueName, consumerName); err != nil {
		return err
	}
	defer consumer.Close()

	for {
		select {
		case <-timer:
			err = errors.New("消费者接收超时")
			return err

		default:

			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			message, err = consumer.Receive(ctx)
			if err != nil {
				consumer.Close()
				m.Close()
				if consumer, err = m.InitSubscribe(consumerName, queueName); err != nil {
					return err
				}
				continue
			}

			if message == nil {
				continue
			}

			err = handlerFunc(string(message.Payload()))

			consumer.Ack(message)
			timer = time.After(time.Duration(timeWait) * time.Minute)
		}
	}
}
