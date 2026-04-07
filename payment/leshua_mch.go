package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/xuanlingzi/go-admin-core/logger"
)

// -------- 商户入驻管理 API --------

// PostMchJSON 商户管理 API 通用 POST（对外导出，供 API 层直接透传）
// 签名：lepos+key+data 的 MD5 Base64，URL 前缀走 CollectAddr
func (m *Leshua) PostMchJSON(endpoint string, data interface{}) (map[string]interface{}, error) {
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


// UploadMchPicture 上传图片到乐刷商户管理平台（multipart，不走 PostMchJSON）
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

// RegisterMerchant 商户入驻（进件）—— service 层调用，保留
func (m *Leshua) RegisterMerchant(data map[string]interface{}) (map[string]interface{}, error) {
	return m.PostMchJSON("/apiv2/merchant/register", data)
}

// QueryMchAudit 商户审核状态查询 —— service 层调用，保留
func (m *Leshua) QueryMchAudit(merchantId string) (map[string]interface{}, error) {
	return m.PostMchJSON("/apiv2/merchant/audit_qry", map[string]string{"merchantId": merchantId})
}

// UpdateMerchant 商户信息修改 —— service 层调用，保留
func (m *Leshua) UpdateMerchant(merchantId string, data map[string]interface{}) (map[string]interface{}, error) {
	data["merchantId"] = merchantId
	return m.PostMchJSON("/apiv2/merchant/update", data)
}

// ApplyWechatSubjectVerify 微信实名认证申请 —— service 层调用，保留
func (m *Leshua) ApplyWechatSubjectVerify(merchantId string, microBizType string, confirmMchidList string) (map[string]interface{}, error) {
	data := map[string]string{"merchantId": merchantId}
	if microBizType != "" {
		data["microBizType"] = microBizType
	}
	if confirmMchidList != "" {
		data["confirmMchidList"] = confirmMchidList
	}
	return m.PostMchJSON("/apiv2/wechat/subject/apply", data)
}

// QueryWechatSubjectVerify 微信实名认证状态查询 —— service 层调用，保留
func (m *Leshua) QueryWechatSubjectVerify(merchantId string) (map[string]interface{}, error) {
	return m.PostMchJSON("/apiv2/wechat/subject/query", map[string]string{"merchantId": merchantId})
}

// ApplyAlipayVerification 支付宝实名认证申请 —— service 层调用，保留
func (m *Leshua) ApplyAlipayVerification(merchantId string, confirmMchidList string) (map[string]interface{}, error) {
	data := map[string]string{"merchantId": merchantId}
	if confirmMchidList != "" {
		data["confirmMchidList"] = confirmMchidList
	}
	return m.PostMchJSON("/apiv2/zfbVerify/apply", data)
}

// QueryAlipayVerification 支付宝实名认证状态查询 —— service 层调用，保留
func (m *Leshua) QueryAlipayVerification(merchantId string, businessCode string, applymentId string) (map[string]interface{}, error) {
	data := map[string]string{"merchantId": merchantId}
	if businessCode != "" {
		data["businessCode"] = businessCode
	}
	if applymentId != "" {
		data["applymentId"] = applymentId
	}
	return m.PostMchJSON("/apiv2/zfbVerify/queryApplyStatus", data)
}

// CancelWechatSubjectVerify 微信实名认证撤销 —— service 层调用，保留
func (m *Leshua) CancelWechatSubjectVerify(merchantId string) (map[string]interface{}, error) {
	return m.PostMchJSON("/apiv2/wechat/subject/cancel", map[string]string{"merchantId": merchantId})
}

// CancelAlipayVerification 支付宝实名认证撤销 —— service 层调用，保留
func (m *Leshua) CancelAlipayVerification(merchantId string, businessCode string, applymentId string) (map[string]interface{}, error) {
	data := map[string]string{"merchantId": merchantId}
	if businessCode != "" {
		data["businessCode"] = businessCode
	}
	if applymentId != "" {
		data["applymentId"] = applymentId
	}
	return m.PostMchJSON("/apiv2/zfbVerify/revoke", data)
}

// WechatSubjectPreCheck 微信实名认证前置检查 —— service 层调用，保留
func (m *Leshua) WechatSubjectPreCheck(merchantId string) (map[string]interface{}, error) {
	return m.PostMchJSON("/apiv2/wechat/subject/preCheck", map[string]string{"merchantId": merchantId})
}
