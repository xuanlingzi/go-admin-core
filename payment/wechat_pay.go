package payment

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/cipher/decryptors"
	"github.com/wechatpay-apiv3/wechatpay-go/core/cipher/encryptors"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	log "github.com/xuanlingzi/go-admin-core/logger"
	"os"
	"path/filepath"
	"time"
)

var _wechat_pay *core.Client

type WeChatPay struct {
	client         *core.Client
	merchantId     string
	appId          string
	apiKey         string
	serialNo       string
	privateKeyPath string
	certPath       string
}

// GetWeChatPayClient 获取BlockChain客户端
func GetWeChatPayClient() *core.Client {
	return _wechat_pay
}

func NewWeChatPay(client *core.Client, merchantId string, appId string, apiKey string, serialNo string, privateKeyPath string, certPath string) (*WeChatPay, error) {
	ctx := context.Background()
	if client == nil {
		privateKey, err := utils.LoadPrivateKeyWithPath(privateKeyPath)
		if err != nil {
			return nil, err
		}

		client, err = core.NewClient(
			ctx,
			option.WithWechatPayAutoAuthCipher(merchantId, serialNo, privateKey, apiKey),
			option.WithWechatPayCipher(
				encryptors.NewWechatPayEncryptor(downloader.MgrInstance().GetCertificateVisitor(merchantId)),
				decryptors.NewWechatPayDecryptor(privateKey),
			),
		)
		if err != nil {
			return nil, err
		}
	}

	c := &WeChatPay{
		client:         client,
		merchantId:     merchantId,
		appId:          appId,
		apiKey:         apiKey,
		serialNo:       serialNo,
		privateKeyPath: privateKeyPath,
		certPath:       certPath,
	}
	return c, nil
}

// Close 关闭连接
func (m *WeChatPay) Close() {
	if m.client != nil {
		log.Info("Closing connection to WeChatPay server")
		m.client = nil
	}
}

func (*WeChatPay) String() string {
	return "wechat_pay"
}

func (m *WeChatPay) GetPrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := utils.LoadPrivateKeyWithPath(m.privateKeyPath)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func (m *WeChatPay) GetCertificates() (*x509.Certificate, error) {
	var err error
	ctx := context.Background()
	dl, err := downloader.NewCertificateDownloaderWithClient(ctx, m.client, m.apiKey)
	if err != nil {
		return nil, err
	}
	err = dl.DownloadCertificates(ctx)
	if err != nil {
		return nil, err
	}
	storePath, _ := filepath.Split(m.certPath)
	for serialNo, certContent := range dl.ExportAll(ctx) {
		m.certPath = filepath.Join(storePath, fmt.Sprintf("wechatpay_%v.pem", serialNo))

		f, err := os.Create(m.certPath)
		if err != nil {
			return nil, fmt.Errorf("创建证书文件`%v`失败：%v", m.certPath, err)
		}

		_, err = f.WriteString(certContent + "\n")
		if err != nil {
			return nil, fmt.Errorf("写入证书到`%v`失败: %v", m.certPath, err)
		}

		log.Infof("写入证书到`%v`成功\n", m.certPath)
	}

	certificate, err := utils.LoadCertificateWithPath(m.certPath)
	if err != nil {
		return nil, err
	}
	return certificate, nil
}

func (m *WeChatPay) PrePay(orderId string, amount int64, payerOpenId string, attach string, description string, expireAt time.Time, callbackAddr string) (string, error) {

	api := jsapi.JsapiApiService{Client: m.client}
	req := jsapi.PrepayRequest{
		Appid:       core.String(m.appId),
		Mchid:       core.String(m.merchantId),
		Description: core.String(description),
		OutTradeNo:  core.String(orderId),
		Attach:      core.String(attach),
		NotifyUrl:   core.String(callbackAddr),
		TimeExpire:  core.Time(expireAt),
		Amount: &jsapi.Amount{
			Total: core.Int64(amount),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(payerOpenId),
		},
	}
	resp, _, err := api.Prepay(context.Background(), req)
	if err != nil {
		e := err.(*core.APIError)
		return "", fmt.Errorf("微信支付预支付错误, %s", e.Message)
	}
	log.Infof("微信支付预支付成功, %s", *resp.PrepayId)

	return *resp.PrepayId, nil
}

func (m *WeChatPay) Refund(orderId string, transactionId string, refundId string, reason string, amount int64, total int64, currency string, callbackAddr string) (string, error) {

	api := refunddomestic.RefundsApiService{Client: m.client}
	req := refunddomestic.CreateRequest{
		TransactionId: core.String(transactionId),
		OutTradeNo:    core.String(orderId),
		OutRefundNo:   core.String(refundId),
		Reason:        core.String(reason),
		NotifyUrl:     core.String(callbackAddr),
		FundsAccount:  refunddomestic.REQFUNDSACCOUNT_AVAILABLE.Ptr(),
		Amount: &refunddomestic.AmountReq{
			Currency: core.String(currency),
			Refund:   core.Int64(amount),
			Total:    core.Int64(total),
		},
	}
	resp, _, err := api.Create(context.Background(), req)
	if err != nil {
		e := err.(*core.APIError)
		return "", fmt.Errorf("微信支付退款错误, %s", e.Message)
	}

	jsonString, err := json.Marshal(resp)
	log.Infof("微信支付退款成功, %s", *resp)

	return string(jsonString), nil
}

// GetClient 暴露原生client
func (m *WeChatPay) GetClient() *core.Client {
	return m.client
}
