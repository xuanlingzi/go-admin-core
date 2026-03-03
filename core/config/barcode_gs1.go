package config

type BarcodeGS1 struct {
	Addr        string `json:"addr" yaml:"addr"`
	AccessToken string `json:"access_token" yaml:"access_token"`
	PageSize    int    `json:"page_size" yaml:"page_size"`
	TimeoutSec  int    `json:"timeout_sec" yaml:"timeout_sec"`
}

var GS1Config = new(BarcodeGS1)
