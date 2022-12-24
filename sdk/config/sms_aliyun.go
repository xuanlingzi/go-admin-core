package config

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

var _aliyun_sms *dysmsapi20170525.Client

// GetAliyunClient 获取sms客户端
func GetAliyunClient() *dysmsapi20170525.Client {
	return _aliyun_sms
}

type Aliyun struct {
	AccessId     string `json:"access_id" yaml:"access_id"`
	AccessSecret string `json:"access_secret" yaml:"access_secret"`
	Region       string `json:"region" yaml:"region"`
	Signature    string `json:"signature" yaml:"signature"`
}

func (e Aliyun) GetAliyunOptions() *openapi.Config {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(e.AccessId),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(e.AccessSecret),
		// 访问的域名
		Endpoint: tea.String(e.Region),
	}
	return config
}
