package config

type RtcZego struct {
	Addr             string `yaml:"addr" json:"addr"`
	AppID            uint32 `yaml:"app_id" json:"app_id"`
	ServerSecret     string `yaml:"server_secret" json:"server_secret"`
	AppSign          string `yaml:"app_sign" json:"app_sign"`
	SignatureVersion string `yaml:"signature_version" json:"signature_version"`
	TimeoutSec       int    `yaml:"timeout_sec" json:"timeout_sec"`
	TokenExpireSec   int64  `yaml:"token_expire_sec" json:"token_expire_sec"`
	CallbackSecret   string `yaml:"callback_secret" json:"callback_secret"`
	CallbackSkewSec  int64  `yaml:"callback_skew_sec" json:"callback_skew_sec"`
}
