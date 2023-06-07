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
	log "github.com/xuanlingzi/go-admin-core/logger"
	"os"
	"path/filepath"
	"time"
)

type WeChatPay struct {
	client      *core.Client
	merchantId  string
	appId       string
	apiKey      string
	serialNo    string
	certPath    string
	privateKey  *rsa.PrivateKey
	certificate *x509.Certificate
}

func NewWeChatPay(client *core.Client, merchantId string, appId string, apiKey string, serialNo string, certPath string, privateKey *rsa.PrivateKey, certificate *x509.Certificate) (*WeChatPay, error) {
	var err error
	ctx := context.Background()
	if client == nil {
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
		client:      client,
		merchantId:  merchantId,
		appId:       appId,
		apiKey:      apiKey,
		serialNo:    serialNo,
		certPath:    certPath,
		privateKey:  privateKey,
		certificate: certificate,
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

func (m *WeChatPay) DownloadCertificates() error {
	var err error
	ctx := context.Background()
	dl, err := downloader.NewCertificateDownloaderWithClient(ctx, m.client, m.apiKey)
	if err != nil {
		return err
	}
	err = dl.DownloadCertificates(ctx)
	if err != nil {
		return err
	}
	for serialNo, certContent := range dl.ExportAll(ctx) {
		fileLocation := filepath.Join(m.certPath, fmt.Sprintf("wechatpay_%v.pem", serialNo))

		f, err := os.Create(fileLocation)
		if err != nil {
			return fmt.Errorf("创建证书文件`%v`失败：%v", fileLocation, err)
		}

		_, err = f.WriteString(certContent + "\n")
		if err != nil {
			return fmt.Errorf("写入证书到`%v`失败: %v", fileLocation, err)
		}

		log.Infof("写入证书到`%v`成功\n", fileLocation)
	}

	return nil
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
