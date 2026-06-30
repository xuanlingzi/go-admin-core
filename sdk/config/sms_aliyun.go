package config

type AliyunSms struct {
	Addr         string `json:"addr" yaml:"addr"`
	AccessId     string `json:"access_id" yaml:"access_id"`
	AccessSecret string `json:"access_secret" yaml:"access_secret"`
	Region       string `json:"region" yaml:"region"`
	Signature    string `json:"signature" yaml:"signature"`
}
