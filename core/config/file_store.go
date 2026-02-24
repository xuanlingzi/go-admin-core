package config

type FileStore struct {
	Tencent *TencentFile `json:"tencent" yaml:"tencent"`
	Aliyun  *AliyunFile  `json:"aliyun" yaml:"aliyun"`
	Huawei  *HuaweiFile  `json:"huawei" yaml:"huawei"`
}

var FileStoreConfig = new(FileStore)
