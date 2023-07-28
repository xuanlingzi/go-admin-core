package file_store

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var _aliyunFile *oss.Client

func GetAliyunFileClient() *oss.Client {
	return _aliyunFile
}

type AliyunFileClient struct {
	client *oss.Client
	bucket string
}

func NewAliyunFile(client *oss.Client, accessId, accessSecret, bucket, endpoint string) *AliyunFileClient {
	var err error
	if client == nil {
		client, err = oss.New(endpoint, accessId, accessSecret)
		if err != nil {
			panic(err)
		}
	}

	r := &AliyunFileClient{
		client: client,
		bucket: bucket,
	}
	return r
}

func (m *AliyunFileClient) String() string {
	return "aliyun_file"
}

func (m *AliyunFileClient) Check() bool {
	return m.client != nil
}

func (m *AliyunFileClient) Close() {
}

func (m *AliyunFileClient) Upload(ctx context.Context, name, fileLocation string) (string, error) {
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
func (m *AliyunFileClient) GetClient() interface{} {
	return m.client
}
