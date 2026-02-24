package config

type BrokerConnectOptions struct {
	MintAddr     string `json:"mint_addr" yaml:"mint_addr"`
	CallbackAddr string `yaml:"callback_addr" json:"callback_addr"`
	ClientId     string `yaml:"client_id" json:"client_id"`
	FilePath     string `yaml:"file_path" json:"file_path"`
	Addr         string `yaml:"addr" json:"addr"`
}
