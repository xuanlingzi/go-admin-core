package haibo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xuanlingzi/go-admin-core/logger"
)

// requestTimeout 单次出站请求的连接与响应等待超时（Requirement 1.9 / 16.3）。
const requestTimeout = 10 * time.Second

// contentTypeJSON 海博要求的请求 Content-Type（Requirement 1.1）。
const contentTypeJSON = "application/json;charset=UTF-8"

// HaiboResponse 为海博开放平台通信层响应（Requirement 1.7）。
// 仅 Code == "0" 表示通信成功；Result 为业务系统响应（再次序列化的 JSON）。
type HaiboResponse struct {
	Code       string          // "0" 仅表示通信成功
	Msg        string          // 海博返回的提示信息
	Result     json.RawMessage // 业务系统响应（再次序列化）
	ResponseId string          // 海博返回的 responseId，便于在海博平台自助查询
	Raw        string          // 原始响应体，便于排查
}

// CallLogger 出站调用日志钩子。
//
// 定义在 SDK 内以避免 go-admin-core 反向依赖 retail-core；由上层（DI 装配处）
// 注入具体实现，将每次出站调用落库到 haibo_message_log（Requirement 16.1, 16.3）。
type CallLogger interface {
	LogOutbound(info *OutboundCall)
}

// OutboundCall 描述一次出站调用的全量可观测信息（Requirement 16.1, 16.3）。
type OutboundCall struct {
	TargetURL    string    // 目标接口完整地址 domain+path
	Path         string    // 接口路径
	RequestBody  string    // 出站请求体（外层 JSON）
	ResponseBody string    // 原始响应体
	Code         string    // 海博响应 code
	Msg          string    // 海博响应 msg
	ResponseId   string    // 海博响应 responseId
	StartedAt    time.Time // 调用发起时间戳
	DurationMs   int64     // 调用耗时（毫秒）
	Err          error     // 通信失败原因（成功为 nil）
}

// HaiboClient 封装全部出站 HTTP 调用与签名，对业务层暴露统一的 Invoke。
type HaiboClient struct {
	appId      string
	secret     string
	domain     string // 由环境决定：prod=https://api-open.hiboos.com，test=https://p-open.hiboos.com
	client     *http.Client
	callLogger CallLogger
}

// NewHaiboClient 构建海博出站客户端，HTTP 超时固定为 10 秒（Requirement 1.9）。
func NewHaiboClient(appId, secret, domain string) *HaiboClient {
	return &HaiboClient{
		appId:  appId,
		secret: secret,
		domain: strings.TrimRight(domain, "/"),
		client: &http.Client{Timeout: requestTimeout},
	}
}

// String 返回 appId，满足适配器注册的 String() 约定。
func (c *HaiboClient) String() string {
	return c.appId
}

// SetCallLogger 注入出站调用日志钩子（在 DI 装配处接线到 HaiboMessageLog）。
func (c *HaiboClient) SetCallLogger(l CallLogger) {
	c.callLogger = l
}

// Invoke 组装系统级参数、序列化 body、签名、POST、解析响应。
//
// 流程（Requirement 1.1-1.4, 1.7-1.11, 16.1, 16.3）：
//  1. 将 bizParams 序列化为 body JSON 字符串。
//  2. 组装外层请求体 {appId, signMethod=md5, timestamp(13 位毫秒), sign, body}。
//  3. 以 Content-Type: application/json;charset=UTF-8 POST 到 domain+path。
//  4. 解析响应 code/msg/result/responseId。
//
// 失败语义：
//   - 连接失败 / 超时 / 非 2xx / 响应非 JSON / 缺少 code → 通信失败，返回 error 且不返回业务结果。
//   - code != "0" → 通信失败，记录 msg/responseId，返回带这些字段的响应与 error。
//
// 无论成功失败，调用结束都会回调 CallLogger（若已注入）。
func (c *HaiboClient) Invoke(path string, bizParams interface{}) (*HaiboResponse, error) {
	startedAt := time.Now()
	targetURL := c.domain + path

	bodyBytes, err := json.Marshal(bizParams)
	if err != nil {
		return nil, fmt.Errorf("序列化业务参数失败: %w", err)
	}
	body := string(bodyBytes)
	timestamp := strconv.FormatInt(startedAt.UnixMilli(), 10)
	sign := CalcSign(c.appId, body, signMethodMD5, timestamp, c.secret)

	requestPayload := map[string]string{
		"appId":      c.appId,
		"signMethod": signMethodMD5,
		"timestamp":  timestamp,
		"sign":       sign,
		"body":       body,
	}
	requestBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}
	requestBody := string(requestBytes)

	// 调用结束统一回调日志钩子。
	call := &OutboundCall{
		TargetURL:   targetURL,
		Path:        path,
		RequestBody: requestBody,
		StartedAt:   startedAt,
	}
	defer func() {
		call.DurationMs = time.Since(startedAt).Milliseconds()
		if c.callLogger != nil {
			c.callLogger.LogOutbound(call)
		}
	}()

	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(requestBytes))
	if err != nil {
		call.Err = fmt.Errorf("创建请求失败: %w", err)
		return nil, call.Err
	}
	req.Header.Set("Content-Type", contentTypeJSON)

	resp, err := c.client.Do(req)
	if err != nil {
		// 连接失败 / 超时 → 通信失败，不返回业务结果（Requirement 1.10, 16.3）。
		call.Err = fmt.Errorf("海博接口请求失败: %w", err)
		logger.Errorf("[海博] POST %s 通信失败: %s", targetURL, err.Error())
		return nil, call.Err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		call.Err = fmt.Errorf("读取海博响应失败: %w", err)
		return nil, call.Err
	}
	raw := string(respBytes)
	call.ResponseBody = raw
	logger.Infof("[海博] POST %s 响应(HTTP %d): %s", targetURL, resp.StatusCode, raw)

	// 非 2xx → 通信失败，不返回业务结果（Requirement 1.10）。
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		call.Err = fmt.Errorf("海博接口返回非 2xx 状态码(%d): %s", resp.StatusCode, raw)
		return nil, call.Err
	}

	// 解析响应；非 JSON 或缺少 code → 通信失败并保留原始响应（Requirement 1.11）。
	var fields map[string]json.RawMessage
	if err = json.Unmarshal(respBytes, &fields); err != nil {
		call.Err = fmt.Errorf("海博响应无法解析为 JSON: %s", raw)
		return nil, call.Err
	}
	rawCode, ok := fields["code"]
	if !ok {
		call.Err = fmt.Errorf("海博响应缺少 code 字段: %s", raw)
		return nil, call.Err
	}

	result := &HaiboResponse{
		Code:       rawToString(rawCode),
		Msg:        rawToString(fields["msg"]),
		Result:     normalizeRaw(fields["result"]),
		ResponseId: rawToString(fields["responseId"]),
		Raw:        raw,
	}
	call.Code = result.Code
	call.Msg = result.Msg
	call.ResponseId = result.ResponseId

	// 仅 code == "0" 视为通信成功；否则标记通信失败并记录 msg/responseId（Requirement 1.8）。
	if result.Code != "0" {
		call.Err = fmt.Errorf("海博接口通信失败: code=%s msg=%s responseId=%s", result.Code, result.Msg, result.ResponseId)
		return result, call.Err
	}

	return result, nil
}

// rawToString 将 JSON 原始字段转换为字符串，兼容字符串与数值两种表示。
func rawToString(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" {
		return ""
	}
	if strings.HasPrefix(s, "\"") {
		var str string
		if err := json.Unmarshal(raw, &str); err == nil {
			return str
		}
	}
	return s
}

// normalizeRaw 归一化可能为空的 result 字段：null / 空 → nil。
func normalizeRaw(raw json.RawMessage) json.RawMessage {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" {
		return nil
	}
	return raw
}
