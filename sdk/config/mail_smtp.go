package config

import (
	"net/smtp"
)

var _smtp *smtp.Client

// GetSmtpClient 获取smtp客户端
func GetSmtpClient() *smtp.Client {
	return _smtp
}

// SetSmtpClient 设置smtp客户端
func SetSmtpClient(c *smtp.Client) {
	if _smtp != nil && _smtp != c {
		_smtp = nil
	}
	_smtp = c
}

type SmtpConnectOptions struct {
	Addr     string `yaml:"addr" json:"addr"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	From     string `yaml:"from" json:"from"`
}

func (e SmtpConnectOptions) GetSmtpOptions() (string, string, string, string) {
	return e.Addr, e.Username, e.Password, e.From
}
