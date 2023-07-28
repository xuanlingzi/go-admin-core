package config

type Moderation struct {
	Tencent *TencentAudit `json:"tencent" yaml:"tencent"`
	Aliyun  *AliyunAudit  `json:"aliyun" yaml:"aliyun"`
	Huawei  *HuaweiAudit  `json:"huawei" yaml:"huawei"`
}

var ModerationConfig = new(Moderation)
