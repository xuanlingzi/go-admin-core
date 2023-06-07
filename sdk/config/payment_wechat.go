package config

import (
	"crypto/rsa"
	"crypto/x509"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var _wechat_pay *core.Client

// GetWeChatPayClient 获取BlockChain客户端
func GetWeChatPayClient() *core.Client {
	return _wechat_pay
}

// SetWeChatPayClient 设置WeChat客户端
func SetWeChatPayClient(c *core.Client) {
	if _wechat_pay != nil && _wechat_pay != c {
		_wechat_pay = nil
	}
	_wechat_pay = c
}

type WeChatPayOption struct {
	MerchantId     string `yaml:"merchant_id" json:"merchant_id"`
	AppId          string `yaml:"app_id" json:"app_id"`
	PrivateKeyPath string `yaml:"private_key_path" json:"private_key_path"`
	SerialNumber   string `yaml:"serial_no" json:"serial_no"`
	ApiV3Key       string `yaml:"api_v3_key" json:"api_v3_key"`
	CallbackAddr   string `yaml:"callback_addr" json:"callback_addr"`
	WeChatCertPath string `yaml:"wechat_cert_path" json:"wechat_cert_path"`
}

func (e WeChatPayOption) GetPrivateKey() *rsa.PrivateKey {
	privateKey, err := utils.LoadPrivateKeyWithPath(e.PrivateKeyPath)
	if err != nil {
		return nil
	}
	return privateKey
}

func (e WeChatPayOption) GetWeChatCert() *x509.Certificate {
	certificate, err := utils.LoadCertificateWithPath(e.WeChatCertPath)
	if err != nil {
		return nil
	}
	return certificate
}

func (e WeChatPayOption) GetWeChatPayOption() *core.ClientOption {
	o := option.WithWechatPayAutoAuthCipher(e.MerchantId, e.SerialNumber, e.GetPrivateKey(), e.ApiV3Key)
	return &o
}
