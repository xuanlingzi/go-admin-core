package config

type Secret struct {
	AppId         string `json:"app_id" yaml:"app_id"`
	AppSecret     string `json:"app_secret" yaml:"app_secret"`
	AesSecret     string `json:"aes_secret,omitempty" yaml:"aes_secret"`
	RsaPubKey     string `json:"rsa_pub_key" yaml:"rsa_pub_key"`
	RsaPrivateKey string `json:"rsa_private_key" yaml:"rsa_private_key"`
}

var SecretConfig = new(Secret)
