package config

type Payment struct {
	WeChatPay *WeChatPayOption `yaml:"wechat" json:"wechat"`
	Leshua    *Leshua          `yaml:"leshua" json:"leshua"`
}

var PaymentConfig = new(Payment)
