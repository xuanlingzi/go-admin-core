package payment

import (
	"crypto/x509"
	"time"
)

type AdapterPaymentService interface {
	String() string
	GetCertificates() (*x509.Certificate, error)
	PrePay(orderId string, amount int64, payerOpenId string, attach string, description string, expireAt time.Time, callbackAddr string) (string, error)
	Refund(orderId string, transactionId string, refundId string, reason string, amount int64, total int64, currency string, callbackAddr string) (string, error)
	QueryOrder(orderId string) (string, error)
	GetClient() interface{}
}

type AdapterLeshuaService interface {
	String() string
	Close() error
	PayByAuthCode(merchantID string, thirdOrderID string, authCode string, amountFen int64, body string, shopNo string, posNo string, terminalType string, terminalID string, goodsDetail string, sceneInfo string) (map[string]string, error)
	QueryOrder(merchantID, thirdOrderID string) (map[string]string, error)
	CloseOrder(merchantID, thirdOrderID, leshuaOrderID string) (map[string]string, error)
	CollectTerminalID(merchantID, serialNum, existingDeviceID string) (string, error)
	DeregisterTerminal(merchantID, deviceID string) (string, error)
	QueryTerminalReport(merchantID, serialNum, deviceID string) (map[string]interface{}, error)
	VerifyNotifySign(payload map[string]string) error
	VerifyResponseSign(payload map[string]string) error
	ParseNotifyXML(raw []byte) (map[string]string, error)
	GetClient() interface{}
}
