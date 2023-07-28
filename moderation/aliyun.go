package moderation

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
)

var _aliyunAudit *green.Client

func GetAliyunAuditClient() *green.Client {
	return _aliyunAudit
}

type AliyunAuditClient struct {
	client      *green.Client
	callbackUrl string
}

func NewAliyunAudit(client *green.Client, region string, accessId, accessSecret, callbackUrl string) *AliyunAuditClient {
	var err error
	if client == nil {
		client, err = green.NewClientWithAccessKey(region, accessId, accessSecret)
		if err != nil {
			panic("Failed to setup Aliyun Audit client.")
		}
	}

	r := &AliyunAuditClient{
		client:      client,
		callbackUrl: callbackUrl,
	}
	return r
}

func (m *AliyunAuditClient) String() string {
	return "aliyun_audit"
}

func (m *AliyunAuditClient) Check() bool {
	return m.client != nil
}

func (m *AliyunAuditClient) Close() {

}

func (rc *AliyunAuditClient) AuditText(content string, suggestion *string, label *string, detail *string) error {

	return nil
}

func (rc *AliyunAuditClient) AuditImage(url string, suggestion *string, label *string, detail *string) error {

	return nil
}

func (rc *AliyunAuditClient) AuditVideo(url string, frame int32) error {

	return nil
}

// GetClient 暴露原生client
func (m *AliyunAuditClient) GetClient() interface{} {
	return m.client
}
