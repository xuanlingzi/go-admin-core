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
	PublishOnQueue(exchangeName, exchangeType, queueName, key, tag string, body interface{}) error
	SubscribeToQueue(exchangeName, exchangeType, queueName, key, tag string, durableQueue bool, consumerExclusive bool, handlerFunc AmqpConsumerFunc) error
}

type AmqpConsumerFunc func([]byte) error
