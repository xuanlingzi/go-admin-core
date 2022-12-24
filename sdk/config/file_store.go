package config

type FileStore struct {
	TencentCos *TencentCos `json:"tencent_cos" yaml:"tencent_cos"`
	AliyunOss  *AliyunOss  `json:"aliyun_oss" yaml:"aliyun_oss"`
}

var FileStoreConfig = new(FileStore)
