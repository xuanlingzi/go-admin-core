package config

type Sms struct {
	Tencent *Tencent `json:"tencent" yaml:"tencent"`
	Aliyun  *Aliyun  `json:"aliyun" yaml:"aliyun"`
}

var SmsConfig = new(Sms)
