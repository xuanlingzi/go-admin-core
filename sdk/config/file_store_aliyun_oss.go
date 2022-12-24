package config

type AliyunOss struct {
	AccessId     string `json:"access_id" yaml:"access_id"`
	AccessSecret string `json:"access_secret" yaml:"access_secret"`
	Bucket       string `json:"bucket" yaml:"bucket"`
	Endpoint     string `json:"endpoint" yaml:"endpoint"`
}
