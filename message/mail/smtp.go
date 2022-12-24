package mail

import (
	"crypto/tls"
	"fmt"

	"net"
	"net/smtp"
	"sync"
)

var smtpLock sync.Mutex
var _smtp *smtp.Client

type SmtpClient struct {
	client *smtp.Client
	host   string
	port   string
	auth   smtp.Auth
	from   string
}

// GetSmtpClient 获取smtp客户端
func GetSmtpClient() *smtp.Client {
	return _smtp
}

func NewSmtpClient(client *smtp.Client, addr string, username, password, from string) *SmtpClient {
	var err error

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil
	}

	auth := smtp.PlainAuth("", username, password, host)

	c := &SmtpClient{
		client: client,
		host:   host,
		port:   port,
		auth:   auth,
		from:   from,
	}
	return c
}

func (m *SmtpClient) Setup() error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%v:%v", m.host, m.port), tlsConfig)
	if err != nil {
		return err
	}

	m.client, err = smtp.NewClient(conn, m.host)
	if err != nil {
		return err
	}
	if err = m.client.Auth(m.auth); err != nil {
		return err
	}

	return err
}

func (*SmtpClient) String() string {
	return "smtp"
}

func (m *SmtpClient) Send(addresses []string, body []byte) error {
	defer smtpLock.Unlock()
	smtpLock.Lock()

	var err error
	if err = m.Setup(); err != nil {
		return err
	}
	defer m.Close()

	if err = m.client.Mail(m.from); err != nil {
		return err
	}
	for _, address := range addresses {
		if err = m.client.Rcpt(address); err != nil {
			return err
		}
	}

	w, err := m.client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	m.client.Quit()
	return err
}

func (m *SmtpClient) Close() {
	if m.client != nil {
		m.client.Close()
	}
}

// GetClient 暴露原生client
func (m *SmtpClient) GetClient() *smtp.Client {
	return m.client
}
