package config

type WeChat struct {
	Platforms *[]WeChatOption `json:"platforms" yaml:"platforms"`
}

var WeChatConfig = new(WeChat)

type WeChatOption struct {
	Scope        string `json:"scope" yaml:"scope"`
	Addr         string `json:"addr" yaml:"addr"`
	AppId        string `json:"app_id" yaml:"app_id"`
	AppSecret    string `json:"app_secret" yaml:"app_secret"`
	AesKey       string `json:"aes_key" yaml:"aes_key"`
	Token        string `json:"token" yaml:"token"`
	CallbackAddr string `json:"callback_addr" yaml:"callback_addr" yaml:"callback_addr"`
}
