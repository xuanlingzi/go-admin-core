package amqp

import (
	"encoding/json"
	"fmt"
	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/gogap/errors"
	"github.com/xuanlingzi/go-admin-core/logger"
	"github.com/xuanlingzi/go-admin-core/message"
	"strings"
	"sync"
	"time"
)

var _rocket mq_http_sdk.MQClient
var rocketLock sync.Mutex

type Rocket struct {
	conn       mq_http_sdk.MQClient
	accessId   string
	instanceId string
	namespace  string
}

// GetRocketClient 获取sms客户端
func GetRocketClient() mq_http_sdk.MQClient {
	return _rocket
}

// NewRocket redis模式
func NewRocket(client mq_http_sdk.MQClient, endpoint, accessKey, secretKey, instanceId, namespace string) *Rocket {
	if client == nil {
		client = mq_http_sdk.NewAliyunMQClient(endpoint, accessKey, secretKey, "")
		_rocket = client
	}
	r := &Rocket{
		conn:       client,
		accessId:   accessKey,
		instanceId: instanceId,
		namespace:  namespace,
	}
	return r
}

// Close 关闭连接
func (m *Rocket) Close() {
	if m.conn != nil {
		return
	}
}

// String 字符
func (m *Rocket) String() string {
	return m.accessId
}

// PublishOnQueue 发布消息
func (m *Rocket) PublishOnQueue(exchangeName, exchangeType, queueName, key, tag string, durableQueue bool, body interface{}) error {

	var err error
	// Topic所属的实例ID，在消息队列RocketMQ版控制台创建。
	// 若实例有命名空间，则实例ID必须传入；若实例无命名空间，则实例ID传入null空值或字符串空值。实例的命名空间可以在消息队列RocketMQ版控制台的实例详情页面查看。
	instanceId := m.instanceId
	if m.namespace == "" {
		instanceId = ""
	}
	producer := m.conn.GetProducer(instanceId, queueName)

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

	msg := mq_http_sdk.PublishMessageRequest{
		MessageBody: string(b),
		MessageTag:  tag,
		Properties:  map[string]string{},
	}
	resp, err := producer.PublishMessage(msg)
	if err != nil {
		logger.Errorf("Error to send message: %v", err.Error())
		return err
	}

	logger.Infof("A message was sent to queue id: %v, %v: %v", resp.MessageId, queueName, body)
	return err
}

func (m *Rocket) SubscribeToQueue(exchangeName, exchangeType, queueName, consumerName string, tag string, durableQueue bool, consumerExclusive bool, handlerFunc message.AmqpConsumerFunc) error {
	defer func() {
		panic("Rocket consumer stop")
	}()

	var err error
	instanceId := m.instanceId
	if m.namespace == "" {
		instanceId = ""
	}
	consumer := m.conn.GetConsumer(instanceId, queueName, consumerName, tag)

	for {
		endChan := make(chan int)
		respChan := make(chan mq_http_sdk.ConsumeMessageResponse)
		errChan := make(chan error)
		go func() {
			select {
			case resp := <-respChan:
				{
					// 处理业务逻辑
					var handles []string
					fmt.Printf("Consume %d messages---->\n", len(resp.Messages))
					for _, v := range resp.Messages {
						handles = append(handles, v.ReceiptHandle)
						fmt.Printf("\tMessageID: %s, PublishTime: %d, MessageTag: %s\n"+
							"\tConsumedTimes: %d, FirstConsumeTime: %d, NextConsumeTime: %d\n"+
							"\tBody: %s\n"+
							"\tProps: %s\n",
							v.MessageId, v.PublishTime, v.MessageTag, v.ConsumedTimes,
							v.FirstConsumeTime, v.NextConsumeTime, v.MessageBody, v.Properties)

						err = handlerFunc([]byte(v.MessageBody))
					}

					// NextConsumeTime前若不确认消息消费成功，则消息会重复消费
					// 消息句柄有时间戳，同一条消息每次消费拿到的都不一样
					ack := consumer.AckMessage(handles)
					if ack != nil {
						// 某些消息的句柄可能超时了会导致确认不成功
						fmt.Println(ack)
						if detail, ok := ack.(errors.ErrCode).Context()["Detail"].([]mq_http_sdk.ErrAckItem); ok {
							for _, errAckItem := range detail {
								fmt.Printf("\tErrorHandle:%s, ErrorCode:%s, ErrorMsg:%s\n",
									errAckItem.ErrorHandle, errAckItem.ErrorCode, errAckItem.ErrorMsg)
							}
							time.Sleep(time.Duration(3) * time.Second)
						}
					} else {
						fmt.Printf("Ack ---->\n\t%s\n", handles)
					}

					endChan <- 1
				}
			case err := <-errChan:
				{
					// 没有消息
					if strings.Contains(err.(errors.ErrCode).Error(), "MessageNotExist") {
						fmt.Println("\nNo new message, continue!")
					} else {
						fmt.Println(err)
						time.Sleep(time.Duration(3) * time.Second)
					}
					endChan <- 1
				}
			case <-time.After(35 * time.Second):
				{
					fmt.Println("Timeout of consumer message ??")
					endChan <- 1
				}
			}
		}()

		// 长轮询消费消息
		// 长轮询表示如果topic没有消息则请求会在服务端挂住3s，3s内如果有消息可以消费则立即返回
		consumer.ConsumeMessage(respChan, errChan,
			3, // 一次最多消费3条(最多可设置为16条)
			3, // 长轮询时间3秒（最多可设置为30秒）
		)
		<-endChan
	}

	return err
}
