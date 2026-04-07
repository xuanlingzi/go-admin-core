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
	"mime/multipart"
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
		"operationId":   "02",
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
	dataMap, ok := r.Data.(map[string]interface{})
	if !ok {
		return nil
	}
	list, _ := dataMap["terminalInfo"].([]interface{})
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
	RespCode    string      `json:"respCode"`
	RespMsg     string      `json:"respMsg"`
	ReqSerialNo string      `json:"reqSerialNo"`
	Version     string      `json:"version"`
	Data        interface{} `json:"data"`
}

func (r *leshuaCollectResp) deviceID() string {
	if r == nil || r.Data == nil {
		return ""
	}
	dataMap, ok := r.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	for _, key := range []string{"deviceId", "device_id"} {
		v, ok := dataMap[key]
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

// -------- 分账接口（全部使用聚合签名） --------

const (
	LedgerApplyPath     = "/api/share-merchant/elec-contract-accredit" // 商户开通分账申请（电子协议）
	LedgerQueryPath     = "/api/share-merchant/accreditQuery"          // 商户开通分账结果查询
	LedgerBindPath      = "/api/share-merchant/bind"                   // 分账关系绑定
	LedgerUnbindPath    = "/api/share-merchant/unbind"                 // 分账关系解绑
	LedgerQueryBindPath = "/api/share-merchant/queryBind"              // 分账关系绑定查询
)

// ApplyLedger 商户开通分账申请（电子协议）
// 文档: /api/share-merchant/elec-contract-accredit
func (m *Leshua) ApplyLedger(merchantID string, sharePercent string, insertFlag int, feeRate int, ledgerMethod int, authTypes string, callbackUrl string) (map[string]interface{}, error) {
	bizData := map[string]interface{}{
		"merchantId": merchantID,
		"baseInfo": map[string]interface{}{
			"shareRole": 0, // 目前只允许商户授权
		},
	}

	otherInfo := map[string]interface{}{}
	if sharePercent != "" {
		otherInfo["sharePercent"] = sharePercent
	}
	if authTypes != "" {
		otherInfo["authTypes"] = authTypes
	} else {
		otherInfo["authTypes"] = "1,2" // 默认使用手机+银行卡验证
	}
	if callbackUrl != "" {
		otherInfo["callBackUrl"] = callbackUrl
	}
	if len(otherInfo) > 0 {
		bizData["otherInfo"] = otherInfo
	}

	if insertFlag > 0 {
		bizData["insertFlag"] = insertFlag
	}
	if feeRate > 0 {
		bizData["feeRate"] = feeRate
	}
	if ledgerMethod > 0 {
		bizData["ledgerMethod"] = ledgerMethod
	}

	addr := m.CollectAddr + LedgerApplyPath
	raw, err := m.postAggregateJSON(addr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("分账开通申请失败: [%s] %s", raw.RespCode, raw.RespMsg)
	}
	if raw.Data == nil {
		return nil, fmt.Errorf("分账开通返回数据为空")
	}
	if dataMap, ok := raw.Data.(map[string]interface{}); ok {
		return dataMap, nil
	}
	return nil, fmt.Errorf("返回数据并非JSON对象")
}

// QueryLedgerStatus 商户开通分账结果查询
// 文档: /api/share-merchant/accreditQuery
func (m *Leshua) QueryLedgerStatus(merchantID string) (map[string]interface{}, error) {
	bizData := map[string]interface{}{
		"merchantId": merchantID,
	}

	addr := m.CollectAddr + LedgerQueryPath
	raw, err := m.postAggregateJSON(addr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("分账状态查询失败: [%s] %s", raw.RespCode, raw.RespMsg)
	}
	if raw.Data == nil {
		return nil, fmt.Errorf("分账状态查询返回数据为空")
	}
	if dataMap, ok := raw.Data.(map[string]interface{}); ok {
		return dataMap, nil
	}
	return nil, fmt.Errorf("返回数据并非JSON对象")
}

// BindLedgerReceiver 分账关系绑定
// 文档: /api/share-merchant/bind
func (m *Leshua) BindLedgerReceiver(merchantID1, merchantID2 string, protocolPic string, remark string) (map[string]interface{}, error) {
	bizData := map[string]interface{}{
		"merchantId1": merchantID1,
		"merchantId2": merchantID2,
	}
	if protocolPic != "" {
		bizData["protocolPic"] = protocolPic
	}
	if remark != "" {
		bizData["remark"] = remark
	}

	addr := m.CollectAddr + LedgerBindPath
	raw, err := m.postAggregateJSON(addr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("分账关系绑定失败: [%s] %s", raw.RespCode, raw.RespMsg)
	}
	result := map[string]interface{}{
		"respCode": raw.RespCode,
		"respMsg":  raw.RespMsg,
	}
	if raw.Data != nil {
		if dataMap, ok := raw.Data.(map[string]interface{}); ok {
			for k, v := range dataMap {
				result[k] = v
			}
		} else {
			result["data"] = raw.Data
		}
	}
	return result, nil
}

// UnbindLedgerReceiver 分账关系解绑
// 文档: /api/share-merchant/unbind
func (m *Leshua) UnbindLedgerReceiver(merchantID1, merchantID2 string, remark string) (map[string]interface{}, error) {
	bizData := map[string]interface{}{
		"merchantId1": merchantID1,
		"merchantId2": merchantID2,
	}
	if remark != "" {
		bizData["remark"] = remark
	}

	addr := m.CollectAddr + LedgerUnbindPath
	raw, err := m.postAggregateJSON(addr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("分账关系解绑失败: [%s] %s", raw.RespCode, raw.RespMsg)
	}
	result := map[string]interface{}{
		"respCode": raw.RespCode,
		"respMsg":  raw.RespMsg,
	}
	if raw.Data != nil {
		if dataMap, ok := raw.Data.(map[string]interface{}); ok {
			for k, v := range dataMap {
				result[k] = v
			}
		} else {
			result["data"] = raw.Data
		}
	}
	return result, nil
}

// QueryBindRelation 分账关系绑定查询
// 文档: /api/share-merchant/queryBind
func (m *Leshua) QueryBindRelation(merchantID1 string) (map[string]interface{}, error) {
	bizData := map[string]interface{}{
		"merchantId1": merchantID1,
	}

	addr := m.CollectAddr + LedgerQueryBindPath
	raw, err := m.postAggregateJSON(addr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("分账绑定查询失败: [%s] %s", raw.RespCode, raw.RespMsg)
	}
	result := map[string]interface{}{
		"respCode": raw.RespCode,
		"respMsg":  raw.RespMsg,
	}
	if raw.Data != nil {
		result["data"] = raw.Data
	}
	return result, nil
}

// -------- 订单分账接口（聚合签名） --------

const (
	SplitApplyPath       = "/api/share-merchant/multi-apply"
	SplitQueryPath       = "/api/share-merchant/multi-query"
	SplitCancelPath      = "/api/share-merchant/cancel"
	SplitRefundPath      = "/api/share-merchant/refund"
	SplitRefundQueryPath = "/api/share-merchant/refundQuery"
)

// postAggregateJSON 使用聚合常规签名发起 JSON 请求
// 签名方式: Base64(md5("lepos" + key + dataJSON).toLowerCase())
func (m *Leshua) postAggregateJSON(addr string, bizData interface{}) (*leshuaCollectResp, error) {
	dataBytes, err := json.Marshal(bizData)
	if err != nil {
		return nil, fmt.Errorf("构造请求参数失败: %w", err)
	}
	dataJSON := string(dataBytes)
	form := url.Values{}
	form.Set("agentId", m.CollectAgentID)
	form.Set("version", "1.0")
	form.Set("reqSerialNo", buildCollectReqSerialNo())
	form.Set("data", dataJSON)
	form.Set("sign", calcCollectSign(dataJSON, m.TradeKey))

	logger.Infof("[乐刷分账] POST %s\n请求data: %s", addr, dataJSON)

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
	logger.Infof("[乐刷分账] 响应(HTTP %d): %s", resp.StatusCode, string(body))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP 异常(%d): %s", resp.StatusCode, string(body))
	}
	var result leshuaCollectResp
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	return &result, nil
}

// ApplyOrderSplit 订单分账申请
func (m *Leshua) ApplyOrderSplit(merchantID, leshuaOrderID, thirdOrderID, thirdRoyaltyID string, shareDetail []map[string]interface{}, remark string) (map[string]interface{}, error) {
	bizData := map[string]interface{}{
		"merchantId":     merchantID,
		"leshuaOrderId":  leshuaOrderID,
		"thirdOrderId":   thirdOrderID,
		"thirdRoyaltyId": thirdRoyaltyID,
		"shareDetail":    shareDetail,
	}
	if remark != "" {
		bizData["Remark"] = remark
	}

	addr := m.CollectAddr + SplitApplyPath
	raw, err := m.postAggregateJSON(addr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("分账申请网关失败: [%s] %s", raw.RespCode, raw.RespMsg)
	}
	if raw.Data == nil {
		return nil, fmt.Errorf("分账申请返回数据为空")
	}
	if dataMap, ok := raw.Data.(map[string]interface{}); ok {
		return dataMap, nil
	}
	return nil, fmt.Errorf("返回数据并非JSON对象")
}

// QueryOrderSplit 订单分账结果查询
func (m *Leshua) QueryOrderSplit(merchantID, leshuaOrderID string, allRoyaltyFlag int, leshuaRoyaltyID, thirdRoyaltyID string) (map[string]interface{}, error) {
	bizData := map[string]interface{}{
		"merchantId":     merchantID,
		"leshuaOrderId":  leshuaOrderID,
		"allRoyaltyFlag": allRoyaltyFlag,
	}
	if leshuaRoyaltyID != "" {
		bizData["leshuaRoyaltyId"] = leshuaRoyaltyID
	}
	if thirdRoyaltyID != "" {
		bizData["thirdRoyaltyId"] = thirdRoyaltyID
	}

	addr := m.CollectAddr + SplitQueryPath
	raw, err := m.postAggregateJSON(addr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("分账查询网关失败: [%s] %s", raw.RespCode, raw.RespMsg)
	}
	if raw.Data == nil {
		return nil, fmt.Errorf("分账查询返回数据为空")
	}
	if dataMap, ok := raw.Data.(map[string]interface{}); ok {
		return dataMap, nil
	}
	return nil, fmt.Errorf("返回数据并非JSON对象")
}

// CancelOrderSplit 订单分账撤销
func (m *Leshua) CancelOrderSplit(merchantID, leshuaOrderID, leshuaRoyaltyID, thirdRoyaltyID string) (map[string]interface{}, error) {
	bizData := map[string]interface{}{
		"merchantId":     merchantID,
		"leshuaOrderId":  leshuaOrderID,
		"thirdRoyaltyId": thirdRoyaltyID,
	}
	if leshuaRoyaltyID != "" {
		bizData["leshuaRoyaltyId"] = leshuaRoyaltyID
	}

	addr := m.CollectAddr + SplitCancelPath
	raw, err := m.postAggregateJSON(addr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("分账撤销网关失败: [%s] %s", raw.RespCode, raw.RespMsg)
	}
	if raw.Data == nil {
		return nil, fmt.Errorf("分账撤销返回数据为空")
	}
	if dataMap, ok := raw.Data.(map[string]interface{}); ok {
		return dataMap, nil
	}
	return nil, fmt.Errorf("返回数据并非JSON对象")
}

// RefundOrderSplit 分账交易退款
func (m *Leshua) RefundOrderSplit(merchantID, thirdOrderID, leshuaOrderID, thirdRefundID string, refundAmount int64, refundMode string, thirdRoyaltyID string, refundDetails []map[string]interface{}, notifyUrl string) (map[string]interface{}, error) {
	bizData := map[string]interface{}{
		"merchantId":    merchantID,
		"thirdOrderId":  thirdOrderID,
		"leshuaOrderId": leshuaOrderID,
		"thirdRefundId": thirdRefundID,
		"refundAmount":  refundAmount,
	}
	if refundMode != "" {
		bizData["refundMode"] = refundMode
	}
	if thirdRoyaltyID != "" {
		bizData["thirdRoyaltyId"] = thirdRoyaltyID
	}
	if len(refundDetails) > 0 {
		bizData["refundDetails"] = refundDetails
	}
	if notifyUrl != "" {
		bizData["notifyUrl"] = notifyUrl
	}

	addr := m.CollectAddr + SplitRefundPath
	raw, err := m.postAggregateJSON(addr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("分账退款网关失败: [%s] %s", raw.RespCode, raw.RespMsg)
	}
	if raw.Data == nil {
		return nil, fmt.Errorf("分账退款返回数据为空")
	}
	if dataMap, ok := raw.Data.(map[string]interface{}); ok {
		return dataMap, nil
	}
	return nil, fmt.Errorf("返回数据并非JSON对象")
}

// QueryRefundOrderSplit 分账退款查询
func (m *Leshua) QueryRefundOrderSplit(merchantID, leshuaOrderID, thirdOrderID, thirdRefundID, leshuaRefundID string) (map[string]interface{}, error) {
	bizData := map[string]interface{}{
		"merchantId":    merchantID,
		"leshuaOrderId": leshuaOrderID,
		"thirdOrderId":  thirdOrderID,
		"thirdRefundId": thirdRefundID,
	}
	if leshuaRefundID != "" {
		bizData["leshuaRefundId"] = leshuaRefundID
	}

	addr := m.CollectAddr + SplitRefundQueryPath
	raw, err := m.postAggregateJSON(addr, bizData)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.RespCode) != "000000" {
		return nil, fmt.Errorf("分账退款查询网关失败: [%s] %s", raw.RespCode, raw.RespMsg)
	}
	if raw.Data == nil {
		return nil, fmt.Errorf("分账退款查询返回数据为空")
	}
	if dataMap, ok := raw.Data.(map[string]interface{}); ok {
		return dataMap, nil
	}
	return nil, fmt.Errorf("返回数据并非JSON对象")
}

// -------- 商户入驻管理 API --------

// postMchJSON 商户管理 API 通用 POST，签名：lepos+key+data 的 MD5 Base64。
func (m *Leshua) postMchJSON(endpoint string, data interface{}) (map[string]interface{}, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化业务数据失败: %w", err)
	}
	dataStr := string(dataBytes)

	reqSerialNo := buildCollectReqSerialNo()
	form := url.Values{}
	form.Set("agentId", m.CollectAgentID)
	form.Set("version", "2.0")
	form.Set("reqSerialNo", reqSerialNo)
	form.Set("sign", calcCollectSign(dataStr, m.TradeKey))
	form.Set("data", dataStr)

	requestURL := strings.TrimRight(m.CollectAddr, "/") + endpoint
	req, err := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	start := time.Now()
	resp, err := m.client.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		logger.Errorf("[乐刷商户] POST %s FAILED (%.0fms) err=%s", endpoint, float64(elapsed.Milliseconds()), err.Error())
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	logger.Infof("[乐刷商户] POST %s %d (%.0fms) body=%s", endpoint, resp.StatusCode, float64(elapsed.Milliseconds()), string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP状态码异常: %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %w", err)
	}
	if respCode, _ := result["respCode"].(string); respCode != "000000" {
		respMsg, _ := result["respMsg"].(string)
		return result, fmt.Errorf("乐刷业务失败[%s]: %s", respCode, respMsg)
	}
	return result, nil
}

// UploadMchPicture 上传图片到乐刷商户管理平台（multipart，不走 postMchJSON）
func (m *Leshua) UploadMchPicture(fileData []byte, fileName string) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	dataStr := "{}"
	_ = writer.WriteField("agentId", m.CollectAgentID)
	_ = writer.WriteField("version", "2.0")
	_ = writer.WriteField("reqSerialNo", buildCollectReqSerialNo())
	_ = writer.WriteField("sign", calcCollectSign(dataStr, m.TradeKey))
	_ = writer.WriteField("data", dataStr)

	part, err := writer.CreateFormFile("media", fileName)
	if err != nil {
		return "", fmt.Errorf("创建文件字段失败: %w", err)
	}
	if _, err = part.Write(fileData); err != nil {
		return "", fmt.Errorf("写入文件数据失败: %w", err)
	}
	writer.Close()

	requestURL := strings.TrimRight(m.CollectAddr, "/") + "/apiv2/picture/upload"
	req, err := http.NewRequest("POST", requestURL, &buf)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := m.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取上传响应失败: %w", err)
	}
	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析上传响应失败: %w", err)
	}
	if respCode, _ := result["respCode"].(string); respCode != "000000" {
		respMsg, _ := result["respMsg"].(string)
		return "", fmt.Errorf("乐刷业务失败[%s]: %s", respCode, respMsg)
	}
	dataMap, _ := result["data"].(map[string]interface{})
	photoUrl, _ := dataMap["photoUrl"].(string)
	if photoUrl == "" {
		return "", fmt.Errorf("上传成功但未返回图片URL")
	}
	return photoUrl, nil
}

func (m *Leshua) QueryRegion(parentCode string) (map[string]interface{}, error) {
	data := map[string]string{}
	if parentCode != "" {
		data["parentCode"] = parentCode
	}
	return m.postMchJSON("/data/area", data)
}

func (m *Leshua) QueryMcc(parentCode string) (map[string]interface{}, error) {
	data := map[string]string{}
	if parentCode != "" {
		data["parentCode"] = parentCode
	}
	return m.postMchJSON("/data/mcc", data)
}

func (m *Leshua) QueryBankBranch(bankName string, cityCode string) (map[string]interface{}, error) {
	data := map[string]string{}
	if bankName != "" {
		data["bankName"] = bankName
	}
	if cityCode != "" {
		data["cityCode"] = cityCode
	}
	return m.postMchJSON("/data/bankbranch2", data)
}

func (m *Leshua) RegisterMerchant(data map[string]interface{}) (map[string]interface{}, error) {
	return m.postMchJSON("/apiv2/merchant/register", data)
}

func (m *Leshua) QueryMchAudit(merchantId string) (map[string]interface{}, error) {
	return m.postMchJSON("/apiv2/merchant/audit_qry", map[string]string{"merchantId": merchantId})
}

func (m *Leshua) UpdateMerchant(merchantId string, data map[string]interface{}) (map[string]interface{}, error) {
	payload := cloneMchPayload(data)
	payload["merchantId"] = merchantId
	return m.postMchJSON("/apiv2/merchant/update", payload)
}

func (m *Leshua) QuerySettlementStatus(merchantId string) (map[string]interface{}, error) {
	return m.postMchJSON("/apiv2/risk-work-order/querySettlementStatus", map[string]string{"merchantId": merchantId})
}

func (m *Leshua) QueryMchFeeRate(merchantId string) (map[string]interface{}, error) {
	return m.postMchJSON("/apiv2/merchant/fee_qry", map[string]string{"merchantId": merchantId})
}

func (m *Leshua) QuerySubMerchant(merchantId string, channel string) (map[string]interface{}, error) {
	return m.postMchJSON("/apiv2/submch/query", map[string]string{
		"merchantId": merchantId,
		"channel":    channel,
	})
}

func (m *Leshua) OpenMerchant(merchantId string, data map[string]interface{}) (map[string]interface{}, error) {
	payload := cloneMchPayload(data)
	payload["merchantId"] = merchantId
	return m.postMchJSON("/apiv2/merchant/open", payload)
}

func (m *Leshua) UpdateMchIdCardInfo(merchantId string, data map[string]interface{}) (map[string]interface{}, error) {
	payload := cloneMchPayload(data)
	payload["merchantId"] = merchantId
	return m.postMchJSON("/apiv2/merchant/updateIdCardInfo", payload)
}

func (m *Leshua) UpdateMchContactInfo(merchantId string, data map[string]interface{}) (map[string]interface{}, error) {
	payload := cloneMchPayload(data)
	payload["merchantId"] = merchantId
	return m.postMchJSON("/apiv2/merchant/updateContactInfo", payload)
}

func (m *Leshua) SaveOrUpdateWxAuthFurtherInfo(merchantId string, data map[string]interface{}) (map[string]interface{}, error) {
	payload := cloneMchPayload(data)
	payload["merchantId"] = merchantId
	return m.postMchJSON("/apiv2/merchant-wx-auth-further-info/saveOrUpdateInfo", payload)
}

func (m *Leshua) UpdateMchShortName(merchantId string, channel string, data map[string]interface{}) (map[string]interface{}, error) {
	payload := cloneMchPayload(data)
	payload["merchantId"] = merchantId

	endpoint := ""
	switch channel {
	case "weixin":
		endpoint = "/apiv2/merchant/merchantUpdateShortname"
	case "alipay":
		endpoint = "/apiv2/submch/syncZfbMsg"
	default:
		return nil, fmt.Errorf("不支持的渠道类型: %s", channel)
	}

	return m.postMchJSON(endpoint, payload)
}

func (m *Leshua) WechatSubjectPreCheck(merchantId string) (map[string]interface{}, error) {
	return m.postMchJSON("/apiv2/wechat/subject/preCheck", map[string]string{"merchantId": merchantId})
}

func (m *Leshua) SetWechatPayConfig(merchantId string, data map[string]interface{}) (map[string]interface{}, error) {
	payload := cloneMchPayload(data)
	payload["merchantId"] = merchantId
	return m.postMchJSON("/apiv2/wechat/wxpayconfig", payload)
}

func (m *Leshua) ReReportSubMerchant(merchantId string, channel string) (map[string]interface{}, error) {
	return m.postMchJSON("/apiv2/submch/reReport", map[string]string{
		"merchantId": merchantId,
		"channel":    channel,
	})
}

func (m *Leshua) ApplyMerchantAcqProtocol(merchantId string) (map[string]interface{}, error) {
	return m.postMchJSON("/apiv2/merchantAcqProtocol/apply", map[string]string{"merchantId": merchantId})
}

func (m *Leshua) QueryMerchantAcqProtocol(merchantId string, contractId string) (map[string]interface{}, error) {
	payload := map[string]string{"merchantId": merchantId}
	if contractId != "" {
		payload["contractId"] = contractId
	}
	return m.postMchJSON("/apiv2/merchantAcqProtocol/signQuery", payload)
}

func (m *Leshua) ApplyWechatSubjectVerify(merchantId string, microBizType string, confirmMchidList string) (map[string]interface{}, error) {
	data := map[string]string{"merchantId": merchantId}
	if microBizType != "" {
		data["microBizType"] = microBizType
	}
	if confirmMchidList != "" {
		data["confirmMchidList"] = confirmMchidList
	}
	return m.postMchJSON("/apiv2/wechat/subject/apply", data)
}

func (m *Leshua) QueryWechatSubjectVerify(merchantId string) (map[string]interface{}, error) {
	return m.postMchJSON("/apiv2/wechat/subject/query", map[string]string{"merchantId": merchantId})
}

func (m *Leshua) ApplyAlipayVerification(merchantId string, confirmMchidList string) (map[string]interface{}, error) {
	data := map[string]string{"merchantId": merchantId}
	if confirmMchidList != "" {
		data["confirmMchidList"] = confirmMchidList
	}
	return m.postMchJSON("/apiv2/zfbVerify/apply", data)
}

func (m *Leshua) QueryAlipayVerification(merchantId string, businessCode string, applymentId string) (map[string]interface{}, error) {
	data := map[string]string{"merchantId": merchantId}
	if businessCode != "" {
		data["businessCode"] = businessCode
	}
	if applymentId != "" {
		data["applymentId"] = applymentId
	}
	return m.postMchJSON("/apiv2/zfbVerify/queryApplyStatus", data)
}

func (m *Leshua) CancelWechatSubjectVerify(merchantId string) (map[string]interface{}, error) {
	return m.postMchJSON("/apiv2/wechat/subject/cancel", map[string]string{"merchantId": merchantId})
}

func (m *Leshua) CancelAlipayVerification(merchantId string, businessCode string, applymentId string) (map[string]interface{}, error) {
	data := map[string]string{"merchantId": merchantId}
	if businessCode != "" {
		data["businessCode"] = businessCode
	}
	if applymentId != "" {
		data["applymentId"] = applymentId
	}
	return m.postMchJSON("/apiv2/zfbVerify/revoke", data)
}

func cloneMchPayload(data map[string]interface{}) map[string]interface{} {
	payload := make(map[string]interface{}, len(data)+1)
	for key, value := range data {
		payload[key] = value
	}
	return payload
}
