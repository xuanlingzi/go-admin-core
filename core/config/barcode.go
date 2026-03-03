package config

type Barcode struct {
	GS1    *BarcodeGS1    `json:"gs1" yaml:"gs1"`
	Aliyun *BarcodeAliyun `json:"aliyun" yaml:"aliyun"`
}

var BarcodeConfig = new(Barcode)
