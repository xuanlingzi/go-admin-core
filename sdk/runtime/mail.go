package runtime

import "github.com/xuanlingzi/go-admin-core/message"

type Mail struct {
	prefix string
	mail   message.AdapterMail
}

// String string输出
func (e *Mail) String() string {
	if e.mail == nil {
		return ""
	}
	return e.mail.String()
}

// Send val by announces
func (e *Mail) Send(address []string, body []byte) error {
	return e.mail.Send(address, body)
}
