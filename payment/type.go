package payment

import "time"

type AdapterPaymentService interface {
	String() string
	DownloadCertificates() error
	PrePay(orderId string, amount int64, payerOpenId string, attach string, description string, expireAt time.Time, callbackAddr string) (string, error)
	Refund(orderId string, transactionId string, refundId string, reason string, amount int64, total int64, currency string, callbackAddr string) (string, error)
}
