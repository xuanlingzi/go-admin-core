package config

type Leshua struct {
	PaymentAddr    string `json:"payment_addr" yaml:"payment_addr"`
	TradeKey       string `json:"trade_key" yaml:"trade_key"`
	NotifyKey      string `json:"notify_key" yaml:"notify_key"`
	SignType       string `json:"sign_type" yaml:"sign_type"`
	NotifyURL      string `json:"notify_url" yaml:"notify_url"`
	CollectAddr    string `json:"collect_addr" yaml:"collect_addr"`
	CollectAgentID string `json:"collect_agent_id" yaml:"collect_agent_id"`
}

var LeshuaConfig = new(Leshua)
