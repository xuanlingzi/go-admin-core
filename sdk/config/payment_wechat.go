package config

type WeChatPayOption struct {
	MerchantId     string `yaml:"merchant_id" json:"merchant_id"`
	AppId          string `yaml:"app_id" json:"app_id"`
	PrivateKeyPath string `yaml:"private_key_path" json:"private_key_path"`
	SerialNumber   string `yaml:"serial_no" json:"serial_no"`
	ApiV3Key       string `yaml:"api_v3_key" json:"api_v3_key"`
	CallbackAddr   string `yaml:"callback_addr" json:"callback_addr"`
	WeChatCertPath string `yaml:"wechat_cert_path" json:"wechat_cert_path"`
}
