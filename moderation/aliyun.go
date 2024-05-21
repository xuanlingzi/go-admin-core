package moderation

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
)

var _aliyunAudit *green.Client

func GetAliyunAuditClient() *green.Client {
	return _aliyunAudit
}

type AliyunAuditClient struct {
	client      *green.Client
	accessId    string
	callbackUrl string
}

func NewAliyunAudit(client *green.Client, accessId, accessSecret, region, callbackUrl string) *AliyunAuditClient {
	var err error
	if client == nil {
		client, err = green.NewClientWithAccessKey(region, accessId, accessSecret)
		if err != nil {
			panic(fmt.Sprintf("Aliyun audit init error: %v", err))
		}
		_aliyunAudit = client
	}

	r := &AliyunAuditClient{
		client:      client,
		accessId:    accessId,
		callbackUrl: callbackUrl,
	}
	return r
}

func (m *AliyunAuditClient) String() string {
	return m.accessId
}

func (m *AliyunAuditClient) Check() bool {
	return m.client != nil
}

func (m *AliyunAuditClient) Close() {

}

func (rc *AliyunAuditClient) AuditText(content string, result *int, label *string, score *int, detail *string, jobId *string) error {

	return nil
}

func (rc *AliyunAuditClient) AuditImage(url string, fileSize int, result *int, label *string, score *int, detail *string, jobId *string) error {

	return nil
}

func (rc *AliyunAuditClient) AuditVideo(url string, frame int32, jobId *string) error {

	return nil
}

func (rc *AliyunAuditClient) AuditResult(body *[]byte, result *int, label *string, score *int, detail *string, jobId *string) error {

	return nil
}

// GetClient 暴露原生client
func (m *AliyunAuditClient) GetClient() interface{} {
	return m.client
}
