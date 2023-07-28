package config

type HuaweiAudit struct {
	AccessKey   string `json:"access_key" yaml:"access_key"`
	SecretKey   string `json:"secret_key" yaml:"secret_key"`
	Region      string `json:"region" yaml:"region"`
	CallbackUrl string `json:"callback_url" yaml:"callback_url"`
}
