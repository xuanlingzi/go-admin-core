package config

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/url"
)

type TencentCos struct {
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Region    string `json:"region" yaml:"region"`
	CosUrl    string `json:"cos_url" yaml:"cos_url"`
	CiUrl     string `json:"ci_url" yaml:"ci_url"`
}

func (f *TencentCos) GetFileStoreOptions() (*cos.BaseURL, *cos.AuthorizationTransport) {
	cosUrl, err := url.Parse(f.CosUrl)
	if err != nil {
		return nil, nil
	}
	ciUrl, err := url.Parse(f.CiUrl)
	if err != nil {
		return nil, nil
	}
	base := &cos.BaseURL{
		BucketURL: cosUrl,
		CIURL:     ciUrl,
	}
	transport := &cos.AuthorizationTransport{
		SecretID:  f.AccessKey,
		SecretKey: f.SecretKey,
	}
	return base, transport
}
