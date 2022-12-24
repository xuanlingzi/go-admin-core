package config

import (
	"github.com/xuanlingzi/go-admin-core/message"
	"github.com/xuanlingzi/go-admin-core/message/mail"
)

type Mail struct {
	Smtp *SmtpConnectOptions `json:"smtp" yaml:"smtp"`
}

var MailConfig = new(Mail)

// Setup 构造邮件配置
func (e Mail) Setup() (message.AdapterMail, error) {
	if e.Smtp != nil {
		addr, username, password, from := e.Smtp.GetSmtpOptions()
		r := mail.NewSmtpClient(GetSmtpClient(), addr, username, password, from)
		if _smtp == nil {
			_smtp = r.GetClient()
		}
		return r, nil
	}
	return nil, nil
}
