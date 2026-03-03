package config

// BarcodeAliyun 阿里云市场API配置（条码查询等）
type BarcodeAliyun struct {
	Addr       string `json:"addr" yaml:"addr"`               // 网关地址，例如 https://lhsptmxxcx.market.alicloudapi.com
	AppCode    string `json:"app_code" yaml:"app_code"`       // 阿里云市场 AppCode
	TimeoutSec int    `json:"timeout_sec" yaml:"timeout_sec"` // 请求超时秒数，默认 10
}
