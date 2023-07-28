package config

type AliyunAudit struct {
	AccessId     string `json:"access_id" yaml:"access_id"`
	AccessSecret string `json:"access_secret" yaml:"access_secret"`
	Region       string `json:"region" yaml:"region"`
	CallbackUrl  string `json:"callback_url" yaml:"callback_url"`
}
