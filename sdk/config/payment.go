package config

type Payment struct {
	WeChatPay *WeChatPayOption `yaml:"wechat" json:"wechat"`
}

var PaymentConfig = new(Payment)
