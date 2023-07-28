package file_store

import (
	"context"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/xuanlingzi/go-admin-core/logger"
)

var _huaweiFile *obs.ObsClient

func GetHuaweiFileClient() *obs.ObsClient {
	return _huaweiFile
}

type HuaweiFileClient struct {
	client *obs.ObsClient
	bucket string
}

func NewHuaweiFile(client *obs.ObsClient, accessKey, secretKey, bucket, endpoint string) *HuaweiFileClient {
	var err error
	if client == nil {
		client, err = obs.New(accessKey, secretKey, endpoint)
		if err != nil {
			panic(err)
		}
	}

	r := &HuaweiFileClient{
		client: client,
		bucket: bucket,
	}

	return r
}

func (rc *HuaweiFileClient) String() string {
	return "huawei_file"
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
	response, err := rc.client.PutFile(&input)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			return fileLocation, fmt.Errorf("OBS Upload Error: %v, %v", obsError.StatusCode, obsError.Message)
		}
		return fileLocation, fmt.Errorf("OBS Upload Error: %v", err)
	}
	logger.Infof("OBS Upload Response: %v", response)

	return response.ObjectUrl, nil
}

// GetClient 暴露原生client
func (rc *HuaweiFileClient) GetClient() interface{} {
	return rc.client
}
