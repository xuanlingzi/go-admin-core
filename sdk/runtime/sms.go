package runtime

import "github.com/xuanlingzi/go-admin-core/message"

type Sms struct {
	prefix string
	sms    message.AdapterSms
}

// String string输出
func (e *Sms) String() string {
	if e.sms == nil {
		return ""
	}
	return e.sms.String()
}

// Send val by announces
func (e *Sms) Send(addresses []string, template string, params map[string]string) error {
	return e.sms.Send(addresses, template, params)
}

func (e *Sms) GetClient() interface{} {
	return e.sms.GetClient()
}
