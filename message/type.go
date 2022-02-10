package message

const (
	PrefixKey = "__host"
)

type AdapterMail interface {
	String() string
	Send(address []string, body []byte) error
}

type AdapterSms interface {
	String() string
	Send(addresses []string, template string, params []string) error
}

type AdapterAmqp interface {
	String() string
	PublishOnQueue(queueName string, body []byte) error
	SubscribeToQueue(queueName string, consumerName string, handlerFunc AmqpConsumerFunc) error
}

type AmqpConsumerFunc func([]byte) error
