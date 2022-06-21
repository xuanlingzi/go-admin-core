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

type AdapterThirdParty interface {
	String() string
	GetToken(merchant, platform string) string
	GetConnectUrl(state, scope, redirectUrl string, popUp bool) (string, error)
	GetAccessToken(force bool) (string, error)
	GetUserAccessToken(code, state string) (string, error)
	GetUserInfo(accessToken, openId string) (string, error)
	SendTemplateMessage(openId, templateId, url string, data []byte) (string, error)
}
