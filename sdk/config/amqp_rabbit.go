package config

type Rabbit struct {
	Addr      string `json:"addr" yaml:"addr"`
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
}
