package moderation

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
)

var _tencentAudit *cos.Client

func GetTencentCosClient() *cos.Client {
	return _tencentAudit
}

type TencentAuditClient struct {
	client *cos.Client
}

func NewTencentAudit(client *cos.Client, base *cos.BaseURL, transport *cos.AuthorizationTransport) *TencentAuditClient {
	if client == nil {
		client = cos.NewClient(base, &http.Client{
			Transport: transport,
		})
	}

	r := &TencentAuditClient{
		client: client,
	}
	return r
}

func (rc *TencentAuditClient) String() string {
	return "tencent_audit"
}

func (rc *TencentAuditClient) Check() bool {
	return rc.client != nil
}

func (rc *TencentAuditClient) Close() {

}

func (rc *TencentAuditClient) AuditText(content string, suggestion *string, label *string, detail *string) error {

	return nil
}

func (rc *TencentAuditClient) AuditImage(url string, suggestion *string, label *string, detail *string) error {

	return nil
}

func (rc *TencentAuditClient) AuditVideo(url string, frame int32) error {

	return nil
}

// GetClient 暴露原生client
func (rc *TencentAuditClient) GetClient() interface{} {
	return rc.client
}
