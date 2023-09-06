package config

type TencentFile struct {
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	CosUrl    string `json:"cos_url" yaml:"cos_url"`
	CiUrl     string `json:"ci_url" yaml:"ci_url"`
	Region    string `json:"region" yaml:"region"`
}
