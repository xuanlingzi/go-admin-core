package config

type WeChat struct {
	Open *WeChatOption `json:"open,omitempty" yaml:"open"`
	Mp   *WeChatOption `json:"mp,omitempty" yaml:"mp"`
}

var WeChatConfig = new(WeChat)

type WeChatOption struct {
	CallbackAddr string `json:"callback_addr" yaml:"callback_addr" yaml:"callback_addr"`
	AppId        string `json:"app_id" yaml:"app_id"`
	AppSecret    string `json:"app_secret" yaml:"app_secret"`
	AesKey       string `json:"aes_key" yaml:"aes_key"`
	Token        string `json:"token" yaml:"token"`
}
