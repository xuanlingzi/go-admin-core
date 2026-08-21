package config

type WeChat struct {
	Platforms *[]WeChatOption `json:"platforms" yaml:"platforms"`
	// CryptoKey 解密 wechat_platform 表密文列（app_secret/aes_key/token）的口令，
	// 与 rktv-console 的 WECHAT_CRYPTO_KEY 一致。为空则不启用「读库」，回退 yml。
	CryptoKey string `json:"crypto_key" yaml:"crypto_key"`
}

var WeChatConfig = new(WeChat)

type WeChatOption struct {
	Key          string `json:"key" yaml:"key"`
	Scope        string `json:"scope" yaml:"scope"`
	Addr         string `json:"addr" yaml:"addr"`
	AppId        string `json:"app_id" yaml:"app_id"`
	AppSecret    string `json:"app_secret" yaml:"app_secret"`
	AesKey       string `json:"aes_key" yaml:"aes_key"`
	Token        string `json:"token" yaml:"token"`
	CallbackAddr string `json:"callback_addr" yaml:"callback_addr"`
}
