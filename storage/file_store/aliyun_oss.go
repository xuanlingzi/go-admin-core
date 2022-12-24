package file_store

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var _aliyun_oss *oss.Client

func GetAliyunOssClient() *oss.Client {
	return _aliyun_oss
}

type AliyunOssClient struct {
	client *oss.Client
	bucket string
}

func NewAliyunOss(client *oss.Client, endpoint, accessId, accessSecret, bucket string) *AliyunOssClient {
	var err error
	if client == nil {
		client, err = oss.New(endpoint, accessId, accessSecret)
		if err != nil {
			panic("Failed to setup Aliyun OSS client.")
		}
	}

	r := &AliyunOssClient{
		client: client,
		bucket: bucket,
	}
	return r
}

func (m *AliyunOssClient) String() string {
	return "aliyun_oss"
}

func (m *AliyunOssClient) Check() bool {
	return m.client != nil
}

func (m *AliyunOssClient) Close() {
}

func (m *AliyunOssClient) Put(ctx context.Context, name, fileLocation string) (string, error) {
	bucket, err := m.client.Bucket(m.bucket)
	if err != nil {
		return "", err
	}
	err = bucket.PutObjectFromFile(name, fileLocation, nil)
	if err != nil {
		return "", err
	}
	return name, nil
}

// GetClient 暴露原生client
func (m *AliyunOssClient) GetClient() interface{} {
	return m.client
}
