package config

type HuaweiFile struct {
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Bucket    string `json:"bucket" yaml:"bucket"`
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	Region    string `json:"region" yaml:"region"`
}
