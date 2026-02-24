package config

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
