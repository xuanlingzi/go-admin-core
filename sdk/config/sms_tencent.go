package config

import (
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

var _tencentSms *sms.Client

// GetTencentSmsClient 获取sms客户端
func GetTencentSmsClient() *sms.Client {
	return _tencentSms
}

type TencentSms struct {
	SecretId  string `json:"secret_id" yaml:"secret_id"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	AppId     string `json:"app_id" yaml:"app_id"`
	AppKey    string `json:"app_key" yaml:"app_key"`
	Region    string `json:"region" yaml:"region"`
	Addr      string `json:"addr" yaml:"addr"`
	Alg       string `json:"alg" yaml:"alg"`
	Signature string `json:"signature" yaml:"signature"`
}
