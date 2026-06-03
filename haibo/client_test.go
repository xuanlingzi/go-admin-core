package haibo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// thirteenDigits 校验 timestamp 为 13 位毫秒级整数（Requirement 1.3）。
var thirteenDigits = regexp.MustCompile(`^\d{13}$`)

// captureLogger 是 CallLogger 的测试桩，记录最后一次出站调用，用于断言日志钩子被回调。
type captureLogger struct {
	last  *OutboundCall
	count int
}

func (l *captureLogger) LogOutbound(info *OutboundCall) {
	l.last = info
	l.count++
}

// newTestClient 构造指向 mock server 的客户端，secret 固定便于断言。
func newTestClient(domain string) *HaiboClient {
	return NewHaiboClient("4de67a45194f12345678968195fbc", goldenSecret, domain)
}

// TestInvoke_Success 正常成功（code="0"）：解析 code/msg/result/responseId（Requirement 1.7）。
func TestInvoke_Success(t *testing.T) {
	var gotContentType, gotMethod string
	var gotPayload map[string]string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		gotMethod = r.Method
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotPayload)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"0","msg":"ok","result":"{\"currentStock\":5}","responseId":"resp-123"}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	logger := &captureLogger{}
	client.SetCallLogger(logger)

	resp, err := client.Invoke("/api/stock/get", map[string]any{"stationNo": "1234"})
	if err != nil {
		t.Fatalf("期望成功，却返回错误: %v", err)
	}
	if resp.Code != "0" {
		t.Errorf("Code = %q, want \"0\"", resp.Code)
	}
	if resp.Msg != "ok" {
		t.Errorf("Msg = %q, want \"ok\"", resp.Msg)
	}
	if resp.ResponseId != "resp-123" {
		t.Errorf("ResponseId = %q, want \"resp-123\"", resp.ResponseId)
	}
	if string(resp.Result) != `"{\"currentStock\":5}"` {
		t.Errorf("Result = %q, 未原样保留业务结果", string(resp.Result))
	}

	// 校验出站协议约定（Requirement 1.1-1.4）。
	if gotMethod != http.MethodPost {
		t.Errorf("HTTP 方法 = %q, want POST", gotMethod)
	}
	if gotContentType != contentTypeJSON {
		t.Errorf("Content-Type = %q, want %q", gotContentType, contentTypeJSON)
	}
	for _, field := range []string{"appId", "signMethod", "timestamp", "sign", "body"} {
		if _, ok := gotPayload[field]; !ok {
			t.Errorf("请求体缺少系统级参数 %q", field)
		}
	}
	if gotPayload["signMethod"] != "md5" {
		t.Errorf("signMethod = %q, want \"md5\"", gotPayload["signMethod"])
	}
	if !thirteenDigits.MatchString(gotPayload["timestamp"]) {
		t.Errorf("timestamp = %q, want 13 位毫秒整数", gotPayload["timestamp"])
	}
	// body 必须是业务参数序列化后的 JSON 字符串。
	if !strings.Contains(gotPayload["body"], `"stationNo":"1234"`) {
		t.Errorf("body = %q, 未包含序列化后的业务参数", gotPayload["body"])
	}
	// sign 必须与按相同入参重算的结果一致。
	wantSign := CalcSign(gotPayload["appId"], gotPayload["body"], "md5", gotPayload["timestamp"], goldenSecret)
	if gotPayload["sign"] != wantSign {
		t.Errorf("sign = %q, want %q", gotPayload["sign"], wantSign)
	}

	// 日志钩子必须被回调一次且记录成功信息（Requirement 16.1）。
	if logger.count != 1 {
		t.Fatalf("CallLogger 回调次数 = %d, want 1", logger.count)
	}
	if logger.last.Err != nil {
		t.Errorf("成功调用的日志 Err 应为 nil, got %v", logger.last.Err)
	}
	if logger.last.Code != "0" || logger.last.ResponseId != "resp-123" {
		t.Errorf("日志未记录正确的 code/responseId: %+v", logger.last)
	}
	if logger.last.DurationMs < 0 {
		t.Errorf("日志耗时异常: %d", logger.last.DurationMs)
	}
}

// TestInvoke_BusinessCodeNonZero code≠0：标记通信失败并记录 msg/responseId（Requirement 1.8）。
func TestInvoke_BusinessCodeNonZero(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"500","msg":"签名错误","responseId":"resp-err-9"}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	logger := &captureLogger{}
	client.SetCallLogger(logger)

	resp, err := client.Invoke("/api/order/picked", map[string]any{"hbOrderId": "X1"})
	if err == nil {
		t.Fatal("code≠0 时期望返回错误，却为 nil")
	}
	// 仍应返回带 msg/responseId 的响应以供记录（Requirement 1.8）。
	if resp == nil {
		t.Fatal("code≠0 时应返回带 msg/responseId 的响应，却为 nil")
	}
	if resp.Code != "500" || resp.Msg != "签名错误" || resp.ResponseId != "resp-err-9" {
		t.Errorf("响应字段不正确: %+v", resp)
	}
	if logger.last == nil || logger.last.Err == nil {
		t.Error("通信失败应在日志中记录 Err")
	}
}

// TestInvoke_ConnectionFailure 连接失败：通信失败且不返回业务结果（Requirement 1.10）。
//
// 注：client 的 http 超时固定为 10s 且不可注入，直接测真实超时会让用例耗时 10s。
// 连接失败与超时走同一「c.client.Do 返回 error → 通信失败不返回业务结果」分支，
// 因此用指向已关闭端口的地址在毫秒级内触发同一失败语义，覆盖该需求而不阻塞 10s。
func TestInvoke_ConnectionFailure(t *testing.T) {
	// 启动后立即关闭，获得一个保证无人监听的地址。
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	closedURL := srv.URL
	srv.Close()

	client := newTestClient(closedURL)
	logger := &captureLogger{}
	client.SetCallLogger(logger)

	start := time.Now()
	resp, err := client.Invoke("/api/stock/get", map[string]any{"a": 1})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("连接失败时期望返回错误，却为 nil")
	}
	if resp != nil {
		t.Errorf("连接失败时不应返回业务结果，got %+v", resp)
	}
	if elapsed > 2*time.Second {
		t.Errorf("连接失败用例耗时过长(%v)，应在毫秒级返回", elapsed)
	}
	if logger.last == nil || logger.last.Err == nil {
		t.Error("连接失败应在日志中记录 Err")
	}
}

// TestInvoke_Non2xxStatus 非 2xx：通信失败且不返回业务结果（Requirement 1.10）。
func TestInvoke_Non2xxStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"code":"0"}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	logger := &captureLogger{}
	client.SetCallLogger(logger)

	resp, err := client.Invoke("/api/order/picked", map[string]any{"a": 1})
	if err == nil {
		t.Fatal("非 2xx 时期望返回错误，却为 nil")
	}
	if resp != nil {
		t.Errorf("非 2xx 时不应返回业务结果，got %+v", resp)
	}
	if !strings.Contains(err.Error(), strconv.Itoa(http.StatusBadGateway)) {
		t.Errorf("错误信息应包含状态码 %d: %v", http.StatusBadGateway, err)
	}
	if logger.last == nil || logger.last.Err == nil {
		t.Error("非 2xx 应在日志中记录 Err")
	}
}

// TestInvoke_NonJSONBody 响应非 JSON：通信失败并保留原始响应（Requirement 1.11）。
func TestInvoke_NonJSONBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><body>502 Bad Gateway</body></html>`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	logger := &captureLogger{}
	client.SetCallLogger(logger)

	resp, err := client.Invoke("/api/stock/get", map[string]any{"a": 1})
	if err == nil {
		t.Fatal("响应非 JSON 时期望返回错误，却为 nil")
	}
	if resp != nil {
		t.Errorf("响应非 JSON 时不应返回业务结果，got %+v", resp)
	}
	// 原始响应应被保留以便排查（记录在日志的 ResponseBody）。
	if logger.last == nil || !strings.Contains(logger.last.ResponseBody, "Bad Gateway") {
		t.Errorf("应保留原始响应内容用于排查, got %+v", logger.last)
	}
}

// TestInvoke_MissingCodeField 缺少 code 字段：通信失败并保留原始响应（Requirement 1.11）。
func TestInvoke_MissingCodeField(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"msg":"missing code","responseId":"r1"}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)

	resp, err := client.Invoke("/api/stock/get", map[string]any{"a": 1})
	if err == nil {
		t.Fatal("缺少 code 字段时期望返回错误，却为 nil")
	}
	if resp != nil {
		t.Errorf("缺少 code 字段时不应返回业务结果，got %+v", resp)
	}
	if !strings.Contains(err.Error(), "code") {
		t.Errorf("错误信息应说明缺少 code 字段: %v", err)
	}
}

// TestInvoke_NumericCodeField 兼容 code 为数值类型（如 0 而非 "0"）的响应解析。
func TestInvoke_NumericCodeField(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"msg":"","responseId":"r2"}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)

	resp, err := client.Invoke("/api/stock/get", map[string]any{"a": 1})
	if err != nil {
		t.Fatalf("数值 code=0 应视为成功, got err=%v", err)
	}
	if resp.Code != "0" {
		t.Errorf("Code = %q, want \"0\"", resp.Code)
	}
}
