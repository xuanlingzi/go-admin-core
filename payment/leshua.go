package payment

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/xuanlingzi/go-admin-core/logger"
)

var _leshua *Leshua

const (
	PaymentLePosPayPath = "/cgi-bin/lepos_pay_gateway.cgi"
	CollectQueryPath    = "/apiv2/terminal/queryMerchantChannelTerminals"
	CollectRegisterPath = "/apiv2/terminal/collect"
)

// Leshua adapter 封装所有对乐刷网关的 HTTP 调用
type Leshua struct {
	PaymentAddr    string
	TradeKey       string
	NotifyKey      string
	SignType       string
	NotifyURL      string
	CollectAddr    string
	CollectAgentID string
	client         *http.Client
}

// NewLeshua 构建 Leshua adapter 并保存到包级变量
func NewLeshua(paymentAddr string, tradeKey string, notifyKey string, signType string, notifyURL string, collectAddr string, collectAgentID string) *Leshua {
	l := &Leshua{
		PaymentAddr:    paymentAddr,
		TradeKey:       tradeKey,
		NotifyKey:      notifyKey,
		SignType:       signType,
		NotifyURL:      notifyURL,
		CollectAddr:    collectAddr,
		CollectAgentID: collectAgentID,
		client:         &http.Client{Timeout: 8 * time.Second},
	}
	_leshua = l
	return l
}

func (*Leshua) String() string {
	return "leshua"
}

// GetClient 暴露原生 http.Client，满足 AdapterLeshuaService 接口
func (m *Leshua) GetClient() interface{} {
	return m.client
}

// Close 释放连接
func (m *Leshua) Close() error {
	if m.client != nil {
		m.client.CloseIdleConnections()
		m.client = nil
	}
	return nil
}

// -------- 支付网关调用 --------

// PayByAuthCode 发起付款码支付，返回原始响应 map，由 Service 层解析
func (m *Leshua) PayByAuthCode(merchantID string, thirdOrderID string, authCode string, amountFen int64, body string, shopNo string, posNo string, terminalType string, terminalID string, goodsDetail string, sceneInfo string) (map[string]string, error) {
	terminalInfo, _ := json.Marshal(map[string]string{
		"device_type": terminalType,
		"device_id":   terminalID,
	})
	params := map[string]string{
		"service":        "upload_authcode",
		"sign_type":      "MD5",
		"auth_code":      authCode,
		"merchant_id":    merchantID,
		"third_order_id": thirdOrderID,
		"amount":         strconv.FormatInt(amountFen, 10),
		"nonce_str":      m.nonce(16),
		"body":           body,
		"shop_no":        shopNo,
		"pos_no":         posNo,
		"terminal_info":  string(terminalInfo),
	}
	if goodsDetail != "" {
		params["goods_detail"] = goodsDetail
	}
	if sceneInfo != "" {
		params["scene_info"] = sceneInfo
	}
	if m.NotifyURL != "" {
		params["notify_url"] = m.NotifyURL
	}
	params["sign"] = m.calcTradeSign(params)

	paymentAddr := m.PaymentAddr + PaymentLePosPayPath
	return m.postXML(paymentAddr, params)
}

// QueryOrder 查询订单状态
func (m *Leshua) QueryOrder(merchantID, thirdOrderID string) (map[string]string, error) {
	params := map[string]string{
		"service":        "query_order",
		"sign_type":      "MD5",
		"merchant_id":    merchantID,
		"third_order_id": thirdOrderID,
		"nonce_str":      m.nonce(16),
	}
	params["sign"] = m.calcTradeSign(params)
	paymentAddr := m.PaymentAddr + PaymentLePosPayPath
	return m.postXML(paymentAddr, params)
}

// CloseOrder 关闭订单；leshuaOrderID 优先，否则用 thirdOrderID
func (m *Leshua) CloseOrder(merchantID, thirdOrderID, leshuaOrderID string) (map[string]string, error) {
	params := map[string]string{
		"service":     "close_order",
		"sign_type":   "MD5",
		"merchant_id": merchantID,
		"nonce_str":   m.nonce(16),
	}
	if leshuaOrderID != "" {
		params["leshua_order_id"] = leshuaOrderID
	} else {
		params["third_order_id"] = thirdOrderID
	}
	params["sign"] = m.calcTradeSign(params)
	paymentAddr := m.PaymentAddr + PaymentLePosPayPath
	return m.postXML(paymentAddr, params)
}

// VerifyNotifySign 验证异步通知签名
func (m *Leshua) VerifyNotifySign(payload map[string]string) error {
	key := m.NotifyKey
	if strings.TrimSpace(key) == "" {
		key = m.TradeKey
	}
	return m.verifyResponseSign(payload, key)
}

// VerifyResponseSign 验证同步响应签名
func (m *Leshua) VerifyResponseSign(payload map[string]string) error {
	return m.verifyResponseSign(payload, m.TradeKey)
}

// ParseNotifyXML 解析异步通知 XML 为 map
func (m *Leshua) ParseNotifyXML(raw []byte) (map[string]string, error) {
	return parseLeshuaXMLToMap(raw)
}

// -------- 终端采集 --------

// CollectTerminalID 采集/更新终端设备号
func (m *Leshua) CollectTerminalID(merchantID, serialNum, existingDeviceID string) (string, error) {

	operationID := "00"
	if existingDeviceID != "" {
		operationID = "01"
	}
	bizData := map[string]string{
		"merchantId":    merchantID,
		"operationId":   operationID,
		"terminalState": "00",
		"terminalType":  "11",
	}
	if serialNum != "" {
		bizData["serialNum"] = serialNum
	}
	if existingDeviceID != "" {
		bizData["deviceId"] = existingDeviceID
	}

	collectAddr := m.CollectAddr + CollectRegisterPath
	raw, err := m.postCollect(collectAddr, bizData)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(raw.RespCode) != "0" {
		return "", fmt.Errorf("终端采集失败: %s %s", raw.RespMsg, raw.RespCode)
	}
	deviceID := raw.deviceID()
	if deviceID == "" {
		return "", fmt.Errorf("终端采集未返回 device_id")
	}
	return deviceID, nil
}

// DeregisterTerminal 停用/注销终端
func (m *Leshua) DeregisterTerminal(merchantID, deviceID string) (string, error) {
	bizData := map[string]string{
		"merchantId":    merchantID,
		"operationId":   "01",
		"terminalState": "01",
		"terminalType":  "11",
		"deviceId":      deviceID,
	}

	collectAddr := m.CollectAddr + CollectRegisterPath
	raw, err := m.postCollect(collectAddr, bizData)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(raw.RespCode) != "0" {
		return "", fmt.Errorf("终端注销失败: %s %s", raw.RespMsg, raw.RespCode)
	}
	return strings.TrimSpace(raw.deviceID()), nil
}

// QueryTerminalReport 查询终端入网上报结果
func (m *Leshua) QueryTerminalReport(merchantID, serialNum, deviceID string) (map[string]interface{}, error) {
	queryAddr := m.CollectAddr + CollectQueryPath
	bizData := map[string]string{"merchantId": merchantID}
	if serialNum != "" {
		bizData["serialNum"] = serialNum
	}
	if deviceID != "" {
		bizData["deviceId"] = deviceID
	}
	raw, err := m.postCollect(queryAddr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("终端查询失败: %s %s", raw.RespMsg, raw.RespCode)
	}
	info := raw.pickTerminalInfo(serialNum, deviceID)
	if info == nil {
		return nil, fmt.Errorf("终端查询无匹配结果")
	}
	return info, nil
}

// pickTerminalInfo 从响应 Data["terminalInfo"] 中按 deviceID > serialNum > 首条的优先级匹配
func (r *leshuaCollectResp) pickTerminalInfo(serialNum, deviceID string) map[string]interface{} {
	if r == nil || r.Data == nil {
		return nil
	}
	list, _ := r.Data["terminalInfo"].([]interface{})
	var items []map[string]interface{}
	for _, raw := range list {
		if item, ok := raw.(map[string]interface{}); ok {
			items = append(items, item)
		}
	}
	for _, item := range items {
		if deviceID != "" && strings.EqualFold(cast.ToString(item["deviceId"]), deviceID) {
			return item
		}
	}
	for _, item := range items {
		if serialNum != "" && strings.EqualFold(cast.ToString(item["serialNum"]), serialNum) {
			return item
		}
	}
	if len(items) > 0 {
		return items[0]
	}
	return nil
}

// -------- 内部 HTTP 调用 --------

func (m *Leshua) postXML(addr string, params map[string]string) (map[string]string, error) {
	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}
	requestBody := form.Encode()
	logger.Infof("[乐刷] POST %s\n请求Body: %s", addr, requestBody)

	req, err := http.NewRequest(http.MethodPost, addr, strings.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	logger.Infof("[乐刷] 响应(HTTP %d):\n%s", resp.StatusCode, string(body))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP 异常(%d): %s", resp.StatusCode, string(body))
	}
	result, err := parseLeshuaXMLToMap(body)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	return result, nil
}

type leshuaCollectResp struct {
	RespCode    string                 `json:"respCode"`
	RespMsg     string                 `json:"respMsg"`
	ReqSerialNo string                 `json:"reqSerialNo"`
	Version     string                 `json:"version"`
	Data        map[string]interface{} `json:"data"`
}

func (r *leshuaCollectResp) deviceID() string {
	if r == nil || r.Data == nil {
		return ""
	}
	for _, key := range []string{"deviceId", "device_id"} {
		v, ok := r.Data[key]
		if !ok {
			continue
		}
		switch val := v.(type) {
		case string:
			if s := strings.TrimSpace(val); s != "" {
				return s
			}
		case float64:
			return strconv.FormatInt(int64(val), 10)
		case int:
			return strconv.Itoa(val)
		}
	}
	return ""
}

func (m *Leshua) postCollect(addr string, bizData map[string]string) (*leshuaCollectResp, error) {
	dataBytes, err := json.Marshal(bizData)
	if err != nil {
		return nil, fmt.Errorf("构造请求参数失败: %w", err)
	}
	dataJSON := string(dataBytes)
	form := url.Values{}
	form.Set("agentId", m.CollectAgentID)
	form.Set("version", "2.0")
	form.Set("reqSerialNo", buildCollectReqSerialNo())
	form.Set("data", dataJSON)
	form.Set("sign", calcCollectSign(dataJSON, m.TradeKey))

	req, err := http.NewRequest(http.MethodPost, addr, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP 异常(%d): %s", resp.StatusCode, string(body))
	}
	var result leshuaCollectResp
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	return &result, nil
}

// -------- 签名工具 --------

func (m *Leshua) calcTradeSign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || strings.TrimSpace(params[k]) == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf strings.Builder
	for _, k := range keys {
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(strings.TrimSpace(params[k]))
		buf.WriteString("&")
	}
	buf.WriteString("key=")
	buf.WriteString(m.TradeKey)
	h := md5.New()
	h.Write([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func (m *Leshua) verifyResponseSign(data map[string]string, key string) error {
	sign := strings.TrimSpace(data["sign"])
	if sign == "" {
		return fmt.Errorf("响应缺少签名")
	}
	excluded := map[string]bool{"sign": true, "leshua": true, "resp_code": true}
	keys := make([]string, 0, len(data))
	for k := range data {
		if !excluded[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	var buf strings.Builder
	for _, k := range keys {
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(data[k])
		buf.WriteString("&")
	}
	buf.WriteString("key=")
	buf.WriteString(key)
	h := md5.New()
	h.Write([]byte(buf.String()))
	calc := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	if !strings.EqualFold(sign, calc) {
		return fmt.Errorf("签名验证失败")
	}
	return nil
}

func (m *Leshua) nonce(n int) string {
	return uuid.New().String()[:n]
}

// -------- 公共工具 --------

func parseLeshuaXMLToMap(raw []byte) (map[string]string, error) {
	d := xml.NewDecoder(bytes.NewReader(raw))
	result := make(map[string]string)
	var tagName string
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch token := t.(type) {
		case xml.StartElement:
			tagName = token.Name.Local
		case xml.CharData:
			if tagName != "" && tagName != "xml" {
				// 空值字段也必须保留，乐刷验签规则：值为空的参数不剔除
				result[tagName] = strings.TrimSpace(string(token))
			}
		case xml.EndElement:
			tagName = ""
		}
	}
	return result, nil
}

func buildCollectReqSerialNo() string {
	now := time.Now()
	millis := now.Nanosecond() / int(time.Millisecond)
	seq := now.UnixNano() % 100000
	return fmt.Sprintf("%s%03d%05d", now.Format("20060102150405"), millis, seq)
}

func calcCollectSign(dataJSON, key string) string {
	sum := md5.Sum([]byte("lepos" + key + dataJSON))
	md5Str := hex.EncodeToString(sum[:])
	return base64.StdEncoding.EncodeToString([]byte(md5Str))
}
