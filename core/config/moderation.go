package config

type Moderation struct {
	Tencent *TencentAudit `json:"tencent" yaml:"tencent"`
}

var ModerationConfig = new(Moderation)
