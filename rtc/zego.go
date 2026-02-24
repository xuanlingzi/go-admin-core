package rtc

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	ZegoDefaultAddr             = "https://rtc-api.zego.im"
	ZegoDefaultSignatureVersion = "2.0"
)

// Zego adapter 封装即构服务端 API 的通用调用能力。
type Zego struct {
	Addr             string
	AppID            string
	ServerSecret     string
	SignatureVersion string
	client           *http.Client
}

// NewZego 构建 Zego adapter。
func NewZego(addr, appID, serverSecret, signatureVersion string, timeoutSec int) *Zego {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		addr = ZegoDefaultAddr
	}
	signatureVersion = strings.TrimSpace(signatureVersion)
	if signatureVersion == "" {
		signatureVersion = ZegoDefaultSignatureVersion
	}
	if timeoutSec <= 0 {
		timeoutSec = 8
	}

	return &Zego{
		Addr:             strings.TrimRight(addr, "/"),
		AppID:            strings.TrimSpace(appID),
		ServerSecret:     strings.TrimSpace(serverSecret),
		SignatureVersion: signatureVersion,
		client: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		},
	}
}

func (*Zego) String() string {
	return "zego"
}

func (m *Zego) Close() error {
	if m.client != nil {
		m.client.CloseIdleConnections()
		m.client = nil
	}
	return nil
}

func (m *Zego) GetClient() interface{} {
	return m.client
}

func (m *Zego) Get(action string, query map[string]string) (map[string]interface{}, error) {
	return m.Call(action, http.MethodGet, query, nil)
}

func (m *Zego) Post(action string, payload interface{}) (map[string]interface{}, error) {
	return m.Call(action, http.MethodPost, nil, payload)
}

// Call 按即构规范追加公共参数并完成签名后发起请求。
func (m *Zego) Call(action, method string, query map[string]string, payload interface{}) (map[string]interface{}, error) {
	if m.client == nil {
		return nil, fmt.Errorf("zego client 已关闭")
	}

	action = strings.TrimSpace(action)
	if action == "" {
		return nil, fmt.Errorf("缺少即构 Action")
	}

	method = strings.ToUpper(strings.TrimSpace(method))
	if method == "" {
		method = http.MethodGet
	}
	if method != http.MethodGet && method != http.MethodPost {
		return nil, fmt.Errorf("不支持的HTTP方法: %s", method)
	}

	requestURL, err := m.buildRequestURL(action, query)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if method == http.MethodPost {
		bodyBytes, marshalErr := marshalPayload(payload)
		if marshalErr != nil {
			return nil, marshalErr
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, requestURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("即构请求失败(status=%d): %s", resp.StatusCode, compactBody(raw))
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("即构响应为空")
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, fmt.Errorf("即构响应解析失败: %s, body=%s", err.Error(), compactBody(raw))
	}
	return result, nil
}

func (m *Zego) buildRequestURL(action string, query map[string]string) (string, error) {
	addr := strings.TrimSpace(m.Addr)
	if addr == "" {
		addr = ZegoDefaultAddr
	}

	u, err := url.Parse(addr)
	if err != nil {
		return "", fmt.Errorf("即构地址格式错误: %s", err.Error())
	}

	nonce := strings.ReplaceAll(uuid.New().String(), "-", "")
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	values := u.Query()
	values.Set("Action", action)
	values.Set("AppId", m.AppID)
	values.Set("SignatureNonce", nonce)
	values.Set("SignatureVersion", m.SignatureVersion)
	values.Set("Timestamp", timestamp)
	values.Set("Signature", m.calcSignature(nonce, timestamp))

	for key, value := range query {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		values.Set(key, strings.TrimSpace(value))
	}

	u.RawQuery = values.Encode()
	return u.String(), nil
}

func (m *Zego) calcSignature(nonce, timestamp string) string {
	signData := m.AppID + nonce + m.ServerSecret + timestamp
	sum := md5.Sum([]byte(signData))
	return hex.EncodeToString(sum[:])
}

func marshalPayload(payload interface{}) ([]byte, error) {
	if payload == nil {
		return []byte("{}"), nil
	}
	switch v := payload.(type) {
	case []byte:
		if len(v) == 0 {
			return []byte("{}"), nil
		}
		return v, nil
	case string:
		if strings.TrimSpace(v) == "" {
			return []byte("{}"), nil
		}
		return []byte(v), nil
	default:
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("即构请求体序列化失败: %s", err.Error())
		}
		if len(data) == 0 {
			return []byte("{}"), nil
		}
		return data, nil
	}
}

func compactBody(raw []byte) string {
	body := strings.TrimSpace(string(raw))
	if len(body) > 512 {
		return body[:512] + "..."
	}
	return body
}
