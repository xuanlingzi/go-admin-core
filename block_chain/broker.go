package block_chain

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/xuanlingzi/go-admin-core/tools/utils"
	"net/http"
)

type Broker struct {
	client   *http.Client
	mintAddr string
}

func NewBroker(client *http.Client, mintAddr string) (*Broker, error) {
	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{},
		}
	}
	c := &Broker{
		client:   client,
		mintAddr: mintAddr,
	}
	return c, nil
}

// Close 关闭连接
func (m *Broker) Close() {
	if m.client != nil {
		m.client.CloseIdleConnections()
		m.client = nil
	}
}

func (*Broker) String() string {
	return "broker"
}

func (m *Broker) Send(chain string, content string, callback string) (string, error) {

	var body string
	body, _ = sjson.Set(body, "chain", chain)
	body, _ = sjson.Set(body, "content", content)
	body, _ = sjson.Set(body, "callback", callback)

	url := fmt.Sprintf("%v/send", m.mintAddr)
	resp, err := utils.HttpPost(url, body)
	if err != nil {
		return "", err
	}
	if !gjson.GetBytes(resp, "code").Exists() {
		return "", errors.New(string(body))
	}
	if gjson.GetBytes(resp, "code").Int() != 0 {
		return "", errors.New(gjson.GetBytes(resp, "message").String())
	}

	return string(resp), nil
}

func (m *Broker) Status(chain string, content string, hash string) (string, error) {

	var body string
	body, _ = sjson.Set(body, "chain", chain)
	body, _ = sjson.Set(body, "content", content)
	body, _ = sjson.Set(body, "hash", hash)

	url := fmt.Sprintf("%v/status", m.mintAddr)
	resp, err := utils.HttpPost(url, body)
	if err != nil {
		return "", err
	}
	if !gjson.GetBytes(resp, "code").Exists() {
		return "", errors.New(string(body))
	}
	if gjson.GetBytes(resp, "code").Int() != 0 {
		return "", errors.New(gjson.GetBytes(resp, "message").String())
	}

	return string(resp), nil
}

// GetClient 暴露原生client
func (m *Broker) GetClient() *http.Client {
	return m.client
}
