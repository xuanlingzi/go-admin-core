package config

type Sms struct {
	Tencent *TencentSms `json:"tencent" yaml:"tencent"`
	Aliyun  *AliyunSms  `json:"aliyun" yaml:"aliyun"`
}

var SmsConfig = new(Sms)
