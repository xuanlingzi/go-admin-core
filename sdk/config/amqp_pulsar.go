package config

type Pulsar struct {
	Addr      string `yaml:"addr"`
	AppId     string `yaml:"appId"`
	Namespace string `yaml:"namespace"`
	Role      string `yaml:"role"`
	SecretKey string `yaml:"secretKey"`
}
