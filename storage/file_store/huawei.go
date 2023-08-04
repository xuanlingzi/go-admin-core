package file_store

import (
	"context"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

var _huaweiFile *obs.ObsClient

func GetHuaweiFileClient() *obs.ObsClient {
	return _huaweiFile
}

type HuaweiFileClient struct {
	client   *obs.ObsClient
	accessId string
	bucket   string
	endpoint string
}

func NewHuaweiFile(client *obs.ObsClient, accessKey, secretKey, bucket, endpoint string) *HuaweiFileClient {
	var err error
	if client == nil {
		client, err = obs.New(accessKey, secretKey, endpoint)
		if err != nil {
			panic(fmt.Sprintf("Huawei file store init error: %s", err.Error()))
		}
		_huaweiFile = client
	}

	r := &HuaweiFileClient{
		client:   client,
		accessId: accessKey,
		bucket:   bucket,
		endpoint: endpoint,
	}

	return r
}

func (rc *HuaweiFileClient) String() string {
	return rc.accessId
}

func (rc *HuaweiFileClient) Check() bool {
	return rc.client != nil
}

func (rc *HuaweiFileClient) Close() {
	rc.client.Close()
}

func (rc *HuaweiFileClient) Upload(ctx context.Context, name, fileLocation string) (string, error) {
	input := obs.PutFileInput{}
	input.Bucket = rc.bucket
	input.Key = name
	input.SourceFile = fileLocation
	_, err := rc.client.PutFile(&input)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			return fileLocation, fmt.Errorf("OBS Upload Error: %v, %v", obsError.StatusCode, obsError.Message)
		}
		return fileLocation, fmt.Errorf("OBS Upload Error: %v", err)
	}

	return fmt.Sprintf("https://%s.%s/%s", rc.bucket, rc.endpoint, name), nil
}

// GetClient 暴露原生client
func (rc *HuaweiFileClient) GetClient() interface{} {
	return rc.client
}
