package config

import (
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
)

var _aliyunSms *dysmsapi20170525.Client

// GetAliyunSmsClient 获取sms客户端
func GetAliyunSmsClient() *dysmsapi20170525.Client {
	return _aliyunSms
}

type AliyunSms struct {
	AccessId     string `json:"access_id" yaml:"access_id"`
	AccessSecret string `json:"access_secret" yaml:"access_secret"`
	Region       string `json:"region" yaml:"region"`
	Signature    string `json:"signature" yaml:"signature"`
}
