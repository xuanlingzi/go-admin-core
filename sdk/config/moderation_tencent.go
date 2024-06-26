package config

type TencentAudit struct {
	AccessKey   string `json:"access_key" yaml:"access_key"`
	SecretKey   string `json:"secret_key" yaml:"secret_key"`
	Region      string `json:"region" yaml:"region"`
	CosUrl      string `json:"cos_url" yaml:"cos_url"`
	CiUrl       string `json:"ci_url" yaml:"ci_url"`
	CallbackUrl string `json:"callback_url" yaml:"callback_url"`
}
