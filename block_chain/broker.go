package block_chain

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"net/http"
	"strings"
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
	resp, err := httpPost(url, body)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (m *Broker) Status(chain string, content string, hash string) (string, error) {

	var body string
	body, _ = sjson.Set(body, "chain", chain)
	body, _ = sjson.Set(body, "content", content)
	body, _ = sjson.Set(body, "hash", hash)

	url := fmt.Sprintf("%v/status", m.mintAddr)
	resp, err := httpPost(url, body)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func httpGet(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if !gjson.GetBytes(body, "code").Exists() {
		return "", errors.New(string(body))
	}
	if gjson.GetBytes(body, "code").Int() != 0 {
		return "", errors.New(gjson.GetBytes(body, "message").String())
	}

	return string(body), nil
}

func httpPost(url, content string) (string, error) {
	response, err := http.Post(url, "application/json", strings.NewReader(content))
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if !gjson.GetBytes(body, "code").Exists() {
		return "", errors.New(string(body))
	}
	if gjson.GetBytes(body, "code").Int() != 0 {
		return "", errors.New(gjson.GetBytes(body, "message").String())
	}

	return string(body), nil
}

// GetClient 暴露原生client
func (m *Broker) GetClient() *http.Client {
	return m.client
}
