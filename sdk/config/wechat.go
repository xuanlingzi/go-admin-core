package config

type WeChat struct {
	Platforms *[]WeChatOption `json:"platforms" yaml:"platforms"`
	// ActivityPlatformKey 覆盖 activity 入口实际使用的平台标识，留空表示按入口原样解析。
	// 这是双服务号切换的灰度开关：填 rktv 即让 C 端回退到主服务号，
	// 改这一项 + 重启即可切换或回滚，无需回退代码。
	ActivityPlatformKey string `json:"activity_platform_key" yaml:"activity_platform_key"`
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
