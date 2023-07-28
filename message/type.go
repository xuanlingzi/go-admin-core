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
	Close()
	Send(addresses []string, template string, params map[string]string) error
	GetClient() interface{}
}

type AdapterAmqp interface {
	String() string
	PublishOnQueue(queueName string, body string, tag string) error
	SubscribeToQueue(queueName string, consumerName string, tag string, handlerFunc AmqpConsumerFunc) error
}

type AmqpConsumerFunc func(string) error
