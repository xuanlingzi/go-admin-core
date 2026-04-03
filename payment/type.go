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
	// 分账开通接口（聚合签名）
	ApplyLedger(merchantID string, sharePercent string, insertFlag int, feeRate int, ledgerMethod int, authTypes string, callbackUrl string) (map[string]interface{}, error)
	QueryLedgerStatus(merchantID string) (map[string]interface{}, error)
	// 分账关系接口（聚合签名）
	BindLedgerReceiver(merchantID1, merchantID2 string, protocolPic string, remark string) (map[string]interface{}, error)
	UnbindLedgerReceiver(merchantID1, merchantID2 string, remark string) (map[string]interface{}, error)
	QueryBindRelation(merchantID1 string) (map[string]interface{}, error)
	// 订单分账接口（聚合签名）
	ApplyOrderSplit(merchantID, leshuaOrderID, thirdOrderID, thirdRoyaltyID string, shareDetail []map[string]interface{}, remark string) (map[string]interface{}, error)
	QueryOrderSplit(merchantID, leshuaOrderID string, allRoyaltyFlag int, leshuaRoyaltyID, thirdRoyaltyID string) (map[string]interface{}, error)
	CancelOrderSplit(merchantID, leshuaOrderID, leshuaRoyaltyID, thirdRoyaltyID string) (map[string]interface{}, error)
	RefundOrderSplit(merchantID, thirdOrderID, leshuaOrderID, thirdRefundID string, refundAmount int64, refundMode string, thirdRoyaltyID string, refundDetails []map[string]interface{}, notifyUrl string) (map[string]interface{}, error)
	QueryRefundOrderSplit(merchantID, leshuaOrderID, thirdOrderID, thirdRefundID, leshuaRefundID string) (map[string]interface{}, error)
	GetClient() interface{}
}
